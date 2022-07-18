import { Router } from "express";
import { parse, parseSync, resync, stringify, stringifySync } from "subtitle";
import osdb, { SearchInformations } from "../lib/osdb";
import { getTorrent, Torrent } from "../postgres/movies";
import {
	createTorrentSubtitle,
	getMediaInformations,
	getTorrentSubtitle,
	getTorrentSubtitles,
	MediaInformations,
	Subtitle,
} from "../postgres/subtitles";
import { mkdir, writeFile } from "fs/promises";

const router = Router();

router.get("/torrent/:torrentId/subtitles", async function (req, res) {
	// Sanitize torrentId
	const torrentId = parseInt(req.params.torrentId);
	if (isNaN(torrentId) || torrentId < 1) {
		return res.status(400).send({ error: "Invalid torrent ID" });
	}

	// ... check that the torrent exists
	let torrent: Torrent;
	try {
		torrent = await getTorrent(torrentId);
		// IF the torrent has no Media attached to it, abort the request with an empty response
		if (!torrent.media_id) {
			return res.status(200).send({ subtitles: [] });
		}
	} catch (err) {
		return res.status(404).send({ error: "The torrent does not exists" });
	}

	// Check if there already is subtitles for the torrent
	const subtitles = await getTorrentSubtitles(torrentId);
	if (subtitles.length > 0) {
		return res.status(200).send({
			subtitles: subtitles.map((subtitle) => ({
				id: parseInt(subtitle.id),
				lang: subtitle.lang,
			})),
		});
	}

	// Get the media informations to query OpenSubtitles
	let mediaInformations: MediaInformations & { name: string };
	try {
		mediaInformations = await getMediaInformations(torrent.media_id);
	} catch (err) {
		return res.status(404).send({ error: "No media found for the torrent" });
	}

	// ... avoid sending requests to OSDB if we're blocked (quota exceeded)
	if (!osdb.resolveBlock) {
		return res.status(200).send({ subtitles: [], error: "Quota exceeded" });
	}

	// ... and search OpenSubtitles
	const searchParams: SearchInformations = {
		imdb_id: mediaInformations.imdb_id,
		tmdb_id: mediaInformations.tmdb_id,
		year: mediaInformations.year,
		media_name: mediaInformations.name,
		torrent_name: torrent.name,
	};
	console.log("Searching subtitles for torrent", torrent.id, searchParams);
	const osResults = await osdb.search(searchParams);
	if (osResults === false) {
		return res.status(200).send({ subtitles: [] });
	}

	// If there is results, select, download them and uncompress them
	const collected: { id: number; lang: string }[] = [];
	const langs = ["en", "fr"];
	for (const lang of langs) {
		const firstResult = osResults.find((subtitle) => subtitle.attributes.language == lang);
		if (firstResult) {
			console.log("downloading lang", lang, "for", torrent.name);
			const downloadedSubtitle = await osdb.download({
				fileId: firstResult.attributes.files[0].file_id,
				torrentId: torrent.id,
				lang,
			});
			// If the download succeeded, convert it and save it
			if (downloadedSubtitle) {
				const convertedSubtitle = stringifySync(parseSync(downloadedSubtitle.buffer), {
					format: "WebVTT",
				});
				const path = `./.cache/subtitles/${torrent.id}/${lang}.vtt`;
				await mkdir(path.split("/").slice(0, -1).join("/"), { recursive: true });
				await writeFile(path, convertedSubtitle);
				const created = await createTorrentSubtitle({
					torrent_id: torrent.id,
					lang,
					path,
				});
				collected.push({ id: parseInt(created.rows[0].id), lang });
			} else if (osdb.blocked) {
				break;
			}
		}
	}

	return res.status(200).send({ subtitles: collected });
});

router.get("/subtitle/:subtitleId", async function (req, res) {
	// Sanitize subtitleId
	const subtitleId: number = parseInt(req.params.subtitleId);

	// ... check that the subtitle exists
	let subtitle: Subtitle;
	try {
		subtitle = await getTorrentSubtitle(subtitleId);
	} catch (err) {
		return res.status(404).send({ error: "The subtitle does not exists" });
	}

	// Read and return the file
	res.contentType("text/vtt");
	return res.sendFile(subtitle.path);
});

export default router;
