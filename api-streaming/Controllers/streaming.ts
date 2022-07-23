import { Router } from 'express';
import TailFile from '@logdna/tail-file';
import fs from 'fs';
import { download, DownloadInfo, localTorrents } from '../lib/downloader';
import pump from 'pump';
import {
	deleteTorrent,
	getTorrent,
	refreshTorrentLastAccess,
	updateTorrent
} from '../postgres/torrents';
import { Readable } from 'stream';
import { stat } from 'fs/promises';

const router = Router();

router.get('/torrent/:torrentId/stream', async function (req, res) {
	// Sanitize Torrent id
	const torrentId: number = parseInt(req.params.torrentId);
	if (isNaN(torrentId) || torrentId < 0) {
		res.status(400).send({ error: 'Invalid torrent id' });
	}

	// Check if the torrent exists
	const torrent = await getTorrent(torrentId);
	if (torrent === undefined) {
		return res.status(404).send({ error: 'The torrent does not exists' });
	}

	// Update torrent last access to delay deletion
	try {
		await refreshTorrentLastAccess(torrentId);
	} catch (err) {
		return res.status(500).send({ error: 'Internal server error' });
	}

	// Check the Accept header to know if we can send mkv to the client if he can handle it
	const acceptHeader = req.headers.accept;
	const acceptMkv =
		(acceptHeader && (acceptHeader.indexOf('video/mkv') >= 0 || acceptHeader == '*/*')) == true;

	// Download the torrent
	let downloadResult: DownloadInfo;
	try {
		downloadResult = await download(torrent, acceptMkv);
	} catch (error) {
		console.error('Failed to download torrent', error);
		const corruptedMessages = [
			'Magnet or torrent corrupted',
			'No movie found in the torrent',
			'Failed to parse torrent',
			'No magnet or torrent url available for this torrent'
		];
		// ! If the torrent is corrupted or invalid, delete it to avoid it being displayed for the users
		if (error instanceof Error && corruptedMessages.indexOf(error.message) >= 0) {
			console.log(`Deleting corrupt Torrent #${torrent.id}`);
			await deleteTorrent(torrent.id);
		}
		return res.status(500).send({ error: 'Failed to download torrent' });
	}
	if (downloadResult.path == null) {
		return res.status(500).send();
	}

	// Extract some informations from the result to return the file
	const extension = downloadResult.path.split('.').pop()!;
	const nativeExtensions = ['mp4', 'webm'];
	const sendNative =
		(nativeExtensions.indexOf(downloadResult.originalExtension) >= 0 ||
			(downloadResult.originalExtension == 'mkv' && acceptMkv)) &&
		// HEVC and x265 files can't be read on Chromium browsers
		// -- So even if the client accept mkv it can't display it
		!torrent.name.match(/hevc|x\s?265/i) &&
		downloadResult.length > 0;

	/*console.log(
		'acceptMkv',
		acceptMkv,
		'extension',
		extension,
		'downloadResult.originalExtension',
		downloadResult.originalExtension,
		'downloadResult.downloaded',
		downloadResult.downloaded,
		'sendNative',
		sendNative,
		'using path',
		downloadResult.path,
		'has stream ?',
		downloadResult.stream ? 'yes' : 'no'
	);*/

	// * Downloading
	// * This can either be the original file if the torrent has an accepted file for the client
	// * -- or a transcode result in webm format
	// * The result is an used file path in either way, which is tailed
	if (!downloadResult.downloaded) {
		let stream: Readable | TailFile;
		if (sendNative && downloadResult.stream) {
			console.log('Using TorrentStream file stream');
			stream = downloadResult.stream;
		} else {
			console.log('Opening tail stream for file still being downloaded');
			// TODO: Bigger idle timeout ?
			stream = new TailFile(downloadResult.path, { startPos: 0 })
				.on('tail_error', (err) => {
					console.error('TailFile had an error!', err);
				})
				.on('error', (err) => {
					console.error('A TailFile stream error was likely encountered', err);
				});
			(stream as TailFile).start();
		}

		// When sending the file that existed in the torrent we already have all of the size informations ...
		if (sendNative) {
			console.log('Sending native torrent file');
			const realSize = downloadResult.length;
			console.log('Sending headers', {
				'Content-Length': realSize,
				// "Content-Range": `bytes 0-${realSize - 1}/${realSize}`,
				// "Accept-Ranges": "bytes",
				'Content-Type': `video/${downloadResult.originalExtension}`
			});
			res.writeHead(200, {
				'Content-Length': realSize,
				// "Content-Range": `bytes 0-${realSize - 1}/${realSize}`,
				// "Accept-Ranges": "bytes",
				'Content-Type': `video/${downloadResult.originalExtension}`
			});
		}
		// ... else the file comes from a transcode and we don't know the final file size
		else {
			console.log('Sending headers', {
				// "Content-Range": `bytes 0-`,
				// "Accept-Ranges": "bytes",
				'Content-Type': `video/${extension}`
			});
			res.writeHead(200, {
				// "Content-Range": `bytes 0-`,
				// "Accept-Ranges": "bytes",
				'Content-Type': `video/${extension}`
			});
		}

		// Send the file read to the response
		console.log('Sending file still being downloaded');
		// Avoid closing torrent-stream stream, only close tailed streams
		// -- since it also mark the torrent as "completed" if closed
		if ('quit' in stream) {
			return pump(stream, res, () => {
				if ('quit' in stream) {
					stream.quit();
				}
			});
		} else {
			return stream.pipe(res);
		}
	}
	// * File is completed already
	// * Accept ranges on the file
	else {
		// Check if the completed file exists
		try {
			await stat(downloadResult.path);
		} catch (error) {
			// -- else delete the entry from the database and sadly return an error
			console.error("Completed file doesn't exists");
			await updateTorrent(torrentId, null, false, null);
			return res.status(500).send({ error: 'Completed file was deleted, refresh the player' });
		}

		const size = downloadResult.length;

		// Check and support ranges if the client supports them
		const range = req.headers['range'];
		let start = 0;
		let end = size - 1;
		// Handle invalid ranges with a 416
		console.log('Sending completed torrent with ranges', range);
		if (range) {
			const parts = range!.slice(6).split('-');
			if (parts.length >= 1 && parts.length < 3) {
				const partialStart = parseInt(parts[0]);
				const partialEnd = parseInt(parts[1]);
				if (!isNaN(partialStart) && partialStart >= 0 && partialStart < size) {
					start = partialStart;
				} else {
					console.error('Invalid range start', start, size);
					return res.status(416).header('Content-Range', `*/${size}`).send();
				}
				if (!isNaN(partialEnd)) {
					if (partialEnd > start && partialEnd <= size) {
						end = partialEnd;
					} else {
						console.error('Invalid range end');
						return res.status(416).header('Content-Range', `*/${size}`).send();
					}
				}
			} else {
				console.error('Invalid ranges format');
				return res.status(416).header('Content-Range', `*/${size}`).send();
			}
		}

		// If the browser has no support for range we just returns a 200 with no ranges
		let returnCode = 200;
		const headers: Record<string, string | number> = {
			'Content-Length': size,
			'Content-Type': `video/${extension}`
		};
		if (range) {
			returnCode = 206;
			headers['Content-Range'] = `bytes ${start}-${end - 1}/${size}`;
			headers['Accept-Ranges'] = 'bytes';
		}
		console.log('Sending headers', headers);
		res.writeHead(returnCode, headers);

		// Send the file read to the response
		console.log('Sending completed torrent file');
		const file = fs.createReadStream(downloadResult.path, { start, end });
		pump(file, res);
	}
});

// Check the torrent status for
// -- completion status
// -- download progress
// -- transcode progress
router.get('/torrent/:torrentId/status', async function (req, res) {
	// Sanitize Torrent id
	const torrentId: number = parseInt(req.params.torrentId);
	if (isNaN(torrentId) || torrentId < 0) {
		res.status(400).send({ error: 'Invalid torrent id' });
	}

	// Check if the torrent exists
	const torrent = await getTorrent(torrentId);
	if (torrent === undefined) {
		return res.status(404).send({ error: 'The torrent does not exists' });
	}

	// If the torrent is already completed there is no informations to get
	if (torrent.downloaded) {
		return res.status(200).send({ status: 'complete' });
		// -- else check if there is a LocalTorrent instance, that means the torrent is being downloaded/transcoded
	} else if (localTorrents[torrentId]) {
		const localTorrent = localTorrents[torrentId];
		const response: {
			status: 'ongoing';
			download?: { completed: number; total: number };
			encoding?: { processed?: string; fps?: number; completeDuration: string };
		} = { status: 'ongoing' };
		if (localTorrent.movieFile && localTorrent.engine?.swarm) {
			response.download = {
				completed: localTorrent.engine.swarm.downloaded,
				total: localTorrent.movieFile.length
			};
		}
		if (localTorrent.ffmpegProgress) {
			response.encoding = localTorrent.ffmpegProgress;
		}
		return res.status(200).send(response);
	}

	// -- else if there is no LocalTorrent instance, the torrent is idle, no download or transcode
	return res.status(200).send({ status: 'idle' });
});

export default router;
