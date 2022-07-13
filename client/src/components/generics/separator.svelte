<!-- ========================= SCRIPT -->
<script lang="ts">
	export let content: string;
</script>

<!-- ========================= HTML -->
<div class="flex flex-row justify-center items-center w-full h-4">
	<div class="bar" />
	<p class="text-white px-2 -translate-y-[2px]">{content}</p>
	<div class="bar" />
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	.bar {
		@apply inline-block border-t border-white/60 h-[1px] w-full;
	}
</style>
