<!-- https://www.flaticon.com/search?author_id=361&style_id=&type=standard&word=eye -->

<!-- ========================= SCRIPT -->
<script lang="ts">
	export let open = false;
	$: img_src = open ? '/open-eye.png' : '/closed-eye.png';
	$: img_alt = open ? 'open eye' : 'closed eye';

	function handle_click() {
		open = !open;
	}
</script>

<!-- ========================= HTML -->
<button
	class="absolute right-2 p-1 translate-y-1"
	on:click|preventDefault={handle_click}
    tabindex="-1"
>
	<img class="" src={img_src} width="24px" height="24px" alt={img_alt} />
</button>
