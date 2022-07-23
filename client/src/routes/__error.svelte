<!-- ========================= SCRIPT -->
<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';
	import { i18n } from '$lib/i18n';

	export const load: Load = async ({ session, error, status }) => {
		await i18n(session);
		await waitLocale();
		return {
			props: {
				status,
				error
			}
		};
	};
</script>

<script lang="ts">
	import { onMount } from 'svelte';
	import { session } from '$app/stores';
	import { waitLocale, _ } from 'svelte-i18n';
	import { browser } from '$app/env';
	import Background from '$components/animations/Background.svelte';

	export let status: number;
	// export let error: Error | null;

	let background: Background;
	onMount(() => {
		if (browser && background) {
			background.start();
		}
	});
</script>

<!-- ========================= HTML -->
<div class="relative w-full h-auto pb-4 bg-black text-white">
	<Background bind:this={background} palette={['rgb(147, 197, 253)', 'red', 'white	']} />
	<div class="relative flex flex-col justify-center items-center w-full h-full">
		<h1 class="text-6xl block">
			{#if status == 404}
				{$_('error.not_found')}
			{:else}
				{status}
			{/if}
		</h1>
		<a
			href={$session.user ? '/search' : '/login'}
			class="text-md p-2 mt-2 text-white bg-blue-500 rounded-md hover:bg-blue-600 hover:shadow-blue-600 hover:shadow-sm transition-all"
		>
			{#if $session.user}
				{$_('error.back_to_search')}
			{:else}
				{$_('error.login_or_register')}
			{/if}
		</a>
	</div>
</div>
