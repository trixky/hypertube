import { connect, disconnect } from './redis/db';
import { DESCRIBE_TEST_client } from './tests/1.client';
import { DESCRIBE_TEST_api_auth } from './tests/2.api-auth-internal';
import { DESCRIBE_TEST_api_user } from './tests/3.api-user';
import { DESCRIBE_TEST_api_position } from './tests/4.api-position';
import { DESCRIBE_TEST_api_media } from './tests/5.api-media';

beforeAll(async () => {
	await connect();
});

jest.setTimeout(15000);
describe('DESCRIBE_TEST_client_external', DESCRIBE_TEST_client);
describe('DESCRIBE_TEST_api_auth', DESCRIBE_TEST_api_auth);
describe('DESCRIBE_TEST_api_user', DESCRIBE_TEST_api_user);
describe('DESCRIBE_TEST_api_position', DESCRIBE_TEST_api_position);
describe('DESCRIBE_TEST_api_media', DESCRIBE_TEST_api_media);

afterAll(async () => {
	await disconnect();
});
