<!-- ========================= SCRIPT -->
<script lang="ts">
	import Hd from '../../../src/components/icons/HD.svelte';
	import Icon4K from '../../../src/components/icons/4K.svelte';
	import Icon8K from '../../../src/components/icons/8K.svelte';
	import Sd from '../../../src/components/icons/SD.svelte';

	export let quality: string | undefined = undefined;
	let classes: string = '';
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
