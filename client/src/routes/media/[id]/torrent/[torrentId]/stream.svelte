<!-- ========================= SCRIPT -->
<script lang="ts">
	import { page } from '$app/stores';

	let torrentId = $page.params.torrentId;
</script>

<!-- ========================= HTML -->
<div width="500" height="300">
	<video
		src={`http://localhost:3030/${torrentId}`}
		preload="metadata"
		controls
		autoplay
		muted
		style="width: 100rem"
		crossorigin="anonymous"
	>
		<track
			kind="captions"
			label="English"
			srclang="en"
			src="http://localhost:3030/899/subtitles/en"
			default
		/>
		Sorry, your browser doesn't support embedded videos.
	</video>
</div>

<!-- ========================= CSS -->
<style lang="postcss">
</style>
