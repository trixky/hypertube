import path from 'path';
import Cron from 'node-cron';
import rimraf from 'rimraf';
import env from '../env';
import { getUnusedFiles, markTorrentAsDeleted } from '../postgres/movies';

const DEFAULT_CRON = '0 * * * *'; // Defaults to every hour
const DEFAULT_INTERVAL = '30 days'; // Defaults to file older than 30 days

let cron = env.OLD_FILES_CRON ?? DEFAULT_CRON;
let interval = env.OLD_FILES_INTERVAL ?? DEFAULT_INTERVAL; // Defaults to file older than 30 days

export default function scheduleDeleteFiles() {
	// Handle invalid CRONs
	if (!Cron.validate(cron)) {
		cron = DEFAULT_CRON;
	}

	console.log('Schedule [delete-files]', cron, 'older than', interval);
	let running = false;
	return Cron.schedule(cron, async (now) => {
		if (running) return;
		running = true;
		console.log(`[${now}] Checking unused files`);
		try {
			const unused = await getUnusedFiles(interval);
			if (unused.rows.length > 0) {
				console.log('Deleting', unused.rows.length, 'old torrents');
				// Delete all old torrents
				for (const unusedTorrent of unused.rows) {
					if (unusedTorrent.file_path) {
						// Remove the root folder of the torrent
						// Handle files that may be nested and also avoid leaving empty folders
						// -- e.g ./.cache/movies/nested/file.mkv
						let upperFolder = unusedTorrent.file_path;
						let folder = unusedTorrent.file_path;
						// eslint-disable-next-line no-constant-condition
						while (true) {
							const relativeUpperFolder = path.dirname(folder);
							upperFolder = path.basename(relativeUpperFolder);
							if (upperFolder == 'movies' || upperFolder == '/') {
								break;
							}
							folder = relativeUpperFolder;
						}
						if (upperFolder == 'movies') {
							// Check if they exist first to just ignore already deleted files
							try {
								// Delete the folder and all of it's files
								await new Promise((resolve, reject) => {
									let fullPath = path.resolve(folder);

									// If the file is a transcoded output
									// -- remove the additional extension from the path
									// -- to delete all files
									if (fullPath.endsWith('.webm')) {
										const parts = fullPath.split('.').slice(0, -1);
										fullPath = `${parts.join('.')}*`;
									}
									console.log('Deleting file or folder', fullPath);
									rimraf(fullPath, (error) => {
										if (error) {
											console.log(
												`Failed to delete folder ${folder}, it either does not exists or rights are missing.`
											);
											reject(error);
										} else {
											resolve(true);
										}
									});
								});

								// Delete subtitles
								await new Promise((resolve, reject) => {
									const fullPath = path.resolve(`./.cache/subtitles/${unusedTorrent.id}`);
									console.log('Deleting file or folder', fullPath);
									rimraf(fullPath, (error) => {
										if (error) {
											console.log(
												`Failed to delete folder ${folder}, it either does not exists or rights are missing.`
											);
											reject(error);
										} else {
											resolve(true);
										}
									});
								});
							} catch (error) {
								// Folder does not exists or is not deletable
							}
						}
					}
					await markTorrentAsDeleted(unusedTorrent.id);
				}
			}
		} catch (error) {
			// Revert to default interval on failure, to make sure it's not the error
			interval = DEFAULT_INTERVAL;
			console.error('Failed to check or delete unused files', error);
			console.log('Schedule [delete-files]', cron, 'older than', interval);
		}
		running = false;
	});
}
