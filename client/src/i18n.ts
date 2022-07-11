// import { browser } from '$app/env';
import { init, getLocaleFromNavigator, addMessages, locale } from 'svelte-i18n';
import { add_a_cookie, del_a_cookie, get_a_cookie } from './utils/cookies';
import en from './locales/en.json';
import fr from './locales/fr.json';

export async function i18n() {
	addMessages('en', en);
	addMessages('fr', fr);

	// * localStorage
	/*if (browser) {
		// Always save the locale on update
		let once = false;
		locale.subscribe((value) => {
			if (!once) {
				once = true;
				return;
			}
			if (value) {
				localStorage.setItem('locale', value);
			} else {
				localStorage.removeItem('locale');
			}
		});

		// Check if the user already has a locale saved
		const storageLocale = localStorage.getItem('locale');
		if (storageLocale && storageLocale.match(/^(fr|en)(-\w+)?/)) {
			return init({
				fallbackLocale: 'en',
				initialLocale: storageLocale
			});
		} else {
			localStorage.removeItem('locale');
		}

		// -- else try to guess and save the locale
		await init({
			fallbackLocale: 'en',
			initialLocale: getLocaleFromNavigator()
		});
	}*/

	// * Cookies
	// Always save the locale on update
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

	// Check if the user already has a locale saved
	const cookieLocale = get_a_cookie('locale');
	if (cookieLocale && cookieLocale.match(/^(fr|en)(-\w+)?/)) {
		return init({
			fallbackLocale: 'en',
			initialLocale: cookieLocale
		});
	} else {
		del_a_cookie('locale');
	}

	// -- else try to guess and save the locale
	await init({
		fallbackLocale: 'en',
		initialLocale: getLocaleFromNavigator()
	});
}
