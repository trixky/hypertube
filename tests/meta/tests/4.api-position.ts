import supertest from 'supertest';
import { shared } from '../shared';
import { get_config } from '../config';

const API_POSITION_request = supertest(
	`${get_config().API_POSITION_addresse}:${get_config().API_POSITION_port}`
);

export function DESCRIBE_TEST_api_position() {
	// ---------------------------------------- POST 400 /v1/position/{torrent_id} (position missing)
	it('POST 400 /v1/position/{torrent_id} (position missing)', (done) => {
		API_POSITION_request.post('/v1/position/1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/position/{torrent_id} (position corrupted #1)
	it('POST 400 /v1/position/{torrent_id} (position corrupted #1)', (done) => {
		API_POSITION_request.post('/v1/position/1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send({
				position: '-1'
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 404 /v1/position/{torrent_id} (torrent id missing)
	it('POST 404 /v1/position/{torrent_id} (torrent id missing)', (done) => {
		API_POSITION_request.post('/v1/position/')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send({
				position: 1
			})
			.expect(404)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/position/{torrent_id} (torrent id corrupted #1)
	it('POST 400 /v1/position/{torrent_id} (torrent id corrupted #1', (done) => {
		API_POSITION_request.post('/v1/position/-1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send({
				position: 1
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/position/{torrent_id} (torrent id corrupted #2)
	it('POST 400 /v1/position/{torrent_id} (torrent id corrupted #2', (done) => {
		API_POSITION_request.post('/v1/position/x')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send({
				position: 1
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/position/{torrent_id} (position corrupted #2)
	it('POST 400 /v1/position/{torrent_id} (position corrupted #2)', (done) => {
		API_POSITION_request.post('/v1/position/1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send({
				position: 'x'
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 401 /v1/position/{torrent_id} (token cookie missing)
	it('POST 401 /v1/position/{torrent_id} (token cookie missing)', (done) => {
		API_POSITION_request.post('/v1/position/1')
			.send({
				position: 1
			})
			.expect(401)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/position/{torrent_id} (corrupted torrent id #1)
	it('POST 400 /v1/position/{torrent_id} (corrupted torrent id #1)', (done) => {
		API_POSITION_request.post('/v1/position/-1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send({
				position: 1
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 400 /v1/position/{torrent_id} (corrupted torrent id #2)
	it('POST 400 /v1/position/{torrent_id} (corrupted torrent id #2)', (done) => {
		API_POSITION_request.post('/v1/position/x')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send({
				position: 1
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 200 /v1/position/{torrent_id}
	it('POST 200 /v1/position/{torrent_id}', (done) => {
		API_POSITION_request.post('/v1/position/1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send({
				position: 1
			})
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 401 /v1/position/{torrent_id} (token cookie missing)
	it('GET 401 /v1/position/{torrent_id} (token cookie missing)', (done) => {
		API_POSITION_request.get('/v1/position/-1').send().expect(401).end(done);
	});
	// ---------------------------------------- GET 400 /v1/position/{torrent_id} (corrupted torrent id #1)
	it('GET 400 /v1/position/{torrent_id} (corrupted torrent id #1)', (done) => {
		API_POSITION_request.get('/v1/position/-1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/position/{torrent_id} (corrupted torrent id #2)
	it('GET 400 /v1/position/{torrent_id} (corrupted torrent id #2)', (done) => {
		API_POSITION_request.get('/v1/position/x')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/position/{torrent_id}
	it('GET 200 /v1/position/{torrent_id}', (done) => {
		API_POSITION_request.get('/v1/position/1')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.expect((res) => {
				if (res?.body?.position == undefined) throw new Error('userinfo missing');
			})
			.end(done);
	});
}
