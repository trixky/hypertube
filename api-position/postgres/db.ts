import { Client } from 'pg';

export let client: Client;

export async function connect() {
	client = new Client();
	await client.connect(); // read db configuration by default from environment variables
	await client.query('SELECT $1::text as message', ['Connected to Postgres !']);
}
