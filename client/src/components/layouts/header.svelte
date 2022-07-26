<!-- ========================= SCRIPT -->
<script lang="ts">
	import Logo from './logo.svelte';
	import * as cookies from '$utils/cookies';
	import { locale, _ } from 'svelte-i18n';
	import { page, session } from '$app/stores';
	import { scale } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import Logout from '$components/icons/Logout.svelte';

	$: user = $session.user;

	function setLocale(event: Event & { currentTarget: EventTarget & HTMLSelectElement }) {
		locale.set(event.currentTarget.value);
	}

	const logout = () => {
		cookies.deleteUserCookies();
		goto('/login');
		$session.token = undefined;
		$session.user = undefined;
	};
</script>

<!-- ========================= HTML -->
<header
	class="relative z-10 w-full bg-black flex flex-col md:flex-row flex-nowrap items-stretch md:h-[7rem] border-b-stone-200 border-b pb-4 md:pb-0"
>
	<div class="relative md:hidden w-full p-4 text-center">
		<Logo />
	</div>
	<div
		class="relative md:w-1/2 flex justify-center md:justify-start items-center text-white mt-2 px-4 md:pr-0"
	>
		{#if user && !($page.url.pathname == '/search')}
			<a class="border border-blue-400 rounded-md text-white p-2" href="/search" transition:scale>
				{$_('navigation.search')}
			</a>
		{/if}
	</div>
	<div class="hidden md:block absolute top-1/2 -translate-y-1/2 left-1/2 -translate-x-1/2">
		<Logo />
	</div>
	<div class="relative md:w-1/2 md:flex justify-end items-center px-4 md:pl-0">
		{#if user}
			<div class="flex flex-col">
				<div title={$_('language')}>
					<select
						name="language"
						id="language"
						class="p-1 bg-transparent rounded-md text-white border border-gray-400"
						on:change={setLocale}
					>
						<option value="en">English</option>
						<option value="fr">Francais</option>
					</select>
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
		{/if}
	</div>
</header>
