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
				status: 302,
				redirect: '/login'
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
				status: 302,
				redirect: '/login'
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
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';
	import { fade } from 'svelte/transition';
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
	import LazyLoad from '$components/lazy/LazyLoad.svelte';
	import { imageUrl } from '$utils/image';
	import Eye from '$components/icons/Eye.svelte';
	import ChevronDown from '$components/icons/ChevronDown.svelte';
	import ChevronUp from '$components/icons/ChevronUp.svelte';
	import { genres, getGenres, type Genre } from '$stores/genres';
	import Times from '$components/icons/Times.svelte';
	import { accordion } from '$directives/accordion';
	import type { Result } from '$types/Media';

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

	// * Infinite loader
	// Obser the Load More card when it's visible and loadMore if the user can see it
	let observer: IntersectionObserver;
	function onIntersectionEvent(entries: IntersectionObserverEntry[]) {
		if (loading || $results.length == 0 || $totalResults == $results.length) {
			return;
		}
		for (const entry of entries) {
			if (entry.isIntersecting) {
				loadMore();
			}
		}
	}
	if (browser) {
		observer = new IntersectionObserver(onIntersectionEvent, { threshold: 0 });
	}
	let loader: HTMLElement | undefined;
	let observing: HTMLElement | undefined;
	$: {
		if (loader) {
			observer.observe(loader);
			observing = loader;
		} else if (observing) {
			observer.unobserve(observing);
			observing = undefined;
		}
	}

	// Genres
	let genresOpen = false;
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
			onScroll();
		} catch (error) {
			loadMoreError = true;
		}
	}

	function toggleSort() {
		loadMoreError = false;
		search.toggleSort();
		search.execute();
	}

	function onScroll() {
		if (browser) {
			if (
				$loadingMore ||
				$results.length == 0 ||
				$totalResults == $results.length ||
				loadMoreError
			) {
				return;
			}
			const element = document.documentElement;
			const offset = element.scrollHeight - element.clientHeight - element.scrollTop;
			if (offset <= 100) {
				loadMore();
			}
		}
	}

	onMount(async () => {
		if (!$search.hasResults) {
			await search.execute();
		} else {
			onScroll();
		}
		if (browser) {
			window.addEventListener('scroll', onScroll, { passive: true });
			window.addEventListener('resize', onScroll, { passive: true });
		}
	});

	onDestroy(() => {
		if (browser) {
			window.removeEventListener('scroll', onScroll);
			window.removeEventListener('resize', onScroll);
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
		<div class="relative text-white border-t border-blue-500" use:accordion={genresOpen}>
			<div class="flex items-center flex-wrap p-4 pb-0">
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
		<div class="result-wrapper p-4">
			{#each $results as result, index (result.id)}
				{@const cover = result.thumbnail ? imageUrl(result.thumbnail) : '/no_cover.png'}
				<LazyLoad
					tag="a"
					href={`/media/${result.id}`}
					class="relative result overflow-hidden h-[268px] w-40 min-h-[268px] mx-auto"
				>
					<div
						class="cover"
						class:opacity-80={result.watched}
						style={`background-image: url(${cover})`}
						in:fade={{ duration: 150, delay: (index - $search.startAt) * 10 }}
					>
						{#if result.rating}
							{@const rating = Math.round(result.rating * 10) / 10}
							<div class="rating">
								<div class="flex justify-between items-center w-full">
									<div class="stars" style={`--rating: ${rating};`} />
									<div class="text-sm">{rating}/10</div>
								</div>
							</div>
						{/if}
					</div>
					<div
						class="text-white font-bold truncate"
						title={result.userTitle ? result.userTitle : result.title}
					>
						{result.userTitle ? result.userTitle : result.title}
					</div>
					{#if result.year}
						<div class="text-white text-sm opacity-80">{result.year}</div>
					{/if}
					{#if result.watched}
						<div class="absolute bottom-1 right-1 text-white">
							<Eye />
						</div>
					{/if}
				</LazyLoad>
			{/each}
			{#if $totalResults != $results.length}
				<div
					bind:this={loader}
					class="result overflow-hidden min-h-[14rem] w-40 mx-auto cursor-pointer text-white"
					class:opacity-50={loading}
					on:click={loadMore}
				>
					{#if loading}
						<div class="flex w-full h-full justify-center items-center text-2xl text-white">
							{$_('search.loading_more')}
							<Spinner size={32} />
						</div>
					{:else}
						<div class="flex w-full h-full justify-center items-center text-2xl text-white">
							{$_('search.load_more')}
						</div>
					{/if}
				</div>
			{/if}
		</div>
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

	.result-wrapper {
		@apply grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 2xl:grid-cols-9 3xl:grid-cols-11 gap-x-1 gap-y-2;
	}

	.cover {
		@apply h-56 w-40 rounded-md overflow-hidden cursor-pointer transition-all bg-center bg-cover relative;
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
