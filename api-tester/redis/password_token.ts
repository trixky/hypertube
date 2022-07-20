import { client } from '../redis/db';

export default async function get_password_token(user_id: number): Promise<string | undefined> {
	const keys: Array<string> = await client.sendCommand([
		'KEYS',
		'password_token.*.' + user_id.toString()
	]);

	return keys.length > 0 ? keys[0].slice(15).slice(0, -(user_id.toString().length + 1)) : undefined;
}
