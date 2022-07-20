import express from "express";
import { sanitize_string_to_positive_integer } from "../utils/sanitize"
import { get_position } from "../postgres/positions"

export default async function get_position_handler(req: express.Request, res: express.Response) {
    const torrend_id: string = req.params.torrent_id

    let sanitized_torrent_id: number
    try {
        sanitized_torrent_id = sanitize_string_to_positive_integer(torrend_id)
    } catch {
        return res.status(400).send("corrupted torrent id")
    }
    
    const position = await get_position(res.locals.user_id, sanitized_torrent_id)

    res.json({position: position.position})
}