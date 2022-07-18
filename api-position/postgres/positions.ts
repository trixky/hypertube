import { client } from './db'
import * as models from '../models/positions'

export async function save_position(position: models.Position): Promise<any> {
    try {
        await client.query("INSERT INTO positions (user_id, torrent_id, position) \
        VALUES ($1, $2, $3) \
        ON CONFLICT ON CONSTRAINT unique_user_torrent_relation DO UPDATE \
        SET position = EXCLUDED.position;", [position.user_id, position.torrent_id, position.position])
    } catch (err: any) {
        throw (err)
    }
}

export async function get_position(torrend_id: number): Promise<models.Position> {
    try {
        const res = await client.query("SELECT * FROM positions WHERE id = $1;", [torrend_id])

        if (res.rows.length != 1)
            throw (new Error('zero or too many positions finded'))

        return res.rows[0]
    } catch (err: any) {
        throw (err)
    }
}
