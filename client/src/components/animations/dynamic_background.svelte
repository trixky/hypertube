<!-- ========================= SCRIPT -->
<script lang="ts">
	import background_movies from './background_movies';
	import Movie from './movie.svelte';

	let w: number = 0;
	let h: number = 0;

	$: scale = w < h ? h / 700 : w / 1000;
</script>

<!-- ========================= HTML -->
<div class="linear-background" />
<div
	class="movie-container"
	style="transform: scale({scale}) translateY({-900 + h / 2}px) translateX({-w / 2}px);"
	bind:clientWidth={w}
	bind:clientHeight={h}
>
	<div class="movies-container">
		{#each background_movies as movies_line, index}
			<div
				class="movie-line"
				style="animation-direction: {index % 2 === 0 ? 'alternate' : 'alternate-reverse'};"
			>
				{#each movies_line as movie}
					<Movie {movie} {index} />
				{/each}
			</div>
		{/each}
	</div>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.movie-container {
		@apply absolute w-full h-full bg-black;
	}

	.movies-container {
		@apply absolute w-full h-full bg-black;
		/* overflow: scroll; */
		transform: rotateX(20deg) rotateY(20deg) rotateZ(5deg) scale(1);
	}

	.linear-background {
		@apply absolute w-full h-full z-30;
		background: radial-gradient(rgba(68, 0, 255, 0.2) 5%, black);
	}

	.movie-line {
		white-space: nowrap;
		animation: x-translation 60s linear 0s infinite;
	}

	@keyframes x-translation {
		0% {
			transform: translateX(-400px);
		}
		100% {
			transform: translateX(400px);
		}
	}
</style>
