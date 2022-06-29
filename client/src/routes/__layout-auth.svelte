<!-- ========================= SCRIPT -->
<script lang="ts">
    import '../app.css';
</script>

<!-- ========================= HTML -->
<main class="h-full linea-gradient bg-no-repeat bg-cover bg-center">
    <slot />
</main>

<!-- ========================= CSS -->
<style lang="postcss">
    .linea-gradient {
        background: radial-gradient(rgba(68, 0, 255, 0.2) 5%, black), url('/background.jpg');
    }
</style>