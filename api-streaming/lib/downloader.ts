import torrentStream from 'torrent-stream'; // https://github.com/mafintosh/torrent-stream#readme
import { Torrent, updateTorrent } from '../postgres/torrents';
import parseTorrent from 'parse-torrent';
import ffmpeg from 'fluent-ffmpeg';
import { sleep } from '../utils/time';
import { CACHE_PATH, CACHE_PATH_MOVIES } from './cache';
import { Readable } from 'stream';

const EXTENSION_mkv = 'mkv';
const EXTENSION_mp4 = 'mp4';
const EXTENSION_webm = 'webm';

const TRANSCODE_OUTPUT = EXTENSION_webm;

type TorrentFile = TorrentStream.TorrentFile & {
	extension: string | null;
};

interface LocalTorrent {
	file_path: string;
	downloaded: boolean;
	extension: string | null;
	corrupted: boolean;
	length: number;
	downloadStarted: boolean;
	downloadComplete: boolean;
	ffmpegClosed: boolean;
	movieFile?: TorrentFile;
}

export interface DownloadInfo {
	path: string | null;
	downloaded: boolean;
	originalExtension: string;
	length: number;
	stream?: Readable;
}

export const localTorrents: Record<number, LocalTorrent> = {};

function fileExtension(path: string): string {
	return path.split('.').pop()!;
}

function generateFullPath(file_path: string, transcoded: boolean): string {
	return (
		CACHE_PATH_MOVIES +
		file_path +
		(fileExtension(file_path) === EXTENSION_mkv && !transcoded ? '.' + TRANSCODE_OUTPUT : '')
	);
}

function convertTorrentToMagnet(torrent_url: string): Promise<string> {
	return new Promise((resolve) => {
		parseTorrent.remote(torrent_url, (err, parsedTorrent) => {
			if (err) throw err;
			if (!parsedTorrent) throw new Error('Failed to parse torrent');
			resolve(parseTorrent.toMagnetURI(parsedTorrent));
		});
	});
}

async function getMagnet(torrent: Torrent): Promise<string> {
	let magnet = torrent.magnet;

	// extract or generate the magnet
	if (magnet == null) {
		if (torrent.torrent_url === null)
			throw new Error('no magnet or torrent url available for this torrent');
		magnet = await convertTorrentToMagnet(torrent.torrent_url);
	}

	return magnet;
}

function selectMovieFromTorrent(engine: TorrentStream.TorrentEngine): Promise<TorrentFile> {
	return new Promise((resolve, reject) => {
		engine.on('ready', async function () {
			const validFiles = engine.files.filter((file) => {
				const extension = fileExtension(file.name);
				return (
					extension == EXTENSION_mkv || extension == EXTENSION_mp4 || extension == EXTENSION_webm
				);
			});
			validFiles.sort((a, b) => b.length - a.length);
			if (validFiles.length > 0) {
				const selectedFile = validFiles[0] as TorrentFile;
				selectedFile.extension = fileExtension(selectedFile.name);
				return resolve(selectedFile);
			}
			return reject();
		});
	});
}

function transcodeMovieFile(
	engine: TorrentStream.TorrentEngine,
	torrentId: number,
	movieFile: TorrentFile
) {
	const localFilePath = generateFullPath(movieFile.path, true);
	const localOutputPath = localFilePath + '.' + TRANSCODE_OUTPUT;

	const stream = movieFile.createReadStream();

	ffmpeg()
		.input(stream)
		.inputFormat('matroska')
		// * mp4
		// .audioCodec("aac")
		// .videoCodec("libx264")
		// .videoBitrate(1)
		// .outputOptions("-preset veryfast")
		// .outputOptions("-crf 50")
		// .outputOptions("-movflags +frag_keyframe+separate_moof+omit_tfhd_offset+empty_moov+faststart")
		// .outputFormat("mp4")
		// * webm
		.audioCodec('libvorbis')
		.videoCodec('libvpx')
		.outputOptions('-movflags +frag_keyframe+separate_moof+omit_tfhd_offset+empty_moov')
		.outputFormat('webm')
		// *
		// .outputOptions("-vf scale=-1:101")
		// .outputOptions("-qp 0")
		.outputOptions('-crf 14')
		.outputOptions('-b:v 5000K')
		// .videoBitrate(1)
		// .outputOptions('-preset veryfast')
		// .outputOptions('-crf 50')
		.outputOptions('-threads 4')
		.outputOptions('-flags low_delay')
		.on('ffmpeg: start', () => {
			console.log('start');
		})
		.on('progress', (progress: { timemark: string }) => {
			localTorrents[torrentId].downloadStarted = true;
			console.log(`ffmpeg: progress > ${progress.timemark} for ${localOutputPath}`);
		})
		.on('end', async () => {
			console.log('ffmpeg: Finished processing');

			localTorrents[torrentId].downloaded = true;
			localTorrents[torrentId].ffmpegClosed = true;

			await updateTorrent(torrentId, movieFile.path, true, movieFile.length);

			// Destroy the engine only when the transcode is completed
			// -- to avoid killing the stream while transcoding
			// -- The torrent *should* be complete since the transcode need the whole file
			if (localTorrents[torrentId].downloadComplete) {
				engine.destroy(() => {
					console.log('TorrentEngine destroyed');
				});
				delete localTorrents[torrentId];
			}
		})
		.on('error', (err: Error) => {
			console.log(`ffmpeg: ERROR > ${err.message}`);
			localTorrents[torrentId].ffmpegClosed = true;

			// Don't forget to destroy the engine if the download was already completed
			if (localTorrents[torrentId].downloadComplete) {
				engine.destroy(() => {
					console.log('TorrentEngine destroyed');
				});
			}
			delete localTorrents[torrentId];
		})
		.output(localOutputPath)
		.run();
}

async function waitForFileToExist(torrentId: number): Promise<unknown> {
	let downloadStarted: boolean | undefined = false;

	while (downloadStarted === undefined || downloadStarted === false) {
		if (!localTorrents[torrentId]) {
			throw new Error('Torrent failed before it started');
		}
		downloadStarted = localTorrents[torrentId].downloadStarted;
		await sleep(500);
	}

	return;
}

export async function download(torrent: Torrent, acceptMkv: boolean): Promise<DownloadInfo> {
	let needTranscode = false;

	// eslint-disable-next-line no-async-promise-executor
	return new Promise(async (resolve, reject) => {
		// Check if the torrent is already downloaded
		if (torrent.downloaded) {
			console.log('Resolving with completed Torrent');
			if (torrent.file_path == null || torrent.length == null) {
				reject(new Error('no path or size for the downloaded movie'));
				return;
			}

			return resolve(<DownloadInfo>{
				downloaded: true,
				path: generateFullPath(torrent.file_path, acceptMkv),
				originalExtension: fileExtension(torrent.file_path),
				length: parseInt(torrent.length)
			});
		}

		// Check if torrent is already known in local
		if (localTorrents[torrent.id]) {
			console.log('Resolving with LocalTorrent entry');
			const localTorrent = localTorrents[torrent.id];

			// check if torrent is corrupted in local
			if (localTorrent.corrupted) {
				reject();
				return null;
			}

			// wait the torrents start to download
			if (!localTorrent.downloadStarted) {
				console.log("LocalTorrent entry download hasn't started");
				await waitForFileToExist(torrent.id);
			}
			const originalExtension = fileExtension(localTorrent.file_path);
			return resolve(<DownloadInfo>{
				path: generateFullPath(localTorrent.file_path, acceptMkv),
				downloaded: localTorrent.downloaded,
				originalExtension: originalExtension,
				length: localTorrent.length,
				stream: localTorrent.movieFile?.createReadStream()
			});
		}

		// Create a local torrent entry
		console.log('Created LocalTorrent entry');
		localTorrents[torrent.id] = {
			file_path: '',
			extension: null,
			downloaded: false,
			corrupted: false,
			length: 0,
			downloadStarted: false,
			ffmpegClosed: false,
			downloadComplete: false
		};

		// get the magnet of the torrent
		let magnet: string;
		try {
			console.log('Using magnet or torrent');
			magnet = await getMagnet(torrent);
		} catch {
			console.log('Magnet or torrent corrupted');
			// If the torrent is corrupted
			localTorrents[torrent.id].corrupted = true;
			reject(new Error('magnet/torrent corrupted'));
			return;
		}

		// Start a Torrent engine
		const engine = torrentStream(magnet, {
			connections: 100,
			uploads: 1,
			tmp: CACHE_PATH,
			path: CACHE_PATH_MOVIES,
			verify: true,
			tracker: true // false ?
		});

		// Select the good file in the torrent (using extension and size)
		let movieFile: TorrentFile;
		try {
			movieFile = await selectMovieFromTorrent(engine);
		} catch {
			console.log('No movie found in the torrent');
			engine.destroy(() => {
				console.log('TorrentEngine destroyed');
			});
			reject(new Error('No movie found in the torrent'));
			return;
		}

		// Update local torrent state
		localTorrents[torrent.id].file_path = movieFile.path;
		localTorrents[torrent.id].extension = movieFile.extension;
		localTorrents[torrent.id].length = movieFile.length;
		localTorrents[torrent.id].ffmpegClosed = movieFile.extension != EXTENSION_mkv;

		// Start downloading and transcode if necessary
		if (movieFile.extension == EXTENSION_mkv) {
			console.log('Resolving download with transcode');
			needTranscode = true;
			transcodeMovieFile(engine, torrent.id, movieFile);

			await waitForFileToExist(torrent.id);
			resolve({
				path: generateFullPath(movieFile.path, acceptMkv),
				downloaded: false,
				originalExtension: fileExtension(movieFile.path),
				length: movieFile.length
			});
		}
		// A native view will simply return a stream of the download torrent file
		else {
			console.log('Resolving download with native file from torrent');
			resolve({
				path: generateFullPath(movieFile.path, acceptMkv),
				downloaded: false,
				originalExtension: fileExtension(movieFile.path),
				length: movieFile.length,
				stream: movieFile.createReadStream()
			});
		}

		engine.on('download', (index: string) => {
			// Notify the torrent start to be readable
			localTorrents[torrent.id].downloadStarted = true;
			console.log(`Download piece ${index} for ${movieFile.path}`);
		});

		// When the download is complete
		// -- Destroy the engine if no transcode is being done
		engine.on('idle', async () => {
			console.log('Torrent download complete', torrent.name);
			const localTorrent = localTorrents[torrent.id];
			if (localTorrent) {
				localTorrent.downloadComplete = true;
				// If no transcode was needed the torrent is updated with the final values
				if (!needTranscode) {
					console.log('Saving native torrent completion');
					await updateTorrent(torrent.id, localTorrent.file_path, true, localTorrent.length);
					delete localTorrents[torrent.id];
				}
			}
			// If the local torrent errored out or ffmpeg also errored out
			// -- destroy the torrent engine
			if (!localTorrent || localTorrent.ffmpegClosed) {
				console.log('Closing engine since ffmpeg is already closed');
				engine.destroy(() => {
					console.log('TorrentEngine destroyed');
				});
				delete localTorrents[torrent.id];
			}
		});
	});
}
