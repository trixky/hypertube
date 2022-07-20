import { client } from './db';

export type Torrent = {
	id: number;
	name: string;
	media_id: number | null;
	torrent_url: string | null;
	magnet: string | null;
	file_path: string | null;
	downloaded: boolean | null;
};

export async function getTorrent(id: number) {
	const res = await client.query<Torrent>(
		'SELECT id, name, media_id, torrent_url, magnet, file_path, downloaded \
        FROM torrents \
        WHERE id = $1;',
		[id]
	);

	if (res.rows.length == 1) {
		return res.rows[0];
	} else {
		throw 'no torrent found with this id';
	}
}

export async function updateTorrent(
	id: number,
	file_path: string | null,
	downloaded: boolean | null,
	length: number | null
): Promise<boolean> {
	const update_strings = [];
	const update_values: Array<string | number | boolean> = [id];

	let arg_nbr = 2;

	if (file_path != null) {
		update_strings.push('file_path = $' + arg_nbr++);
		update_values.push(file_path);
	}
	if (downloaded != null) {
		update_strings.push('downloaded = $' + arg_nbr++);
		update_values.push(downloaded);
	}
	if (length != null) {
		update_strings.push('length = $' + arg_nbr++);
		update_values.push(length);
	}

	const res: { rowCount: number } = await client.query(
		`UPDATE torrents \
        SET ${update_strings.join(',')} \
        WHERE id = $1;`,
		update_values
	);

	return res.rowCount > 0;
}
