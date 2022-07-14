import type { GetSession, Handle } from '@sveltejs/kit';
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
	return event.locals.user
		? {
				locale: event.locals.locale,
				token: event.locals.token,
				user: event.locals.user
		  }
		: {};
};
