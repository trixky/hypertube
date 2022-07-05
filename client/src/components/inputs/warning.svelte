<!-- ========================= SCRIPT -->
<script lang="ts">
	import { tweened } from 'svelte/motion';
	import { cubicOut } from 'svelte/easing';
	import { uppercase_first_character } from '../../utils/str';

	export let content = '';
	export let centered = false;

	let height_warning = 0;
	let progress_warning = tweened(height_warning, {
		duration: 200,
		easing: cubicOut
	});

	$: progress_warning.set(content.length ? height_warning : 0);
</script>

<!-- ========================= HTML -->
<div
	class:text-center={centered}
	class="relative overflow-hidden"
	style="height: {$progress_warning}px;"
>
	<p class="text-red-300 text-xs  pt-2 whitespace-pre-line px-2" bind:offsetHeight={height_warning}>
		{uppercase_first_character(content)}
	</p>
</div>
