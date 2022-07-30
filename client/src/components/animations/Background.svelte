<!-- ========================= SCRIPT -->
<script lang="ts">
	import { onDestroy } from 'svelte';

	import { linear } from 'svelte/easing';
	import { fade, fly } from 'svelte/transition';

	// * Props

	export let palette: string[];
	$: paletteLength = palette.length;
	export let enabled = false;

	// * Logic

	let backgroundHeight = 0;

	function randomNumber(minInc: number, maxExcl: number) {
		return Math.random() * (maxExcl - minInc) + minInc;
	}

	export let nbLines = 5;
	let lines: {
		id: number;
		visible: boolean;
		left: number;
		height: number;
		color: string;
		duration: number;
		timeout: number;
	}[] = [];
	let restartTimeouts: number[] = new Array(nbLines).fill(0);

	export function start() {
		enabled = true;
		for (let index = 0; index < nbLines; index++) {
			let line = {
				id: index,
				visible: false,
				left: 0,
				height: 0,
				color: '',
				duration: 0,
				timeout: 0
			};
			lines.push(line);
			line.timeout = setTimeout(() => {
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
		line.timeout = setTimeout(function () {
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
			restartTimeouts[index] = setTimeout(() => {
				resetLine(index);
			}, randomNumber(0, 500));
		}
	}

	onDestroy(() => {
		for (const line of lines) {
			clearTimeout(line.timeout);
		}
		for (const restartTimeout of restartTimeouts) {
			clearTimeout(restartTimeout);
		}
	});
</script>

<!-- ========================= HTML -->
<div
	bind:clientHeight={backgroundHeight}
	class="absolute top-0 right-0 bottom-0 left-0 overflow-hidden text-white"
>
	{#each lines as line (line.id)}
		{#if line.visible}
			<div
				class="absolute top-0 w-1 rounded-sm will-change-transform"
				style={`left: ${line.left}px; height: ${line.height}px; background-color: ${line.color}`}
				in:fade={{ duration: 0 }}
				out:fly={{ y: backgroundHeight, duration: line.duration, delay: 0, easing: linear }}
				on:introend={removeLine.bind(null, line.id)}
				on:outroend={resetLine.bind(null, line.id)}
			/>
		{/if}
	{/each}
</div>
