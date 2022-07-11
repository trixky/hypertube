<!-- ========================= SCRIPT -->
<script lang="ts">
	import Spinner from '../../../components/animations/spinner.svelte';
	import { uppercase_first_character } from '../../../utils/str';

	export let label: string;
	export let value = '';
	export let no_value = false;
	export let centered = false;
	export let can_be_empty = false;
</script>

<!-- ========================= HTML -->
<div class:centered class="my-0">
	<div class="inline-block px-2 align-bottom">
		<label class="inline-block" for="username">{uppercase_first_character(label)}</label>
	</div>
	<div class="inline-block px-2 align-bottom text-gray-400">
		{#if value.length || no_value || can_be_empty}
			<p class="truncate">{value}</p>
		{:else}
			<Spinner inversed />
		{/if}
	</div>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
    div.centered {
        @apply flex justify-center align-middle content-center;
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
