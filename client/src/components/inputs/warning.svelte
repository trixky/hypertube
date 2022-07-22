<!-- ========================= SCRIPT -->
<script lang="ts">
	import { tweened } from 'svelte/motion';
	import { cubicOut } from 'svelte/easing';
	import { uppercase_first_character } from '$utils/str';

	export let content = '';
	export let centered = false;
	export let color = 'gray'; // gray (default) / green / red / blue / orange

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
	<p
		class:text-gray-300={color == 'gray'}
		class:text-red-300={color == 'red'}
		class:text-green-300={color == 'green'}
		class:text-blue-300={color == 'blue'}
		class:text-orange-300={color == 'orange'}
		class="text-xs pb-2 whitespace-pre-line px-2 transition-colors"
		bind:offsetHeight={height_warning}
	>
		{uppercase_first_character(content)}
	</p>
</div>
