<!-- ========================= SCRIPT -->
<script lang="ts">
	export let title: string;
</script>

<!-- ========================= HTML -->
<div class="my-auto w-full z-[200]">
	<div
		class="mx-auto relative bg-black/80 rounded-2xl p-8 pb-6 w-[92%] max-w-[400px] flex justify-center items-center my-6 flex-col h-fit"
	>
		<h1 class="my-2 text-2xl text-white first-letter:capitalize">{title}</h1>
		<slot />
	</div>
</div>
