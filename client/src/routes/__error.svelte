<!-- ========================= SCRIPT -->
<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';

	export const load: Load = async ({ session, error, status }) => {
		return {
			props: {
				title: `${status}: ${error?.message}`
			}
		};
	};
</script>

<script lang="ts">
	import { session } from '$app/stores';

	export let title: string;
</script>

<!-- ========================= HTML -->
<div class="flex flex-col justify-center items-center w-full h-auto pb-4 bg-black text-white">
	<h1 class="block">{title}</h1>
	<p>
		{#if $session.user}
			<a href="/search">Search</a>
		{:else}
			<a href="/login">Login</a>
		{/if}
	</p>
</div>
