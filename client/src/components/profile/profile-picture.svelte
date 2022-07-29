<!-- ========================= SCRIPT -->
<script lang="ts">
	import { _ } from 'svelte-i18n';

    export let user_id: number;
</script>

<!-- ========================= HTML -->
<a href="/users/{user_id}">
	<img
		class="inline bg-blue-300 rounded-[10px]"
		src="/user.png"
		width="20px"
		height="20px"
		alt={$_('auth.my_profile')}
	/>
</a>

<!-- ========================= CSS -->
<style lang="postcss">
</style>
