import supertest from 'supertest';
import { shared } from '../shared';
import { get_config } from '../config';

const API_MEDIA_request = supertest(
	`${get_config().API_MEDIA_addresse}:${get_config().API_MEDIA_port}`
);

export function DESCRIBE_TEST_api_media() {
	// ---------------------------------------- GET 403 /v1/media/genres (user logged out)
	it('GET 403 /v1/media/genres (user logged out)', (done) => {
		API_MEDIA_request.get('/v1/media/genres').send().expect(403).end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/genres
	it('GET 200 /v1/media/genres', (done) => {
		API_MEDIA_request.get('/v1/media/genres')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});

	// ---------------------------------------- GET 403 /v1/media/{media_id}/get (user logged out)
	it('GET 403 /v1/media/{media_id}/get (user logged out)', (done) => {
		API_MEDIA_request.get('/v1/media/1/get').send().expect(403).end(done);
	});
	// ---------------------------------------- GET 403 /v1/media/{media_id}/get (user logged out)
	it('GET 403 /v1/media/{media_id}/get (user logged out)', (done) => {
		API_MEDIA_request.get('/v1/media/9999999/get').send().expect(403).end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/{media_id}/get (exists)
	it('GET 200 /v1/media/{media_id}/get (exists)', (done) => {
		API_MEDIA_request.get('/v1/media/1/get')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/{media_id}/get (exists)
	it("GET 404 /v1/media/{media_id}/get (doesn't exists)", (done) => {
		API_MEDIA_request.get('/v1/media/9999999/get')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(404)
			.end(done);
	});

	// ---------------------------------------- GET 403 /v1/media/search (user logged out)
	it('GET 403 /v1/media/search (user logged out)', (done) => {
		API_MEDIA_request.get('/v1/media/search').send().expect(403).end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (homepage)
	it('GET 200 /v1/media/search (homepage)', (done) => {
		API_MEDIA_request.get('/v1/media/search')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (empty query)
	it('GET 200 /v1/media/search (empty query)', (done) => {
		API_MEDIA_request.get('/v1/media/search?query=')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (query)
	it('GET 200 /v1/media/search (query)', (done) => {
		API_MEDIA_request.get('/v1/media/search?query=hello')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/media/search (empty year)
	it('GET 400 /v1/media/search (empty year)', (done) => {
		API_MEDIA_request.get('/v1/media/search?year=')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/media/search (invalid year)
	it('GET 400 /v1/media/search (invalid year)', (done) => {
		API_MEDIA_request.get('/v1/media/search?year=deux-mille')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (year)
	it('GET 200 /v1/media/search (year)', (done) => {
		API_MEDIA_request.get('/v1/media/search?year=1997')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/media/search (empty rating)
	it('GET 400 /v1/media/search (empty rating)', (done) => {
		API_MEDIA_request.get('/v1/media/search?rating=')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/media/search (invalid rating)
	it('GET 400 /v1/media/search (invalid rating)', (done) => {
		API_MEDIA_request.get('/v1/media/search?rating=deux-mille')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (rating)
	it('GET 200 /v1/media/search (rating)', (done) => {
		API_MEDIA_request.get('/v1/media/search?rating=6.5')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/media/search (empty genre_ids)
	it('GET 400 /v1/media/search (empty genre_ids)', (done) => {
		API_MEDIA_request.get('/v1/media/search?genre_ids=')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/media/search (invalid genre_ids)
	it('GET 400 /v1/media/search (invalid genre_ids)', (done) => {
		API_MEDIA_request.get('/v1/media/search?genre_ids=deux-mille')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (genre_ids)
	it('GET 200 /v1/media/search (genre_ids)', (done) => {
		API_MEDIA_request.get('/v1/media/search?genre_ids=6')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (multiple genre_ids)
	it('GET 200 /v1/media/search (multiple genre_ids)', (done) => {
		API_MEDIA_request.get('/v1/media/search?genre_ids=6&genre_ids=7')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (empty sort_by)
	it('GET 200 /v1/media/search (empty sort_by)', (done) => {
		API_MEDIA_request.get('/v1/media/search?sort_by=')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (invalid sort_by column)
	it('GET 200 /v1/media/search (invalid sort_by column)', (done) => {
		API_MEDIA_request.get('/v1/media/search?sort_by=column')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (sort_by)
	it('GET 200 /v1/media/search (sort_by)', (done) => {
		API_MEDIA_request.get('/v1/media/search?sort_by=id')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (empty sort_order)
	it('GET 200 /v1/media/search (empty sort_order)', (done) => {
		API_MEDIA_request.get('/v1/media/search?sort_order=')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (invalid sort_order direction)
	it('GET 200 /v1/media/search (invalid sort_order direction)', (done) => {
		API_MEDIA_request.get('/v1/media/search?sort_order=both')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (sort_order)
	it('GET 200 /v1/media/search (sort_order)', (done) => {
		API_MEDIA_request.get('/v1/media/search?sort_order=asc')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/media/search (empty page)
	it('GET 400 /v1/media/search (empty page)', (done) => {
		API_MEDIA_request.get('/v1/media/search?page=')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 400 /v1/media/search (invalid page)
	it('GET 400 /v1/media/search (invalid page)', (done) => {
		API_MEDIA_request.get('/v1/media/search?page=two')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (page)
	it('GET 200 /v1/media/search (page)', (done) => {
		API_MEDIA_request.get('/v1/media/search?page=2')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/search (full)
	it('GET 200 /v1/media/search (full)', (done) => {
		API_MEDIA_request.get(
			'/v1/media/search?query=h&year=1997&rating=5.0&genres_ids=1&sort_by=duration&sort_order=desc&page=1'
		)
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});

	// ---------------------------------------- GET 403 /v1/media/{media_id}/comments (user logged out)
	it('GET 403 /v1/media/{media_id}/comments (user logged out)', (done) => {
		API_MEDIA_request.get('/v1/media/1/comments').send().expect(403).end(done);
	});
	// ---------------------------------------- GET 403 /v1/media/{media_id}/comments (user logged out)
	it('GET 403 /v1/media/{media_id}/comments (user logged out)', (done) => {
		API_MEDIA_request.get('/v1/media/9999999/comments').send().expect(403).end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/{media_id}/comments (exists)
	it('GET 200 /v1/media/{media_id}/comments (exists)', (done) => {
		API_MEDIA_request.get('/v1/media/1/comments')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- GET 200 /v1/media/{media_id}/comments (exists)
	it("GET 404 /v1/media/{media_id}/comments (doesn't exists)", (done) => {
		API_MEDIA_request.get('/v1/media/9999999/comments')
			.set('Cookie', `token=${shared.random_user_token};`)
			.send()
			.expect(404)
			.end(done);
	});

	// ---------------------------------------- POST 403 /v1/media/{media_id}/comment (user logged out)
	it('POST 403 /v1/media/{media_id}/comment (user logged out)', (done) => {
		API_MEDIA_request.post('/v1/media/1/comment')
			.set('Content-Type', 'application/json')
			.send({
				content: 'Hello'
			})
			.expect(403)
			.end(done);
	});
	// ---------------------------------------- POST 403 /v1/media/{media_id}/comment (user logged out)
	it('POST 403 /v1/media/{media_id}/comment (user logged out)', (done) => {
		API_MEDIA_request.post('/v1/media/9999999/comment')
			.set('Content-Type', 'application/json')
			.send({
				content: 'Hello'
			})
			.expect(403)
			.end(done);
	});
	// ---------------------------------------- POST 200 /v1/media/{media_id}/comment (exists)
	it('POST 200 /v1/media/{media_id}/comment (exists)', (done) => {
		API_MEDIA_request.post('/v1/media/1/comment')
			.set('Cookie', `token=${shared.random_user_token};`)
			.set('Content-Type', 'application/json')
			.send({
				content: 'Hello'
			})
			.expect(200)
			.end(done);
	});
	// ---------------------------------------- POST 200 /v1/media/{media_id}/comment (exists)
	it("POST 404 /v1/media/{media_id}/comment (doesn't exists)", (done) => {
		API_MEDIA_request.post('/v1/media/9999999/comment')
			.set('Cookie', `token=${shared.random_user_token};`)
			.set('Content-Type', 'application/json')
			.send({
				content: 'Hello'
			})
			.expect(404)
			.end(done);
	});
	// ---------------------------------------- POST 200 /v1/media/{media_id}/comment (comment too short)
	it('POST 400 /v1/media/{media_id}/comment (comment too short)', (done) => {
		API_MEDIA_request.post('/v1/media/1/comment')
			.set('Cookie', `token=${shared.random_user_token};`)
			.set('Content-Type', 'application/json')
			.send({
				content: 'H'
			})
			.expect(400)
			.end(done);
	});
	// ---------------------------------------- POST 200 /v1/media/{media_id}/comment (comment too long)
	it('POST 400 /v1/media/{media_id}/comment (comment too long)', (done) => {
		API_MEDIA_request.post('/v1/media/1/comment')
			.set('Cookie', `token=${shared.random_user_token};`)
			.set('Content-Type', 'application/json')
			.send({
				content: new Array(1000).fill('H').join('h')
			})
			.expect(400)
			.end(done);
	});
}
