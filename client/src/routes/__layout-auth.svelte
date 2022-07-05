<!-- ========================= SCRIPT -->
<script lang="ts">
	import '../app.css';
</script>

<!-- ========================= HTML -->
<main
	class="h-full relative flex content-center justify-center linear-gradient bg-no-repeat bg-cover bg-center"
>
	<slot />
</main>

<!-- ========================= CSS -->
<style lang="postcss">
	.linear-gradient {
		background: radial-gradient(rgba(68, 0, 255, 0.2) 5%, black), url('/background.jpg');
	}
</style>
