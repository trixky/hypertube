<!-- ========================= SCRIPT -->
<script lang="ts">
	import { _ } from 'svelte-i18n';

	export let destination = '?';
	export let logo = '';
	export let disabled = false;
	export let href: string;
</script>

<!-- ========================= HTML -->
<a {href}>
	<button
		class="inline-block bg-white py-2 px-3 rounded-full duration-[0.35s] hover:bg-gray-300 disabled:opacity-30"
		{disabled}
	>
		{$_('auth.continue_with')}
		<img class="inline-block ml-[1px]" alt="{destination} logo" src={logo} height="20" width="20" />
	</button>
</a>
