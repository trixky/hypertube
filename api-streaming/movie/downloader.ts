import torrentStream from "torrent-stream"; // https://github.com/mafintosh/torrent-stream#readme
import { get_one, update_one } from "../postgres/movies";
import parseTorrent from "parse-torrent";
import ffmpeg from "fluent-ffmpeg";
import fs from "fs";
import { sleep } from "../utils/time";

const cache_path = "./.cache";
const cache_path_movies = cache_path + "/movies/";
const extension_mp4 = ".mp4";
const extension_webm = ".webm";
const extension_mkv = ".mkv";

interface DBTorrent {
	torrent_url: string | null;
	magnet: string | null;
	file_path: string | null;
	downloaded: boolean | null;
	is_mkv: boolean | null;
}

interface LocalTorrent {
	file_path: string;
	downloaded: boolean;
	is_mkv: boolean;
	corrupted: boolean;
}

interface SelectedFile {
	name: string;
	path: string;
	length: number;
	createReadStream: Function;
	is_mkv: boolean;
}

const local_torrents = new Map<number, LocalTorrent>();

function is_mp4(path: string): boolean {
	const file_name_length = path.length;

	return Boolean(path.length > 4 && path.slice(file_name_length - 4, file_name_length) === extension_mp4);
}

function is_mkv(path: string): boolean {
	const file_name_length = path.length;

	return Boolean(path.length > 4 && path.slice(file_name_length - 4, file_name_length) === extension_mkv);
}

function generate_full_path(file_path: string, block_extention_adding: boolean): string {
	return cache_path_movies + file_path + (is_mkv(file_path) && !block_extention_adding ? extension_webm : "");
}

function sanitize_torrent_id(torrent_id: string): number {
	const sanitized_torrent_id = parseInt(torrent_id);

	if (isNaN(sanitized_torrent_id)) throw new Error("torrent id corrupted");
	if (sanitized_torrent_id < 0) throw new Error("torrent id need to be positive");

	return sanitized_torrent_id;
}

function torrent_to_magnet(torrent_url: string): Promise<string> {
	return new Promise((resolve) => {
		parseTorrent.remote(torrent_url, (err: Error, parsedTorrent: any) => {
			if (err) throw err;
			resolve(parseTorrent.toMagnetURI(parsedTorrent));
		});
	});
}

async function get_magnet(db_torrent_infos: DBTorrent): Promise<string> {
	let magnet = db_torrent_infos.magnet;

	// extract or generate the magnet
	if (magnet == null) {
		if (db_torrent_infos.torrent_url === null) throw new Error("no magnet or torrent url available for this torrent");
		magnet = await torrent_to_magnet(db_torrent_infos.torrent_url);
	}

	return magnet;
}

function get_movie_file_from_engine(engine: any): Promise<SelectedFile> {
	return new Promise((resolve, reject) => {
		engine.on("ready", async function () {
			engine.files.every(function (file: SelectedFile) {
				const file_is_mp4 = is_mp4(file.name);
				const file_is_mkv = file_is_mp4 ? false : is_mkv(file.name);

				if (file_is_mp4 || file_is_mkv) {
					file.is_mkv = file_is_mkv;
					resolve(file);
					return false;
				}
				return true;
			});
			reject();
		});
	});
}

function start_download_mp4(selected_file: SelectedFile) {
	selected_file.createReadStream();
}

function start_download_mkv(torrent_id: number, selected_file: SelectedFile) {
	let local_file_path = generate_full_path(selected_file.path, true);
	let local_file_path_mp4 = local_file_path + extension_webm;

	const stream = selected_file.createReadStream();

	ffmpeg()
		.input(stream)
		.inputFormat("matroska")
		// * mp4
		// .audioCodec("aac")
		// .videoCodec("libx264")
		// .outputOptions("-movflags frag_keyframe+empty_moov")
		// .outputFormat("mp4")
		// * webm
		.audioCodec("libvorbis")
		.videoCodec("libvpx-vp9")
		.videoBitrate(20)
		.outputOptions("-vf scale=-1:101")
		.outputOptions("-preset veryfast")
		.outputOptions("-crf 50")
		.outputOptions("-movflags frag_keyframe+empty_moov")
		.outputFormat("webm")
		// *
		.on("ffmpeg: start", () => {
			console.log("start");
		})
		.on("progress", (progress: { timemark: string }) => {
			console.log(`ffmpeg: progress > ${progress.timemark}`);
		})
		.on("end", () => {
			console.log("ffmpeg: Finished processing");

			local_torrents.set(torrent_id, <LocalTorrent>{
				file_path: selected_file.path,
				downloaded: true,
				is_mkv: true,
				corrupted: false,
			});
			update_one(torrent_id, selected_file.path, true);
		})
		.on("error", (err: Error) => {
			console.log(`ffmpeg: ERROR > ${err.message}`);
			local_torrents.delete(torrent_id);
		})
		.output(local_file_path_mp4)
		.run();
}

async function wait_file_path(torrent_id: number): Promise<string> {
	let file_path: string | undefined = undefined;

	while (file_path === undefined || file_path === "") {
		await sleep(1000);
		file_path = local_torrents.get(torrent_id)?.file_path;
	}

	return file_path;
}

export async function download(torrent_id: string, want_mkv: boolean): Promise<string | null> {
	return new Promise(async (resolve, reject) => {
		// sanitize the torrent id
		let sanitized_torrent_id: number;
		try {
			sanitized_torrent_id = sanitize_torrent_id(torrent_id);
		} catch (err) {
			reject(err);
			return null;
		}

		// check if torrent is known in local
		if (local_torrents.has(sanitized_torrent_id)) {
			const local_torrent_info = local_torrents.get(sanitized_torrent_id);

			// check if torrent is corrupted in local
			if (local_torrent_info == undefined || local_torrent_info.corrupted) {
				reject();
				return null;
			}

			resolve(generate_full_path(await wait_file_path(sanitized_torrent_id), want_mkv));
			return;
		}
		local_torrents.set(sanitized_torrent_id, <LocalTorrent>{
			file_path: "",
			is_mkv: false,
			downloaded: false,
			corrupted: false,
		});

		// get the torrent infos from db
		let db_torrent_infos: DBTorrent = <DBTorrent>{};
		try {
			const res = await get_one(sanitized_torrent_id);
			db_torrent_infos.downloaded = res.downloaded;
			db_torrent_infos.file_path = res.file_path;
			db_torrent_infos.is_mkv = null;
			db_torrent_infos.magnet = res.magnet;
			db_torrent_infos.torrent_url = res.torrent_url;
		} catch (err) {
			reject(err);
			return;
		}

		// check if the torrent is already downloaded
		if (db_torrent_infos.downloaded) {
			if (db_torrent_infos.file_path == null) {
				reject(new Error("no path for the downloaded movie from db"));
				return;
			}
			resolve(generate_full_path(db_torrent_infos.file_path, want_mkv));
			return;
		}

		// get the magnet of the torrent
		let magnet: string;
		try {
			magnet = await get_magnet(db_torrent_infos);
		} catch {
			// if the torrent is corrupted
			local_torrents.set(sanitized_torrent_id, <LocalTorrent>{
				file_path: "",
				is_mkv: false,
				downloaded: false,
				corrupted: true,
			});
			reject(new Error("magnet/torrent corrupted"));
			return;
		}

		// start the torrent engine
		const engine = torrentStream(magnet, {
			connections: 100,
			uploads: 10, // 0 ?
			tmp: cache_path,
			path: cache_path_movies,
			verify: true,
			tracker: true, // false ?
		});

		// select the good file in the torrent (using extensions)
		let movie_file: SelectedFile;
		try {
			movie_file = await get_movie_file_from_engine(engine);
			resolve(generate_full_path(movie_file.path, want_mkv));
		} catch {
			engine.destroy(() => {});
			reject(new Error("no movie finded in the torrent"));
			return;
		}

		// save torrent/movie in local
		local_torrents.set(sanitized_torrent_id, <LocalTorrent>{
			file_path: movie_file.path,
			is_mkv: false,
			downloaded: false,
			corrupted: false,
		});

		// start download
		if (movie_file.is_mkv) {
			start_download_mkv(sanitized_torrent_id, movie_file);
		} else {
			start_download_mp4(movie_file);
		}

		engine.on("download", (index: string) => {
			console.log(`state for: ${index}`);
		});

		engine.on("idle", async () => {
			engine.destroy(() => {}); // ?
		});
	});
}
