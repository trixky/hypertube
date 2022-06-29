<!-- ========================= SCRIPT -->
<script lang="ts">
	import Spinner from '../animations/spinner.svelte';
	import { fade } from 'svelte/transition';

	export let name: string = '?';
	let loading: boolean = false;

	function toto() {
		loading = !loading;
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
	<p class="inline-block">{name.charAt(0).toUpperCase() + name.slice(1)}</p>
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
