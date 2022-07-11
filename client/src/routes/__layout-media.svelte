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
	import Header from '../components/layouts/header.svelte';
</script>

<!-- ========================= HTML -->

<Header />
<main class="search min-h-[90%] flex flex-row content-start justify-start relative">
	<slot />
</main>
