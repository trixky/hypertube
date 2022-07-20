import { config } from 'dotenv';

const env = config({ override: true });
export default env.parsed as unknown as {
	API_STREAMING_PORT: number;

	PGUSER: string;
	PGHOST: string;
	PGPASSWORD: number;
	PGDATABASE: string;
	PGPORT: number;

	OSDB_API_KEY: string;
	OSDB_USERNAME: string;
	OSDB_PASSWORD: string;

	OLD_FILES_INTERVAL: string;
	OLD_FILES_CRON: string;
};
