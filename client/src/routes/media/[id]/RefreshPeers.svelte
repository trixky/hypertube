<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	export type RefreshResult = { torrentId: number; seed: number; leech: number };
</script>

<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';
	import { _ } from 'svelte-i18n';
	import Spinner from '$components/animations/spinner.svelte';

	export let mediaId: number;
	let loading = false;

	const dispatch = createEventDispatcher();
	async function refresh() {
		loading = true;
		const response = await fetch(`http://localhost:7072/v1/media/${mediaId}/refresh`, {
			method: 'GET',
			credentials: 'include',
			headers: { accept: 'application/json' }
		});
		if (response.ok) {
			const lines = await await response.text();
			let results: RefreshResult[] = [];
			for (const line of lines.trim().split('\n')) {
				try {
					const data = JSON.parse(line) as {
						result: RefreshResult;
					};
					if (
						data.result?.torrentId != undefined &&
						typeof data.result?.seed == 'number' &&
						typeof data.result?.leech == 'number'
					) {
						results.push(data.result);
					}
				} catch (error) {
					console.error('Failed to read line in sream response', error);
				}
			}
			dispatch('refresh', results);
		}
		loading = false;
	}
</script>

<!-- ========================= HTML -->

{#if loading}
	<div
		in:fade={{ duration: 350 }}
		out:fade={{ duration: 80 }}
		id="spinner"
		class="absolute inline-block -translate-x-full transition-all duration-[0.35s] opacity-50"
		style={`--tw-translate-x: calc(-100% - 4px);`}
	>
		<Spinner size={16} />
	</div>
{/if}
<button
	class="inline-block text-sm transition-all disabled:opacity-50"
	class:hover:text-opacity-90={!loading}
	disabled={loading}
	on:click={refresh}
>
	{$_('media.refresh_peers')}
</button>

<!-- ========================= CSS -->
<style lang="postcss">
</style>
