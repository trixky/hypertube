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

export type MediaProps = {
	media: Result;
	torrents: {
		id: number;
		name: string;
		size?: string | null;
		leech: number;
		seed: number;
		quality?: string;
		hover: boolean;
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
	comments: {
		id: number | string;
		user: {
			id: number | string;
			name: string;
		};
		date: string;
		content: string;
	}[];
};
