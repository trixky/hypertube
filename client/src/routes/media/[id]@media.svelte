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
			quality?: string;
			hover: boolean;
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

	export const load: Load = async ({ params, fetch }) => {
		const url = browser
			? `http://localhost:7072/v1/media/${params.id}/get`
			: `http://api-media:7072/v1/media/${params.id}/get`;
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
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import { linear } from 'svelte/easing';
	import { browser } from '$app/env';
	import { goto } from '$app/navigation';
	// @ts-expect-error No types for quantize
	import quantize from 'quantize';
	import ArrowLeft from '../../../src/components/icons/ArrowLeft.svelte';
	import Play from '../../../src/components/icons/Play.svelte';
	import LazyLoad from '../../../src/components/lazy/LazyLoad.svelte';
	import Spinner from '../../../src/components/animations/spinner.svelte';
	import QualityIcon from './QualityIcon.svelte';

	/// @ts-expect-error media is given as a prop
	export let props: MediaProps;
	let { media, torrents, staffs, actors } = props;

	// Find quality for torrents
	for (const torrent of torrents) {
		if (/sd|720p?|(hq)?\s*cam(\s*rip)?/i.exec(torrent.name)) {
			torrent.quality = 'sd';
		} else if (/hd|1080p?/i.exec(torrent.name)) {
			torrent.quality = 'hd';
		} else if (/2160p?|4k/i.exec(torrent.name)) {
			torrent.quality = '4k';
		} else if (/8k/i.exec(torrent.name)) {
			torrent.quality = '8k';
		}
	}

	const cover = media.thumbnail ? `http://localhost:7260${media.thumbnail}` : '/no_cover.png';

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

	// Add full thumbnail to staffs and actors
	const cleanStaffs = staffs.map((staff) => ({
		id: staff.id,
		name: staff.name,
		thumbnail: `http://localhost:7260${staff.thumbnail}`,
		roles: [staff.role!]
	}));
	const cleanActors = actors.map((actor) => ({
		id: actor.id,
		name: actor.name,
		thumbnail: `http://localhost:7260${actor.thumbnail}`,
		characters: [actor.character!]
	}));

	// Background animation
	let palette: string[] = [];
	let paletteLength = palette.length;

	function randomNumber(minInc: number, maxExcl: number) {
		return Math.random() * (maxExcl - minInc) + minInc;
	}

	const nbLines = 10;
	let lines: {
		id: number;
		visible: boolean;
		left: number;
		height: number;
		color: string;
		duration: number;
	}[] = [];
	function startBackground() {
		for (let index = 0; index < nbLines; index++) {
			lines.push({ id: index, visible: false, left: 0, height: 0, color: '', duration: 0 });
			setTimeout(() => {
				resetLine(index);
			}, randomNumber(0, 500));
		}
	}
	function removeLine(index: number) {
		const line = lines[index];
		line.visible = false;
		lines = lines;
	}
	function resetLine(index: number) {
		const line = lines[index];
		line.left = Math.round(randomNumber(0, window.outerWidth));
		line.height = Math.round(randomNumber(32, 64));
		line.color = palette[Math.round(randomNumber(0, paletteLength))];
		line.duration = Math.round(randomNumber(1500, 3500));
		setTimeout(function () {
			line.visible = true;
			lines = lines;
		}, randomNumber(100, 500));
	}

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

	function seedColor(amount: number) {
		if (amount == 0) {
			return 'text-red-600';
		}
		if (amount <= 10) {
			return 'text-orange-600';
		}
		return 'text-green-600';
	}

	// Image average color
	// @source https://stackoverflow.com/a/49837149
	let gradientColor: string | undefined;
	let backgroundHeight: number = 0;
	function extractPalette(image: HTMLImageElement) {
		var context = document.createElement('canvas').getContext('2d');
		if (!context) {
			return undefined;
		}
		context.imageSmoothingEnabled = true;
		context.drawImage(image, 0, 0, image.width, image.height);

		// Extract pixels RGB channels as an array of pixel data
		const pixels = context.getImageData(0, 0, image.width, image.height).data;
		const imageData: [number, number, number][] = [];
		for (var i = 0; i < pixels.length; i += 4) {
			let rgb: [number, number, number] = [pixels[i], pixels[i + 1], pixels[i + 2]];
			let a = pixels[i + 3];
			// If pixel is mostly opaque and not white
			if (typeof a === 'undefined' || a >= 125) {
				if (
					!(rgb[0] > 250 && rgb[1] > 250 && rgb[2] > 250) &&
					!(rgb[0] < 30 && rgb[1] < 30 && rgb[2] < 30)
				) {
					imageData.push(rgb);
				}
			}
		}

		// Extract a color palette
		const rawPalette: [number, number, number][] = quantize(imageData, 5, 10).palette();
		palette = rawPalette.map((color) => {
			return `rgb(${color[0]}, ${color[1]}, ${color[2]})`;
		});
		paletteLength = palette.length;
		startBackground();

		// Clamp each channels to 150 to avoid bright colors
		let color = [rawPalette[0][0], rawPalette[0][1], rawPalette[0][2]];
		let difference = color.reduce((carry, value) => {
			if (value > 70) {
				return Math.max(carry, value - 70);
			}
			return carry;
		}, 0);
		color = color.map((color) => Math.max(0, color - difference)) as [number, number, number];
		gradientColor = `rgb(${color[0]}, ${color[1]}, ${color[2]})`;
		loadingGradient = false;
	}

	const comments: {
		id: number;
		user: {
			id: number;
			name: string;
		};
		date: Date;
		content: string;
	}[] = [
		{
			id: 1,
			user: {
				id: 1,
				name: 'ncolomer'
			},
			date: new Date(),
			content:
				'Lorem ipsum dolor sit, amet consectetur adipisicing elit. Provident non debitis enim autem dolor in consequatur odit, nisi nemo nesciunt cumque eligendi obcaecati. Expedita impedit sit animi nam aliquam quasi?'
		},
		{
			id: 2,
			user: {
				id: 2,
				name: 'mcolomer'
			},
			date: new Date(),
			content:
				'Lorem ipsum dolor sit, amet consectetur adipisicing elit. Provident non debitis enim autem dolor in consequatur odit, nisi nemo nesciunt cumque eligendi obcaecati. Expedita impedit sit animi nam aliquam quasi?'
		}
	];

	// Refresh Peers
	let refreshingPeers = false;
	async function refreshPeers() {
		refreshingPeers = true;
		const response = await fetch(`http://localhost:7072/v1/media/${media.id}/refresh`, {
			method: 'GET',
			headers: { accept: 'application/json' }
		});
		if (response.ok) {
			const lines = await await response.text();
			for (const line of lines.trim().split('\n')) {
				try {
					const data = JSON.parse(line) as {
						result: { torrentId: number; seed: number; leech: number };
					};
					if (
						data.result?.torrentId != undefined &&
						typeof data.result?.seed == 'number' &&
						typeof data.result?.leech == 'number'
					) {
						let torrent = torrents.find((torrent) => torrent.id == data.result.torrentId);
						if (torrent) {
							torrent.seed = data.result.seed;
							torrent.leech = data.result.leech;
						}
					}
				} catch (error) {
					console.error('Failed to read line in sream response', error);
				}
			}
			// Try to sort again on a simple metric
			torrents.sort((a, b) => b.seed - a.seed);
			torrents = torrents;
		}
		refreshingPeers = false;
	}

	let loadingGradient = true;
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
				useBackground = `http://localhost:7260${useBackground}`;
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
				useImage = `http://localhost:7260${useImage}`;
				const image = new Image();
				image.setAttribute('crossOrigin', 'anonymous');
				image.src = useImage;
				image.addEventListener('load', () => {
					extractPalette(image);
				});
			} else {
				loadingGradient = false;
			}
		}
	});
</script>

<!-- ========================= HTML -->
<div class="flex flex-col w-full h-auto bg-black">
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
			<a href="/search" on:click={goBack} class="cursor-pointer"><ArrowLeft /> Go Back</a>
		</div>
		<div
			class="relative flex flex-col md:flex-row justify-center items-center w-11/12 md:w-4/5 lg:w-1/2 mx-auto py-10"
		>
			<img
				src={cover}
				alt={`${userFavoriteTitle} Cover`}
				in:fade={{ duration: 150, delay: 50 }}
				class="h-[500px] rounded-md flex-grow-0 "
			/>
			<div
				class="md:ml-8 max-w-full md:max-w-[348px] lg:max-w-[612px] xl:max-w-[720px] text-white transition-all"
			>
				<div class="text-3xl mt-4 lg:mt-0">{userFavoriteTitle.title}</div>
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
					<div class="text-lg mt-4 mb-2">Actors</div>
					<ol class="flex w-full pb-4 overflow-x-auto overflow-y-hidden">
						{#each cleanActors as actor (actor.id)}
							<LazyLoad tag="li" class="mr-6 last:mr-0 w-24 flex-shrink-0 h-full">
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
					<div class="text-lg mt-4 mb-2">Staffs</div>
					<ol class="flex w-full pb-4 overflow-x-auto overflow-y-hidden">
						{#each cleanStaffs as staff (staff.id)}
							<LazyLoad tag="li" class="mr-6 last:mr-0 w-24 flex-shrink-0 h-full">
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
		<div
			bind:clientHeight={backgroundHeight}
			class="absolute top-0 right-0 bottom-0 left-0 overflow-hidden text-white"
		>
			{#each lines as line (line.id)}
				{#if line.visible}
					<div
						class="absolute top-0 w-1 rounded-sm"
						style={`left: ${line.left}px; height: ${line.height}px; background-color: ${line.color}`}
						in:fade={{ duration: 0 }}
						out:fly={{ y: backgroundHeight, duration: line.duration, delay: 0, easing: linear }}
						on:introend={removeLine.bind(null, line.id)}
						on:outroend={resetLine.bind(null, line.id)}
					/>
				{/if}
			{/each}
		</div>
		<div class="w-11/12 md:w-4/5 lg:w-1/2 mx-auto text-white my-4 flex-grow relative">
			<div>
				<h1 class="flex justify-between items-center text-2xl mb-4">
					<span>Torrents</span>
					<div class="inline-block relative opacity-80">
						{#if refreshingPeers}
							<div
								in:fade={{ duration: 350 }}
								out:fade={{ duration: 80 }}
								id="spinner"
								class="absolute inline-block -translate-x-full transition-all duration-[0.35s] opacity-50"
								style={`--tw-translate-x: calc(-100% - 4px);`}
							>
								<Spinner size={16} />
							</div>
						{/if}
						<button
							class="inline-block text-sm transition-all disabled:opacity-50"
							class:hover:text-opacity-90={!refreshingPeers}
							disabled={refreshingPeers}
							on:click={refreshPeers}
						>
							Refresh Peers
						</button>
					</div>
				</h1>
				{#if torrents.length > 0}
					<div class="w-full">
						{#each torrents as torrent (torrent.id)}
							<div
								class="flex flex-col xl:flex-row xl:items-center w-full mb-4 last:mb-0 xl:mb-2 my-2 bg-black bg-opacity-80"
							>
								<div class="hidden xl:inline-block">
									<QualityIcon quality={torrent.quality} class="mr-2" />
								</div>
								<div class="flex-grow truncate" title={torrent.name}>
									{torrent.name}
								</div>
								{#if torrent.size}
									<div class="hidden xl:block flex-shrink-0 opacity-80">
										{torrent.size}
									</div>
								{/if}
								<div class="xl:hidden">
									<div class="inline-block xl:hidden">
										<QualityIcon quality={torrent.quality} class="mr-1" />
									</div>
									{#if torrent.size}
										Size: {torrent.size} &#x2022;
									{/if}
									Seed: <span class={`${seedColor(torrent.seed)}`}>{torrent.seed}</span> &#x2022;
									Leech:
									<span class=" text-red-600">{torrent.leech}</span>
								</div>
								<div class="hidden xl:block mx-4 flex-shrink-0 min-w-[3rem] text-center">
									<span class={`${seedColor(torrent.seed)}`}>{torrent.seed}</span> /
									<span class="text-red-600">{torrent.leech}</span>
								</div>
								<button
									class="flex-shrink-0 p-[2px] mt-2 xl:mt-0 rounded-md font-bold border border-stone-400 hover:border-transparent transition-all relative overflow-hidden"
									on:mouseenter={() => (torrent.hover = true)}
									on:mouseleave={() => (torrent.hover = false)}
								>
									{#if torrent.hover}
										<div class="loader" transition:fade />
									{/if}
									<div
										class="flex items-center w-full h-full px-2 py-1 rounded-md relative overflow-hidden bg-black hover:bg-stone-900 transition-all text-blue-400"
									>
										<Play />
										<div class="inline-block flex-grow text-white">Watch</div>
									</div>
								</button>
							</div>
						{/each}
					</div>
				{:else}
					<div>No torrents for this media, yet !</div>
				{/if}
			</div>
			<div class="my-4">
				<h1 class="text-2xl mb-4">
					Comments {#if comments.length > 0}
						({comments.length})
					{/if}
				</h1>
				{#if comments.length > 0}
					{#each comments as comment (comment.id)}
						<div class="comment" class:self={comment.user.id == 1}>
							{#if comment.user.id == 1}
								<div class="bordered" />
							{/if}
							<div class="comment-header">
								<div>
									<span class="opacity-60 text-sm">#{comment.id}</span>
									<span class="font-bold">{comment.user.name}</span>
								</div>
								<div class="text-sm">{comment.date}</div>
							</div>
							<div class="comment-content">{comment.content}</div>
						</div>
					{/each}
				{:else}
					<div>No comments on this media, yet, be the first one !</div>
				{/if}
			</div>
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

	.loader {
		@apply absolute top-0 right-0 bottom-0 left-0;
		background: rgb(170, 50, 201);
		background: linear-gradient(90deg, rgb(170, 50, 201) 0%, rgba(107, 139, 176, 1) 100%);
		background-size: 300% 300%;
		background-position: 0 50%;
		animation: move-background 1s alternate infinite;
	}

	.comment {
		@apply mb-4 p-2 border border-stone-400 rounded-md bg-stone-900 relative;
	}

	.comment.self {
		@apply border-transparent overflow-hidden;
		padding: 1px;
	}

	.comment.self .comment-header {
		@apply p-2 pb-0 bg-stone-900 rounded-t-md;
	}
	.comment.self .comment-content {
		@apply p-2 pt-0 bg-stone-900 rounded-b-md;
	}

	.comment.self .bordered {
		@apply absolute top-0 right-0 bottom-0 left-0;
		background: rgb(170, 50, 201);
		background: linear-gradient(to bottom right, rgb(170, 50, 201) 0%, rgba(107, 139, 176, 1) 100%);
	}

	.comment-header {
		@apply flex justify-between w-full relative;
	}

	.comment-content {
		@apply relative;
	}

	.comment-content::before {
		@apply block w-full mb-1;
		content: '';
		height: 1px;
		background: linear-gradient(
			to right,
			rgba(0, 0, 0, 0) 25%,
			rgba(255, 255, 255, 0.8) 50%,
			rgba(0, 0, 0, 0) 75%
		);
	}
</style>
