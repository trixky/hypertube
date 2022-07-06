<!-- ========================= SCRIPT -->
<script lang="ts">
	import { fly } from 'svelte/transition';
	import { clickOutside } from '../../../src/directives/clickOutside';
	import ChevronDown from '../../../src/components/icons/ChevronDown.svelte';
	import ChevronUp from '../../../src/components/icons/ChevronUp.svelte';
	import { loaded, loading, genres } from '../../../src/stores/genres';
	import { onMount } from 'svelte';
	import { search } from '../../../src/stores/search';
	import Times from '../../../src/components/icons/Times.svelte';

	// * Popup logic
	export let disabled = false;
	let classes: string = '';
	let show = false;
	let isOpen = false;
	let wrapper: HTMLElement;
	let offset: [number, number] = [0, 0];

	function close() {
		show = false;
	}

	function hide() {
		isOpen = false;
	}

	function toggle() {
		if (!show) {
			offset = [wrapper.offsetLeft, wrapper.offsetTop];
			isOpen = true;
		}
		show = !show;
	}

	// * Value logic
	let selected: number[] = [];
	function onInput() {
		search.setGenres(selected);
		search.execute();
	}

	function clear() {
		toggle();
		selected = [];
		search.setGenres(selected);
		search.execute();
	}

	onMount(() => {
		if (!$loaded) {
			genres.load();
		}
	});

	export { classes as class };
</script>

<!-- ========================= HTML -->
<div
	bind:this={wrapper}
	class={`wrapper ${classes}`}
	class:absolute={isOpen}
	class:opacity-50={disabled}
	style={`left: ${offset[0] - 8}px; top: ${offset[1]}px`}
	on:click={toggle}
>
	<div class="flex justify-between items-center p-2" use:clickOutside on:clickOutside={close}>
		<span>
			Genres
			{#if selected.length > 0}
				({selected.length})
			{/if}
		</span>
		{#if isOpen}
			<ChevronDown />
		{:else}
			<ChevronUp />
		{/if}
	</div>
	{#if show}
		<div class="details" transition:fly={{ duration: 150 }} on:outroend={hide}>
			<div
				class="flex items-center p-2 border-b last:border-b-0 last:rounded-b-md border-slate-400 text-blue-600"
				on:click={clear}
			>
				<Times /> Clear
			</div>
			{#each $genres as genre (genre.id)}
				<div class="inline-block p-2 border-b last:border-b-0 last:rounded-b-md border-slate-400">
					<input
						type="checkbox"
						class="inline-block"
						name="genres"
						id={genre.name}
						bind:group={selected}
						value={genre.id}
						on:input={onInput}
					/>
					<label for={genre.name} class="inline-block flex-grow">{genre.name}</label>
				</div>
			{:else}
				{#if $loading}
					Loading...
				{:else}
					No genres, yet !
				{/if}
			{/each}
		</div>
	{/if}
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.wrapper {
		@apply inline-block min-w-[10rem] rounded-md bg-transparent border border-slate-400 bg-slate-900 text-white cursor-pointer;
		transition: height 150ms ease-in-out, width 150ms ease-in-out;
	}

	.details {
		@apply flex flex-col text-white border-t border-slate-400 overflow-hidden;
	}

	.details > div:nth-child(even) {
		@apply bg-slate-800 bg-opacity-80;
	}
</style>
