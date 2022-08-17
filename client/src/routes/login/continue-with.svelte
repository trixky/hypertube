<!-- ========================= SCRIPT -->
<script lang="ts">
	import { _ } from 'svelte-i18n';

	export let destination = '?';
	export let logo = '';
	export let disabled = false;
	export let href: string;

	const demo_mode = import.meta.env.VITE_DEMO_MODE;

	$: _disabled = disabled || demo_mode != 'false'
</script>

<!-- ========================= HTML -->
<a class="w-full text-center overflow-hidden" href={_disabled ? '#' : href}>
	<button
		class:cursor-not-allowed={_disabled}
		class="inline-block bg-white w-full py-2 px-3 rounded-full duration-[0.35s] hover:bg-gray-300 disabled:opacity-30 text-center"
		disabled={_disabled}
	>
		<div class="flex justify-center items-center">
			<p class="truncate pr-2 hidden md:block">{$_('auth.continue_with')}</p>
			<img
				class="w-[20px] h-[20px] "
				alt="{destination} logo"
				src={logo}
				height="20"
				width="20"
			/>
		</div>
	</button>
</a>