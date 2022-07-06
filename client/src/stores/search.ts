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
			list.push(...results);
			return set(list);
		}
	};
}
export const results = resultsStore();
export const totalResults = writable<number>(0);

type SortColumns = 'id' | 'year' | 'duration';
type SortOrder = 'ASC' | 'DESC';
type SearchStore = {
	query?: string | null;
	page: number;
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
	if (params.genres.length > 0) {
		for (const genre of params.genres) {
			result.push(`${encodeURIComponent('genre_ids[]')}=${encodeURIComponent(genre)}`);
		}
	}
	return result.join('&');
}

export function searchStore() {
	const store: SearchStore = {
		query: '',
		page: 1,
		genres: [],
		sortBy: 'id',
		sortOrder: 'DESC'
	};
	const { subscribe, set, update } = writable<SearchStore>(store);

	let searchTimeout = 0;
	return {
		subscribe,
		set,
		update,
		setQuery(query: string | undefined | null) {
			if (query) {
				store.query = query;
			} else {
				store.query = '';
			}
			return set(store);
		},
		resetPage() {
			store.page = 1;
			return set(store);
		},
		nextPage() {
			store.page += 1;
			return set(store);
		},
		setYear(year: number | null | undefined) {
			if (year) {
				store.year = year;
			} else {
				store.year = undefined;
			}
			return set(store);
		},
		setRating(rating: number | null | undefined) {
			if (rating) {
				store.rating = rating;
			} else {
				store.rating = undefined;
			}
			return set(store);
		},
		addGenre(genre: number) {
			store.genres.push(genre);
			return set(store);
		},
		removeGenre(genre: number) {
			const index = store.genres.indexOf(genre);
			if (index >= 0) {
				store.genres.splice(index, 1);
			}
			return set(store);
		},
		setGenres(genres: number[]) {
			store.genres = genres;
			return set(store);
		},
		setSortBy(column: SortColumns) {
			if (column == 'id' || column == 'year' || column == 'duration') {
				store.sortBy = column;
			}
			return set(store);
		},
		setSortOrder(order: SortOrder) {
			if (order == 'ASC' || order == 'DESC') {
				store.sortOrder = order;
			}
			return set(store);
		},
		execute() {
			clearTimeout(searchTimeout);
			searchTimeout = setTimeout(async () => {
				searching.set(true);
				loadingMore.set(false);

				// Reset form
				store.page = 1;
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
				}

				searching.set(false);
			}, 200);
		},
		async loadMore() {
			loadingMore.set(true);

			// Reset form
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
