<!-- ========================= SCRIPT -->
<script lang="ts">
	import { profilePicture } from '../../stores/profile-picture';

	export let user_id: number;
	export let tranlsate = false;
	export let size = 20;
</script>

<!-- ========================= HTML -->
{#if $profilePicture}
	<div
		class="inline-block rounded-full bg-center bg-cover"
		style={`background-image: url("http://localhost:5180/v1/user/${user_id}/picture?${$profilePicture}"); width: ${size}px; height: ${size}px; ${
			tranlsate ? `transform: translateY(${size / 4}px);` : ''
		}`}
	/>
{/if}
