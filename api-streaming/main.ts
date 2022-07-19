import express from "express";
import cors from "cors";
import { config } from "dotenv";
import middleware_auth from "./middlewares/auth";
import bodyParser from "body-parser";
import cookieParser from "cookie-parser";
import { connect as connectPostgres } from "./postgres/db";
import { connect as connectRedis } from "./redis/db";
import subtitlesRouter from "./Controllers/subtitles";
import streamRouter from "./Controllers/streaming";

config();

(async () => {
	if (!process.env.OSDB_API_KEY || !process.env.OSDB_USERNAME || !process.env.OSDB_PASSWORD) {
		throw new Error("Missing required env keys");
	}

	// *

	try {
		await connectPostgres();
		console.log("connected to pg !");
	} catch (err: any) {
		console.log("failed to connect to pg: ", err);
		return;
	}

	try {
		await connectRedis();
		console.log("connected to redis !");
	} catch (err: any) {
		console.log("failed to connect to redis: ", err);
		return;
	}

	// *

	const app = express();

	console.log("**************************************** START");
	console.log("**************************************** START");
	console.log("**************************************** START");

	app.use(
		cors({
			origin: ["http://localhost:4040"],
			credentials: true,
		})
	);

	app.use(cookieParser());
	app.use(middleware_auth);
	app.use(bodyParser.json());

	app.use("/", subtitlesRouter);
	app.use("/", streamRouter);

	app.listen(3030);
})();
