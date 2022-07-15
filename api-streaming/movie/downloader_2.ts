const torrentStream = require('torrent-stream'); // https://github.com/mafintosh/torrent-stream#readme
const { get_one } = require('../postgres/movies');
const parseTorrent = require('parse-torrent')

const cache_path = './.cache'
const cache_path_movies = cache_path + '/movies'
const extension_mp4 = '.mp4'
const extension_mkv = '.mkv'

interface LocalTorrent {
    file_path: string;
    is_mkv: boolean;
    downloaded: boolean;
    corrupted: boolean;
}

function generate_full_path(local_torrent_info: LocalTorrent): string {
    return cache_path_movies + local_torrent_info.file_path + (local_torrent_info.is_mkv ? extension_mp4 : '');
}

const local_torrents = new Map<number, LocalTorrent>()

function sanitize_torrent_id(torrent_id: string): number {
    const sanitized_torrent_id = parseInt(torrent_id)

    if (isNaN(sanitized_torrent_id))
        throw new Error('torrent id corrupted')
    if (sanitized_torrent_id < 0)
        throw new Error('torrent id need to be positive')

    return sanitized_torrent_id
}

function torrent_to_magnet(torrent_url: string): Promise<string> {
    return new Promise((resolve) => {
        parseTorrent.remote(torrent_url, { timeout: 60 * 1000 }, (err: Error, parsedTorrent: any) => {
            if (err) throw err
            resolve(parseTorrent.toMagnetURI(parsedTorrent))
        })
    })
}

async function get_magnet(torrent_id: number): Promise<string> {
    // get torrent infos from postgres
    let res: {
        torrent_url: string | null,
        magnet: string | null,
        file_path: string | null,
        downloaded: boolean | null,
    } = await get_one(torrent_id)

    let magnet = res.magnet

    // extract or generate the magnet
    if (res.magnet != null) {
        magnet = res.magnet
    } else if (res.torrent_url != null) {
        magnet = await torrent_to_magnet(res.torrent_url)
    } else {
        throw new Error('no magnet or torrent url available for this torrent')
    }

    return magnet
}

function get_movie_file_from_engine(engine: any): Promise<string> {
    return new Promise(resolve )
    engine.on('ready', async function () {
        console.log("_______________________ 8")

        // select the good file in the torrent (using extensions)
        engine.files.every(function (file: { name: string; path: string, length: number, createReadStream: Function }) {
            console.log("_______________________ 9: ", file.name)
            const file_name_length = file.name.length
        })
    })
}

export async function download_2(torrent_id: string): Promise<string | null> {
    // sanitize the torrent id
    let sanitized_torrent_id: number
    try {
        sanitized_torrent_id = sanitize_torrent_id(torrent_id)
    } catch {
        return null
    }

    // check if torrent is known in local
    if (local_torrents.has(sanitized_torrent_id)) {
        const local_torrent_info = local_torrents.get(sanitized_torrent_id)

        if (local_torrent_info == undefined) return null
        // check if torrent is corrupted in local
        if (local_torrent_info.corrupted) return null

        return generate_full_path(local_torrent_info)
    }

    // get the magnet of the torrent
    let magnet: string;
    try {
        magnet = await get_magnet(sanitized_torrent_id)
    } catch {
        // if the torrent is corrupted
        local_torrents.set(sanitized_torrent_id, <LocalTorrent>{
            file_path: '',
            is_mkv: false,
            downloaded: false,
            corrupted: true,
        })

        return null
    }

    // start the torrent engine
    const engine = torrentStream(magnet, {
        connections: 100,
        uploads: 10, // 0 ?
        tmp: cache_path,
        path: cache_path_movies,
        verify: true,
        tracker: true, // false ?
    })

    const movie_file = get_movie_file_from_engine(engine)

    return null
}