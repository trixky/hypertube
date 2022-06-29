<!-- ========================= SCRIPT -->
<script lang="ts">
	export let alone = false
</script>

<!-- ========================= HTML -->
<h1 class:alone class="text-blue-300 font-bold logo-title font-['Holtwood One SC'] text-5xl">Hypertube</h1>

<!-- ========================= CSS -->
<style lang="postcss">
	.logo-title {
		text-shadow: 1px 1px 2px red, 0 0 1em blue, 0 0 0.2em blue;
	}

	.alone {
		@apply absolute z-10 left-1/2 -translate-x-1/2 -top-6;
	}
</style>