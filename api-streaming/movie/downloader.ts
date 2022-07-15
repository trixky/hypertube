const torrentStream = require('torrent-stream'); // https://github.com/mafintosh/torrent-stream#readme
const ffmpeg = require('fluent-ffmpeg')
const fs = require('fs');
const { get_one, update_one } = require('../postgres/movies');
const parseTorrent = require('parse-torrent')

const cache_path = './cache'
const cache_path_movies = './cache/movies'

const downloaded_torrent_movies = new Map<number, { file_path: string | null, downloaded: boolean }>(); // movie_id // false: in progress, true: finished

export function generate_full_paths(file_path: string): [string, string] {
    const basic_path = cache_path_movies + '/' + file_path
    return [basic_path, basic_path + '.mp4']
}

function torrent_to_magnet(torrent_url: string): Promise<string> {
    return new Promise((resolve) => {
        parseTorrent.remote(torrent_url, { timeout: 60 * 1000 }, (err: Error, parsedTorrent: any) => {
            if (err) throw err
            resolve(parseTorrent.toMagnetURI(parsedTorrent))
        })
    })
}

export function download(id: number): Promise<string | null> {
    return new Promise(async resolve => {
        // check if torrent is not already downloaded from local map
        if (downloaded_torrent_movies.has(id) && downloaded_torrent_movies.get(id)?.downloaded == true) return downloaded_torrent_movies.get(id)?.file_path || null

        console.log("_______________________ 0")

        // get torrent infos from postgres
        let res: {
            torrent_url: string | null,
            magnet: string | null,
            file_path: string | null,
            downloaded: boolean | null,
        } = await get_one(id)

        let file_path = res.file_path

        console.log("_______________________ 1")

        // check if torrent is not already downloaded from postgres
        if (res.downloaded === true) {
            downloaded_torrent_movies.set(id, {
                file_path: file_path,
                downloaded: true,
            })
            resolve(file_path)
        }

        console.log("_______________________ 2")

        console.log("----------------- va ton delete ?")
        // check if torrent not start to be downloaded from postgres and is unknown from local map
        if (!downloaded_torrent_movies.has(id) && res.downloaded === false && file_path != null) {
            // delete the started downloaded movie

            const full_paths = generate_full_paths(file_path)
            fs.unlink(full_paths[0], () => { })
            fs.unlink(full_paths[1], () => { })
        }

        console.log("_______________________ 3")

        let magnet = res.magnet

        // extract or generate the magnet
        if (res.magnet != null) {
            console.log("_______________________ 3.1 magnet")
            magnet = res.magnet
        } else if (res.torrent_url != null) {
            console.log("_______________________ 3.2 torrent_url [" + res.torrent_url + "]")
            magnet = await torrent_to_magnet(res.torrent_url)
        } else {
            console.log("_______________________ 3.3 rien ...")
            throw new Error('no magnet or torrent url available for this torrent')
        }

        console.log("_______________________ 4")

        // save torrent status in local map
        downloaded_torrent_movies.set(id, {
            file_path: file_path,
            downloaded: false,
        })

        console.log("_______________________ 5")

        // start to download the torrent movie
        var engine = torrentStream(magnet, {
            connections: 100,
            uploads: 10, // 0 ?
            tmp: cache_path,
            path: cache_path_movies,
            verify: true,
            tracker: true, // false ?
        })

        console.log("_______________________ 6")

        let selected_file: { name: string; path: string, length: number } | null = null;

        console.log("_______________________ 7")

        engine.on('ready', async function () {
            console.log("_______________________ 8")

            // select the good file in the torrent (using extensions)
            engine.files.every(function (file: { name: string; path: string, length: number, createReadStream: Function }) {
                console.log("_______________________ 9: ", file.name)
                const file_name_length = file.name.length

                // if file have a correct extension length

                if (file_name_length > 4) {
                    console.log("_______________________ 10")
                    const extension = file.name.slice(file_name_length - 4, file_name_length)
                    // if file is a mp4
                    if (extension === ".mp4") {
                        console.log("_______________________ 11 mp4")
                        selected_file = file
                        file.createReadStream();
                        return false
                    } else if (extension === ".mkv") {
                        console.log("_______________________ 11.1 mkv")
                        // if file is a mkv
                        selected_file = file

                        // create a stream for transcode mkv to mp4
                        let local_file_path = cache_path_movies + '/' + file.path
                        let writeStream = fs.createWriteStream(local_file_path + '.mp4');

                        console.log("_______________________ 11.2 mkv")

                        const stream = file.createReadStream();
                        ffmpeg()
                            .input(stream)
                            .inputFormat('matroska')
                            .audioCodec('aac')
                            .videoCodec('libx264')
                            .outputOptions('-movflags frag_keyframe+empty_moov')
                            .outputFormat('mp4')
                            .on('start', () => {
                                console.log('start')
                            })
                            .on('progress', (progress: { timemark: string }) => {
                                console.log(`progress: ${progress.timemark}`)
                            })
                            .on('end', () => {
                                console.log('Finished processing')
                                downloaded_torrent_movies.set(id, {
                                    file_path: file_path,
                                    downloaded: true,
                                })
                                fs.unlink(file_path, () => { })
                            })
                            .on('error', (err: Error) => {
                                console.log(`ERROR: ${err.message}`)
                                downloaded_torrent_movies.delete(id)
                                fs.unlink(file_path, () => { })
                            })
                            .pipe(writeStream)
                        return false
                    } else {
                        file.createReadStream();
                    }
                    console.log("-------------- rien")
                    return true
                }
            })

            if (selected_file == null) {
                console.log("************** BOOOOOOOM NANNN")
                engine.destroy()
                throw "no .mp4 or .mkv finded"
            }

            console.log("_______________________ 12")
            await update_one(id, selected_file.path, false)

            console.log("_______________________ 13")
            engine.on("download", (index: string) => {
                console.log(`state for: ${index}`);
                resolve(selected_file?.path || null)
            });

            console.log("_______________________ 14")
            engine.on("idle", async () => {
                console.log('downloaded');
                if (selected_file != null) await update_one(id, selected_file.path, true)
                resolve(selected_file?.path || null)
            });

            engine.once("destroyed", () => engine.removeAllListeners());
        })

        return file_path
    })
}
