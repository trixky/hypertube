<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	import type { Load } from '@sveltejs/kit';
	import { waitLocale } from 'svelte-i18n';
	import { i18n } from '$lib/i18n';

	export const load: Load = async ({ params, session }) => {
		// Auth pages require to be logged out
		if (session.user) {
			return {
				status: 302,
				redirect: '/search'
			};
		}
		await i18n(params);
		await waitLocale();
		return {};
	};
</script>

<script lang="ts">
	import '../app.css';
</script>

<!-- ========================= HTML -->
<main
	class="h-full relative flex content-center justify-center linear-gradient bg-no-repeat bg-cover bg-center"
>
	<slot />
</main>

<!-- ========================= CSS -->
<style lang="postcss">
	.linear-gradient {
		background: radial-gradient(rgba(68, 0, 255, 0.2) 5%, black), url('/background.jpg') no-repeat;
		background-size: cover;
		background-position: center;
	}
</style>
