<!-- ========================= SCRIPT -->
<script lang="ts">
	import Logo from './logo.svelte';
	import * as cookies from '$utils/cookies';
	import { locale, _ } from 'svelte-i18n';
	import Logout from '$components/icons/Logout.svelte';
	import { page, session } from '$app/stores';
	import { scale } from 'svelte/transition';
	import { goto } from '$app/navigation';

	$: user = $session.user!;

	const logout = () => {
		$session.token = undefined;
		$session.user = undefined;
		cookies.del_a_cookie(cookies.labels.token);
		cookies.del_a_cookie(cookies.labels.user_info);
		goto('/login');
	};
</script>

<!-- ========================= HTML -->
<header
	class="relative z-10 w-full bg-black flex flex-col md:flex-row flex-nowrap items-stretch md:h-[7rem] border-b-stone-200 border-b pb-4 md:pb-0"
>
	<div class="md:hidden w-full p-4 text-center">
		<Logo />
	</div>
	<div
		class="md:w-1/2 flex justify-center md:justify-start items-center text-white mt-2 px-4 md:pr-0"
	>
		{#if !($page.url.pathname == '/search')}
			<a class="border border-blue-400 rounded-md text-white p-2" href="/search" transition:scale>
				{$_('navigation.search')}
			</a>
		{/if}
	</div>
	<div class="hidden md:block absolute top-1/2 -translate-y-1/2 left-1/2 -translate-x-1/2">
		<Logo />
	</div>
	<div class="md:w-1/2 md:flex justify-end items-center px-4 md:pl-0">
		<div class="flex flex-col">
			<div title={$_('language')}>
				{#if $locale == 'fr'}
					<button on:click={() => locale.set('en')} class="text-white">English</button>
				{:else}
					<button on:click={() => locale.set('fr')} class="text-white">Francais</button>
				{/if}
			</div>

			<div>
				<a href={`/users/${user.id}`}>
					<p class="text-white inline">{user.username}</p>
					<img
						class="invert inline-block -translate-y-[1px] translate-x-1"
						src="/user.png"
						width="16px"
						height="16px"
						alt={$_('auth.my_profile')}
					/>
				</a>
			</div>

			<button
				class="flex items-center text-red-500 border border-red-100 py-1 px-2 mt-2 rounded-md hover:bg-red-700 transition-all hover:shadow-md shadow-red-900 hover:text-white"
				on:click|preventDefault={logout}
			>
				<Logout /> <span class="hover:text-white transition-all">{$_('logout')}</span>
			</button>
		</div>
	</div>
</header>
