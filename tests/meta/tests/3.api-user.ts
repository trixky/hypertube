import supertest from 'supertest';
import { shared } from '../shared';
import { get_config } from '../config';

const API_USER_request = supertest(
	`${get_config().API_USER_addresse}:${get_config().API_USER_port}`
);

export function DESCRIBE_TEST_api_user() {
	// ---------------------------------------- GET 404 /v1/me (token missing)
	it('GET 404 /v1/me (token missing)', (done) => {
		API_USER_request.get('/v1/me' + shared.random_user_token)
			.send()
			.expect(404)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/me (token corrupted)
	it('GET 400 /v1/me (token corrupted)', (done) => {
		API_USER_request.get('/v1/me?token=' + shared.random_user_token_corrupted)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 200 /me
	it('GET 200 /me', (done) => {
		API_USER_request.get('/v1/me?token=' + shared.random_user_token)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- PATCH 403 /v1/me (token cookie missing)
	it('PATCH 403 /v1/me (token cookie missing)', (done) => {
		API_USER_request.patch('/v1/me')
			.send({
				username: 'new_username',
				firstname: 'new_firstname',
				lastname: 'new_lastname',
				email: 'new_email',
				current_password: shared.default_new_password_1,
				new_password: shared.default_new_password_2
			})
			.expect(403)
			.end(done);
	});
	// ---------------------------------------- PATCH 400 /v1/me (corrupted email)
	it('PATCH 400 /v1/me (corrupted email)', (done) => {
		const url_search_params = new URLSearchParams();

		url_search_params.append('email', shared.random_user_email_corrupted);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- PATCH 400 /v1/me (corrupted username)
	it('PATCH 400 /v1/me (corrupted username)', (done) => {
		const url_search_params = new URLSearchParams();

		url_search_params.append('username', 'x');

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- PATCH 400 /v1/me (corrupted firstname)
	it('PATCH 400 /v1/me (corrupted firstname)', (done) => {
		const url_search_params = new URLSearchParams();

		url_search_params.append('firstname', 'x');

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- PATCH 400 /v1/me (corrupted lastname)
	it('PATCH 400 /v1/me (corrupted lastname)', (done) => {
		const url_search_params = new URLSearchParams();

		url_search_params.append('lastname', 'x');

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- PATCH 400 /v1/me (corrupted current/new_password)
	it('PATCH 400 /v1/me (corrupted current/new_password)', (done) => {
		const url_search_params = new URLSearchParams();

		url_search_params.append('current_password', shared.default_password_corrupted);
		url_search_params.append('new_password', shared.default_new_password_corrupted_1);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- PATCH 200 /v1/me (missing new_password)
	it('PATCH 200 /v1/me (missing new_password)', (done) => {
		const url_search_params = new URLSearchParams();

		url_search_params.append('current_password', shared.default_new_password_2);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- PATCH 403 /v1/me (missing current_password)
	it('PATCH 403 /v1/me (missing current_password)', (done) => {
		const url_search_params = new URLSearchParams();

		url_search_params.append('new_password', shared.default_new_password_2);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(403)
			.end(done);
	});
	// ---------------------------------------- PATCH 200 /v1/me (no infos)
	it('PATCH 200 /v1/me (no infos)', (done) => {
		API_USER_request.patch('/v1/me')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- PATCH 200 /v1/me (only email)
	it('PATCH 200 /v1/me (only email)', (done) => {
		shared.random_user_userInfo!.email += 'x';

		const url_search_params = new URLSearchParams();

		url_search_params.append('email', shared.random_user_userInfo!.email);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- PATCH 200 /v1/me (only username)
	it('PATCH 200 /v1/me (only username)', (done) => {
		shared.random_user_userInfo!.username += 'x';

		const url_search_params = new URLSearchParams();

		url_search_params.append('username', shared.random_user_userInfo!.username);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- PATCH 200 /v1/me (only firstname)
	it('PATCH 200 /v1/me (only firstname)', (done) => {
		shared.random_user_userInfo!.firstname += 'x';

		const url_search_params = new URLSearchParams();

		url_search_params.append('firstname', shared.random_user_userInfo!.firstname);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- PATCH 200 /v1/me (only lastname)
	it('PATCH 200 /v1/me (only lastname)', (done) => {
		shared.random_user_userInfo!.lastname += 'x';

		const url_search_params = new URLSearchParams();

		url_search_params.append('lastname', shared.random_user_userInfo!.lastname);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- PATCH 200 /v1/me (only password)
	it('PATCH 200 /v1/me (only password)', (done) => {
		const url_search_params = new URLSearchParams();

		url_search_params.append('current_password', shared.default_new_password_1);
		url_search_params.append('new_password', shared.default_new_password_2);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- PATCH 200 /v1/me (full infos)
	it('PATCH 200 /v1/me (full infos)', (done) => {
		shared.random_user_userInfo!.email += 'x';
		shared.random_user_userInfo!.username += 'x';
		shared.random_user_userInfo!.firstname += 'x';
		shared.random_user_userInfo!.lastname += 'x';

		const url_search_params = new URLSearchParams();

		url_search_params.append('username', shared.random_user_userInfo!.username);
		url_search_params.append('firstname', shared.random_user_userInfo!.firstname);
		url_search_params.append('lastname', shared.random_user_userInfo!.lastname);
		url_search_params.append('email', shared.random_user_userInfo!.email);
		url_search_params.append('current_password', shared.default_new_password_2);
		url_search_params.append('new_password', shared.default_new_password_3);

		API_USER_request.patch('/v1/me' + '?' + url_search_params.toString())
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 403 /v1/user (token cookie missing)
	it('GET 403 /v1/user (token cookie missing)', (done) => {
		API_USER_request.get('/v1/user?id=0').send().expect(403).end(done);
	});
	// ---------------------------------------- GET 404 /v1/user (corrupted user id #1)
	it('GET 404 /v1/user (corrupted user id #1)', (done) => {
		API_USER_request.get('/v1/user?id=-1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(404)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/user (corrupted user id #2)
	it('GET 400 /v1/user (corrupted user id #2)', (done) => {
		API_USER_request.get('/v1/user?id=x2')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/user
	it('GET 200 /v1/user', (done) => {
		API_USER_request.get('/v1/user?id=' + shared.random_user_userInfo!.id)
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.expect((res) => {
				if (res?.body?.userInfo == undefined) throw new Error('userinfo missing');
			})
			.end(done);
	});
}
