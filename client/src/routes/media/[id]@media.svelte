<!-- ========================= SCRIPT -->
<script context="module" lang="ts">
	import type { Load } from '@sveltejs/kit';

	type Media = {
		id: number;
		overview: string;
	};

	type MediaProps = {
		media: {
			id: number;
			description: string;
			duration?: number | null;
			genres: string[];
			names: { lang: string; title: string }[];
			rating?: number | null;
			thumbnail?: string | null;
			background?: string | null;
			type: string;
			year?: number | null;
		};
		torrents: {
			id: number;
			leech: number;
			name: string;
			seed: number;
		}[];
		staffs: {
			id: number;
			name: string;
			thumbnail?: string | null;
			role?: string | null;
		}[];
		actors: {
			id: number;
			name: string;
			thumbnail?: string | null;
			character?: string | null;
		}[];
	};

	export const load: Load = async ({ params, fetch, session, stuff }) => {
		const url = browser
			? `http://localhost:7072/v1/media/get/${params.id}`
			: `http://api-search:7072/v1/media/get/${params.id}`;
		const response = await fetch(url, {
			method: 'GET',
			headers: { accept: 'application/json' }
		});

		let props: MediaProps | false = false;
		if (response.ok) {
			props = (await response.json()) as MediaProps;
		}

		return {
			status: response.status,
			redirect: !response.ok ? '/search' : undefined,
			props: { props }
		};
	};
</script>

<script lang="ts">
	import { page } from '$app/stores';
	import { browser } from '$app/env';

	/// @ts-expect-error media is given as a prop
	export let props: MediaProps;
	const { media, torrents, staffs, actors } = props;
</script>

<!-- ========================= HTML -->
<div>
	<a href="/search">Back</a>
	<h1>Media</h1>
	<div>{media.id}</div>
	<div>{media.description}</div>
	<div>{media.duration}</div>
	<div>{media.genres.join(',')}</div>
	<div>{media.names.map((e) => `${e.lang}=${e.title}`).join(',')}</div>
	<div>{media.rating}</div>
	<div>{media.thumbnail}</div>
	<div>{media.background}</div>
	<div>{media.type}</div>
	<div>{media.year}</div>
	<h1>Torrents</h1>
	{#each torrents as torrent (torrent.id)}
		<div>{torrent.id}</div>
		<div>{torrent.leech}</div>
		<div>{torrent.seed}</div>
		<div>{torrent.name}</div>
	{/each}
	<h1>Staffs</h1>
	{#each staffs as staff (staff.id + (staff.role ?? ''))}
		<div>{staff.id}</div>
		<div>{staff.name}</div>
		<div>{staff.role}</div>
		<div>{staff.thumbnail}</div>
	{/each}
	<h1>Actors</h1>
	{#each actors as actor (actor.id + (actor.character ?? ''))}
		<div>{actor.id}</div>
		<div>{actor.name}</div>
		<div>{actor.character}</div>
		<div>{actor.thumbnail}</div>
	{/each}
</div>

<!-- ========================= CSS -->
<style lang="postcss">
</style>
