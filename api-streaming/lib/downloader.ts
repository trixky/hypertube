import torrentStream from 'torrent-stream'; // https://github.com/mafintosh/torrent-stream#readme
import { getTorrent, updateTorrent } from '../postgres/movies';
import parseTorrent from 'parse-torrent';
import ffmpeg from 'fluent-ffmpeg';
import { sleep } from '../utils/time';

const cache_path = './.cache';
const cache_path_movies = cache_path + '/movies/';

const EXTENSION_mkv = 'mkv';
const EXTENSION_mp4 = 'mp4';
const EXTENSION_webm = 'webm';

const TRANSCODE_OUTPUT = EXTENSION_webm;

interface DBTorrent {
	torrent_url: string | null;
	magnet: string | null;
	file_path: string | null;
	downloaded: boolean | null;
	length: number | null;
}

interface LocalTorrent {
	file_path: string;
	downloaded: boolean;
	original_extension: string | null;
	corrupted: boolean;
	length: number;
	started: boolean;
	download_complete: boolean;
	ffmpeg_closed: boolean;
}

type TorrentFile = TorrentStream.TorrentFile & {
	original_extension: string | null;
};

export interface DownloadInfo {
	path: string | null;
	downloaded: boolean;
	original_extension: string;
	length: number;
}

const local_torrents = new Map<number, LocalTorrent>();

function get_extension(path: string): string | null {
	const extension = path.split('.').pop();

	if (extension != EXTENSION_mkv && extension != EXTENSION_mp4 && extension != EXTENSION_webm)
		return null;

	return extension;
}

function generate_full_path(file_path: string, block_extention_adding: boolean): string {
	return (
		cache_path_movies +
		file_path +
		(get_extension(file_path) === EXTENSION_mkv && !block_extention_adding
			? '.' + TRANSCODE_OUTPUT
			: '')
	);
}

function sanitize_torrent_id(torrent_id: string): number {
	const sanitized_torrent_id = parseInt(torrent_id);

	if (isNaN(sanitized_torrent_id)) throw new Error('torrent id corrupted');
	if (sanitized_torrent_id < 0) throw new Error('torrent id need to be positive');

	return sanitized_torrent_id;
}

function torrent_to_magnet(torrent_url: string): Promise<string> {
	return new Promise((resolve) => {
		parseTorrent.remote(torrent_url, (err, parsedTorrent) => {
			if (err) throw err;
			if (!parsedTorrent) throw new Error('Failed to parse torrent');
			resolve(parseTorrent.toMagnetURI(parsedTorrent));
		});
	});
}

async function get_magnet(db_torrent_infos: DBTorrent): Promise<string> {
	let magnet = db_torrent_infos.magnet;

	// extract or generate the magnet
	if (magnet == null) {
		if (db_torrent_infos.torrent_url === null)
			throw new Error('no magnet or torrent url available for this torrent');
		magnet = await torrent_to_magnet(db_torrent_infos.torrent_url);
	}

	return magnet;
}

function get_movie_file_from_engine(engine: TorrentStream.TorrentEngine): Promise<TorrentFile> {
	return new Promise((resolve, reject) => {
		engine.on('ready', async function () {
			engine.files.every((file) => {
				const original_extension = get_extension(file.name);
				if (original_extension != null) {
					const torrentFile = file as TorrentFile;
					torrentFile.original_extension = original_extension;
					resolve(torrentFile);
					return false;
				}
				return true;
			});
			reject();
		});
	});
}

function start_download_mp4_or_webm(
	engine: TorrentStream.TorrentEngine,
	selected_file: TorrentFile
) {
	selected_file.createReadStream();
}

function start_download_mkv(
	engine: TorrentStream.TorrentEngine,
	torrent_id: number,
	selected_file: TorrentFile
) {
	const started_saved = false;
	const local_file_path = generate_full_path(selected_file.path, true);
	const local_file_path_webm = local_file_path + '.' + TRANSCODE_OUTPUT;

	const stream = selected_file.createReadStream();

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
		// .videoBitrate(1)
		// .outputOptions("-preset veryfast")
		// .outputOptions("-crf 50")
		.outputOptions('-movflags +frag_keyframe+separate_moof+omit_tfhd_offset+empty_moov')
		.outputFormat('webm')
		// *
		// .outputOptions("-vf scale=-1:101")
		// .outputOptions("-qp 0")
		.outputOptions('-crf 18')
		.on('ffmpeg: start', () => {
			console.log('start');
		})
		.on('progress', (progress: { timemark: string }) => {
			if (started_saved === false) {
				// notify the torrent start to be readable
				local_torrents.set(torrent_id, <LocalTorrent>{
					file_path: selected_file.path,
					downloaded: false,
					original_extension: selected_file.original_extension,
					corrupted: false,
					length: selected_file.length,
					started: true,
					ffmpeg_closed: false
				});
			}
			console.log(`ffmpeg: progress > ${progress.timemark} for ${local_file_path_webm}`);
		})
		.on('end', () => {
			console.log('ffmpeg: Finished processing');

			local_torrents.set(torrent_id, <LocalTorrent>{
				file_path: selected_file.path,
				downloaded: true,
				original_extension: selected_file.original_extension,
				corrupted: false,
				length: selected_file.length,
				started: true,
				ffmpeg_closed: true
			});
			updateTorrent(torrent_id, selected_file.path, true, selected_file.length);
			// Destroy the engine only when the transcode is completed
			// -- to avoid killing the stream while transcoding
			if (local_torrents.get(torrent_id)?.download_complete) {
				engine.destroy(() => {
					//
				});
			}
		})
		.on('error', (err: Error) => {
			console.log(`ffmpeg: ERROR > ${err.message}`);
			// Don't forget to destroy the engine if the download was already completed
			if (local_torrents.get(torrent_id)?.download_complete) {
				engine.destroy(() => {
					//
				});
			}
			local_torrents.delete(torrent_id);
		})
		.output(local_file_path_webm)
		.run();
}

async function wait_file_start_to_download(torrent_id: number): Promise<unknown> {
	let started: boolean | undefined = false;

	while (started === undefined || started === false) {
		started = local_torrents.get(torrent_id)?.started;

		await sleep(1000);
	}

	return;
}

export async function download(torrent_id: string, want_mkv: boolean): Promise<DownloadInfo> {
	let need_transcode = false;
	const started_saved = false;

	// eslint-disable-next-line no-async-promise-executor
	return new Promise(async (resolve, reject) => {
		// sanitize the torrent id
		let sanitized_torrent_id: number;
		try {
			sanitized_torrent_id = sanitize_torrent_id(torrent_id);
		} catch (err) {
			reject(err);
			return null;
		}

		// check if torrent is known in local
		console.log('############################ 1.1');
		if (local_torrents.has(sanitized_torrent_id)) {
			const local_torrent_info = local_torrents.get(sanitized_torrent_id);

			// check if torrent is corrupted in local
			if (local_torrent_info == undefined || local_torrent_info.corrupted) {
				reject();
				return null;
			}

			// wait the torrents start to download
			await wait_file_start_to_download(sanitized_torrent_id);
			resolve(<DownloadInfo>{
				path: generate_full_path(local_torrent_info.file_path, want_mkv),
				downloaded: local_torrent_info.downloaded,
				original_extension: get_extension(local_torrent_info.file_path),
				length: local_torrent_info.length
			});
			return;
		}

		// set/block the torrent in local
		console.log('############################ 1.2');
		local_torrents.set(sanitized_torrent_id, <LocalTorrent>{
			file_path: '',
			original_extension: null,
			downloaded: false,
			corrupted: false,
			length: 0,
			started: false,
			ffmpeg_closed: false
		});

		// get the torrent infos from db
		console.log('############################ 1.3');
		const db_torrent_infos: DBTorrent = <DBTorrent>{};
		try {
			const res = await getTorrent(sanitized_torrent_id);
			db_torrent_infos.downloaded = res.downloaded;
			db_torrent_infos.file_path = res.file_path;
			db_torrent_infos.magnet = res.magnet;
			db_torrent_infos.torrent_url = res.torrent_url;
		} catch (err) {
			reject(err);
			return;
		}

		// check if the torrent is already downloaded
		console.log('############################ 1.4');
		if (db_torrent_infos.downloaded) {
			if (db_torrent_infos.file_path == null) {
				reject(new Error('no path for the downloaded movie from db'));
				return;
			}

			resolve(<DownloadInfo>{
				path: generate_full_path(db_torrent_infos.file_path, want_mkv),
				downloaded: db_torrent_infos.downloaded,
				original_extension: get_extension(db_torrent_infos.file_path),
				length: db_torrent_infos.length
			});
			return;
		}

		// get the magnet of the torrent
		console.log('############################ 1.5');
		let magnet: string;
		try {
			magnet = await get_magnet(db_torrent_infos);
		} catch {
			// if the torrent is corrupted
			local_torrents.set(sanitized_torrent_id, <LocalTorrent>{
				file_path: '',
				original_extension: null,
				downloaded: false,
				corrupted: true,
				length: 0,
				started: false,
				ffmpeg_closed: false
			});
			reject(new Error('magnet/torrent corrupted'));
			return;
		}

		// start the torrent engine
		console.log('############################ 1.6');
		const engine = torrentStream(magnet, {
			connections: 100,
			uploads: 1,
			tmp: cache_path,
			path: cache_path_movies,
			verify: true,
			tracker: true // false ?
		});

		// select the good file in the torrent (using extensions)
		console.log('############################ 1.7');
		let movie_file: TorrentFile;
		try {
			movie_file = await get_movie_file_from_engine(engine);

			(async () => {
				await wait_file_start_to_download(sanitized_torrent_id);

				resolve(<DownloadInfo>{
					path: generate_full_path(movie_file.path, want_mkv),
					downloaded: false,
					original_extension: get_extension(movie_file.path),
					length: movie_file.length
				});
			})();
		} catch {
			engine.destroy(() => {
				//
			});
			reject(new Error('no movie found in the torrent'));
			return;
		}

		// save torrent/movie in local
		console.log('############################ 1.8');
		local_torrents.set(sanitized_torrent_id, <LocalTorrent>{
			file_path: movie_file.path,
			original_extension: movie_file.original_extension,
			downloaded: false,
			corrupted: false,
			length: movie_file.length,
			started: false,
			ffmpeg_closed: movie_file.original_extension != EXTENSION_mkv
		});

		// start download
		console.log('############################ 1.9');
		if (movie_file.original_extension == EXTENSION_mkv) {
			need_transcode = true;
			start_download_mkv(engine, sanitized_torrent_id, movie_file);
		} else {
			start_download_mp4_or_webm(engine, movie_file);
		}

		engine.on('download', (index: string) => {
			// notify the torrent start to be readable
			if (started_saved === false && need_transcode === false)
				local_torrents.set(sanitized_torrent_id, <LocalTorrent>{
					file_path: movie_file.path,
					original_extension: movie_file.original_extension,
					downloaded: false,
					corrupted: false,
					length: movie_file.length,
					started: true,
					ffmpeg_closed: local_torrents.get(sanitized_torrent_id)!.ffmpeg_closed
				});
			console.log(`Download piece ${index} for ${movie_file.path}`);
		});

		engine.on('idle', async () => {
			console.log('############################ 1.10');
			const localTorrent = local_torrents.get(sanitized_torrent_id);
			if (localTorrent) {
				localTorrent.download_complete = true;
				if (!need_transcode) {
					await updateTorrent(
						sanitized_torrent_id,
						localTorrent.file_path,
						true,
						localTorrent.length
					);
				}
			}
			if (!localTorrent || localTorrent.ffmpeg_closed) {
				engine.destroy(() => {
					//
				});
			}
		});
	});
}
