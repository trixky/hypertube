import TailFile from "@logdna/tail-file";
import express from "express";
import fs from "fs";
import { download, DownloadInfo } from "./movie/downloader";
import { connect } from "./postgres/db";
import pump from "pump";

const app = express();

async function main() {
	await connect();

	app.get("/:torrent_id", async function (req, res) {
		const torrentId: number = parseInt(req.params.torrent_id);

		if (isNaN(torrentId) || torrentId < 0) {
			res.status(400).send();
		}

		const acceptHeader = req.headers.accept;
		const acceptMkv = (acceptHeader && (acceptHeader.indexOf("video/mkv") >= 0 || acceptHeader == "*/*")) == true;

		let downloadResult: DownloadInfo;
		try {
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 0.1");
			downloadResult = await download(`${torrentId}`, acceptMkv);
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 0.2");
		} catch {
			res.status(500).send();
			return;
		}

		if (downloadResult.path == null) {
			return res.status(500).send();
		}

		console.log("on s'en sort aveeeec download result: ", downloadResult);
		const extension = downloadResult.path.split(".").pop()!;
		const nativeExtensions = ["mp4", "webm"];
		const sendNative =
			(nativeExtensions.indexOf(downloadResult.original_extension) >= 0 ||
				(downloadResult.original_extension == "mkv" && acceptMkv)) &&
			downloadResult.length;

		console.log("acceptMkv", acceptMkv);
		console.log("extension", extension);
		console.log("sendNative", sendNative);
		console.log("using path", downloadResult.path);

		// * Downloading
		// * This can either be the original file if the torrent has an accepted file for the client
		// * -- or a transcode result in webm format
		// * The result is an used file path in either way, which is tailed
		if (!downloadResult.downloaded) {
			const file = new TailFile(downloadResult.path, {
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
				const realSize = downloadResult.length;
				res.writeHead(200, {
					"Content-Length": realSize,
					// "Content-Range": `bytes 0-${realSize - 1}/${realSize}`,
					// "Accept-Ranges": "bytes",
					"Content-Type": `video/${extension}`,
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
			pump(file, res, () => {
				file.quit();
			});
		}
		// * File is completed already
		// * Accept ranges on the file
		else {
			const size = sendNative ? downloadResult.length : fs.statSync(downloadResult.path).size;

			// Check and support ranges if the client supports them
			const range = req.headers["range"];
			let start = 0;
			let end = size - 1;
			// Handle invalid ranges with a 416
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.0");
			if (range) {
				var parts = range!.slice(6).split("-");
				if (parts.length >= 1 && parts.length < 3) {
					var partialStart = parseInt(parts[0]);
					var partialEnd = parseInt(parts[1]);
					if (!isNaN(partialStart) && partialStart >= 0 && partialStart < size) {
						start = partialStart;
					} else {
						console.error("Invalid range start", start, size);
						return res.status(416).header("Content-Range", `*/${size}`).send();
					}
					if (!isNaN(partialEnd)) {
						if (partialEnd > start && partialEnd <= size) {
							end = partialEnd;
						} else {
							console.error("Invalid range end");
							return res.status(416).header("Content-Range", `*/${size}`).send();
						}
					}
				} else {
					console.error("Invalid ranges format");
					return res.status(416).header("Content-Range", `*/${size}`).send();
				}
			}

			// If the browser has no support for range we just returns a 200 with no ranges
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.1");
			let returnCode = 200;
			let headers: Record<string, any> = {
				"Content-Length": size,
				"Content-Type": `video/${extension}`,
			};
			if (range) {
				returnCode = 206;
				headers["Content-Range"] = `bytes ${start}-${end - 1}/${size}`;
				headers["Accept-Ranges"] = "bytes";
			}
			res.writeHead(returnCode, headers);

			// Send the file read to the response
			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.2");
			const file = fs.createReadStream(downloadResult.path, { start, end });
			pump(file, res, () => {
				file.close();
			});

			console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2.3");
		}
	});

	app.listen(3030);
}

console.log("**************************************** START");
console.log("**************************************** START");
console.log("**************************************** START");

main();
