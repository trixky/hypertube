import { writable } from 'svelte/store';
import { addUserTitle } from '$utils/media';
import type { Result } from 'src/types/Media';
import { browser } from '$app/env';
import { apiMedia } from '$utils/api';

export const searching = writable(false);
export const loadingMore = writable(false);
function resultsStore() {
	let list: Result[] = [];
	const { set, subscribe, update } = writable<Result[]>([]);

	return {
		subscribe,
		set,
		update,
		setResults(results: Result[]) {
			list = results;
			return set(results);
		},
		append(results: Result[]) {
			// Avoid duplicate entries on DB update within pages
			for (const result of results) {
				if (!list.find((media) => media.id == result.id)) {
					list.push(result);
				}
			}
			return set(list);
		}
	};
}
export const results = resultsStore();
export const totalResults = writable<number>(0);

type SortColumns = 'id' | 'name' | 'year' | 'duration';
type SortOrder = 'ASC' | 'DESC';
type SearchStore = {
	hasResults: boolean;
	query?: string | null;
	page: number;
	startAt: number;
	year?: number | null;
	rating?: number | null;
	genres: number[];
	sortBy: SortColumns;
	sortOrder: SortOrder;
};

export const baseUrl = apiMedia(`/v1/media/search`);

export async function executeSearch(
	url: string,
	sender: (info: RequestInfo, init?: RequestInit | undefined) => Promise<Response> = fetch,
	session?: App.Session
) {
	const res = await sender(url.toString(), {
		method: 'GET',
		credentials: 'include',
		headers: {
			accept: 'application/json',
			cookie: !browser ? `token=${session?.token}; locale=${session?.locale}` : ''
		}
	});
	if (res.ok) {
		const body = (await res.json()) as {
			page: number;
			results: number;
			totalResults: number;
			medias: Result[];
		};
		return {
			response: res,
			results: body.medias.map(addUserTitle),
			totalResults: body.totalResults
		};
	}
	return {
		response: res,
		results: [],
		totalResults: 0
	};
}

function buildParams(store: SearchStore): string {
	const params = new URLSearchParams([
		['page', `${store.page}`],
		['sort_by', `${store.sortBy}`],
		['sort_order', `${store.sortOrder}`]
	]);
	if (store.query && store.query != '') {
		params.append('query', encodeURIComponent(store.query));
	}
	if (store.year) {
		params.append('year', encodeURIComponent(store.year));
	}
	if (store.rating) {
		params.append('rating', encodeURIComponent(store.rating));
	}
	for (const genre of store.genres) {
		params.append('genre_ids', encodeURIComponent(genre));
	}
	return params.toString();
}

export function searchStore() {
	const store: SearchStore = {
		hasResults: false,
		query: '',
		page: 1,
		startAt: 0,
		genres: [],
		sortBy: 'id',
		sortOrder: 'DESC'
	};
	const { subscribe, set, update } = writable<SearchStore>(store);

	function buildURL(to: string): URL {
		const url = new URL(to);
		url.search = buildParams(store);
		return url;
	}

	return {
		subscribe,
		set,
		update,
		setHasResults(value: boolean) {
			store.hasResults = value;
			return set(store);
		},
		setGenres(genres: number[]) {
			store.genres = genres;
			return set(store);
		},
		toggleSort() {
			if (store.sortOrder == 'ASC') {
				store.sortOrder = 'DESC';
			} else {
				store.sortOrder = 'ASC';
			}
			return set(store);
		},
		async execute() {
			searching.set(true);
			loadingMore.set(false);

			// Reset form
			store.page = 1;
			store.startAt = 0;
			set(store);
			results.setResults([]);

			// Send request
			const {
				response,
				results: searchResults,
				totalResults: searchTotalResults
			} = await executeSearch(buildURL(baseUrl).toString());
			if (response.ok) {
				results.setResults(searchResults);
				totalResults.set(searchTotalResults);
			}

			searching.set(false);
			return response.status;
		},
		async loadMore() {
			loadingMore.set(true);

			// Reset form
			store.startAt = store.page * 20;
			store.page = store.page + 1;
			set(store);

			// Send request
			const res = await fetch(buildURL(baseUrl), {
				method: 'GET',
				credentials: 'include',
				headers: { accept: 'application/json' }
			});
			if (res.ok) {
				const body = (await res.json()) as {
					page: number;
					results: number;
					totalResults: number;
					medias: Result[];
				};
				results.append(body.medias.map(addUserTitle));
			}

			loadingMore.set(false);
		}
	};
}
export const search = searchStore();
