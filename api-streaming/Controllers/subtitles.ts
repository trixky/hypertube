import { mkdir, stat, writeFile } from 'fs/promises';
import { resolve } from 'path';
import { Router } from 'express';
import { parseSync, stringifySync } from 'subtitle';
import { getTorrent } from '../postgres/torrents';
import {
	createTorrentSubtitle,
	deleteTorrentSubtitle,
	getMediaInformations,
	getTorrentSubtitle,
	getTorrentSubtitles,
	MediaInformations
} from '../postgres/subtitles';
import osdb, { OSDBSubtitle, SearchInformations } from '../lib/osdb';
import { CACHE_PATH_SUBTITLES } from '../lib/cache';

const router = Router();

router.get('/torrent/:torrentId/subtitles', async function (req, res) {
	// Sanitize torrentId
	const torrentId = parseInt(req.params.torrentId);
	if (isNaN(torrentId) || torrentId < 1) {
		return res.status(400).send({ error: 'Invalid torrent ID' });
	}

	// ... check that the torrent exists
	const torrent = await getTorrent(torrentId);
	if (!torrent) {
		return res.status(404).send({ error: 'The torrent does not exists' });
	}
	// IF the torrent has no Media attached to it, abort the request with an empty response
	if (!torrent.media_id) {
		return res.status(200).send({ subtitles: [] });
	}

	// Check if there already is subtitles for the torrent
	try {
		const subtitles = await getTorrentSubtitles(torrentId);
		if (subtitles.length > 0) {
			return res.status(200).send({
				subtitles: subtitles.map((subtitle) => ({
					id: parseInt(subtitle.id),
					lang: subtitle.lang
				}))
			});
		}
	} catch (err) {
		return res.status(500).send({ error: 'Failed to get subtitles' });
	}

	// Get the media informations to query OpenSubtitles
	let mediaInformations: MediaInformations & { name: string };
	try {
		mediaInformations = await getMediaInformations(torrent.media_id);
	} catch (err) {
		return res.status(404).send({ error: 'No media found for the torrent' });
	}

	// ... avoid sending requests to OSDB if we're blocked (quota exceeded)
	if (!osdb.resolveBlock) {
		return res.status(200).send({ subtitles: [], error: 'Quota exceeded' });
	}

	// ... and search OpenSubtitles
	const searchParams: SearchInformations = {
		imdb_id: mediaInformations.imdb_id,
		tmdb_id: mediaInformations.tmdb_id,
		year: mediaInformations.year,
		media_name: mediaInformations.name,
		torrent_name: torrent.name
	};
	console.log('Searching subtitles for torrent', torrent.id, searchParams);
	let osResults: false | OSDBSubtitle[];
	try {
		osResults = await osdb.search(searchParams);
		if (osResults === false) {
			return res.status(200).send({ subtitles: [] });
		}
	} catch (error) {
		return res.status(500).send({ error: 'Failed to fetch subtitles' });
	}

	// If there is results, select, download them and uncompress them
	const collected: { id: number; lang: string }[] = [];
	const langs = ['en', 'fr'];
	for (const lang of langs) {
		const firstResult = osResults.find((subtitle) => subtitle.attributes.language == lang);
		if (firstResult) {
			console.log('downloading lang', lang, 'for', torrent.name);
			let downloadedSubtitle:
				| false
				| {
						extension: string;
						buffer: string;
				  } = false;
			try {
				downloadedSubtitle = await osdb.download({
					fileId: firstResult.attributes.files[0].file_id,
					torrentId: torrent.id,
					lang
				});
			} catch (error) {
				console.error('Failed to download subtitles for lang', lang, torrent.name);
				continue;
			}
			// If the download succeeded, convert it and save it
			if (downloadedSubtitle) {
				try {
					const convertedSubtitle = stringifySync(parseSync(downloadedSubtitle.buffer), {
						format: 'WebVTT'
					});
					const folder = `${CACHE_PATH_SUBTITLES}${torrent.id}`;
					const path = `${folder}/${lang}.vtt`;
					await mkdir(folder, { recursive: true });
					await writeFile(path, convertedSubtitle);
					const created = await createTorrentSubtitle({
						torrent_id: torrent.id,
						lang,
						path
					});
					collected.push({ id: parseInt(created.rows[0].id), lang });
				} catch (error) {
					console.error('Failed to save subtitles for lang', lang, torrent.name);
				}
			} else if (osdb.blocked) {
				break;
			}
		}
	}

	return res.status(200).send({ subtitles: collected });
});

router.get('/subtitles/:subtitleId', async function (req, res) {
	// Sanitize subtitleId
	const subtitleId: number = parseInt(req.params.subtitleId);

	// ... check that the subtitle exists
	const subtitle = await getTorrentSubtitle(subtitleId);
	if (!subtitle) {
		return res.status(404).send({ error: 'The subtitle does not exists' });
	}

	// Delete files that are saved but doesn't exist
	try {
		await stat(resolve(subtitle.path));
	} catch (error) {
		deleteTorrentSubtitle(subtitle.id);
		return res.status(404).send({ error: 'The subtitle does not exists' });
	}

	// Read and return the file
	res.contentType('text/vtt');
	return res.sendFile(resolve(subtitle.path));
});

export default router;
