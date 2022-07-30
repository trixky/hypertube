import { config } from 'dotenv';

const env = config({ override: true });
export default env.parsed as unknown as {
	API_PICTURE_PORT: number;

	PGUSER: string;
	PGHOST: string;
	PGPASSWORD: number;
	PGDATABASE: string;
	PGPORT: number;

	REDIS_HOST: string;
	REDIS_PORT: number;
};
