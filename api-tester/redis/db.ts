import { createClient } from 'redis';
import { sleep } from '../utils/time';
import { get_config } from '../config';

let connected = false;
let start_connection = false;

export const client = createClient({
	url: `redis://${get_config().REDIS_addresse}:${get_config().REDIS_port}`
});

export function connect() {
	// eslint-disable-next-line no-async-promise-executor
	return new Promise(async (resolve, reject) => {
		if (connected) return resolve(undefined);
		if (start_connection) {
			while (connected == false) {
				await sleep(1000);
			}
			resolve(undefined);
		}
		start_connection = true;
		client.on('error', (err) => reject(err));
		client.on('ready', () => {
			connected = true;
			resolve(undefined);
		});
		client.connect();
	});
}

export function disconnect() {
	// eslint-disable-next-line no-async-promise-executor
	return new Promise(async (resolve, reject) => {
		if (connected) {
			start_connection = true;
			client.on('error', (err) => reject(err));
			client.on('end', () => {
				connected = false;
				start_connection = false;
				resolve(undefined);
			});
			client.disconnect();
		} else {
			resolve(undefined);
		}
	});
}
