<!-- ========================= SCRIPT -->
<script lang="ts">
	import Spinner from '../animations/spinner.svelte';
	import { fade } from 'svelte/transition';
	import { uppercase_first_character } from '../../utils/str';

	export let name: string = '?';
	export let handler: Function = () => {};

	export let loading: boolean = false;

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
	class="relative mt-6 py-2 px-4 bg-blue-300 block m-auto rounded-sm duration-[0.35s]"
>
	{#if loading}
		<div
			in:fade={{ duration: 350, delay: 50 }}
			out:fade={{ duration: 80 }}
			id="spinner"
			class="absolute inline-block left-4 opacity-0 transition-all duration-[0.35s]"
		>
			<Spinner />
		</div>
	{/if}
	<p class="inline-block">{uppercase_first_character(name)}</p>
</button>

<!-- ========================= CSS -->
<style lang="postcss">
	button.loading {
		@apply pl-12 bg-blue-400;
	}

	button.loading > #spinner {
		@apply opacity-100;
	}

	button:hover {
		@apply bg-blue-400;
	}
</style>
