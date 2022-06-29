<!-- ========================= SCRIPT -->
<script lang="ts">
    import '../app.css';
    import Header from '../components/layouts/header.svelte';
    import Footer from '../components/layouts/footer.svelte';
</script>

<!-- ========================= HTML -->

<Header />
<main class="min-h-[80%]">
    <slot />
</main>
<Footer />
