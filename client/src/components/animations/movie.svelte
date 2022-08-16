<!-- ========================= SCRIPT -->
<script lang="ts">
    export let index: number;
    export let movie: string;

    let flip = false;
</script>

<!-- ========================= HTML -->
<div class="movie inline-block w-[120px] h-[180px]" style="transform: translateX({(index * 40) % 120}px); background-image: url({movie});"  alt="" />

<!-- ========================= CSS -->
<style lang="postcss">
    .movie {
        margin: 0 3px;
        background-repeat: no-repeat;
        background-size: cover;
        background-position: center;
    }
</style>