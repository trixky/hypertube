<!-- ========================= SCRIPT -->
<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';
	import type { MediaProps, MediaTorrent } from '$types/Media';

	export const load: Load = async ({ params, fetch, session }) => {
		const url = apiMedia(`/v1/media/${params.id}/get`);
		const response = await fetch(url, {
			method: 'GET',
			credentials: 'include',
			headers: {
				accept: 'application/json',
				cookie: !browser ? `token=${session.token}; locale=${session.locale}` : ''
			}
		});
		if (!response.ok || response.status >= 400) {
			return {
				status: response.status > 0 ? response.status : 500
			};
		}

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
	import { onMount, tick } from 'svelte';
	import { fade } from 'svelte/transition';
	import { browser } from '$app/env';
	import { _ } from 'svelte-i18n';
	import { addUserTitle } from '$utils/media';
	import ArrowLeft from '$components/icons/ArrowLeft.svelte';
	import LazyLoad from '$components/lazy/LazyLoad.svelte';
	import Background from '$components/animations/Background.svelte';
	import RefreshPeers, { type RefreshResult } from './RefreshPeers.svelte';
	import Comments from './Comments.svelte';
	import { extractPalette } from '$utils/color';
	import Torrent from './Torrent.svelte';
	import Player from './Player.svelte';
	import { apiMedia, imageProxy } from '$utils/api';

	/// @ts-expect-error media is given as a prop
	export let props: MediaProps;
	let { media, torrents, staffs, actors, comments } = props;
	addUserTitle(media);

	const cover = media.thumbnail ? imageProxy(media.thumbnail) : '/no_cover.png';
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

	// Add full thumbnail to staffs and actors
	const cleanStaffs = staffs.map((staff) => ({
		id: staff.id,
		name: staff.name,
		thumbnail: imageProxy(staff.thumbnail!),
		roles: [staff.role!]
	}));
	const cleanActors = actors.map((actor) => ({
		id: actor.id,
		name: actor.name,
		thumbnail: imageProxy(actor.thumbnail!),
		characters: [actor.character!]
	}));

	// Background animation
	let backgroundAnimation: Background;
	let palette: string[] = [];

	// Utility
	/*function goBack(event: Event) {
		event.preventDefault();
		// TODO Avoid exit when opening the media page directly
		if (history.length > 1) {
			history.back();
		} else {
			goto('/search');
		}
	}*/

	// Peers refresh
	function onPeersRefresh(event: CustomEvent<RefreshResult[]>) {
		for (const result of event.detail) {
			let torrent = torrents.find((torrent) => torrent.id == result.torrentId);
			if (torrent) {
				torrent.seed = result.seed;
				torrent.leech = result.leech;
			}
		}
		torrents.sort((a, b) => b.seed - a.seed);
		torrents = torrents;
	}

	// Player
	let infoContainer: HTMLElement;
	let selectedTorrent: MediaTorrent | undefined;
	async function playTorrent(torrent: MediaTorrent) {
		selectedTorrent = torrent;
		// Add watched icon
		if (typeof selectedTorrent.position != 'number') {
			selectedTorrent.position = 0.0;
			torrents = torrents;
		}
		backgroundAnimation.stop();
		scrollPlayerIntoView();
	}
	async function scrollPlayerIntoView() {
		await tick();
		if (infoContainer) {
			infoContainer.scrollIntoView(true);
		}
	}

	function onPlayerClose() {
		selectedTorrent = undefined;
		backgroundAnimation.restart();
	}

	function onPlayerTimeUpdate(event: CustomEvent<number | undefined | null>) {
		if (selectedTorrent) {
			selectedTorrent.position = event.detail;
		}
	}

	// If there was un unrecoverable error, delete the torrent from the list and close the player
	function onPlayerError() {
		torrents = torrents.filter((torrent) => torrent.id != selectedTorrent?.id);
		selectedTorrent = undefined;
		backgroundAnimation.restart();
	}

	// Background gradient
	let loadingGradient = true;
	let gradientColor: string | undefined;
	let background: string | undefined;
	onMount(() => {
		if (browser) {
			// Load background first to have a clean fade-in
			let useBackground = media.background
				? media.background
				: media.thumbnail
				? media.thumbnail
				: undefined;
			if (useBackground) {
				useBackground = imageProxy(useBackground);
				const image = new Image();
				image.setAttribute('crossOrigin', 'anonymous');
				image.src = useBackground;
				image.addEventListener('load', () => {
					background = useBackground;
				});
			}

			// Load the image used for the gradient
			let useImage = media.thumbnail ? media.thumbnail : undefined;
			if (useImage && useImage != '') {
				useImage = imageProxy(useImage);
				const image = new Image();
				image.setAttribute('crossOrigin', 'anonymous');
				image.src = useImage;
				image.addEventListener('load', () => {
					const result = extractPalette(image);
					gradientColor = result?.color;
					if (result?.palette) {
						palette = result.palette;
					}
					if (backgroundAnimation) {
						backgroundAnimation.start();
					}
					loadingGradient = false;
				});
			} else {
				loadingGradient = false;
			}
		}
	});
</script>

<!-- ========================= HTML -->
<svelte:head>
	<title>
		{$_('title.media', { values: { title: media.userTitle ? media.userTitle : media.title } })}
	</title>
</svelte:head>
<div class="media-page flex flex-col w-full h-auto pb-4 bg-black">
	<div class="header min-h-[30rem] flex-grow-0 border-b-stone-200 border-b">
		{#if !loadingGradient}
			{#if background}
				<div
					class="header-image transition-all"
					in:fade|local={{ duration: 250 }}
					style={`background-image: url("${background}")`}
				/>
			{/if}
			{#if gradientColor}
				<div
					class="header-gradient transition-all"
					style={gradientColor ? `--gradient-color: ${gradientColor}` : ''}
				/>
			{/if}
		{:else}
			<div class="header-gradient transition-all" />
		{/if}
		<div
			class="absolute top-0 left-0 m-2 px-2 py-1 text-stone-200 inline-block hover:text-blue-500 transition-colors"
		>
			<a href="/search" class="cursor-pointer">
				<ArrowLeft />
				{$_('media.go_back')}
			</a>
		</div>
		<div
			class="relative flex flex-col md:flex-row justify-center items-center w-11/12 md:w-4/5 lg:w-1/2 mx-auto py-10"
		>
			<img
				src={cover}
				alt={`${media.userTitle ? media.userTitle : media.title} Cover`}
				in:fade|local={{ duration: 150, delay: 50 }}
				class="h-[500px] rounded-md flex-grow-0 "
			/>
			<div
				class="md:ml-8 max-w-full md:max-w-[348px] lg:max-w-[612px] xl:max-w-[720px] text-white transition-all"
			>
				<div class="text-3xl mt-4 lg:mt-0">{media.userTitle ? media.userTitle : media.title}</div>
				{#if media.userTitle}
					<div class="text-xl opacity-80">{media.title}</div>
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
				{#if media.rating}
					{@const rating = Math.round(media.rating * 10) / 10}
					<div>
						<div class="rating">
							<div class="flex items-center w-full">
								<div class="stars h-4" style={`--rating: ${rating};`} />
								<div class="text-sm">{rating}/10</div>
							</div>
						</div>
					</div>
				{/if}
				<div class="text-white mt-4">{media.description}</div>
				{#if cleanActors.length > 0}
					<div class="text-lg mt-4 mb-2">{$_('media.actors')}</div>
					<ol class="flex w-full pb-4 overflow-x-auto overflow-y-hidden">
						{#each cleanActors as actor (actor.id)}
							<LazyLoad tag="li" class="mr-6 last:mr-0 w-24 flex-shrink-0 h-full min-h-[1px]">
								<div
									class="h-24 w-24 rounded-full border-4 border-black border-opacity-80 bg-center bg-cover transition-all"
									style={`background-image: url("${actor.thumbnail}"); ${
										gradientColor ? `border-color: ${gradientColor}` : ''
									}`}
								/>
								<div class="font-medium">{actor.name}</div>
								<div class="opacity-80 text-sm truncate" title={actor.characters.join(', ')}>
									{actor.characters.join(', ')}
								</div>
							</LazyLoad>
						{/each}
					</ol>
				{/if}
				{#if cleanStaffs.length > 0}
					<div class="text-lg mt-4 mb-2">{$_('media.staffs')}</div>
					<ol class="flex w-full pb-4 overflow-x-auto overflow-y-hidden">
						{#each cleanStaffs as staff (staff.id)}
							<LazyLoad tag="li" class="mr-6 last:mr-0 w-24 flex-shrink-0 h-full min-h-[1px]">
								<div
									class="h-24 w-24 rounded-full border-4 border-black border-opacity-80 bg-center bg-cover transition-all"
									style={`background-image: url("${staff.thumbnail}"); ${
										gradientColor ? `border-color: ${gradientColor}` : ''
									}`}
								/>
								<div class="font-medium">{staff.name}</div>
								<div class="opacity-80 text-sm truncate" title={staff.roles.join(', ')}>
									{staff.roles.join(', ')}
								</div>
							</LazyLoad>
						{/each}
					</ol>
				{/if}
			</div>
		</div>
	</div>
	<div bind:this={infoContainer} class="relative flex-grow">
		<Background bind:this={backgroundAnimation} {palette} />
		{#if selectedTorrent}
			<Player
				torrent={selectedTorrent}
				on:open={scrollPlayerIntoView}
				on:focus={scrollPlayerIntoView}
				on:close={onPlayerClose}
				on:timeUpdate={onPlayerTimeUpdate}
				on:error={onPlayerError}
			/>
		{/if}
		<div class="w-11/12 md:w-4/5 lg:w-1/2 mx-auto text-white my-4 flex-grow relative z-10">
			<h1 class="flex justify-between items-center text-2xl mb-4">
				<span>Torrents</span>
				<div class="inline-block relative opacity-80">
					<RefreshPeers mediaId={media.id} on:refresh={onPeersRefresh} />
				</div>
			</h1>
			{#if torrents.length > 0}
				<div class="w-full">
					{#each torrents as torrent (torrent.id)}
						<Torrent
							{torrent}
							on:select={playTorrent.bind(null, torrent)}
							selected={selectedTorrent ? selectedTorrent.id == torrent.id : undefined}
						/>
					{/each}
				</div>
			{:else}
				<div>{$_('media.no_torrents')}</div>
			{/if}
			<Comments mediaId={media.id} list={comments} />
		</div>
	</div>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	@keyframes move-background {
		50% {
			background-position: 100% 50%;
		}
	}

	.header {
		@apply relative;
	}

	.header-image {
		@apply absolute w-full h-full z-0 bg-center bg-cover bg-no-repeat;
		--gradient-color: rgba(0, 0, 0, 0.7);
		background-image: linear-gradient(to bottom right, rgba(0, 0, 0, 0.9), var(--gradient-color));
		transition: background 150ms linear;
	}

	.header-gradient {
		@apply absolute w-full h-full z-0;
		--gradient-color: rgba(0, 0, 0, 0.7);
		background-image: linear-gradient(to bottom right, rgba(0, 0, 0, 0.9), var(--gradient-color));
		transition: background 150ms linear;
	}
</style>
