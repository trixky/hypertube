import express from "express";
import { download, generate_full_paths } from "./movie/downloader"
import { download_2 } from "./movie/downloader_2"
const fs = require('fs');
import { connect } from "./postgres/db"

const app = express();

let tutu = false;

async function main() {
  await connect()

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
    console.log("ASDFASD;LFKJASD;LKFJ")
    console.log("ASDFASD;LFKJASD;LKFJ")
    console.log("ASDFASD;LFKJASD;LKFJ")
    console.log("ASDFASD;LFKJASD;LKFJ")
    console.log("ASDFASD;LFKJASD;LKFJ")
    console.log("ASDFASD;LFKJASD;LKFJ")
    console.log("ASDFASD;LFKJASD;LKFJ")
  })

  app.get("/:torrent_id", async function (req, res) {
    const torrent_id: number = parseInt(req.params.torrent_id)

    if (isNaN(torrent_id) || torrent_id < 0) {
      res.status(400).send();
    }

    let file_path: string | null = null

    if (tutu === false) {
      tutu = true
      try {
        file_path = await download(torrent_id)
      } catch {
        res.status(500).send();
        return
      }
      if (file_path == null) {
        res.status(500).send();
        return
      }
    }

    setTimeout(() => {
      // const full_file_path = generate_full_paths(file_path)[1]
      const full_file_path = "./cache/movies/Green.Ghost.and.the.Masters.of.the.Stone.2022.720p.WEBRip.800MB.x264-GalaxyRG[TGx]/Green.Ghost.and.the.Masters.of.the.Stone.2022.720p.WEBRip.800MB.x264-GalaxyRG.mkv.mp4"
      // const full_file_path = "./cache/movies/cat.mp4"

      console.log("on s'en sort aveeeec full_file_path: " + full_file_path)

      const stat = fs.statSync(full_file_path);
      const total = stat.size;

      console.log("on s'en sort aveeeec total: " + total)
      console.log("req.headers['range'] ======= ???: " + req.headers['range'])

      if (req.headers['range'] && req.headers['range'] != 'bytes=0-') {
        console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 1: " + req.headers.range)
        var range = req.headers['range'];
        var parts = range.replace(/bytes=/, "").split("-");
        var partialstart = parts[0];
        var partialend = parts[1];
        console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 2")

        var start = parseInt(partialstart, 10);
        var end = partialend ? parseInt(partialend, 10) : total - 1;
        var chunksize = (end - start) + 1;
        console.log('RANGE: ' + start + ' - ' + end + ' = ' + chunksize);

        console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 3: " + [start, end])

        var file = fs.createReadStream(full_file_path, { start: start, end: end });
        console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 4")
        res.writeHead(206, { 'Content-Range': 'bytes ' + start + '-' + end + '/' + total, 'Accept-Ranges': 'bytes', 'Content-Length': chunksize, 'Content-Type': 'video/mp4' });
        console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 5")
        file.pipe(res);
      } else {
        console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@ ******* 1")
        console.log('ALL: ' + total);
        res.writeHead(200, { 'Content-Length': total, 'Content-Type': 'video/mp4' });
        // res.writeHead(206, { 'Content-Range': 'bytes ' + 0 + '-' + 150000 + '/' + (total - 1), 'Accept-Ranges': 'bytes', 'Content-Length': 150000 + 1, 'Content-Type': 'video/mp4' }); // >
        console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@ ******* 2")
        var file = fs.createReadStream(full_file_path, { start: 0, end: total });
        console.log("@@@@@@@@@@@@@@@@@@@@@@@@@@@@ ******* 3")
        file.pipe(res);
      }
      // res.send("Hello Worldds");
    }, 15000);

  });


  app.listen(3030);
}

// function main_2() {
//   download_2("TR45")
// }

console.log("**************************************** START")
console.log("**************************************** START")
console.log("**************************************** START")

main()