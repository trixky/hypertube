<!-- ========================= SCRIPT -->
<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';
	import type { MediaProps, MediaTorrent } from '../../../../src/types/Media';

	export const load: Load = async ({ params, fetch, session }) => {
		const url = `http://localhost:7072/v1/media/${params.id}/get`;
		const response = await fetch(url, {
			method: 'GET',
			credentials: 'include',
			headers: {
				accept: 'application/json',
				cookie: !browser ? `token=${session.token}; locale=${session.locale}` : ''
			}
		});

		let props: MediaProps | false = false;
		if (response.ok) {
			props = (await response.json()) as MediaProps;
		}

		const notFound = response.status == 404;
		const forbidden = response.status >= 400 && response.status < 500 && !notFound;

		if (forbidden) {
			return {
				status: 302,
				redirect: '/login'
			};
		}
		return {
			status: response.status,
			redirect: notFound || !response.ok ? '/search' : undefined,
			props: { props }
		};
	};
</script>

<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import { browser } from '$app/env';
	import { goto } from '$app/navigation';
	import { locale, _ } from 'svelte-i18n';
	import { addUserTitle } from '$utils/media';
	import ArrowLeft from '$components/icons/ArrowLeft.svelte';
	import LazyLoad from '$components/lazy/LazyLoad.svelte';
	import Background from '$components/animations/Background.svelte';
	import RefreshPeers, { type RefreshResult } from './RefreshPeers.svelte';
	import { imageUrl } from '$utils/image';
	import Comments from './Comments.svelte';
	import { extractPalette } from '$utils/color';
	import Torrent from './Torrent.svelte';
	import { session } from '$app/stores';

	/// @ts-expect-error media is given as a prop
	export let props: MediaProps;
	let { media, torrents, staffs, actors, comments } = props;
	addUserTitle(media);

	const cover = media.thumbnail ? imageUrl(media.thumbnail) : '/no_cover.png';
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
		thumbnail: imageUrl(staff.thumbnail!),
		roles: [staff.role!]
	}));
	const cleanActors = actors.map((actor) => ({
		id: actor.id,
		name: actor.name,
		thumbnail: imageUrl(actor.thumbnail!),
		characters: [actor.character!]
	}));

	// Background animation
	let backgroundAnimation: Background;
	let palette: string[] = [];

	// Utility
	function goBack(event: Event) {
		event.preventDefault();
		// TODO Avoid exit when opening the media page directly
		if (history.length > 1) {
			history.back();
		} else {
			goto('/search');
		}
	}

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
	const POSITION_DELAY_MS = 15_000; /* 15sec */
	const ERROR_TIMEOUT = 300_000; /* 5min */

	// @source https://stackoverflow.com/a/41776483
	function seekToTime(player: HTMLVideoElement, ts: number) {
		// try and avoid pauses after seeking
		player.pause();
		player.currentTime = ts; // if this is far enough away from current, it implies a "play" call as well...oddly. I mean seriously that is junk.
		// however if it close enough, then we need to call play manually
		// some shenanigans to try and work around this:
		let timer = setInterval(function () {
			if ((player.paused && player.readyState == 4) || !player.paused) {
				player.play();
				clearInterval(timer);
			}
		}, 50);
	}

	let infoContainer: HTMLElement;
	let play: number | undefined;
	let player: HTMLVideoElement | undefined;
	let playMessage: { type: 'info' | 'error'; content: string } | undefined;
	let updateErrored: number | null = null;
	let lastUpdate: number = 0;
	async function playTorrent(torrent: MediaTorrent) {
		backgroundAnimation.stop();

		// Open player
		play = torrent.id;
		playMessage = { type: 'info', content: 'Video is loading...' };
		// Infer media lang from the torrent name
		let mediaLang = 'en';
		if (torrent.name.match(/(?:vost)?fr(?:ench|an[cç]ais)?/i)) {
			mediaLang = 'fr';
		}

		// Wait for the player to exist and bind events
		await tick();
		if (player) {
			// Load subtitles
			try {
				const subtitlesUrl = `http://localhost:3030/torrent/${play}/subtitles`;
				fetch(subtitlesUrl, { method: 'GET' })
					.then(async (response) => {
						if (response.ok) {
							const body = (await response.json()) as
								| { subtitles: { id: number; lang: string }[] }
								| { error: string };
							if ('error' in body) {
								// TODO playMessage: failed to load subtitles
							} else {
								let index = 0;
								for (const subtitle of body.subtitles) {
									const track = document.createElement('track');
									track.kind = 'captions';
									track.label = subtitle.lang == 'fr' ? 'Francais' : 'English';
									track.srclang = subtitle.lang;
									track.src = `http://localhost:3030/subtitles/${subtitle.id}`;
									player?.appendChild(track);
									if (
										!$session.locale!.startsWith(mediaLang) &&
										$session.locale!.startsWith(subtitle.lang)
									) {
										/// @ts-expect-error Track element *has* the showing property
										track.mode = 'showing';
										player!.textTracks[index].mode = 'showing';
									}
									index++;
								}
							}
						}
					})
					.catch((error) => {
						console.error(error);
						// TODO playMessage: failed to load subtitles
					});
			} catch (error) {
				console.error(error);
				// TODO playMessage: failed to load subtitles
			}

			// Bind
			player.addEventListener('loadedmetadata', () => {
				playMessage = undefined;
				if (torrent.position) {
					seekToTime(player!, torrent.position);
				}
			});
			let updatingPosition = false;
			player.addEventListener('timeupdate', (event) => {
				if (!updatingPosition && event.timeStamp - lastUpdate >= POSITION_DELAY_MS) {
					if (!updateErrored || Date.now() - updateErrored >= ERROR_TIMEOUT /* 5min */) {
						lastUpdate = event.timeStamp;
						const currentTime = player?.currentTime;
						updatingPosition = true;
						fetch(`http://localhost:3040/v1/position/${play}`, {
							method: 'POST',
							credentials: 'include',
							headers: {
								'Content-Type': 'application/json'
							},
							body: JSON.stringify({
								position: currentTime
							})
						})
							.then((response) => {
								updatingPosition = false;
								if (!response.ok || response.status >= 400) {
									throw new Error('Response status is not ok');
								}
							})
							.catch((error) => {
								updatingPosition = false;
								console.error(error);
								updateErrored = Date.now();
							});
					} else {
						lastUpdate = event.timeStamp + Date.now() - updateErrored;
					}
				}
			});
		}
		if (infoContainer) {
			infoContainer.scrollIntoView(true);
		}
	}

	function closePlayer() {
		play = undefined;
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
				useBackground = imageUrl(useBackground);
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
				useImage = imageUrl(useImage);
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
<div class="flex flex-col w-full h-auto pb-4 bg-black">
	<div class="header min-h-[30rem] flex-grow-0 border-b-stone-200 border-b">
		{#if !loadingGradient}
			{#if background}
				<div
					class="header-image transition-all"
					in:fade={{ duration: 250 }}
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
			<a href="/search" on:click={goBack} class="cursor-pointer"
				><ArrowLeft /> {$_('media.go_back')}</a
			>
		</div>
		<div
			class="relative flex flex-col md:flex-row justify-center items-center w-11/12 md:w-4/5 lg:w-1/2 mx-auto py-10"
		>
			<img
				src={cover}
				alt={`${media.userTitle ? media.userTitle : media.title} Cover`}
				in:fade={{ duration: 150, delay: 50 }}
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
	<div class="relative flex-grow">
		<Background bind:this={backgroundAnimation} {palette} />
		<div
			bind:this={infoContainer}
			class="w-11/12 md:w-4/5 lg:w-1/2 mx-auto text-white my-4 flex-grow relative"
		>
			<div>
				{#if play}
					<div class="mb-10" transition:fade>
						<video
							bind:this={player}
							class="w-full"
							src={`http://localhost:3030/torrent/${play}/stream`}
							controls
							autoplay
							muted
							crossorigin="use-credentials"
						>
							Sorry, your browser doesn't support embedded videos.
						</video>
						<div class="mt-2 flex justify-between items-center">
							{#if playMessage}
								<span class={`play-message ${playMessage.type}`} transition:fly>
									{playMessage.content}
								</span>
							{:else}
								<span />
							{/if}
							<div>
								<button
									class="p-1 border border-gray-400 text-sm hover:bg-gray-800 rounded-md transition-colors"
									on:click={closePlayer}
								>
									Close
								</button>
							</div>
						</div>
					</div>
				{/if}
			</div>
			<div>
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
								selected={play ? play == torrent.id : undefined}
							/>
						{/each}
					</div>
				{:else}
					<div>{$_('media.no_torrents')}</div>
				{/if}
			</div>
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

	.play-message {
		@apply text-sm;
	}
	.play-message.info {
		@apply text-blue-400;
	}
	.play-message.error {
		@apply text-red-600;
	}
</style>
