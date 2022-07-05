<!-- ========================= SCRIPT -->
<script lang="ts">
	import '../app.css';
	import { browser } from '$app/env';
	import { already_connected } from '../utils/redirect';

	already_connected(browser);
</script>

<!-- ========================= HTML -->
<main class="h-full linea-gradient bg-no-repeat bg-cover bg-center">
	<slot />
</main>

<!-- ========================= CSS -->
<style lang="postcss">
	.linea-gradient {
		background: radial-gradient(rgba(68, 0, 255, 0.2) 5%, black), url('/background.jpg') no-repeat;
		background-size: cover;
		background-position: center;
	}
</style>
