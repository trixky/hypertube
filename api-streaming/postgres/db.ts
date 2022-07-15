import { Client } from 'pg'

export let client: any = {}

export async function connect() {
    client = new Client()
    await client.connect() // read db configuration by default from environment variables
    const res = await client.query('SELECT $1::text as message', ['Connected to Postgres !'])
    console.log(res.rows[0].message) // Hello world!
    // await client.end()
}