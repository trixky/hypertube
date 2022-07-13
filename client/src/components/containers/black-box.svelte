<!-- ========================= SCRIPT -->
<script lang="ts">
	import { uppercase_first_character } from '$utils/str';
	export let title: string;
</script>

<!-- ========================= HTML -->
<div class="my-auto">
	<div
		class="relative bg-black/80 rounded-2xl  p-8 pb-6 w-[380px] flex justify-center items-center my-10 flex-col h-fit"
	>
		<h1 class="my-2 text-2xl text-white">{uppercase_first_character(title)}</h1>

		<slot />
	</div>
</div>
