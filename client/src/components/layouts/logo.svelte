<!-- ========================= SCRIPT -->
<script lang="ts">
	import { page, session } from '$app/stores';
	import { goto } from '$app/navigation';
	
	export let alone = false;

	$: link_to_search = $session.user && !($page.url.pathname == '/search')

	const onClick = () => {
		if (link_to_search) {
			goto("/search")
		}
	};
</script>

<!-- ========================= HTML -->
<!-- {#if onClick != undefined} -->
<h1
	class:alone
	on:click={onClick}
	class="text-blue-300 text-5xl font-bold logo-title cursor-pointer select-none font-['Holtwood One SC']  text-clip"
	class:cursor-pointer={link_to_search}
>
	Hypertube
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
