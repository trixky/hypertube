// import { browser } from '$app/env';
import { init, getLocaleFromNavigator, addMessages, locale } from 'svelte-i18n';
import { add_a_cookie, del_a_cookie, extract_cookie } from '$utils/cookies';
import en from '../locales/en.json';
import fr from '../locales/fr.json';
import { browser } from '$app/env';

export function localeFromCookie(cookies: string) {
	const value = extract_cookie(cookies, 'locale');
	if (value && value.match(/^(fr|en)(-\w+)?/)) {
		return value;
	}
	return 'en';
}

export function chooseLocale(params?: Record<string, string>): string {
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
		return params?.locale || 'en';
	}
	// Locale from params (extracted from cookie in hooks)
	console.log('param', params, params?.locale);
	return params?.locale || 'en';
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

export async function i18n(params?: Record<string, string>) {
	const locale = chooseLocale(params);

	addMessages('en', en);
	addMessages('fr', fr);

	return init({
		fallbackLocale: 'en',
		initialLocale: locale
	});
}
