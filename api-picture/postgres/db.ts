import { Client } from 'pg';

export let client: Client;

export async function connect() {
	client = new Client({
		host: process.env.POSTGRES_HOST,
		port: parseInt(process.env.POSTGRES_PORT ?? '5342'),
		user: process.env.POSTGRES_USER,
		password: process.env.POSTGRES_PASSWORD,
		database: process.env.POSTGRES_DB
	});
	await client.connect(); // read db configuration by default from environment variables
	await client.query('SELECT $1::text as message', ['Connected to Postgres !']);
}
