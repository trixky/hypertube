<!-- ========================= SCRIPT -->
<script lang="ts">
	import Hd from '$components/icons/HD.svelte';
	import Icon4K from '$components/icons/4K.svelte';
	import Icon8K from '$components/icons/8K.svelte';
	import Sd from '$components/icons/SD.svelte';

	export let quality: string | undefined = undefined;
	let classes = '';
	export { classes as class };
</script>

<!-- ========================= HTML -->

{#if quality}
	<div class={`flex-shrink-0 ${classes}`} {...$$restProps}>
		{#if quality == 'hd'}
			<Hd />
		{:else if quality == '4k'}
			<Icon4K />
		{:else if quality == '8k'}
			<Icon8K />
		{:else}
			<Sd />
		{/if}
	</div>
{/if}
