import { browser } from '$app/env';
import { writable } from 'svelte/store';

export type Genre = {
	id: number;
	name: string;
};

export const loaded = writable(false);
export const loading = writable(true);
function genresStore() {
	const { set, subscribe, update } = writable<Genre[]>([]);

	return {
		set,
		subscribe,
		update,
		async load(
			sender: (info: RequestInfo, init?: RequestInit | undefined) => Promise<Response> = fetch
		) {
			loading.set(true);
			set([]);

			const url = browser
				? `http://localhost:7072/v1/media/genres`
				: `http://api-media:7072/v1/media/genres`;
			const res = await sender(url, {
				method: 'GET',
				credentials: 'include',
				headers: { accept: 'application/json' }
			});
			if (res.ok) {
				const body = (await res.json()) as {
					genres: Genre[];
				};
				set(
					body.genres.sort((a, b) => {
						return a.name.localeCompare(b.name);
					})
				);
				loaded.set(true);
			}

			loading.set(true);
		}
	};
}

export const genres = genresStore();
