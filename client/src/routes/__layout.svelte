<!-- ========================= SCRIPT -->
<script lang="ts">
	import '../app.css';
	import Header from '$components/layouts/header.svelte';
</script>

<!-- ========================= HTML -->
<div class="flex flex-col min-h-screen">
	<Header />
	<main class="media-page flex-grow flex flex-row content-start justify-start relative">
		<slot />
	</main>
</div>
