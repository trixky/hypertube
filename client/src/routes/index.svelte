<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	import type { Load } from '@sveltejs/kit';

	export const load: Load = async () => {
		return {
			status: 301,
			redirect: '/search'
		};
	};
</script>
