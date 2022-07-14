<!-- ========================= SCRIPT -->
<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';
	import { fade } from 'svelte/transition';
	import { _ } from 'svelte-i18n';
	import Spinner from '$components/animations/spinner.svelte';
	import { searching, loadingMore, results, totalResults, search } from '$stores/search';
	import Genres from './genres.svelte';
	import SortAsc from '$components/icons/SortAsc.svelte';
	import SortDesc from '$components/icons/SortDesc.svelte';
	import LazyLoad from '$components/lazy/LazyLoad.svelte';
	import { imageUrl } from '$utils/image';

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
			console.log('loadMore', error);
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
<div class="bg-black min-h-[90%] w-full flex-grow">
	<div
		class="flex flex-col md:flex-row items-center w-full sticky top-0 p-4 bg-black z-10 border-b-2 border-blue-500"
	>
		<div>
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
			<Genres disabled={loading} class="lg:ml-2" />
		</div>
		<div class="flex-grow" />
		<div class="mt-2 lg:mt-0">
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
					class="result overflow-hidden w-40 min-h-[220px] mx-auto"
				>
					<div
						class="cover"
						in:fade={{ duration: 150, delay: (index - $search.startAt) * 10 }}
						style={`background-image: url(${cover})`}
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
				</LazyLoad>
			{/each}
			{#if $totalResults != $results.length}
				<div
					bind:this={loader}
					class="result overflow-hidden min-h-[14rem] w-40 mx-auto cursor-pointer"
					class:opacity-50={loading}
					on:click={loadMore}
				>
					<div class="flex w-full h-full justify-center items-center text-2xl text-white">
						{$_('search.load_more')}
					</div>
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
		@apply p-2 text-white;
	}

	.result-wrapper {
		@apply grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 2xl:grid-cols-9 3xl:grid-cols-11 gap-x-1 gap-y-2;
	}

	.cover {
		@apply h-56 w-40 rounded-md overflow-hidden cursor-pointer transition-all bg-center bg-cover relative;
	}
</style>
