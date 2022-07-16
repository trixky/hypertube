import type { ExternalFetch, GetSession, Handle } from '@sveltejs/kit';
import { localeFromCookie } from '$lib/i18n';
import { extract_cookie, get_user, labels } from '$utils/cookies';

export const handle: Handle = ({ event, resolve }) => {
	const cookies = event.request.headers.get('cookie');
	if (cookies) {
		event.locals.locale = localeFromCookie(cookies!);
		event.locals.token = extract_cookie(cookies, labels.token);
		event.locals.user = get_user(cookies);
	} else {
		event.locals.locale = 'en';
		event.locals.token = undefined;
		event.locals.user = undefined;
	}

	return resolve(event);
};

export const getSession: GetSession = (event) => {
	const session: App.Session = {
		locale: event.locals.locale
	};
	if (event.locals.user) {
		session.token = event.locals.token;
		session.user = event.locals.user;
	}
	return session;
};

export const externalFetch: ExternalFetch = (request) => {
	// TODO replace by env value
	// Replace SSR requests to the Media API with the "local" URL
	if (request.url.startsWith('http://localhost:7072/')) {
		request = new Request(
			request.url.replace('http://localhost:7072/', 'http://api-media:7072/'),
			request
		);
	}

	return fetch(request);
};