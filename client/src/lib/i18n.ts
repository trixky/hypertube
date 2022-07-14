// import { browser } from '$app/env';
import { get } from 'svelte/store';
import { browser } from '$app/env';
import { init, getLocaleFromNavigator, addMessages, locale } from 'svelte-i18n';
import { add_a_cookie, del_a_cookie, extract_cookie } from '$utils/cookies';
import en from '../locales/en.json';
import fr from '../locales/fr.json';

export function localeFromCookie(cookies: string) {
	const value = extract_cookie(cookies, 'locale');
	if (value && value.match(/^(fr|en)(-\w+)?/)) {
		return value;
	}
	return 'en';
}

export function chooseLocale(session: App.Session): string {
	if (browser) {
		// Saved locale
		const cookieLocale = localeFromCookie(document.cookie);
		if (cookieLocale) {
			return cookieLocale;
		} else {
			del_a_cookie('locale');
		}
		// Locale from the browser
		const locale = getLocaleFromNavigator();
		if (locale) {
			return locale;
		}
		// Default locale
		return session?.locale || 'en';
	}
	// Locale from session (extracted from cookie in hooks)
	return session?.locale || 'en';
}

// * Cookies
// Always save the locale on update
if (browser) {
	let once = false;
	locale.subscribe((value) => {
		if (!once) {
			once = true;
			return;
		}
		if (value) {
			add_a_cookie('locale', value, 365);
		} else {
			del_a_cookie('locale');
		}
	});
}

export async function i18n(session: App.Session) {
	const locale = chooseLocale(session);

	addMessages('en', en);
	addMessages('fr', fr);

	return init({
		fallbackLocale: 'en',
		initialLocale: locale
	});
}
