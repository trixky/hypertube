<!-- ========================= SCRIPT -->
<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';
	import { i18n } from '$lib/i18n';

	export const load: Load = async ({ session, error, status }) => {
		await i18n(session);
		await waitLocale();
		return {
			props: {
				status,
				error
			}
		};
	};
</script>

<script lang="ts">
	import { onMount } from 'svelte';
	import { page, session } from '$app/stores';
	import { waitLocale, _ } from 'svelte-i18n';
	import { browser } from '$app/env';
	import { goto } from '$app/navigation';
	import Background from '$components/animations/Background.svelte';
	import { deleteUserCookies } from '$utils/cookies';
	import Logout from '$components/icons/Logout.svelte';

	export let status: number;
	// export let error: Error | null;

	const refresh = () => {
		location.reload();
	};

	const logout = () => {
		deleteUserCookies();
		$session.token = undefined;
		$session.user = undefined;
		goto('/login');
	};

	$: errorMap = {
		400: $_('error.bad_request'),
		401: $_('error.unauthorized'),
		403: $_('error.forbidden'),
		404: $_('error.not_found'),
		500: $_('error.server_error')
	} as Record<number, string>;

	let background: Background;
	onMount(() => {
		if (browser && background) {
			background.start();
		}
	});
</script>

<!-- ========================= HTML -->
<svelte:head>
	<title>hypertube :: {status}</title>
</svelte:head>
<div class="relative w-full h-auto pb-4 bg-black text-white">
	<Background bind:this={background} palette={['rgb(147, 197, 253)', 'red', 'white	']} />
	<div class="relative flex flex-col justify-center items-center w-full h-full">
		<h1 class="text-6xl block">
			{#if errorMap[status]}
				{errorMap[status]}
			{:else}
				{status}
			{/if}
		</h1>
		<div class="flex mt-4">
			{#if $page.routeId == 'search@logged'}
				<button
					class="flex items-center text-md p-2 text-white bg-blue-500 rounded-md hover:bg-blue-600 hover:shadow-blue-600 hover:shadow-sm transition-all"
					on:click|preventDefault={refresh}
				>
					Refresh
				</button>
			{:else}
				<a
					href={$session.user ? '/search' : '/login'}
					class="flex items-center text-md p-2 text-white bg-blue-500 rounded-md hover:bg-blue-600 hover:shadow-blue-600 hover:shadow-sm transition-all"
				>
					{#if $session.user}
						{$_('error.back_to_search')}
					{:else}
						{$_('error.login_or_register')}
					{/if}
				</a>
			{/if}
			{#if status == 403}
				<button
					class="flex items-center text-white border border-red-600 p-2 ml-2 rounded-md bg-red-700 hover:bg-red-500 transition-all hover:shadow-md shadow-red-900 hover:text-white"
					on:click|preventDefault={logout}
				>
					<Logout /> <span class="hover:text-white transition-all">{$_('logout')}</span>
				</button>
			{/if}
		</div>
	</div>
</div>
