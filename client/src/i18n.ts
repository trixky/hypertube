import { init, getLocaleFromNavigator, addMessages } from 'svelte-i18n';

import en from './locales/en.json';
import fr from './locales/fr.json';

export function i18n() {
	addMessages('en', en);
	addMessages('fr', fr);

	return init({
		fallbackLocale: 'en',
		initialLocale: getLocaleFromNavigator()
	});
}
