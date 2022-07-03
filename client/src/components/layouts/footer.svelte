<!-- ========================= SCRIPT -->
<script lang="ts">
	import { browser } from '$app/env';
	import * as cookies from '../../utils/cookies';

	let connected = false;

	if (browser && cookies.get_a_cookie('token')) {
		connected = true;
	}

    function disconnect() {
        cookies.del_a_cookie('token')

		window.location.href = window.location.origin + '/login';
    }
</script>

<!-- ========================= HTML -->
<footer
	class="w-full bg-black flex flex-col justify-center items-center h-20 min-h-[10%] border-t-stone-200 border-t"
>
	<p class="text-white">footer</p>
	{#if connected}
		<button on:click={disconnect}>
			<p class="text-white">disconnect</p>
		</button>
	{/if}
</footer>
