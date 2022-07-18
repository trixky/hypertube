import express from "express";
import cors from "cors";
import { config } from "dotenv";
import { connect } from "./postgres/db";
import subtitlesRouter from "./Controllers/subtitles";
import streamRouter from "./Controllers/streaming";

config();

(async () => {
	if (!process.env.OSDB_API_KEY || !process.env.OSDB_USERNAME || !process.env.OSDB_PASSWORD) {
		throw new Error("Missing required env keys");
	}

	const app = express();

	await connect();

	console.log("**************************************** START");
	console.log("**************************************** START");
	console.log("**************************************** START");

	app.use(cors());

	app.use("/", subtitlesRouter);
	app.use("/", streamRouter);

	app.listen(3030);
})();
