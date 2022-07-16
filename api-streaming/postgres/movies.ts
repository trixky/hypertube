import { client } from './db'

export async function get_one(id: number): Promise<{
    torrent_url: string | null,
    magnet: string | null,
    file_path: string | null,
    downloaded: boolean | null,
}> {
    console.log("get_one 1")
    let res = await client.query("SELECT torrent_url, magnet, file_path, downloaded \
    FROM torrents \
    WHERE id = $1;", [id])
    console.log("get_one 2")
    
    console.log("get_one 3")
    
    if (res.rows.length == 1) {
        console.log("get_one 4")
        return res.rows[0]
    } else {
        console.log("get_one ERROR")
        throw ('no torrent finded with this id')
    }
}

export async function update_one(id: number, file_path: string | null, downloaded: boolean | null): Promise<boolean> {
    const update_strings = []
    const update_values: Array<string | number | boolean> = [id]

    let arg_nbr = 2

    if (file_path != null) {
        update_strings.push('file_path = $' + arg_nbr++)
        update_values.push(file_path)
    }
    if (downloaded != null) {
        update_strings.push('downloaded = $' + arg_nbr++)
        update_values.push(downloaded)
    }


    let res: { rowCount: number } = await client.query(`UPDATE torrents \
    SET ${update_strings.join(',')} \
    WHERE id = $1;`, update_values)

    return res.rowCount > 0
}