<!-- ========================= SCRIPT -->
<script lang="ts">
	import { _ } from 'svelte-i18n';

    export let user_id: number;
    export let size = 20;
</script>

<!-- ========================= HTML -->
<a href="/users/{user_id}">
	<img
		class="inline rounded-full"
		src="http://localhost:5180/users/{user_id}/picture"
		width="{size}px"
		height="{size}px"
		alt={$_('auth.my_profile')}
	/>
</a>

<!-- ========================= CSS -->
<style lang="postcss">
</style>
