<!-- ========================= SCRIPT -->
<script lang="ts">
	import { onMount } from 'svelte';

	export let tag: string = 'div';

	let container: HTMLElement;
	let intersected = false;

	onMount(() => {
		const observer = new IntersectionObserver(
			(entries: IntersectionObserverEntry[]) => {
				for (const entry of entries) {
					if (entry.isIntersecting && !intersected) {
						intersected = true;
						observer.unobserve(container);
						observer.disconnect();
					}
				}
			},
			{ root: null, rootMargin: '0px 0px 0px 0px', threshold: 0 }
		);
		observer.observe(container);
	});
</script>

<!-- ========================= HTML -->
<svelte:element this={tag} bind:this={container} {...$$restProps}>
	{#if intersected}
		<slot />
	{/if}
</svelte:element>

<!-- ========================= CSS -->
<style lang="postcss">
</style>
