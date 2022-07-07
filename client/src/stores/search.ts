import { writable } from 'svelte/store';

type Result = {
	id: number;
	type: string;
	name: string;
	names: { lang: string; title: string }[];
	genres: string[];
	description: string;
	year: number | null;
	duration: number | null;
	thumbnail: string;
	rating: number | null;
};

export const searching = writable(true);
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

function buildParams(params: SearchStore): string {
	const result: string[] = [
		`page=${params.page}`,
		`sort_by=${params.sortBy}`,
		`sort_order=${params.sortOrder}`
	];
	if (params.query && params.query != '') {
		result.push(`query=${encodeURIComponent(params.query)}`);
	}
	if (params.year) {
		result.push(`year=${encodeURIComponent(params.year)}`);
	}
	if (params.rating) {
		result.push(`rating=${encodeURIComponent(params.rating)}`);
	}
	for (const genre of params.genres) {
		result.push(`${encodeURIComponent('genre_ids')}=${encodeURIComponent(genre)}`);
	}
	return result.join('&');
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

	return {
		subscribe,
		set,
		update,
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
			const params = buildParams(store);
			const res = await fetch(`http://localhost:7072/v1/media/search?${params}`, {
				method: 'GET',
				headers: { accept: 'application/json' }
			});
			if (res.ok) {
				const body = (await res.json()) as {
					page: number;
					results: number;
					totalResults: number;
					medias: Result[];
				};
				results.setResults(
					body.medias.map((media) => {
						media.name = media.names.find((name) => name.lang == '__')!.title;
						return media;
					})
				);
				totalResults.set(body.totalResults);
				store.hasResults = true;
				set(store);
			}

			searching.set(false);
		},
		async loadMore() {
			loadingMore.set(true);

			// Reset form
			store.startAt = store.page * 20;
			store.page = store.page + 1;
			set(store);

			// Send request
			const params = buildParams(store);
			const res = await fetch(`http://localhost:7072/v1/media/search?${params}`, {
				method: 'GET',
				headers: { accept: 'application/json' }
			});
			if (res.ok) {
				const body = (await res.json()) as {
					page: number;
					results: number;
					totalResults: number;
					medias: Result[];
				};
				results.append(
					body.medias.map((media) => {
						media.name = media.names.find((name) => name.lang == '__')!.title;
						return media;
					})
				);
			}

			loadingMore.set(false);
		}
	};
}
export const search = searchStore();
