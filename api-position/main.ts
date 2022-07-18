import { connect as connect_to_pg } from './postgres/db'
import { connect as connect_to_redis } from './redis/db'
import middleware_auth from './middlewares/auth'
import express from "express";
import bodyParser from 'body-parser'
import cookieParser from 'cookie-parser'
import service_router from './controllers/positions'

async function main() {
    try {
        await connect_to_pg()
        console.log("connected to pg !")
    } catch (err: any) {
        console.log("failed to connect to pg: ", err)
        return
    }

    try {
        await connect_to_redis()
        console.log("connected to redis !")
    } catch (err: any) {
        console.log("failed to connect to redis: ", err)
        return
    }

    const app = express();
	app.use(cookieParser());
	app.use(middleware_auth);
    app.use(bodyParser.json())
    
    app.use("/v1/position", service_router)

	app.listen(3040);
}

console.log("---------------------------------------------------- start")

main()