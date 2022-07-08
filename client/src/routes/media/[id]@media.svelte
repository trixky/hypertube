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
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import { quintInOut } from 'svelte/easing';
	import { browser } from '$app/env';
	import { goto } from '$app/navigation';
	import ArrowLeft from '../../../src/components/icons/ArrowLeft.svelte';
	import Play from '../../../src/components/icons/Play.svelte';
	import Hd from '../../../src/components/icons/HD.svelte';
	import Icon4K from '../../../src/components/icons/4K.svelte';
	import Icon8K from '../../../src/components/icons/8K.svelte';
	import Sd from '../../../src/components/icons/SD.svelte';

	/// @ts-expect-error media is given as a prop
	export let props: MediaProps;
	const { media, torrents, staffs, actors } = props;

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

	// Filter and merge actors and staffs once
	const cleanStaffs: {
		id: number;
		name: string;
		thumbnail: string;
		roles: string[];
	}[] = [];
	for (const staff of staffs.filter((staff) => staff.thumbnail && staff.role)) {
		let existingStaff = cleanStaffs.find((existingStaff) => existingStaff.id == staff.id);
		if (existingStaff) {
			existingStaff.roles.push(staff.role!);
		} else {
			existingStaff = {
				id: staff.id,
				name: staff.name,
				thumbnail: staff.thumbnail!,
				roles: [staff.role!]
			};
			cleanStaffs.push(existingStaff);
		}
	}

	const cleanActors: {
		id: number;
		name: string;
		thumbnail: string;
		characters: string[];
	}[] = [];
	for (const actor of actors.filter((actor) => actor.thumbnail && actor.character)) {
		let existingActor = cleanActors.find((existingActor) => existingActor.id == actor.id);
		if (existingActor) {
			existingActor.characters.push(actor.character!);
		} else {
			existingActor = {
				id: actor.id,
				name: actor.name,
				thumbnail: actor.thumbnail!,
				characters: [actor.character!]
			};
			cleanActors.push(existingActor);
		}
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
	function averageRGB(image: HTMLImageElement): [number, number, number] | undefined {
		var context = document.createElement('canvas').getContext('2d');
		if (!context) {
			return undefined;
		}
		context.imageSmoothingEnabled = true;
		context.drawImage(image as HTMLImageElement, 0, 0, 1, 1);
		const i = context.getImageData(0, 0, 1, 1).data;
		let colors = [i[0], i[1], i[2]];
		// Clamp each channels to 150 to avoid bright colors
		let difference = colors.reduce((carry, value) => {
			if (value > 70) {
				return Math.max(carry, value - 70);
			}
			return carry;
		}, 0);
		return colors.map((color) => Math.max(0, color - difference)) as [number, number, number];
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

	let loadingGradient = true;
	let background: string | undefined;
	onMount(() => {
		if (browser) {
			// Load background first to have a clean fade-in
			const useBackground = media.background
				? media.background
				: media.thumbnail
				? media.thumbnail
				: undefined;
			if (useBackground) {
				const image = new Image();
				image.setAttribute('crossOrigin', '');
				image.src = useBackground;
				image.addEventListener('load', () => {
					background = useBackground;
				});
			}

			// Load the image used for the gradient
			const useImage = media.thumbnail ? media.thumbnail : undefined;
			if (useImage && useImage != '') {
				const image = new Image();
				image.setAttribute('crossOrigin', '');
				image.src = useImage;
				image.addEventListener('load', () => {
					let color = averageRGB(image);
					if (color) {
						gradientColor = `rgba(${color.join(',')}, 0.7)`;
					}
					loadingGradient = false;
				});
			} else {
				loadingGradient = false;
			}
		}
	});

	// Background animation
	function randomNumber(minInc: number, maxExcl: number) {
		return Math.random() * (maxExcl - minInc) + minInc;
	}

	const nbLines = 5;
	let colors: string[] = [
		'rgb(96, 165, 250)',
		'rgb(170, 50, 201)',
		'rgb(219, 126, 20)',
		'rgb(16, 121, 39)'
	];
	let lines: {
		id: number;
		visible: boolean;
		left: number;
		height: number;
		color: string;
		duration: number;
	}[] = [];
	if (browser) {
		for (let index = 0; index < nbLines; index++) {
			lines.push({ id: index, visible: false, left: 0, height: 0, color: '', duration: 0 });
			setTimeout(() => {
				resetLine(index);
			}, randomNumber(0, 1000));
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
		line.color = colors[Math.round(randomNumber(0, 4))];
		line.duration = Math.round(randomNumber(1000, 1500));
		setTimeout(function () {
			line.visible = true;
			lines = lines;
		}, randomNumber(100, 500));
	}
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
				<div class="text-white mt-4">{media.description}</div>
				{#if cleanActors.length > 0}
					<div class="text-lg mt-4">Actors</div>
					<ol class="flex pb-4 overflow-x-auto overflow-y-hidden">
						{#each cleanActors as actor (actor.id)}
							<li class="mr-6 last:mr-0 max-w-[6rem]">
								<div
									class="h-24 w-24 rounded-full border-4 border-black border-opacity-80 bg-center bg-cover transition-all"
									style={`background-image: url("${actor.thumbnail}"); ${
										gradientColor ? `border-color: ${gradientColor}` : ''
									}`}
									in:fade={{ duration: 150 }}
								/>
								<div class="font-medium">{actor.name}</div>
								<div class="opacity-80 text-sm truncate" title={actor.characters.join(', ')}>
									{actor.characters.join(', ')}
								</div>
							</li>
						{/each}
					</ol>
				{/if}
				{#if cleanStaffs.length > 0}
					<div class="text-lg mt-4">Staffs</div>
					<ol class="flex pb-4 overflow-x-auto overflow-y-hidden">
						{#each cleanStaffs as staff (staff.id)}
							<li class="mr-6 last:mr-0 max-w-[6rem]">
								<div
									class="h-24 w-24 rounded-full border-4 border-black border-opacity-80 bg-center bg-cover transition-all"
									style={`background-image: url("${staff.thumbnail}"); ${
										gradientColor ? `border-color: ${gradientColor}` : ''
									}`}
									in:fade={{ duration: 150 }}
								/>
								<div class="font-medium">{staff.name}</div>
								<div class="opacity-80 text-sm truncate" title={staff.roles.join(', ')}>
									{staff.roles.join(', ')}
								</div>
							</li>
						{/each}
					</ol>
				{/if}
			</div>
		</div>
	</div>
	<div class="relative flex-grow">
		<div class="absolute top-0 right-0 bottom-0 left-0 overflow-hidden text-white">
			{#each lines as line (line.id)}
				{#if line.visible}
					<div
						class="absolute top-0 w-1 rounded-sm"
						style={`left: ${line.left}px; height: ${line.height}px; background-color: ${line.color}`}
						in:fly={{ duration: 150, y: -line.height, opacity: 1 }}
						out:fly={{ y: window.outerHeight, duration: line.duration, easing: quintInOut }}
						on:introend={removeLine.bind(null, line.id)}
						on:outroend={resetLine.bind(null, line.id)}
					/>
				{/if}
			{/each}
		</div>
		<div class="w-11/12 md:w-4/5 lg:w-1/2 mx-auto text-white my-4 flex-grow relative">
			<div>
				<h1 class="text-2xl mb-4">Torrents</h1>
				{#if torrents.length > 0}
					<div class="w-full">
						{#each torrents as torrent (torrent.id)}
							<div
								class="flex flex-col xl:flex-row xl:items-center w-full my-2 bg-black bg-opacity-80"
							>
								{#if torrent.quality}
									<div class="flex-shrink-0 mr-2">
										{#if torrent.quality == 'hd'}
											<Hd />
										{:else if torrent.quality == '4k'}
											<Icon4K />
										{:else if torrent.quality == '8k'}
											<Icon8K />
										{:else}
											<Sd />
										{/if}
									</div>
								{/if}
								<div class="flex-grow truncate" title={torrent.name}>
									{torrent.name}
								</div>
								{#if torrent.size}
									<div class="hidden xl:block flex-shrink-0 opacity-80">
										{torrent.size}
									</div>
								{/if}
								<div class="xl:hidden">
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
									class="flex flex-shrink-0 items-center p-[2px] rounded-md font-bold border border-stone-400 hover:border-transparent transition-all relative overflow-hidden"
									on:mouseenter={() => (torrent.hover = true)}
									on:mouseleave={() => (torrent.hover = false)}
								>
									{#if torrent.hover}
										<div class="loader" transition:fade />
									{/if}
									<div
										class="w-full h-full px-2 py-1 rounded-md relative overflow-hidden bg-black hover:bg-stone-900 transition-all text-blue-400"
									>
										<Play />
										<div class="inline-block text-white">Watch</div>
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
