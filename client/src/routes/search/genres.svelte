<!-- ========================= SCRIPT -->
<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { fly } from 'svelte/transition';
	import { _ } from 'svelte-i18n';
	import { clickOutside } from '$directives/clickOutside';
	import ChevronDown from '$components/icons/ChevronDown.svelte';
	import ChevronUp from '$components/icons/ChevronUp.svelte';
	import { loaded, loading, genres } from '$stores/genres';
	import { search } from '$stores/search';
	import Times from '$components/icons/Times.svelte';

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

	let maxHeight = 0;
	let details: HTMLElement | undefined;
	async function toggle() {
		if (!show) {
			offset = [wrapper.offsetLeft, wrapper.offsetTop];
			isOpen = true;
		}
		show = !show;
		await tick();
		if (details) {
			let rect = details.getBoundingClientRect();
			console.log('height + offset', rect.top + rect.height, rect);
			console.log('window height', window.outerHeight - 50);
			console.log('set to', window.outerHeight - rect.top - rect.height - 50);
			if (rect.top + rect.height - 50 > window.outerHeight) {
				maxHeight = window.outerHeight - rect.top - 50;
			}
		}
	}

	// * Value logic
	let selected: number[] = [];
	async function onChange() {
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
	style={`left: ${offset[0]}px; top: ${offset[1]}px`}
	on:click={toggle}
>
	<div class="flex justify-between items-center p-2" use:clickOutside on:clickOutside={close}>
		<span>
			{$_('search.form.genres.name')}
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
		<div
			bind:this={details}
			class="details"
			transition:fly={{ duration: 150 }}
			style={maxHeight > 0 ? `max-height: ${maxHeight}px` : ''}
			on:outroend={hide}
		>
			<div
				class="flex items-center p-2 border-b last:border-b-0 last:rounded-b-md border-slate-400 text-blue-600"
				on:click={clear}
			>
				<Times />
				{$_('search.form.genres.clear')}
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
						on:change={onChange}
					/>
					<label for={genre.name} class="inline-block flex-grow">{genre.name}</label>
				</div>
			{:else}
				{#if $loading}
					{$_('search.form.genres.loading')}
				{:else}
					{$_('search.form.genres.no_results')}
				{/if}
			{/each}
		</div>
	{/if}
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.wrapper {
		@apply inline-block min-w-[10rem] w-full md:w-auto rounded-md bg-transparent border border-slate-400 bg-slate-900 text-white cursor-pointer;
		transition: height 150ms ease-in-out, width 150ms ease-in-out;
	}

	.details {
		@apply flex flex-col text-white border-t border-slate-400 overflow-auto;
	}

	.details > div:nth-child(even) {
		@apply bg-slate-800 bg-opacity-80;
	}
</style>
