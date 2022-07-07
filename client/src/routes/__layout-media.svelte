<!-- ========================= SCRIPT -->
<script lang="ts">
	import '../app.css';
	import Header from '../components/layouts/header.svelte';
</script>

<!-- ========================= HTML -->

<Header />
<main class="search min-h-[90%] flex flex-row content-start justify-start relative">
	<slot />
</main>
