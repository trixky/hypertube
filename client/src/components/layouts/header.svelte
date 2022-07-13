<!-- ========================= SCRIPT -->
<script lang="ts">
	import Logo from './logo.svelte';
	import { goto } from '$app/navigation';
	import { browser } from '$app/env';
	import { me_store } from '$stores/me';
	import { disconnect } from '$utils/redirect';
	import * as cookies from '$utils/cookies';
	import { _ } from 'svelte-i18n';

	me_store.refresh_from_cookies();

	function handle_click_on_me() {
		if ($me_store.username.length > 0) {
			goto('/users/' + $me_store.id);
		}
	}

	let connected = false;
	if (browser && cookies.get_a_cookie('token')) {
		connected = true;
	}
</script>

<!-- ========================= HTML -->
<header
	class="relative z-10 w-full bg-black flex justify-center items-center h-20 min-h-[10%] border-b-stone-200 border-b"
>
	<Logo />
	{#if $me_store.username.length > 0}
		<button class="absolute right-8 top-7" on:click={handle_click_on_me}>
			<p class="text-white inline">{$me_store.username}</p>
			<img
				class="invert inline-block -translate-y-[1px] translate-x-1"
				src="/user.png"
				width="16px"
				height="16px"
				alt="user icone"
			/>
		</button>
	{/if}
	{#if connected}
		<button on:click={disconnect}>
			<p class="text-white">{$_('logout')}</p>
		</button>
	{/if}
</header>
