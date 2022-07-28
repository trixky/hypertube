<!-- ========================= SCRIPT -->
<script lang="ts">
	import Spinner from '../animations/spinner.svelte';
	import { fade } from 'svelte/transition';

	export let name = '...';
	export let handler: () => unknown | Promise<unknown>;
	export let disabled = false;

	export let loading = false;

	async function toto() {
		if (!loading) {
			loading = true;
			await handler();
			loading = false;
		}
	}
</script>

<!-- ========================= HTML -->
<button
	on:click|preventDefault={toto}
	class:loading
	class:disabled
	class="relative block mt-6 mb-2 py-2 px-4 bg-blue-300 block m-auto rounded-sm hover:bg-blue-400 duration-[0.35s]"
	{disabled}
>
	{#if loading}
		<div
			in:fade|local={{ duration: 350, delay: 50 }}
			out:fade|local={{ duration: 80 }}
			id="spinner"
			class="absolute inline-block left-4 opacity-0 transition-all duration-[0.35s]"
		>
			<Spinner />
		</div>
	{/if}
	<p class="inline-block capitalize text-black">{name}</p>
</button>

<!-- ========================= CSS -->
<style lang="postcss">
	button.loading {
		@apply pl-12 bg-blue-400;
	}

	button.loading > #spinner {
		@apply opacity-100;
	}

	button.disabled {
		@apply bg-gray-300 !important;
	}
</style>
