import { client } from './db';

export function getUser(userId: number) {
	return client.query<{ extension: string | null }>('SELECT extension FROM users WHERE id = $1', [
		userId
	]);
}

export async function setUserExtension(userId: number, extension: string | null) {
	return client.query('UPDATE users SET extension = $2 WHERE id = $1', [userId, extension]);
}
