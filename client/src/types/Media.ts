export type Result = {
	id: number;
	type: string;
	title: string;
	userTitle?: string;
	names: { lang: string; title: string }[];
	genres: string[];
	description: string;
	year?: number | null;
	duration: number | null;
	thumbnail: string;
	background?: string | null;
	rating: number | null;
};

export type MediaTorrent = {
	id: number;
	name: string;
	size?: string | null;
	leech: number;
	seed: number;
	quality?: string;
	hover: boolean;
};

export type MediaStaff = {
	id: number;
	name: string;
	thumbnail?: string | null;
	role?: string | null;
};

export type MediaActor = {
	id: number;
	name: string;
	thumbnail?: string | null;
	character?: string | null;
};

export type MediaComment = {
	id: number;
	user: {
		id: number;
		name: string;
	};
	date: string;
	content: string;
};

export type MediaProps = {
	media: Result;
	torrents: MediaTorrent[];
	staffs: MediaStaff[];
	actors: MediaActor[];
	comments: MediaComment[];
};
