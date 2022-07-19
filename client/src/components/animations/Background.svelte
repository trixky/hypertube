<!-- ========================= SCRIPT -->
<script lang="ts">
	import { linear } from 'svelte/easing';
	import { fade, fly } from 'svelte/transition';

	// * Props

	export let palette: string[];
	$: paletteLength = palette.length;
	export let enabled: boolean = true;

	// * Logic

	let backgroundHeight: number = 0;

	function randomNumber(minInc: number, maxExcl: number) {
		return Math.random() * (maxExcl - minInc) + minInc;
	}

	const nbLines = 5;
	let lines: {
		id: number;
		visible: boolean;
		left: number;
		height: number;
		color: string;
		duration: number;
	}[] = [];

	export function start() {
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
			if (enabled) {
				line.visible = true;
				lines = lines;
			}
		}, randomNumber(100, 500));
	}

	export function stop() {
		enabled = false;
	}

	export function restart() {
		enabled = true;
		for (let index = 0; index < lines.length; index++) {
			setTimeout(() => {
				resetLine(index);
			}, randomNumber(0, 500));
		}
	}
</script>

<!-- ========================= HTML -->
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
