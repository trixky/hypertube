<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	import type { Load } from '@sveltejs/kit';
	import { waitLocale } from 'svelte-i18n';
	import { i18n } from '../../src/i18n';

	export const load: Load = async () => {
		await i18n();
		await waitLocale();
		return {};
	};
</script>

<script lang="ts">
	import '../app.css';
	import { browser } from '$app/env';
	import { not_connected } from '../utils/redirect';
	import Header from '../components/layouts/header.svelte';
	import Footer from '../components/layouts/footer.svelte';

	// not_connected(browser);
</script>

<!-- ========================= HTML -->

<Header />
<main class="min-h-[80%] flex flex-row content-start justify-start relative">
	<slot />
</main>
<Footer />
