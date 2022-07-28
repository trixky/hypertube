<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	import type { Load } from '@sveltejs/kit';

	// Preload search results and genres
	// -- and insert them on startup in the client
	export const load: Load = async ({ fetch, session }) => {
		const { response: genresResponse, genres } = await getGenres(fetch, session);
		let notFound = genresResponse.status == 404;
		let forbidden = genresResponse.status >= 400 && genresResponse.status < 500 && !notFound;

		if (forbidden) {
			return {
				status: 403
			};
		} else if (notFound) {
			return {
				status: 404
			};
		} else if (genresResponse.status >= 500) {
			return {
				status: 500
			};
		}

		const {
			response: searchResponse,
			results,
			totalResults
		} = await executeSearch(baseUrl, fetch, session);
		notFound = searchResponse.status == 404;
		forbidden = searchResponse.status >= 400 && searchResponse.status < 500 && !notFound;

		if (forbidden) {
			return {
				status: 403
			};
		} else if (notFound) {
			return {
				status: 404
			};
		} else if (searchResponse.status >= 500) {
			return {
				status: 500
			};
		}

		return {
			status: 200,
			props: {
				ssrGenres: genres,
				ssrResults: results,
				ssrTotalResults: totalResults
			}
		};
	};
</script>

<script lang="ts">
	import { onMount } from 'svelte';
	import { _ } from 'svelte-i18n';
	import Spinner from '$components/animations/spinner.svelte';
	import {
		searching,
		loadingMore,
		results,
		totalResults,
		search,
		executeSearch,
		baseUrl
	} from '$stores/search';
	import SortAsc from '$components/icons/SortAsc.svelte';
	import SortDesc from '$components/icons/SortDesc.svelte';
	import ChevronDown from '$components/icons/ChevronDown.svelte';
	import ChevronUp from '$components/icons/ChevronUp.svelte';
	import { genres, getGenres, type Genre } from '$stores/genres';
	import Times from '$components/icons/Times.svelte';
	import type { Result } from '$types/Media';
	import MediaList from '$components/generics/MediaList.svelte';
	import { tweened } from 'svelte/motion';
	import { cubicOut } from 'svelte/easing';

	export let ssrGenres: Genre[];
	genres.set(ssrGenres);
	export let ssrResults: Result[];
	results.setResults(ssrResults);
	search.setHasResults(ssrResults.length > 0);
	searching.set(false);
	export let ssrTotalResults: number;
	totalResults.set(ssrTotalResults);

	let sortColumns: string[] = ['year', 'name', 'duration', 'id'];

	let loadMoreError = false;
	$: loading = $searching || $loadingMore;

	// Genres
	let genresOpen = false;
	let genresHeight = 0;
	let genresAnimations = tweened(genresHeight, {
		duration: 200,
		easing: cubicOut
	});
	$: genresAnimations.set(genresOpen ? genresHeight : 0);
	function toggleGenres() {
		genresOpen = !genresOpen;
	}

	let selected: number[] = [];
	async function onGenresChange() {
		search.setGenres(selected);
		search.execute();
	}

	function clearGenres() {
		toggleGenres();
		selected = [];
		onGenresChange();
	}

	// Store wrapper to update UI
	let searchTimeout = 0;
	function debounceSearch() {
		if (loading) {
			return;
		}

		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(async () => {
			loadMoreError = false;
			search.execute();
		}, 300);
	}

	async function loadMore() {
		if (loading || loadMoreError) {
			return;
		}

		try {
			await search.loadMore();
		} catch (error) {
			loadMoreError = true;
		}
	}

	function toggleSort() {
		loadMoreError = false;
		search.toggleSort();
		search.execute();
	}

	onMount(async () => {
		if (!$search.hasResults) {
			await search.execute();
		}
	});
</script>

<!-- ========================= HTML -->
<svelte:head>
	{#if !$search.query?.length}
		<title>
			{$_('title.search_empty', { values: { query: $search.query } })}
		</title>
	{:else}
		<title>
			{$_('title.search', { values: { query: $search.query } })}
		</title>
	{/if}
</svelte:head>
<div class="media-page bg-black min-h-[90%] w-full flex-grow">
	<div class="w-full sticky top-0 z-10 border-b-2 border-blue-500">
		<div class="flex flex-col md:flex-row items-center p-4">
			<div class="search-bar-bg" />
			<div class="relative">
				<input
					type="text"
					class="input block w-full mb-2 lg:inline-block lg:w-auto lg:mb-0"
					placeholder={$_('search.form.query_placeholder')}
					disabled={loading}
					bind:value={$search.query}
					on:input={debounceSearch}
				/>
				<label for="year" class="lg:ml-4">{$_('search.form.year')}</label>
				<input
					type="number"
					class="input w-20 mb-2 lg:mb-0"
					placeholder={$_('search.form.year')}
					name="year"
					min="0"
					max="9999"
					step="1"
					disabled={loading}
					bind:value={$search.year}
					on:input={debounceSearch}
				/>
				<label for="rating">{$_('search.form.rating')}</label>
				<input
					type="number"
					class="input w-20"
					placeholder={$_('search.form.rating_placeholder')}
					name="rating"
					min="0"
					max="10"
					step="0.1"
					disabled={loading}
					bind:value={$search.rating}
					on:input={debounceSearch}
				/>
				<button class="input w-full lg:w-auto lg:ml-4" on:click={toggleGenres}>
					<span>
						{$_('search.form.genres.name')}
						{#if $search.genres.length > 0}
							({$search.genres.length})
						{/if}
					</span>
					{#if genresOpen}
						<ChevronDown />
					{:else}
						<ChevronUp />
					{/if}
				</button>
			</div>
			<div class="relative flex-grow" />
			<div class="relative mt-2 lg:mt-0">
				<label for="sort">{$_('search.form.sort_by')}</label>
				<select
					class="input"
					name="sort"
					disabled={loading}
					bind:value={$search.sortBy}
					on:input={debounceSearch}
				>
					{#each sortColumns as column (column)}
						<option value={column}>{$_(`search.form.sort_columns.${column}`)}</option>
					{/each}
				</select>
				<div
					class="input inline-block ml-2 cursor-pointer"
					class:opacity-80={loading}
					on:click={toggleSort}
				>
					{#if $search.sortOrder == 'ASC'}
						{$_('asc.short')} <SortAsc />
					{:else}
						{$_('desc.short')} <SortDesc />
					{/if}
				</div>
			</div>
		</div>
		<div
			class="relative text-white border-t border-blue-500 overflow-hidden"
			style="height: {$genresAnimations}px;"
		>
			<div class="flex items-center flex-wrap p-4 pb-0" bind:offsetHeight={genresHeight}>
				<button
					class="inline-flex items-center text-red-500 border border-red-100 py-1 px-2 mb-2 mr-2 rounded-md hover:bg-red-700 transition-all hover:shadow-md shadow-red-900 hover:text-white"
					on:click={clearGenres}
				>
					<Times />
					{$_('search.form.genres.clear')}
				</button>
				{#each $genres as genre}
					<div class="inline-block px-2 mr-2 mb-2">
						<input
							type="checkbox"
							class="hidden peer"
							name="genres"
							id={genre.name}
							bind:group={selected}
							value={genre.id}
							on:change={onGenresChange}
						/>
						<label
							for={genre.name}
							class="inline-block flex-grow px-2 py-1 peer-checked:bg-green-700 bg-slate-600 rounded-md border border-gray-400 peer-checked:border-green-600 transition-colors cursor-pointer"
						>
							{genre.name}
						</label>
					</div>
				{/each}
			</div>
		</div>
	</div>
	{#if $searching}
		<div class="w-full flex justify-center mt-8 text-white">
			<Spinner size={96} />
		</div>
	{:else if $results.length == 0}
		<div class="w-full flex justify-center mt-8">
			<div class="text-5xl text-white">{$_('search.no_results')}</div>
		</div>
	{:else}
		<MediaList list={$results} totalResults={$totalResults} {loadMore} {loading} />
	{/if}
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.input {
		@apply p-2 rounded-md bg-transparent border border-slate-400 bg-slate-900 text-white;
	}

	label {
		@apply px-2 text-white;
	}

	.search-bar-bg {
		@apply absolute inset-0 bg-black bg-opacity-60 backdrop-blur-xl;
	}

	@-moz-document url-prefix() {
		.search-bar-bg {
			@apply bg-opacity-90;
		}
	}
</style>
