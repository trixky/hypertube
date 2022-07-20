import { client } from './db';
import * as models from '../models/positions';

export function save_position(position: models.Position) {
	return client.query(
		'INSERT INTO positions (user_id, torrent_id, position) \
		VALUES ($1, $2, $3) \
		ON CONFLICT ON CONSTRAINT unique_user_torrent_relation DO UPDATE \
		SET position = EXCLUDED.position;',
		[position.user_id, position.torrent_id, position.position]
	);
}

export async function get_position(user_id: number, torrend_id: number): Promise<models.Position> {
	const res = await client.query(
		'SELECT * FROM positions WHERE user_id = $1 AND torrent_id = $2;',
		[user_id, torrend_id]
	);

	if (res.rows.length != 1) throw new Error('zero or too many positions found');

	return res.rows[0];
}
