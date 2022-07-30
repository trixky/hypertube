<!-- ========================= SCRIPT -->
<script lang="ts">
	import Spinner from '$components/animations/spinner.svelte';

	export let label: string;
	export let value = '';
	export let no_value = false;
	export let can_be_empty = false;
</script>

<!-- ========================= HTML -->
<div class="centered my-0 py-2">
	<div class="inline-block px-2 align-bottom">
		<label class="inline-block first-letter:capitalize" style="padding: 0 !important;" for="username">{label}</label>
	</div>
	<div class="inline-block px-2 align-bottom text-gray-400">
		{#if value.length || no_value || can_be_empty}
			<p class="truncate">{value}</p>
		{/if}
	</div>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
	div.centered {
		@apply flex justify-center align-middle content-center items-center;
	}

	div.centered > div {
		@apply w-1/2;
	}

	div.centered > div:first-child {
		@apply text-right;
	}

	div:not(.centered) > div:first-child {
		@apply w-[25%];
	}
	div:not(.centered) > div:nth-child(2) {
		@apply w-[70%];
	}
</style>
