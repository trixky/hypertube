<!-- ========================= SCRIPT -->
<script lang="ts">
	import Spinner from '../../../src/components/animations/spinner.svelte';
	import { onMount, tick } from 'svelte';
	import { fade } from 'svelte/transition';
	import Genres from './genres.svelte';

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

	let value: string = '';
	let results: Result[] = [];
	let page = 1;
	let totalResults = 0;
	let sorts: string[] = ['#', 'Year', 'Duration'];

	let searching = true;
	let smallLoading = false;
	$: loading = searching || smallLoading;
	let searchTimeout = 0;
	function onInput() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(async () => {
			searching = true;
			page = 1;
			results = [];
			let query: Record<string, any> = { page };
			if (value != '') {
				query.query = value;
			}
			let params = Object.keys(query)
				.map((key) => {
					return `${encodeURIComponent(key)}=${encodeURIComponent(query[key])}`;
				})
				.join('&');
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
				results = body.medias.map((media) => {
					media.name = media.names.find((name) => name.lang == '__')!.title;
					return media;
				});
				totalResults = body.totalResults;
			}
			searching = false;
		}, 200);
	}

	async function loadMore() {
		if (loading) {
			return;
		}

		smallLoading = true;
		page = page + 1;
		let query: Record<string, any> = { page };
		if (value != '') {
			query.query = value;
		}
		let params = Object.keys(query)
			.map((key) => {
				return `${encodeURIComponent(key)}=${encodeURIComponent(query[key])}`;
			})
			.join('&');
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
			results = [
				...results,
				...body.medias.map((media) => {
					media.name = media.names.find((name) => name.lang == '__')!.title;
					return media;
				})
			];
		}
		smallLoading = false;

		await tick();
		document.documentElement.scrollTop = document.documentElement.scrollHeight;
	}

	onMount(async () => {
		onInput();
	});
</script>

<!-- ========================= HTML -->
<div class="bg-black min-h-[90%] w-full flex-grow">
	<div class="flex w-full sticky top-0 p-4 bg-black z-10 border-b-2 border-blue-500">
		<div>
			<input
				type="text"
				class="input"
				placeholder="Search"
				disabled={loading}
				bind:value
				on:input={onInput}
			/>
			<label for="year" class="ml-4 ">Year</label>
			<input
				type="text"
				class="input w-14"
				placeholder="Year"
				name="year"
				min="0"
				disabled={loading}
			/>
			<label for="rating">Rating</label>
			<input
				type="text"
				class="input w-14"
				placeholder="Min. Rating"
				name="rating"
				min="0"
				max="10"
				step="0.1"
				disabled={loading}
			/>
			<Genres class="ml-2" />
		</div>
		<div class="flex-grow" />
		<div>
			<label for="sort">Sort By</label>
			<select
				name="sort"
				class="p-2 rounded-mdbg-transparent border border-slate-400 bg-slate-900 text-white"
			>
				<option value="id">ID</option>
				<option value="year">Year</option>
				<option value="duration">Duration</option>
			</select>
		</div>
	</div>
	{#if searching}
		<div class="w-full flex justify-center mt-8 text-white">
			<Spinner size={96} />
		</div>
	{:else if results.length == 0}
		<div class="w-full flex justify-center mt-8">
			<div class="text-5xl text-white">No results !</div>
		</div>
	{:else}
		<div class="result-wrapper p-4">
			{#each results as result (result.id)}
				{@const cover = result.thumbnail ? result.thumbnail : '/no_cover.png'}
				<div class="result overflow-hidden w-40 mx-auto">
					<div class="cover" in:fade={{ duration: 150 }} style={`background-image: url(${cover})`}>
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
				</div>
			{/each}
			{#if totalResults != results.length}
				<div
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
