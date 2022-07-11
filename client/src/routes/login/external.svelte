<!-- ========================= SCRIPT -->
<script lang="ts">
	import ContinueWith from './continue-with.svelte';

	const href_42 = import.meta.env.VITE_42_OAUTH_AUTHORIZE;
	const href_google = import.meta.env.VITE_GOOGLE_OAUTH_AUTHORIZE;

	export let disabled = false;
</script>

<!-- ========================= HTML -->
<div class="py-4 w-full flex justify-between duration-100">
	<ContinueWith href={href_google} logo="google_logo.svg.png" destination="google" {disabled} />
	<ContinueWith href={href_42} logo="42_logo.svg.png" destination="42" {disabled} />
</div>
