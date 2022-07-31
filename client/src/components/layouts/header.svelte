<!-- ========================= SCRIPT -->
<script lang="ts">
	import Logo from './logo.svelte';
	import * as cookies from '$utils/cookies';
	import { locale, _ } from 'svelte-i18n';
	import { session } from '$app/stores';
	import { goto } from '$app/navigation';
	import Logout from '$components/icons/Logout.svelte';
	import ProfilePicture from '../../components/profile/profile-picture.svelte';

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
	<div class="absolute top-0 left-1/2 -translate-x-1/2 mt-7">
		<Logo />
	</div>
	<div class="mb-4 mt-24 mx-6 flex flex-col md:float-right md:mt-4">
		{#if user}
			<div title={$_('language')} class="inline-block -translate-x-[6px]">
				<select
					name="language"
					id="language"
					class="p-1 bg-transparent text-white cursor-pointer"
					value={$locale}
					on:change={setLocale}
				>
					<option value="en" class="text-black" selected={$locale?.startsWith('en')}>English</option>
					<option value="fr" class="text-black" selected={$locale?.startsWith('fr')}>Francais</option>
				</select>
			</div>
			<div class="my-1 inline">
				<a href={`/users/${user.id}`} class="[&>*]:hover:underline">
					<p class="text-white inline pr-1 underline-offset-1">{user.username}</p>
					<a href="/users/{user.id}">
						<ProfilePicture user_id={user.id} tranlsate />
					</a>
				</a>
			</div>
			<button
				class="inline-block w-fit m-auto text-red-500 transition-all rounded-md hover:underline underline-offset-1"
				on:click|preventDefault={logout}
			>
				<span class="transition-all pr-3 ">{$_('logout')}</span><Logout />
			</button>
		{/if}
	</div>
</header>
