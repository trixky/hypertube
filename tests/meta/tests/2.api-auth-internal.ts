import supertest from 'supertest';
import { extract_hex64 } from '../utils/hex64';
import get_password_token from '../redis/password_token';
import { shared } from '../shared';
import { get_config } from '../config';

const API_AUTH_request = supertest(
	`${get_config().API_AUTH_addresse}:${get_config().API_AUTH_port}`
);

export function DESCRIBE_TEST_api_auth() {
	// ---------------------------------------- POST 400 /v1/internal/register (infos missing)
	it('POST 400 /v1/internal/register(infos missing)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({})
			.expect(400)
			.end(done);
	});
	// ----------------------------------------POST 400 /v1/internal/register (email missing)
	it('POST 400 /v1/internal/register(email missing)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				username: 'asdf',
				firstname: 'asdf',
				lastname: 'asdf',
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (email corrupted)
	it('POST 400 /v1/internal/register (email corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email_corrupted,
				username: 'asdf',
				firstname: 'asdf',
				lastname: 'asdf',
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (password missing)
	it('POST 400 /v1/internal/register (password missing)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				username: 'asdf',
				firstname: 'asdf',
				lastname: 'asdf'
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (password corrupted)
	it('POST 400 /v1/internal/register (password corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				username: 'asdf',
				firstname: 'asdf',
				lastname: 'asdf',
				password: shared.default_password_corrupted
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (username missing)
	it('POST 400 /v1/internal/register (username missing)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				firstname: 'asdf',
				lastname: 'asdf',
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (username corrupted)
	it('POST 400 /v1/internal/register (username corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				username: 'a',
				firstname: 'asdf',
				lastname: 'asdf',
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (firstname missing)
	it('POST 400 /v1/internal/register(firstname missing)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				username: 'aasdf',
				lastname: 'asdf',
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (firstname corrupted)
	it('POST 400 /v1/internal/register (firstname corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				username: 'aasdf',
				firstname: 'a',
				lastname: 'asdf',
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (lastname missing)
	it('POST 400 /v1/internal/register (lastname missing)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				firstname: 'asdf',
				username: 'aasdf',
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/register (lastname corrupted)
	it('POST 400 /v1/internal/register(lastname corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				username: 'aasdf',
				firstname: 'asdf',
				lastname: 'a',
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 200 /v1/internal/register
	it('POST 200 /v1/internal/register', (done) => {
		API_AUTH_request.post('/v1/internal/register')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				username: 'asdf',
				firstname: 'asdf',
				lastname: 'asdf',
				password: shared.default_password
			})
			.expect(200)
			.expect((res) => {
				if (res?.body?.token == undefined) throw new Error('token missing');
				if (res?.body?.userInfo == undefined) throw new Error('userinfo missing');
			})
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/login (infos missing)
	it('POST 400 /v1/internal/login (infos missing)', (done) => {
		API_AUTH_request.post('/v1/internal/login')
			.set('Content-type', 'application/json')
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/login (email missing)
	it('POST 400 /v1/internal/login (email missing)', (done) => {
		API_AUTH_request.post('/v1/internal/login')
			.set('Content-type', 'application/json')
			.send({
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/login (email corrupted)
	it('POST 400 /v1/internal/login (email corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/login')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email_corrupted,
				password: shared.default_password
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/login (password missing)
	it('POST 400 /v1/internal/login (password missing)', (done) => {
		API_AUTH_request.post('/v1/internal/login')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/login (password corrupted)
	it('POST 400 /v1/internal/login (password corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/login')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				password: shared.default_password_corrupted
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 200 /v1/internal/login
	it('POST 200 /v1/internal/login', (done) => {
		API_AUTH_request.post('/v1/internal/login')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email,
				password: shared.default_password
			})
			.expect(200)
			.expect((res) => {
				if (res?.body?.token == undefined) throw new Error('token missing');
				if (res?.body?.userInfo == undefined) throw new Error('userinfo missing');

				// save token and userInfo for next tests
				shared.random_user_token = res.body.token;
				shared.random_user_userInfo = extract_hex64(res.body.userInfo);
				if (shared.random_user_userInfo === undefined) {
					throw new Error('userInfo corrupted');
				}
			})
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/recover-password (email missing)
	it('POST 400 /v1/internal/recover-password (email missing)', (done) => {
		API_AUTH_request.post('/v1/internal/recover-password')
			.set('Content-type', 'application/json')
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/recover-password (email corrupted)
	it('POST 400 /v1/internal/recover-password (email corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/recover-password')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email_corrupted
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/recover-password (email missing)
	it('POST 400 /v1/internal/recover-password (email missing)', (done) => {
		API_AUTH_request.post('/v1/internal/recover-password')
			.set('Content-type', 'application/json')
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/recover-password (email corrupted)
	it('POST 400 /v1/internal/recover-password (email corrupted)', (done) => {
		API_AUTH_request.post('/v1/internal/recover-password')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email_corrupted
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 200 /v1/internal/recover-password
	it('POST 200 /v1/internal/recover-password', (done) => {
		API_AUTH_request.post('/v1/internal/recover-password')
			.set('Content-type', 'application/json')
			.send({
				email: shared.random_user_email
			})
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/internal/apply-token-password (infos missing)
	it('POST 400 /v1/internal/apply-token-password (infos missing)', async () => {
		if (shared.random_user_userInfo == null) {
			throw new Error('cant do the test: need userInfo');
		}

		expect(shared.random_user_userInfo).toBeDefined();

		const password_token = await get_password_token(shared.random_user_userInfo!.id);

		expect(password_token).toBeDefined();

		await API_AUTH_request.post('/v1/internal/apply-token-password')
			.set('Content-type', 'application/json')
			.send()
			.expect(400);
	});
	// ---------------------------------------- POST 400 /v1/internal/apply-token-password (password_token missing)
	it('POST 400 /v1/internal/apply-token-password (password_token missing)', async () => {
		if (shared.random_user_userInfo == null) {
			throw new Error('cant do the test: need userInfo');
		}

		expect(shared.random_user_userInfo).toBeDefined();

		const password_token = await get_password_token(shared.random_user_userInfo!.id);

		expect(password_token).toBeDefined();

		await API_AUTH_request.post('/v1/internal/apply-token-password')
			.set('Content-type', 'application/json')
			.send({
				new_password: shared.default_new_password_1
			})
			.expect(400);
	});
	// ---------------------------------------- POST 400 /v1/internal/apply-token-password (password_token corrupted)
	it('POST 400 /v1/internal/apply-token-password (password_token corrupted)', async () => {
		if (shared.random_user_userInfo == null) {
			throw new Error('cant do the test: need userInfo');
		}

		expect(shared.random_user_userInfo).toBeDefined();

		const password_token = await get_password_token(shared.random_user_userInfo!.id);

		expect(password_token).toBeDefined();

		await API_AUTH_request.post('/v1/internal/apply-token-password')
			.set('Content-type', 'application/json')
			.send({
				password_token: shared.default_password_token_corrupted,
				new_password: shared.default_new_password_1
			})
			.expect(400);
	});
	// ---------------------------------------- POST 400 /v1/internal/apply-token-password (new_password missing)
	it('POST 400 /v1/internal/apply-token-password (new_password missing)', async () => {
		if (shared.random_user_userInfo == null) {
			throw new Error('cant do the test: need userInfo');
		}

		expect(shared.random_user_userInfo).toBeDefined();

		const password_token = await get_password_token(shared.random_user_userInfo!.id);

		expect(password_token).toBeDefined();

		await API_AUTH_request.post('/v1/internal/apply-token-password')
			.set('Content-type', 'application/json')
			.send({
				password_token: password_token
			})
			.expect(400);
	});
	// ---------------------------------------- POST 400 /v1/internal/apply-token-password (new_password corrupted)
	it('POST 400 /v1/internal/apply-token-password (new_password corrupted)', async () => {
		if (shared.random_user_userInfo == null) {
			throw new Error('cant do the test: need userInfo');
		}

		expect(shared.random_user_userInfo).toBeDefined();

		const password_token = await get_password_token(shared.random_user_userInfo!.id);

		expect(password_token).toBeDefined();

		await API_AUTH_request.post('/v1/internal/apply-token-password')
			.set('Content-type', 'application/json')
			.send({
				password_token: password_token,
				new_password: shared.default_new_password_corrupted_1
			})
			.expect(400);
	});
	// ---------------------------------------- POST 200 /v1/internal/apply-token-password
	it('POST 200 /v1/internal/apply-token-password', async () => {
		if (shared.random_user_userInfo == null) {
			throw new Error('cant do the test: need userInfo');
		}

		expect(shared.random_user_userInfo).toBeDefined();

		const password_token = await get_password_token(shared.random_user_userInfo!.id);

		expect(password_token).toBeDefined();

		await API_AUTH_request.post('/v1/internal/apply-token-password')
			.set('Content-type', 'application/json')
			.send({
				password_token: password_token,
				new_password: shared.default_new_password_1
			})
			.expect(200);
	});
}
