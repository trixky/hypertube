import { createClient } from 'redis';

export const client = createClient({
	url: `redis://${process.env.REDIS_HOST}:${process.env.REDIS_PORT}`
});

export function connect() {
	return new Promise((resolve, reject) => {
		client.on('error', (err) => reject(err));
		client.on('ready', () => resolve(undefined));
		client.connect();
	});
}
