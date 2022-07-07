<!-- ========================= SCRIPT -->
<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';

	type MediaProps = {
		media: {
			id: number;
			description: string;
			duration?: number | null;
			genres: string[];
			names: { lang: string; title: string }[];
			rating?: number | null;
			thumbnail?: string | null;
			background?: string | null;
			type: string;
			year?: number | null;
		};
		torrents: {
			id: number;
			name: string;
			size?: string | null;
			leech: number;
			seed: number;
		}[];
		staffs: {
			id: number;
			name: string;
			thumbnail?: string | null;
			role?: string | null;
		}[];
		actors: {
			id: number;
			name: string;
			thumbnail?: string | null;
			character?: string | null;
		}[];
	};

	export const load: Load = async ({ params, fetch, session, stuff }) => {
		const url = browser
			? `http://localhost:7072/v1/media/get/${params.id}`
			: `http://api-search:7072/v1/media/get/${params.id}`;
		const response = await fetch(url, {
			method: 'GET',
			headers: { accept: 'application/json' }
		});

		let props: MediaProps | false = false;
		if (response.ok) {
			props = (await response.json()) as MediaProps;
		}

		return {
			status: response.status,
			redirect: !response.ok ? '/search' : undefined,
			props: { props }
		};
	};
</script>

<script lang="ts">
	import { browser } from '$app/env';
	import ArrowLeft from '../../../src/components/icons/ArrowLeft.svelte';
	import { fade } from 'svelte/transition';

	/// @ts-expect-error media is given as a prop
	export let props: MediaProps;
	const { media, torrents, staffs, actors } = props;

	const cover = media.thumbnail ? media.thumbnail : '/no_cover.png';

	const defaultTitle = media.names.find((name) => name.lang == '__')!;
	const userFavoriteTitle = (() => {
		const userLang = 'FR';
		const favoriteTitle = media.names.find((name) => name.lang == userLang);
		if (favoriteTitle) {
			return favoriteTitle;
		}
		return defaultTitle;
	})();

	const durationStr = (() => {
		if (!media.duration) {
			return '';
		}

		// Collect times
		const total = { hours: 0, minutes: 0 };
		let duration = media.duration;
		if (duration >= 60) {
			total.hours = Math.floor(duration / 60);
			duration -= total.hours * 60;
		}
		total.minutes = duration;

		// Convert to string
		if (total.hours > 0) {
			return `${total.hours}h${total.minutes > 0 ? ` ${total.minutes}m` : ''}`;
		}
		return `${total.minutes}m`;
	})();

	// Utility
	function seedColor(amount: number) {
		if (amount == 0) {
			return 'text-red-600';
		}
		if (amount <= 10) {
			return 'text-orange-600';
		}
		return 'text-green-600';
	}
</script>

<!-- ========================= HTML -->
<div class="flex flex-col w-full h-auto bg-black">
	<div
		class="header min-h-[30rem] flex-grow-0 border-b-stone-200 border-b"
		in:fade={{ duration: 150 }}
		style={`background-image: url("${media.background ? media.background : media.thumbnail}")`}
	>
		<div class="header-gradient">
			<div class="m-2 px-2 py-1 bg-white border-2 border-gray-400 rounded-sm inline-block">
				<a href="/search"><ArrowLeft /> Go Back</a>
			</div>
			<div class="flex flex-row w-11/12 md:w-8/12 lg:w-1/2 mx-auto">
				<img
					src={cover}
					alt={`${userFavoriteTitle} Cover`}
					in:fade={{ duration: 150, delay: 50 }}
					class="h-96 rounded-md flex-grow-0 flex-shrink-0"
				/>
				<div class="ml-4 text-white">
					<div class="text-3xl">{userFavoriteTitle.title}</div>
					{#if userFavoriteTitle.lang != '__'}
						<div class="text-xl opacity-80">{defaultTitle.title}</div>
					{/if}
					<div class="text-white mt-4">
						{#if media.year}
							{media.year}
						{/if}
						{#if media.year && media.genres.length > 0}
							<span class="mx-1">&#x2022;</span>
						{/if}
						{#if media.genres.length > 0}
							{media.genres.join(', ')}
						{/if}
						{#if (media.year || media.genres.length > 0) && media.duration}
							<span class="mx-1">&#x2022;</span>
						{/if}
						{#if media.duration}
							{durationStr}
						{/if}
					</div>
					<div class="text-white mt-4">{media.description}</div>
					{#if staffs.length > 0}
						<div class="text-lg mt-4">Staffs</div>
						<div class="max-h-16 overflow-hidden">
							{#each staffs as staff (staff.id + (staff.role ?? ''))}
								<div>{staff.id}</div>
								<div>{staff.name}</div>
								<div>{staff.role}</div>
								<div>{staff.thumbnail}</div>
							{/each}
						</div>
					{/if}
					{#if actors.length > 0}
						<div class="text-lg mt-4">Actors</div>
						<div class="max-h-16 overflow-hidden">
							{#each actors as actor (actor.id + (actor.character ?? ''))}
								<div>{actor.id}</div>
								<div>{actor.name}</div>
								<div>{actor.character}</div>
								<div>{actor.thumbnail}</div>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
	<div class="w-11/12 md:w-8/12 lg:w-1/2 mx-auto text-white mt-4">
		<h1 class="text-2xl">Torrents</h1>
		{#if torrents.length > 0}
			<table class="w-full">
				<thead>
					<tr>
						<td>Name</td>
						<td>Size</td>
						<td>Seed</td>
						<td>Leech</td>
						<td />
					</tr>
				</thead>
				<tbody>
					{#each torrents as torrent (torrent.id)}
						<tr>
							<td class="p-2 pl-0 truncate" title={torrent.name}>{torrent.name}</td>
							<td class="p-2">
								{#if torrent.size}
									{torrent.size}
								{/if}
							</td>
							<td class={`p-2 mx-2 ${seedColor(torrent.seed)}`}>{torrent.seed}</td>
							<td class="p-2 mx-2 text-red-600">{torrent.leech}</td>
							<td class="p-2">Watch</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{:else}
			<div>No torrents for this media, yet !</div>
		{/if}
	</div>
	<div class="w-11/12 md:w-8/12 lg:w-1/2 mx-auto text-white mt-4">
		<h1 class="text-2xl">Comments</h1>
		{#if [].length > 0}
			<div class="max-h-16 overflow-hidden">
				{#each actors as actor (actor.id + (actor.character ?? ''))}
					<div>{actor.id}</div>
					<div>{actor.name}</div>
					<div>{actor.character}</div>
					<div>{actor.thumbnail}</div>
				{/each}
			</div>
		{:else}
			<div>No comments on this media, yet, be the first one !</div>
		{/if}
	</div>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.header {
		@apply relative bg-center bg-cover bg-no-repeat;
		background-image: linear-gradient(to bottom right, rgba(0, 0, 0, 0.9), rgba(0, 0, 0, 0.64));
	}

	.header-gradient {
		@apply h-full w-full;
		background-image: linear-gradient(to bottom right, rgba(0, 0, 0, 0.9), rgba(0, 0, 0, 0.64));
	}
	/*
 */
</style>
