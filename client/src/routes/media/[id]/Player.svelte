<!-- ========================= SCRIPT -->
<script lang="ts">
	import { createEventDispatcher, onMount, tick } from 'svelte';
	import { fade } from 'svelte/transition';
	import type { MediaTorrent } from 'src/types/Media';
	import { session } from '$app/stores';
	import { _ } from 'svelte-i18n';

	const POSITION_DELAY_MS = 15_000; /* 15sec */
	const ERROR_TIMEOUT = 300_000; /* 5min */

	export let torrent: MediaTorrent;

	const dispatch = createEventDispatcher();

	function closePlayer() {
		dispatch('close');
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
				const subtitlesUrl = `http://localhost:3030/torrent/${torrent.id}/subtitles`;
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
								setSubtitleMessage('success', `Loaded ${index} subtitles`);
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
						fetch(`http://localhost:3040/v1/position/${torrent.id}`, {
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
	});
</script>

<!-- ========================= HTML -->
<div class="mb-10" transition:fade>
	<video
		bind:this={player}
		class="w-full"
		src={`http://localhost:3030/torrent/${torrent.id}/stream`}
		controls
		autoplay
		muted
		crossorigin="use-credentials"
	>
		Sorry, your browser doesn't support embedded videos.
	</video>
	<div class="mt-2 flex justify-between items-center">
		<div class="flex flex-col">
			{#if playMessage}
				<span class={`message ${playMessage.type}`} transition:fade>
					{playMessage.content}
				</span>
			{/if}
			{#if subtitleMessage && subtitleMessage.visible}
				<span class={`message ${subtitleMessage.type}`} transition:fade={{ duration: 1000 }}>
					{subtitleMessage.content}
				</span>
			{/if}
		</div>
		<div>
			<button
				class="p-1 border border-gray-400 text-sm hover:bg-gray-800 rounded-md transition-colors mr-2"
				on:click={closePlayer}
			>
				Close
			</button>
			<button
				class="p-2 border border-blue-400 hover:bg-blue-500 rounded-md transition-colors"
				on:click={closePlayer}
			>
				{$_('media.theatre_mode')}
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
