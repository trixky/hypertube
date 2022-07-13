<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	import type { Load } from '@sveltejs/kit';
	import { waitLocale } from 'svelte-i18n';
	import { i18n } from '$lib/i18n';

	export const load: Load = async ({ params }) => {
		await i18n(params);
		await waitLocale();
		return {};
	};
</script>

<script lang="ts">
	import '../app.css';
	import { browser } from '$app/env';
	import { not_connected } from '../utils/redirect';
	import Header from '../components/layouts/header.svelte';

	not_connected(browser);
</script>

<!-- ========================= HTML -->

<Header />
<main class="media-page min-h-[90%] flex flex-row content-start justify-start relative">
	<slot />
</main>
