import { client } from './db';

export type Subtitle = {
	id: string;
	torrent_id: number;
	lang: string;
	path: string;
};

export async function getTorrentSubtitles(torrentId: number) {
	const res = await client.query<Subtitle>(
		'SELECT * FROM torrent_subtitles WHERE torrent_id = $1',
		[torrentId]
	);
	return res.rows;
}

export async function getTorrentSubtitle(id: number) {
	const res = await client.query<Subtitle>('SELECT * FROM torrent_subtitles WHERE id = $1', [id]);

	if (res.rows.length == 1) {
		return res.rows[0];
	}
	return undefined;
}

export function deleteTorrentSubtitle(id: number | string) {
	return client.query<Subtitle>('DELETE FROM torrent_subtitles WHERE id = $1', [id]);
}

export type MediaInformations = {
	imdb_id: number | null;
	tmdb_id: number | null;
	year: number | null;
};

export async function getMediaInformations(mediaId: number) {
	const media = await client.query<MediaInformations>(
		'SELECT imdb_id, tmdb_id, year FROM medias WHERE id = $1',
		[mediaId]
	);
	if (media.rows.length == 0) {
		throw 'no media found with this id';
	}

	const names = await client.query<{ name: string }>(
		"SELECT name FROM media_names WHERE media_id = $1 AND lang = '__'",
		[mediaId]
	);
	if (names.rows.length != 1) {
		throw 'no name found for the media';
	}

	return {
		...media.rows[0],
		name: names.rows[0].name
	};
}

export function createTorrentSubtitle(informations: {
	torrent_id: number;
	lang: string;
	path: string;
}) {
	return client.query<Subtitle>(
		'INSERT INTO torrent_subtitles (torrent_id, lang, path) VALUES ($1, $2, $3) RETURNING id',
		[informations.torrent_id, informations.lang, informations.path]
	);
}
