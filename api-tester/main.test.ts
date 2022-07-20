import { DESCRIBE_TEST_client } from "./1.client"
import { DESCRIBE_TEST_api_auth } from "./2.api-auth-internal"
import { DESCRIBE_TEST_api_user } from "./3.api-user"
import { DESCRIBE_TEST_api_position } from "./4.api-position"

jest.setTimeout(15000);

describe('DESCRIBE_TEST_client_external', DESCRIBE_TEST_client);
describe('DESCRIBE_TEST_api_auth', DESCRIBE_TEST_api_auth);
describe('DESCRIBE_TEST_api_user', DESCRIBE_TEST_api_user);
describe('DESCRIBE_TEST_api_position', DESCRIBE_TEST_api_position);