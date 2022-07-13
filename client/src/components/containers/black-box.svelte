<!-- ========================= SCRIPT -->
<script lang="ts">
	export let title: string;
</script>

<!-- ========================= HTML -->
<div class="my-auto">
	<div
		class="relative bg-black/80 rounded-2xl p-8 pb-6 w-[380px] flex justify-center items-center my-10 flex-col h-fit"
	>
		<h1 class="my-2 text-2xl text-white capitalize">{title}</h1>

		<slot />
	</div>
</div>
