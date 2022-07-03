<!-- ========================= SCRIPT -->
<script lang="ts">
	import ContinueWith from './continue-with.svelte';

	export let disabled = false;
</script>

<!-- ========================= HTML -->
<div class="py-4 w-full flex justify-between duration-100">
	<ContinueWith logo="google_logo.svg.png" destination="google" disabled={disabled}/>
	<ContinueWith logo="42_logo.svg.png" destination="42" disabled={disabled}/>
</div>