<!-- ========================= SCRIPT -->
<script lang="ts">
	import { page, session } from '$app/stores';

	export let alone = false;

	$: link_to_search = $session.user && !($page.url.pathname == '/search');
</script>

<!-- ========================= HTML -->
<!-- {#if onClick != undefined} -->
<h1
	class:alone
	class="text-blue-300 text-5xl font-bold logo-title cursor-pointer select-none font-['Holtwood One SC'] text-clip"
	class:cursor-pointer={link_to_search}
>
	{#if link_to_search}
		<a href="/search"> Hypertube </a>
	{:else}
		Hypertube
	{/if}
</h1>

<!-- ========================= CSS -->
<style lang="postcss">
	.logo-title {
		text-shadow: 1px 1px 2px red, 0 0 1em blue, 0 0 0.2em blue;
	}

	.alone {
		@apply absolute left-1/2 -translate-x-1/2 -top-6;
	}
</style>
