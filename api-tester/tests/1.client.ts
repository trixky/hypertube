import supertest from 'supertest';
import { get_config } from '../config';

const CLIENT_request = supertest(`${get_config().CLIENT_addresse}:${get_config().CLIENT_port}`);

export function DESCRIBE_TEST_client() {
	// ---------------------------------------- GET 200 /login
	it('GET 200 /login', (done) => {
		CLIENT_request.get('/login').expect(200).end(done);
	});
	// ---------------------------------------- GET 200 /register
	it('GET 200 /register', (done) => {
		CLIENT_request.get('/register').expect(200).end(done);
	});
	// ---------------------------------------- GET 200 /recover/ask
	it('GET 200 /recover/ask', (done) => {
		CLIENT_request.get('/recover/ask').expect(200).end(done);
	});
	// ---------------------------------------- GET 200 /recover/apply
	it('GET 200 /recover/apply', (done) => {
		CLIENT_request.get('/recover/apply').expect(200).end(done);
	});
	// ---------------------------------------- GET 404 /nothing
	it('GET 404 /nothing', (done) => {
		CLIENT_request.get('/nothing').expect(404).end(done);
	});
	// ---------------------------------------- GET 302 /
	it('GET 302 /', (done) => {
		CLIENT_request.get('/').expect(302).expect('Location', '/login').end(done);
	});
	// ----------------------------------------  GET 302 /search
	it('GET 302 /search', (done) => {
		CLIENT_request.get('/search').expect(302).expect('Location', '/login').end(done);
	});
	// ---------------------------------------- GET 302 /users/0
	it('GET 302 /users/0', (done) => {
		CLIENT_request.get('/users/0').expect(302).expect('Location', '/login').end(done);
	});
	// ---------------------------------------- GET 302 /media/1151
	it('GET 302 /media/1151', (done) => {
		CLIENT_request.get('/media/1151').expect(302).expect('Location', '/login').end(done);
	});
}
