<!-- ========================= SCRIPT -->
<script lang="ts">
	import { browser } from '$app/env';
	import * as cookies from '$utils/cookies';
	import { disconnect } from '$utils/redirect';
	import { _ } from 'svelte-i18n';

	let connected = false;

	if (browser && cookies.get_a_cookie('token')) {
		connected = true;
	}
</script>

<!-- ========================= HTML -->
<footer
	class="w-full bg-black flex flex-col justify-center items-center h-20 min-h-[10%] border-t-stone-200 border-t"
>
	{#if connected}
		<button on:click={disconnect}>
			<p class="text-white">{$_('logout')}</p>
		</button>
	{/if}
</footer>
