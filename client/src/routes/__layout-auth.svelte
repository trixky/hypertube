<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	import type { Load } from '@sveltejs/kit';
	import { waitLocale } from 'svelte-i18n';
	import { i18n } from '$lib/i18n';
	import DynamicBackground from '../components/animations/dynamic_background.svelte';

	export const load: Load = async ({ session }) => {
		// Auth pages require to be logged out
		if (session.user) {
			return {
				status: 302,
				redirect: '/search'
			};
		}
		await i18n(session);
		await waitLocale();
		return {};
	};
</script>

<script lang="ts">
	import '../app.css';
</script>

<!-- ========================= HTML -->
<main
	class="h-full relative flex content-center justify-center bg-no-repeat bg-cover bg-center"
>
	<DynamicBackground />
	<slot />
</main>
