<!-- ========================= SCRIPT -->
<script lang="ts">
	export let content: string;
</script>

<!-- ========================= HTML -->
<div class="flex w-full">
	<hr />
	<p class="text-white px-2">{content}</p>
	<hr />
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	hr {
		@apply border-white/60 w-full translate-y-3 mx-1;
	}
</style>
