import { browser } from '$app/env';
import { writable } from 'svelte/store';

export type Genre = {
	id: number;
	name: string;
};

export async function getGenres(
	sender: (info: RequestInfo, init?: RequestInit | undefined) => Promise<Response> = fetch,
	session?: App.Session
) {
	const url = browser
		? `http://localhost:7072/v1/media/genres`
		: `http://api-media:7072/v1/media/genres`;
	const res = await sender(url, {
		method: 'GET',
		credentials: 'include',
		headers: {
			accept: 'application/json',
			cookie: !browser ? `token=${session?.token}; locale=${session?.locale}` : ''
		}
	});
	if (res.ok && res.status >= 200 && res.status < 300) {
		const body = (await res.json()) as {
			genres: Genre[];
		};
		return {
			response: res,
			genres: body.genres.sort((a, b) => {
				return a.name.localeCompare(b.name);
			})
		};
	}
	return {
		response: res,
		genres: []
	};
}

export const loaded = writable(false);
export const loading = writable(true);
function genresStore() {
	const { set, subscribe, update } = writable<Genre[]>([]);

	return {
		set,
		subscribe,
		update
	};
}

export const genres = genresStore();
