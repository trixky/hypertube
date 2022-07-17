import { connect as connect_to_pg } from './postgres/db'
import { connect as connect_to_redis } from './redis/db'
import express from "express";

import { get_position_handler } from './handlers/get_position.handler'

async function main() {
    try {
        await connect_to_pg()
    } catch (err: any) {
        console.log("failed to connect to pg: ", err)
        return
    }

    try {
        await connect_to_redis()
    } catch (err: any) {
        console.log("failed to connect to redis: ", err)
        return
    }

    const app = express();

    app.get("/position/:torrent_id", async function (req, res) {

    })
}

main()