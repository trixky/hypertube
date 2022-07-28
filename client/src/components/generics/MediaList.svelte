<!-- ========================= SCRIPT -->
<script lang="ts">
	import { browser } from '$app/env';
	import Spinner from '$components/animations/spinner.svelte';
	import Eye from '$components/icons/Eye.svelte';
	import LazyLoad from '$components/lazy/LazyLoad.svelte';
	import type { Result } from '$types/Media';
	import { imageProxy } from '$utils/api';
	import { onDestroy } from 'svelte';
	import { _ } from 'svelte-i18n';
	import { fade } from 'svelte/transition';

	export let list: Result[];
	export let totalResults: number;
	export let fadeDelay = 0;
	export let loadMore: (() => void) | undefined = undefined;

	export let loading = false;
	let loader: HTMLElement | undefined;

	// * Infinite loader
	// Obser the Load More card when it's visible and loadMore if the user can see it
	let observer: IntersectionObserver;
	function onIntersectionEvent(entries: IntersectionObserverEntry[]) {
		if (!loadMore || loading || list.length == 0 || totalResults == list.length) {
			return;
		}
		for (const entry of entries) {
			if (entry.isIntersecting) {
				loadMore();
			}
		}
	}
	if (browser && loadMore) {
		observer = new IntersectionObserver(onIntersectionEvent, { threshold: 0 });
	}
	let observing: HTMLElement | undefined;
	$: {
		if (loadMore) {
			if (loader) {
				observer.observe(loader);
				observing = loader;
			} else if (observing) {
				observer.unobserve(observing);
				observing = undefined;
			}
		}
	}

	onDestroy(() => {
		if (observer) {
			observer.disconnect();
		}
	});
</script>

<!-- ========================= HTML -->
<div class="result-wrapper p-4">
	{#each list as result, index (result.id)}
		{@const cover = result.thumbnail ? imageProxy(result.thumbnail) : '/no_cover.png'}
		<LazyLoad
			tag="a"
			href={`/media/${result.id}`}
			class="relative result overflow-hidden h-[268px] w-40 min-h-[268px] mx-auto"
		>
			<div
				class="cover"
				class:opacity-80={result.watched}
				style={`background-image: url(${cover})`}
				in:fade={{ duration: 150, delay: (index - fadeDelay) * 10 }}
			>
				{#if result.rating}
					{@const rating = Math.round(result.rating * 10) / 10}
					<div class="rating">
						<div class="flex justify-between items-center w-full">
							<div class="stars" style={`--rating: ${rating};`} />
							<div class="text-sm">{rating}/10</div>
						</div>
					</div>
				{/if}
			</div>
			<div
				class="text-white font-bold truncate"
				title={result.userTitle ? result.userTitle : result.title}
			>
				{result.userTitle ? result.userTitle : result.title}
			</div>
			{#if result.year}
				<div class="text-white text-sm opacity-80">{result.year}</div>
			{/if}
			{#if result.watched}
				<div class="absolute bottom-1 right-1 text-white">
					<Eye />
				</div>
			{/if}
		</LazyLoad>
	{/each}
	{#if totalResults != list.length}
		<div
			bind:this={loader}
			class="result overflow-hidden min-h-[14rem] w-40 mx-auto cursor-pointer text-white"
			class:opacity-50={loading}
			on:click={loadMore}
		>
			{#if loading}
				<div class="flex w-full h-full justify-center items-center text-2xl text-white">
					{$_('search.loading_more')}
					<Spinner size={32} />
				</div>
			{:else}
				<div class="flex w-full h-full justify-center items-center text-2xl text-white">
					{$_('search.load_more')}
				</div>
			{/if}
		</div>
	{/if}
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.result-wrapper {
		@apply grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-8 3xl:grid-cols-10 gap-x-2 gap-y-2;
	}

	.cover {
		@apply h-56 w-40 rounded-md overflow-hidden cursor-pointer transition-all bg-center bg-cover relative;
	}
</style>
