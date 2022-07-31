<!-- ========================= SCRIPT -->
<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount, tick } from 'svelte';
	import { fade, slide } from 'svelte/transition';
	import type { MediaTorrent } from '$types/Media';
	import { session } from '$app/stores';
	import { _ } from 'svelte-i18n';
	import { readableFileSize } from '$utils/media';
	import { cubicOut } from 'svelte/easing';
	import { apiPosition, apiStreaming } from '$utils/api';

	type StatusResponse =
		| {
				status: 'complete' | 'idle';
		  }
		| {
				status: 'ongoing';
				download?: {
					completed: number;
					total: number;
				};
				encoding?: {
					completeDuration: string;
					fps?: number;
					processed?: string;
				};
		  };

	const POSITION_DELAY_MS = 15_000; /* 15sec */
	const ERROR_TIMEOUT = 300_000; /* 5min */

	export let torrent: MediaTorrent;
	let theatre = false;

	const dispatch = createEventDispatcher();

	let destroyed = false;
	let loading = true;
	function closePlayer() {
		dispatch('close');
		loading = false;
	}

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

	type MessageType = 'success' | 'info' | 'error';
	type Message = { type: MessageType; content: string; visible?: boolean };

	let player: HTMLVideoElement;
	let playMessage: Message | undefined;
	let subtitleMessage: Message | undefined;
	let updateErrored: number | null = null;
	let lastUpdate = 0;

	function setPlayMessage(type: MessageType, content: string) {
		playMessage = { type, content };
	}

	let subtitleMessageTimeout = 0;
	function setSubtitleMessage(type: MessageType, content: string) {
		subtitleMessage = { type, content, visible: true };
		clearTimeout(subtitleMessageTimeout);
		subtitleMessageTimeout = setTimeout(() => {
			if (subtitleMessage) {
				subtitleMessage.visible = false;
			}
		}, 1000);
	}

	function toggleTheatreMode() {
		theatre = !theatre;
	}

	let isPlayerOpen = false;
	let focusTimeout = 0;
	function playerOpened() {
		isPlayerOpen = true;
		dispatch('open');
	}

	function focusPlayer(
		transition?:
			| TransitionEvent
			| (CustomEvent<null> & { currentTarget: EventTarget & HTMLDivElement })
	) {
		if (
			transition &&
			'propertyName' in transition &&
			transition.propertyName == 'background-color'
		) {
			return;
		}
		clearTimeout(focusTimeout);
		dispatch('focus');
		if (!isPlayerOpen) {
			focusTimeout = setTimeout(() => {
				focusPlayer();
			}, 20);
		}
	}

	let statusTimeout = 0;
	function refreshPlayer() {
		clearTimeout(statusTimeout);
		loading = true;
		let currentSrc = player.currentSrc;
		let currentTime = player.currentTime;
		player.src = '';
		player.src = currentSrc;
		seekToTime(player, currentTime);
	}

	// Check torrent status for some informations
	let idledFor = 0;
	async function checkStatus() {
		if (destroyed) return;
		const response = await fetch(apiStreaming(`/torrent/${torrent.id}/status`), {
			method: 'GET',
			credentials: 'include',
			headers: {
				Accept: 'application/json'
			}
		});
		if (response.ok) {
			const body = (await response.json()) as StatusResponse;
			if (body.status == 'ongoing') {
				let message = 'Video is still being processed.';
				if (body.download && body.download.completed != body.download.total) {
					message = `${message}\nDownload: ${readableFileSize(
						body.download.completed
					)} / ${readableFileSize(body.download.total)}`;
				}
				if (body.encoding) {
					message = `${message}\nEncoding: ${body.encoding.processed} / ${body.encoding.completeDuration} (${body.encoding.fps}fps)`;
				}
				setPlayMessage('info', message);
				statusTimeout = setTimeout(checkStatus, 5000);
			} else if (body.status == 'idle') {
				idledFor += 1;
				if (idledFor < 3) {
					statusTimeout = setTimeout(checkStatus, 5000);
					playMessage = undefined;
				} else {
					setPlayMessage(
						'error',
						`Failed to get torrent status: ${response.status}, close and re-open`
					);
				}
			} else {
				playMessage = undefined;
			}
		} else {
			setPlayMessage('error', `Failed to get torrent status: ${response.status}`);
			statusTimeout = setTimeout(checkStatus, 5000);
		}
	}

	let initialStatusTimeout = 0;
	onMount(async () => {
		setPlayMessage('info', 'Video is loading...');

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
				setSubtitleMessage('info', 'Loading subtitles...');
				const subtitlesUrl = apiStreaming(`/torrent/${torrent.id}/subtitles`);
				fetch(subtitlesUrl, { method: 'GET', credentials: 'include' })
					.then(async (response) => {
						if (response.ok) {
							const body = (await response.json()) as
								| { subtitles: { id: number; lang: string }[] }
								| { error: string };
							if ('error' in body) {
								setSubtitleMessage('error', 'Failed to load subtitles...');
							} else {
								let index = 0;
								for (const subtitle of body.subtitles) {
									const track = document.createElement('track');
									track.kind = 'captions';
									track.label = subtitle.lang == 'fr' ? 'Francais' : 'English';
									track.srclang = subtitle.lang;
									track.src = apiStreaming(`/subtitles/${subtitle.id}`);
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
								setSubtitleMessage('success', `Loaded ${index} subtitle${index > 1 ? 's' : ''}`);
							}
						} else {
							throw new Error('Failed to load subtitles');
						}
					})
					.catch((error) => {
						console.error(error);
						setSubtitleMessage('error', 'Failed to load subtitles...');
					});
			} catch (error) {
				console.error(error);
				setSubtitleMessage('error', 'Failed to load subtitles...');
			}

			// Bind
			let bind = false;
			player.addEventListener('loadedmetadata', async () => {
				loading = false;
				playMessage = undefined;
				if (torrent.position) {
					seekToTime(player!, torrent.position);
				}
				clearTimeout(initialStatusTimeout);
				await checkStatus();

				if (!bind) {
					/*let retryCount = 0;
					player.addEventListener('stalled', () => {
						if (retryCount <= 3) {
							refreshPlayer();
							retryCount += 1;
						} else {
							setPlayMessage(
								'error',
								'The player is having some problems, you need to refresh it.'
							);
						}
					});*/
					bind = true;
				}
			});
			initialStatusTimeout = setTimeout(() => {
				checkStatus();
			}, 1000);

			let updatingPosition = false;
			player.addEventListener('timeupdate', (event) => {
				if (!updatingPosition && event.timeStamp - lastUpdate >= POSITION_DELAY_MS) {
					if (!updateErrored || Date.now() - updateErrored >= ERROR_TIMEOUT /* 5min */) {
						lastUpdate = event.timeStamp;
						const currentTime = player?.currentTime;
						if (currentTime && !isNaN(currentTime) && currentTime > 0) {
							updatingPosition = true;
							dispatch('timeUpdate', currentTime);
							fetch(apiPosition(`/v1/position/${torrent.id}`), {
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
						}
					} else {
						lastUpdate = event.timeStamp + Date.now() - updateErrored;
					}
				}
			});

			// Handle torrent errors
			player.addEventListener('error', () => {
				clearTimeout(initialStatusTimeout);
				clearTimeout(statusTimeout);
				clearTimeout(focusTimeout);
				clearTimeout(subtitleMessageTimeout);
				dispatch('error');
			});
		}
	});

	onDestroy(() => {
		destroyed = true;
		clearTimeout(initialStatusTimeout);
		clearTimeout(statusTimeout);
		clearTimeout(focusTimeout);
		clearTimeout(subtitleMessageTimeout);
	});
</script>

<!-- ========================= HTML -->
{#if theatre}
	<div
		class="absolute top-0 right-0 bottom-0 left-0 bg-black bg-opacity-80 select-none pointer-events-none z-20"
		transition:fade|local
	/>
{/if}
<div
	class={`${
		theatre ? 'w-full' : 'w-11/12 md:w-4/5 lg:w-1/2'
	} mx-auto text-white my-4 flex-grow relative mb-10 transition-all z-30`}
	transition:slide|local={{ duration: 1000, easing: cubicOut }}
	on:introstart={focusPlayer}
	on:introend={playerOpened}
	on:transitionstart={focusPlayer}
	on:transitionend={focusPlayer}
>
	<video
		bind:this={player}
		class="w-full max-h-screen"
		src={apiStreaming(`/torrent/${torrent.id}/stream`)}
		controls
		autoplay
		muted
		crossorigin="use-credentials"
	>
		Sorry, your browser doesn't support embedded videos.
	</video>
	<div class="mt-2 flex flex-col md:flex-row justify-between items-center">
		<div class="flex flex-col w-full md:w-auto text-left pl-2">
			{#if playMessage}
				<span class={`message ${playMessage.type} whitespace-pre-line`} transition:fade|local>
					{playMessage.content}
				</span>
			{/if}
			{#if subtitleMessage && subtitleMessage.visible}
				<span class={`message ${subtitleMessage.type}`} transition:fade|local={{ duration: 1000 }}>
					{subtitleMessage.content}
				</span>
			{/if}
		</div>
		<div class="w-full md:w-auto text-right pr-2">
			<button
				class="p-1 border border-gray-400 text-sm hover:bg-gray-800 rounded-md transition-colors mr-2"
				on:click={closePlayer}
			>
				Close
			</button>
			<button
				class="p-1 border border-orange-400 text-sm hover:bg-orange-800 rounded-md transition-colors mr-2 disabled:opacity-50"
				on:click={refreshPlayer}
				disabled={loading}
			>
				Refresh
			</button>
			<button
				class="p-2 border border-blue-400 bg-blue-500 hover:bg-blue-600 rounded-md transition-colors shadow-blue-400 hover:shadow-md"
				on:click={toggleTheatreMode}
			>
				{#if !theatre}
					{$_('media.theatre_mode')}
				{:else}
					{$_('media.disable_theatre_mode')}
				{/if}
			</button>
		</div>
	</div>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.message {
		@apply text-sm;
	}
	.message.success {
		@apply text-green-400;
	}
	.message.info {
		@apply text-blue-400;
	}
	.message.error {
		@apply text-red-600;
	}
</style>
