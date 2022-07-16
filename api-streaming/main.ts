import TailFile from "@logdna/tail-file";
import express from "express";
import fs from "fs";
import { download } from "./movie/downloader";
import { connect } from "./postgres/db";

const app = express();

let tutu = false;

async function main() {
	await connect();

	// download(25, "asdf")
	// download(899)

	// download(22) // marche pas !
	// download(29) // marche pas !
	// download(35) // marche pas !

	// download("magnet:?xt=urn:btih:4551CA5E03147242B45B94E909B43F4B5221B5E0&dn=Death.Hunt.2022.720p.WEBRip.800MB.x264-GalaxyRG&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.cyberia.is%3A6969%2Fannounce&tr=udp%3A%2F%2Fipv4.tracker.harry.lu%3A80%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.birkenwald.de%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.moeking.me%3A6969%2Fannounce&tr=udp%3A%2F%2Fopentor.org%3A2710%2Fannounce&tr=udp%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2970%2Fannounce&tr=https%3A%2F%2Ftracker.foreverpirates.co%3A443%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce")

	// download("magnet:?xt=urn:btih:2b3994ba55cbb4b7d256c5c6438eb329daa1400f&dn=%5BSaizen%5D%20Mahjong%20Soul%20-%20Akagi%20Crossover%20Event%20PV%20%5B1080p-Web%5D%5B4B130D34%5D.mkv&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce")

	// app.get("/", function (req, res) {
	//   res.send("Hello Worldds");
	// });

	app.head("/:torrent_id", () => {
		console.log("ASDFASD;LFKJASD;LKFJ");
		console.log("ASDFASD;LFKJASD;LKFJ");
		console.log("ASDFASD;LFKJASD;LKFJ");
		console.log("ASDFASD;LFKJASD;LKFJ");
		console.log("ASDFASD;LFKJASD;LKFJ");
		console.log("ASDFASD;LFKJASD;LKFJ");
		console.log("ASDFASD;LFKJASD;LKFJ");
	});

	app.get("/:torrent_id", async function (req, res) {
		const torrent_id: number = parseInt(req.params.torrent_id);

		if (isNaN(torrent_id) || torrent_id < 0) {
			res.status(400).send();
		}

		let file_path: string | null = null;
		try {
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 0.1");
			file_path = await download(`${torrent_id}`);
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 0.2");
		} catch {
			res.status(500).send();
			return;
		}

		setTimeout(() => {
			if (file_path == null) {
				res.status(500).send();
				return;
			}

			// const full_file_path = generate_full_paths(file_path)[1]
			// const full_file_path = "./.cache/movies/Green.Ghost.and.the.Masters.of.the.Stone.2022.720p.WEBRip.800MB.x264-GalaxyRG[TGx]/Green.Ghost.and.the.Masters.of.the.Stone.2022.720p.WEBRip.800MB.x264-GalaxyRG.mkv.mp4";
			// const full_file_path = "./.cache/movies/Green.Ghost.and.the.Masters.of.the.Stone.2022.720p.WEBRip.800MB.x264-GalaxyRG[TGx]/Green.Ghost.and.the.Masters.of.the.Stone.2022.720p.WEBRip.800MB.x264-GalaxyRG.mkv.webm";

			console.log("on s'en sort aveeeec full_file_path: " + file_path);

			const stat = fs.statSync(file_path);
			const partialTotal = stat.size;

			console.log("on s'en sort aveeeec total: " + partialTotal);
			console.log("req.headers['range'] ======= ???: " + req.headers["range"]);

			// * Downloading
			// * Assume transcode here
			if (true /* wasTutu */) {
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1.1");
				const file = new TailFile(file_path, {
					startPos: 0,
				})
					.on("tail_error", (err) => {
						console.error("TailFile had an error!", err);
					})
					.on("error", (err) => {
						console.error("A TailFile stream error was likely encountered", err);
					});
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1.2");
				file.start();
				// const file = fs.createReadStream(full_file_path, { start: 0 });
				// let length = 0;
				// const streamReader = new PassThrough();
				// streamReader.on("data", (chunk) => {
				// 	if ("length" in chunk) {
				// 		length += chunk;
				// 		console.log("============> PassTrough saw", chunk.length);
				// 	}
				// });
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1.3");
				res.writeHead(206, {
					// "Content-Length": partialTotal,
					"Content-Range": `bytes 0-`,
					"Accept-Ranges": "bytes",
					"Content-Type": "video/webm",
				});
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1.4", {
					// "Content-Length": partialTotal,
					"Content-Range": `bytes 0-`,
					"Accept-Ranges": "bytes",
					"Content-Type": "video/webm",
				});
				// file.pipe(streamReader);
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1.5");
				// streamReader.pipe(res);
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1.6");
				file.pipe(res);
			}
			// * File is completed already
			// * Accept ranges on the file
			else {
				// TODO No support for ranges ? Send 200 with the whole file
				const range = req.headers["range"];
				let start = 0;
				let end = partialTotal - 1;
				if (range) {
					var parts = range!.slice(6).split("-");
					if (parts.length >= 1) {
						var partialStart = parseInt(parts[0]);
						var partialEnd = parseInt(parts[1]);
						if (isNaN(partialStart) && partialStart >= 0 && partialStart < partialTotal) {
							start = partialStart;
						} else {
							// ! Invalid Range
						}
						if (!isNaN(partialEnd)) {
							if (partialEnd > start && partialEnd <= partialTotal) {
								end = partialEnd;
							} else {
								// ! Invalid Range
							}
						}
					}
				}
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.1");
				res.writeHead(206, {
					"Content-Length": partialTotal,
					"Content-Range": `bytes ${start}-${end - 1}/${partialTotal}`,
					"Accept-Ranges": "bytes",
					"Content-Type": "video/mp4",
				});
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.2", {
					"Content-Length": partialTotal,
					"Content-Range": `bytes ${start}-${end - 1}/${partialTotal}`,
					"Accept-Ranges": "bytes",
					"Content-Type": "video/mp4",
				});
				const file = fs.createReadStream(file_path!, { start, end });
				file.pipe(res);
			}
		}, 10000);
	});

	app.listen(3030);
}

// function main_2() {
//   download_2("TR45")
// }

console.log("**************************************** START");
console.log("**************************************** START");
console.log("**************************************** START");

main();
