import { locale } from 'svelte-i18n';
import type { Result } from '../../src/types/Media';

let userLang: string | undefined = undefined;
let userRegion: string | undefined = undefined;
locale.subscribe((value) => {
	if (value) {
		const values = value.split('-');
		userLang = values[0].toLocaleLowerCase();
		userRegion = values[1]?.toLocaleLowerCase();
	} else {
		userLang = undefined;
		userRegion = undefined;
	}
});

export function addUserTitle(media: Result): Result {
	media.title = media.names.find((name) => name.lang == '__')!.title;
	const favoriteTitle = media.names.find((name) => {
		const titleLocale = name.lang.toLocaleLowerCase();
		return (
			titleLocale == userLang ||
			titleLocale == userRegion ||
			((userLang == 'gb' || userRegion == 'gb') && titleLocale == 'gb')
		);
	});
	media.userTitle = favoriteTitle?.title;
	return media;
}
