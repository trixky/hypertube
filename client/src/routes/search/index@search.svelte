<!-- ========================= SCRIPT -->
<script lang="ts">
	import { fade } from 'svelte/transition';
	import Genres from './genres.svelte';

	let value: string = '';
	let results: string[] = [];
	let sorts: string[] = ['#', 'Year', 'Duration'];

	function onInput() {
		results = [];
		if (value == '') {
			return;
		}
		for (let index = 0; index < Math.round(Math.random() * 10000); index++) {
			results.push(`${Math.round(Math.random() * 100) / 10}`);
		}
	}
</script>

<!-- ========================= HTML -->
<div class="bg-slate-800 min-h-[90%] w-full flex-grow">
	<div class="flex w-full sticky top-0 p-4 bg-slate-800 z-10 border-b-2 border-blue-500">
		<div>
			<input type="text" class="input" placeholder="Search" bind:value on:input={onInput} />
			<label for="year" class="ml-4 ">Year</label>
			<input type="text" class="input w-14" placeholder="Year" name="year" min="0" />
			<label for="rating">Rating</label>
			<input
				type="text"
				class="input w-14"
				placeholder="Min. Rating"
				name="rating"
				min="0"
				max="10"
				step="0.1"
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
	{#if results.length == 0}
		<div class="w-full flex justify-center mt-8">
			<div class="text-5xl text-white">No results !</div>
		</div>
	{:else}
		<div class="result-wrapper p-4">
			{#each results as _}
				<div class="result overflow-hidden w-40 mx-auto">
					<div
						class="cover"
						transition:fade={{ duration: 150 }}
						style={`background-image: url(https://image.tmdb.org/t/p/w500/eHuGQ10FUzK1mdOY69wF5pGgEf5.jpg)`}
					>
						<div class="rating">
							<div class="flex justify-between items-center w-full">
								<div class="stars" style={`--rating: ${_};`} />
								<div class="text-sm">{_}/10</div>
							</div>
						</div>
					</div>
					<div class="text-white font-bold truncate" title={`Finding Nemo ${_}${_}${_}`}>
						Finding Nemo {_}{_}{_}
					</div>
					<div class="text-white text-sm opacity-80">2000</div>
				</div>
			{/each}
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

	/** 134x196 */

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
