import express from 'express';
import cors from 'cors';
import bodyParser from 'body-parser';
import cookieParser from 'cookie-parser';
import middleware_auth from './middlewares/auth';
import { connect as connectPostgres } from './postgres/db';
import env from './env';
import { connect as connectRedis } from './redis/db';
import usersRouter from './Controllers/users';

(async () => {
	if (!env.API_PICTURE_PORT) {
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
			origin: [`http://${process.env.DOMAIN}:${process.env.CLIENT_PORT}`],
			credentials: true
		})
	);

	app.use(cookieParser());
	app.use(middleware_auth);
	app.use(bodyParser.json());

	app.use('/', usersRouter);

	const server = app.listen(env.API_PICTURE_PORT);
	console.log('[api-picture] listening on port', env.API_PICTURE_PORT);

	// *

	function closeGracefully() {
		server.close();
	}

	process.once('SIGINT', closeGracefully);
	process.once('SIGTERM', closeGracefully);
})();
