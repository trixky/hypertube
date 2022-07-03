<!-- ========================= SCRIPT -->
<script lang="ts">
	import { tweened } from 'svelte/motion';
	import { cubicOut } from 'svelte/easing';
	import { uppercase_first_character } from '../../utils/str';

	export let content = 'bad formatted';
	export let centered = false;

	let height_warning = 0;
	let progress_warning = tweened(height_warning, {
		duration: 200,
		easing: cubicOut
	});

	$: progress_warning.set(content.length ? height_warning : 0);
</script>

<!-- ========================= HTML -->
<div class:text-center={centered} class="error_container" style="height: {$progress_warning}px;">
	<p class="error whitespace-pre-line" bind:offsetHeight={height_warning}>{uppercase_first_character(content)}</p>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.error_container {
		@apply relative overflow-hidden;
	}

	p.error {
		@apply text-red-300 text-xs px-2 pt-2;
	}
</style>
