import { createClient } from 'redis';
import { sleep } from '../utils/time'
import { get_config } from '../config'

let connected = false
let start_connection = false

export const client = createClient({
    url: `redis://${get_config().REDIS_addresse}:${get_config().REDIS_port}`
});

export function connect(): Promise<any> {
    return new Promise(async (resolve, reject) => {
        if (connected)
            return resolve(undefined)
        if (start_connection) {
            while (connected == false) {
                await sleep(1000)
            }
            return resolve(undefined)
        }
        start_connection = true
        client.on('error', (err) => reject(err));
        client.on('ready', () => {
            connected = true;
            return resolve(undefined)
        });
        client.connect();
    })
}

export function disconnect(): Promise<any> {
    return new Promise(async (resolve, reject) => {
        if (connected) {
            start_connection = true
            client.on('error', (err) => reject(err));
            client.on('end', () => {
                connected = false;
                start_connection = false
                return resolve(undefined)
            });
            client.disconnect();
        }
        return resolve(undefined)
    })
}
