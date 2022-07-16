import TailFile from "@logdna/tail-file";
import express from "express";
import fs from "fs";
import { download } from "./movie/downloader";
import { connect } from "./postgres/db";
import pump from "pump";

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

		const acceptHeader = req.headers.accept;
		const acceptMkv = (acceptHeader && acceptHeader.indexOf("video/mkv") >= 0) == true;

		let file_path: string | null = null;
		try {
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 0.1");
			file_path = await download(`${torrent_id}`, acceptMkv);
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 0.2");
		} catch {
			res.status(500).send();
			return;
		}

		setTimeout(() => {
			if (file_path == null) {
				return res.status(500).send();
			}

			console.log("on s'en sort aveeeec full_file_path: " + file_path);

			const stat = fs.statSync(file_path);
			const partialTotal = stat.size;

			console.log("on s'en sort aveeeec total: " + partialTotal);
			console.log("req.headers['range'] ======= ???: " + req.headers["range"]);

			// * Downloading
			// * This can either be the original file if the torrent has an accepted file for the client
			// * -- or a transcode result in webm format
			// * The result is an used file path in either way, which is tailed
			if (true /* downloading */) {
				let sendNative = true;
				const file = new TailFile(file_path!, {
					startPos: 0,
				})
					.on("tail_error", (err: any) => {
						console.error("TailFile had an error!", err);
					})
					.on("error", (err: any) => {
						console.error("A TailFile stream error was likely encountered", err);
					});
				file.start();
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1.3");

				// When sending the file that existed in the torrent we already have all of the size informations ...
				if (sendNative) {
					const realSize = 123456789;
					const nativeExtension = "mp4";
					res.writeHead(200, {
						"Content-Length": realSize,
						// "Content-Range": `bytes 0-${realSize - 1}/${realSize}`,
						// "Accept-Ranges": "bytes",
						"Content-Type": `video/${nativeExtension}`,
					});
				}
				// ... else the file comes from a transcode and we don't know the final file size
				// TODO check Transfer-Encoding: chunked for chrome, maybe this works ?
				else {
					res.writeHead(206, {
						"Content-Range": `bytes 0-`,
						"Accept-Ranges": "bytes",
						"Content-Type": "video/webm",
					});
				}

				// Send the file read to the response
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1.5");
				pump(file, res);
			}
			// * File is completed already
			// * Accept ranges on the file
			else {
				const usedFileExtension = "mp4";

				// Check and support ranges if the client supports them
				const range = req.headers["range"];
				let start = 0;
				let end = partialTotal - 1;
				// Handle invalid ranges with a 416
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.0");
				if (range) {
					var parts = range!.slice(6).split("-");
					if (parts.length >= 1 && parts.length < 3) {
						var partialStart = parseInt(parts[0]);
						var partialEnd = parseInt(parts[1]);
						if (!isNaN(partialStart) && partialStart >= 0 && partialStart < partialTotal) {
							start = partialStart;
						} else {
							console.error("Invalid range start", start, partialTotal);
							return res.status(416).header("Content-Range", `*/${partialTotal}`).send();
						}
						if (!isNaN(partialEnd)) {
							if (partialEnd > start && partialEnd <= partialTotal) {
								end = partialEnd;
							} else {
								console.error("Invalid range end");
								return res.status(416).header("Content-Range", `*/${partialTotal}`).send();
							}
						}
					} else {
						console.error("Invalid ranges format");
						return res.status(416).header("Content-Range", `*/${partialTotal}`).send();
					}
				}

				// If the browser has no support for range we just returns a 200 with no ranges
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.1");
				let returnCode = 200;
				let headers: Record<string, any> = {
					"Content-Length": partialTotal,
					"Content-Type": `video/${usedFileExtension}`,
				};
				if (range) {
					returnCode = 206;
					headers["Content-Range"] = `bytes ${start}-${end - 1}/${partialTotal}`;
					headers["Accept-Ranges"] = "bytes";
				}
				res.writeHead(returnCode, headers);

				// Send the file read to the response
				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.2");
				const file = fs.createReadStream(file_path!, { start, end });
				pump(file, res);

				console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.3");
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
