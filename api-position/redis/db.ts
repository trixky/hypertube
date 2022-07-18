import { createClient } from 'redis';


export const client = createClient({
    url: 'redis://redis:6379'
});

export function connect(): Promise<any> {
    return new Promise(async (resolve, reject) => {
        client.on('error', (err) => reject(err));
        client.on('ready', () => resolve(undefined));
        client.connect();
    })
}