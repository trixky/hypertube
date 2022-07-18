<!-- ========================= SCRIPT -->
<script lang="ts">
	import { _ } from 'svelte-i18n';
	import { fade } from 'svelte/transition';
	import QualityIcon from './QualityIcon.svelte';
	import Play from '$components/icons/Play.svelte';
	import type { MediaTorrent } from '../../../../src/types/Media';
	import { createEventDispatcher } from 'svelte';

	export let torrent: MediaTorrent;
	export let selected: boolean | undefined = undefined;

	function seedColor(amount: number) {
		if (amount == 0) {
			return 'text-red-600';
		}
		if (amount <= 10) {
			return 'text-orange-600';
		}
		return 'text-green-600';
	}

	let quality: string = '';
	if (/sd|720p?|(hq)?\s*cam(\s*rip)?/i.exec(torrent.name)) {
		quality = 'sd';
	} else if (/hd|1080p?/i.exec(torrent.name)) {
		quality = 'hd';
	} else if (/2160p?|4k/i.exec(torrent.name)) {
		quality = '4k';
	} else if (/8k/i.exec(torrent.name)) {
		quality = '8k';
	}

	const dispatch = createEventDispatcher();
</script>

<!-- ========================= HTML -->
<div
	class="flex flex-col xl:flex-row xl:items-center w-full mb-4 last:mb-0 xl:mb-2 my-2 bg-black bg-opacity-80 transition-opacity"
	class:opacity-20={selected === false}
	class:opacity-100={selected}
	class:hover:opacity-80={selected !== undefined}
>
	<div class="hidden xl:inline-block">
		<QualityIcon {quality} class="mr-2" />
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
			<QualityIcon {quality} class="mr-1" />
		</div>
		{#if torrent.size}
			{$_('media.size')}: {torrent.size} &#x2022;
		{/if}
		Seed: <span class={`${seedColor(torrent.seed)}`}>{torrent.seed}</span> &#x2022; Leech:
		<span class=" text-red-600">{torrent.leech}</span>
	</div>
	<div class="hidden xl:block mx-4 flex-shrink-0 min-w-[3rem] text-center">
		<span class={`${seedColor(torrent.seed)}`}>{torrent.seed}</span> /
		<span class="text-red-600">{torrent.leech}</span>
	</div>
	<button
		class="flex-shrink-0 p-[2px] mt-2 xl:mt-0 rounded-md font-bold border border-stone-400 hover:border-transparent transition-all relative overflow-hidden"
		class:border-transparent={selected}
		on:mouseenter={() => (torrent.hover = true)}
		on:mouseleave={() => (torrent.hover = false)}
		on:click={() => dispatch('select')}
	>
		{#if torrent.hover || selected}
			<div class="loader" transition:fade />
		{/if}
		<div
			class="flex items-center w-full h-full px-2 py-1 rounded-md relative overflow-hidden bg-black hover:bg-stone-900 transition-all text-blue-400"
		>
			<Play />
			<div class="inline-block flex-grow text-white">
				{#if selected}
					{$_('media.watching')}
				{:else}
					{$_('media.watch')}
				{/if}
			</div>
		</div>
	</button>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.loader {
		@apply absolute top-0 right-0 bottom-0 left-0;
		background: rgb(170, 50, 201);
		background: linear-gradient(90deg, rgb(170, 50, 201) 0%, rgba(107, 139, 176, 1) 100%);
		background-size: 300% 300%;
		background-position: 0 50%;
		animation: move-background 1s alternate infinite;
	}
</style>
