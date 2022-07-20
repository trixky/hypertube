import { fetch } from 'undici';
import { URL, URLSearchParams } from 'url';
import env from '../env';

const ApiKey = env.OSDB_API_KEY || '';

export type SearchInformations = {
	imdb_id?: number | null;
	tmdb_id?: number | null;
	torrent_name: string;
	media_name: string;
	year?: number | null;
};

export type OSDBSubtitle = {
	id: number;
	type: string;
	attributes: {
		subtitle_id: string;
		language: string;
		download_count: number;
		new_download_count: number;
		hearing_impaired: boolean;
		hd: boolean;
		format: string;
		fps: number;
		votes: number;
		points: number;
		ratings: number;
		from_trusted: boolean;
		foreign_parts_only: boolean;
		ai_translated: boolean;
		machine_translated: boolean;
		upload_date: string;
		release: string;
		comments: string;
		legacy_subtitle_id: number;
		uploader: {
			uploader_id: number;
			name: string;
			rank: string;
		};
		feature_details: {
			feature_id: number;
			feature_type: string;
			year: number;
			title: string;
			movie_name: string;
			imdb_id: number;
			tmdb_id: number;
		};
		url: string;
		related_links: {
			label: string;
			url: string;
			img_url: string;
		}[];
		files: {
			file_id: number;
			cd_number: number;
			file_name: string;
		}[];
	};
};

export type SearchResponse = {
	total_pages: number;
	total_count: number;
	page: number;
	data: OSDBSubtitle[];
};

class OSDBClass {
	token: string;
	blocked: boolean;
	reset: Date | null;

	constructor() {
		this.token = '';
		this.blocked = false;
		this.reset = null;
	}

	resolveBlock(): boolean {
		if (this.blocked) {
			if (!this.reset) {
				this.blocked = false;
			} else if (this.reset.getTime() <= Date.now()) {
				this.blocked = false;
				this.reset = null;
			}
		}

		return !this.blocked;
	}

	async login() {
		const response = await fetch('https://api.opensubtitles.com/api/v1/login', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Api-Key': ApiKey
			},
			body: JSON.stringify({
				username: env.OSDB_USERNAME,
				password: env.OSDB_PASSWORD
			})
		});
		if (response.ok) {
			const body = (await response.json()) as {
				user: {
					allowed_downloads: number;
					level: string;
					user_id: number;
					ext_installed: boolean;
					vip: boolean;
				};
				token: string;
				status: number;
			};
			this.token = body.token;
			return true;
		}

		this.token = '';
		return false;
	}

	async search(
		searchInformations: SearchInformations,
		doRetry = true
	): Promise<OSDBSubtitle[] | false> {
		if (!this.token) {
			await this.login();
		}

		// Construct search URL
		// Use the tmdb_id or the imdb_id if there is one, else search by name and year
		const url = new URL('https://api.opensubtitles.com/api/v1/subtitles');
		const params = new URLSearchParams({
			type: 'movie',
			languages: 'en,fr',
			order_by: 'download_count'
		});
		if (searchInformations.tmdb_id) {
			params.append('tmdb_id', `${searchInformations.tmdb_id}`);
		} else if (searchInformations.imdb_id) {
			params.append('imdb_id', `${searchInformations.imdb_id}`);
		} else {
			// if (searchInformations.year) {
			// 	params.append('year', `${searchInformations.year}`)
			// }
			// params.append('query', searchInformations.media_name)
			params.append('query', searchInformations.torrent_name);
		}

		const response = await fetch(`${url}?${params}`, {
			headers: {
				'Content-Type': 'application/json',
				'Api-Key': ApiKey
			}
		});
		if (response.ok) {
			const body = (await response.json()) as SearchResponse;
			if (body.total_count > 0) {
				return body.data;
			}
		} else if (response.status == 401 && doRetry) {
			return this.search(searchInformations, false);
		}

		return false;
	}

	async download(
		downloadInformations: {
			fileId: number;
			lang: string;
			torrentId: number;
		},
		doRetry = true
	): Promise<{ extension: string; buffer: string } | false> {
		if (!this.token) {
			await this.login();
		}

		// Initial request to get the file link
		console.log('Downloading', downloadInformations.fileId);
		const response = await fetch('https://api.opensubtitles.com/api/v1/download', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Api-Key': ApiKey
			},
			body: JSON.stringify({
				file_id: downloadInformations.fileId
			})
		});
		if (response.ok) {
			const body = (await response.json()) as {
				link: string;
				file_name: string;
				requests: number;
				remaining: number;
				message: string;
				reset_time: string;
				reset_time_utc: string;
			};

			// Update block status
			console.log('Remaining quota', body.remaining, 'until', body.reset_time_utc);
			if (body.remaining == 0) {
				this.blocked = true;
				this.reset = new Date(body.reset_time_utc);
			}

			// Then download the file
			const downloadResponse = await fetch(body.link, {
				method: 'GET',
				headers: {
					Authorization: `Bearer ${this.token}`
				}
			});
			if (response.ok) {
				const subtitle = await downloadResponse.text();
				const extension = body.file_name.split('.').pop()!;
				return { extension, buffer: subtitle };
			} else {
				console.error('Failed to download', downloadInformations.fileId);
				this.blocked = true;
				this.reset = new Date();
				this.reset.setMinutes(this.reset.getMinutes() + 5);
			}
		} else if (response.status == 401 && doRetry) {
			return this.download(downloadInformations, false);
		}

		return false;
	}
}

export default new OSDBClass();
