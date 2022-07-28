<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	import type { Load } from '@sveltejs/kit';
	import { waitLocale } from 'svelte-i18n';
	import { i18n } from '$lib/i18n';

	export const load: Load = async ({ session }) => {
		// All other pages require to be logged in
		if (!session.user) {
			return {
				status: 302,
				redirect: '/login'
			};
		}
		await i18n(session);
		await waitLocale();
		return {};
	};
</script>

<script lang="ts">
	import '../app.css';
	import Header from '$components/layouts/header.svelte';
</script>

<!-- ========================= HTML -->
<div class="flex flex-col min-h-screen bg-black">
	<Header />
	<main class="flex-grow flex flex-row content-start justify-start relative">
		<slot />
	</main>
</div>
