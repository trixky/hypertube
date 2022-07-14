import type { Handle } from '@sveltejs/kit';
import { localeFromCookie } from '$lib/i18n';
import { extract_cookie, labels } from '$utils/cookies';

export const handle: Handle = ({ event, resolve }) => {
	const cookies = event.request.headers.get('cookie');
	if (cookies) {
		event.params.locale = localeFromCookie(cookies!);
		event.params.token = extract_cookie(cookies, labels.token) ?? '';
	} else {
		event.params.locale = 'en';
		event.params.token = '';
	}

	return resolve(event);
};

export default handle;
