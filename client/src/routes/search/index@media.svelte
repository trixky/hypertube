<!-- ========================= SCRIPT -->
<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { browser } from '$app/env';
	import { fade } from 'svelte/transition';
	import Spinner from '../../../src/components/animations/spinner.svelte';
	import { searching, loadingMore, results, totalResults, search } from '../../stores/search';
	import Genres from './genres.svelte';
	import SortAsc from '../../../src/components/icons/SortAsc.svelte';
	import SortDesc from '../../../src/components/icons/SortDesc.svelte';

	let columns: { value: string; name: string }[] = [
		{ value: 'year', name: 'Year' },
		{ value: 'name', name: 'Name' },
		{ value: 'duration', name: 'Duration' },
		{ value: 'id', name: 'ID' }
	];

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
			search.execute();
		}, 200);
	}

	async function loadMore() {
		if (loading) {
			return;
		}

		await search.loadMore();
		onScroll();
	}

	function toggleSort() {
		search.toggleSort();
		search.execute();
	}

	function onScroll() {
		if (browser) {
			if ($loadingMore || $results.length == 0 || $totalResults == $results.length) {
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
				placeholder="Search"
				disabled={loading}
				bind:value={$search.query}
				on:input={debounceSearch}
			/>
			<label for="year" class="lg:ml-4">Year</label>
			<input
				type="number"
				class="input w-20 mb-2 lg:mb-0"
				placeholder="Year"
				name="year"
				min="0"
				max="9999"
				step="1"
				disabled={loading}
				bind:value={$search.year}
				on:input={debounceSearch}
			/>
			<label for="rating">Rating</label>
			<input
				type="number"
				class="input w-20"
				placeholder="Min. Rating"
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
			<label for="sort">Sort By</label>
			<select
				class="input"
				name="sort"
				disabled={loading}
				bind:value={$search.sortBy}
				on:input={debounceSearch}
			>
				{#each columns as column (column.name)}
					<option value={column.value}>{column.name}</option>
				{/each}
			</select>
			<div
				class="input inline-block ml-2 cursor-pointer"
				class:opacity-80={loading}
				on:click={toggleSort}
			>
				{#if $search.sortOrder == 'ASC'}
					Asc <SortAsc />
				{:else}
					Desc <SortDesc />
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
			<div class="text-5xl text-white">No results !</div>
		</div>
	{:else}
		<div class="result-wrapper p-4">
			{#each $results as result, index (result.id)}
				{@const cover = result.thumbnail ? result.thumbnail : '/no_cover.png'}
				<a href={`/media/${result.id}`} class="result overflow-hidden w-40 mx-auto">
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
					<div class="text-white font-bold truncate" title={result.name}>
						{result.name}
					</div>
					{#if result.year}
						<div class="text-white text-sm opacity-80">{result.year}</div>
					{/if}
				</a>
			{/each}
			{#if $totalResults != $results.length}
				<div
					bind:this={loader}
					class="result overflow-hidden min-h-[14rem] w-40 mx-auto cursor-pointer"
					class:opacity-50={loading}
					on:click={loadMore}
				>
					<div class="flex w-full h-full justify-center items-center text-2xl text-white	">
						Load More
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
		@apply grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 xl:grid-cols-8 2xl:grid-cols-12 gap-x-1 gap-y-2;
	}

	.cover {
		@apply h-56 w-40 rounded-md overflow-hidden cursor-pointer transition-all bg-center bg-cover relative;
	}

	.result {
		@apply block;
	}

	.result .rating {
		@apply opacity-0 transition-all absolute top-0 right-0 bottom-0 left-0 p-2 flex items-end text-white;
		background: rgb(0, 0, 0);
		background: linear-gradient(360deg, rgba(0, 0, 0, 0.8) 0%, rgba(0, 0, 0, 0) 100%);
	}

	.result .rating .stars {
		@apply inline-block text-sm;
		--percent: calc(var(--rating) / 10 * 100%);
		font-family: Times;
	}

	.result .rating .stars::before {
		@apply bg-clip-text;
		content: '★★★★★';
		letter-spacing: 3px;
		background: linear-gradient(90deg, #fc0 var(--percent), #fff var(--percent));
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	.result:hover .rating {
		@apply opacity-100;
	}
</style>
