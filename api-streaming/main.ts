import express from 'express';
import cors from 'cors';
import bodyParser from 'body-parser';
import cookieParser from 'cookie-parser';
import middleware_auth from './middlewares/auth';
import { connect as connectPostgres } from './postgres/db';
import env from './env';
import { connect as connectRedis } from './redis/db';
import subtitlesRouter from './Controllers/subtitles';
import streamRouter from './Controllers/streaming';
import scheduleDeleteFiles from './jobs/delete-files';

(async () => {
	if (!env.OSDB_API_KEY || !env.OSDB_USERNAME || !env.OSDB_PASSWORD) {
		throw new Error('Missing required env keys');
	}

	// *

	try {
		await connectPostgres();
		console.log('connected to pg !');
	} catch (err) {
		console.log('failed to connect to pg: ', err);
		return;
	}

	try {
		await connectRedis();
		console.log('connected to redis !');
	} catch (err) {
		console.log('failed to connect to redis: ', err);
		return;
	}

	// *

	const app = express();

	app.use(
		cors({
			origin: ['http://localhost:4040'],
			credentials: true
		})
	);

	app.use(cookieParser());
	app.use(middleware_auth);
	app.use(bodyParser.json());

	app.use('/', subtitlesRouter);
	app.use('/', streamRouter);

	app.listen(env.API_STREAMING_PORT);
	console.log('[api-streaming] listening on port', env.API_STREAMING_PORT);

	// *

	scheduleDeleteFiles();
})();
