import { localeFromCookie } from '$lib/i18n';
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = ({ event, resolve }) => {
	const cookieHeader = event.request.headers.get('cookie');
	if (cookieHeader) {
		event.params.locale = localeFromCookie(cookieHeader!);
	} else {
		event.params.locale = 'en';
	}

	return resolve(event);
};

export default handle;
