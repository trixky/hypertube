<!-- ========================= SCRIPT -->
<script lang="ts">
	import Logo from './logo.svelte';
	import * as cookies from '$utils/cookies';
	import { locale, _ } from 'svelte-i18n';
	import { session } from '$app/stores';
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
<header class="relative text-center w-full bg-black">
	<!-- <header
		class="relative w-full bg-black flex flex-col md:flex-row flex-nowrap items-stretch md:h-[7rem] pb-4 md:pb-0"
	> -->
	<!-- <div
		class="relative md:w-1/2 flex justify-center md:justify-start items-center text-white mt-2 px-4 md:pr-0"
	>
		{#if user && !($page.url.pathname == '/search')}
			<a class="border border-blue-400 rounded-md text-white p-2" href="/search" transition:scale>
				{$_('navigation.search')}
			</a>
		{/if}
	</div> -->
	<div class="absolute top-0 left-1/2 -translate-x-1/2 mt-7">
		<Logo />
	</div>
	<div class="mb-4 mt-24 mx-6 flex flex-col md:float-right md:mt-4">
		<!-- <div class="float-right mt-8 mr-10 relative"> -->
		{#if user}
			<!-- <div class="flex flex-col content-end text-right"> -->
			<!-- <div class="flex flex-col content-end"> -->
			<div title={$_('language')} class="inline-block -translate-x-[6px]">
				<select
					name="language"
					id="language"
					class="p-1 bg-transparent text-white cursor-pointer"
					value={$locale}
					on:change={setLocale}
				>
					<option value="en" selected={$locale?.startsWith('en')}>English</option>
					<option value="fr" selected={$locale?.startsWith('fr')}>Francais</option>
				</select>
			</div>
			<div class="my-1 inline">
				<a href={`/users/${user.id}`} class="[&>*]:hover:underline">
					<p class="text-white inline pr-1 underline-offset-1">{user.username}</p>
					<img
						class="invert inline"
						src="/user.png"
						width="16px"
						height="16px"
						alt={$_('auth.my_profile')}
					/>
				</a>
			</div>
			<button
				class="inline text-red-500 transition-all rounded-md hover:underline underline-offset-1"
				on:click|preventDefault={logout}
			>
				<span class="transition-all pr-3 ">{$_('logout')}</span><Logout />
			</button>
			<!-- </div> -->
		{/if}
	</div>
</header>
