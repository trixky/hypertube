interface Config {
	// ----------------------------------- REDIS
	REDIS_addresse: string;
	REDIS_port: number | undefined;
	// ----------------------------------- CLIENT
	CLIENT_addresse: string;
	CLIENT_port: number | undefined;
	// ----------------------------------- API_AUTH
	API_AUTH_addresse: string;
	API_AUTH_port: number | undefined;
	// ----------------------------------- API_USER
	API_USER_addresse: string;
	API_USER_port: number | undefined;
	// ----------------------------------- API_POSITION
	API_POSITION_addresse: string;
	API_POSITION_port: number | undefined;
}

const config = <Config>{
	// ----------------------------------- REDIS
	REDIS_addresse: 'redis',
	REDIS_port: undefined,
	// ----------------------------------- CLIENT
	CLIENT_addresse: 'client',
	CLIENT_port: undefined,
	// ----------------------------------- API_AUTH
	API_AUTH_addresse: 'api-auth',
	API_AUTH_port: undefined,
	// ----------------------------------- API_USER
	API_USER_addresse: 'api-user',
	API_USER_port: undefined,
	// ----------------------------------- API_POSITION
	API_POSITION_addresse: 'api-position',
	API_POSITION_port: undefined
};

export function get_config(): Config {
	if (config.REDIS_port === undefined) {
		read_env();
	}

	return config;
}

export function read_env() {
	// ----------------------------------- REDIS
	config.REDIS_port = parseInt(process.env.REDIS_PORT === undefined ? 'x' : process.env.REDIS_PORT);
	if (isNaN(config.REDIS_port)) throw new Error('REDIS_PORT environment variable is missing');
	// ----------------------------------- CLIENT
	config.CLIENT_port = parseInt(
		process.env.CLIENT_PORT === undefined ? 'x' : process.env.CLIENT_PORT
	);
	if (isNaN(config.CLIENT_port)) throw new Error('CLIENT_PORT environment variable is missing');
	// ----------------------------------- API_AUTH
	config.API_AUTH_port = parseInt(
		process.env.API_AUTH_GRPC_GATEWAY_PORT === undefined
			? 'x'
			: process.env.API_AUTH_GRPC_GATEWAY_PORT
	);
	if (isNaN(config.API_AUTH_port))
		throw new Error('API_AUTH_GRPC_GATEWAY_PORT environment variable is missing');
	// ----------------------------------- API_USER
	config.API_USER_port = parseInt(
		process.env.API_USER_GRPC_GATEWAY_PORT === undefined
			? 'x'
			: process.env.API_USER_GRPC_GATEWAY_PORT
	);
	if (isNaN(config.API_USER_port))
		throw new Error('API_USER_GRPC_GATEWAY_PORT environment variable is missing');
	// ----------------------------------- API_POSITION
	config.API_POSITION_port = parseInt(
		process.env.API_POSITION_PORT === undefined ? 'x' : process.env.API_POSITION_PORT
	);
	if (isNaN(config.API_POSITION_port))
		throw new Error('API_POSITION_PORT environment variable is missing');
}
