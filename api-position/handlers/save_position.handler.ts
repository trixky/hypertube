import express from "express";
import * as models from '../models/positions'
import { save_position } from "../postgres/positions"
import { sanitize_string_to_positive_integer } from "../utils/sanitize"

export default async function save_position_handler(req: express.Request, res: express.Response) {
    const torrend_id: string = req.params.torrent_id

    let sanitized_torrent_id: number
    try {
        sanitized_torrent_id = sanitize_string_to_positive_integer(torrend_id)
    } catch {
        return res.status(400).send("corrupted torrent id")
    }

    let sanitized_position: number
    try {
        sanitized_position = sanitize_string_to_positive_integer(req.body.position)
    } catch {
        return res.status(400).send("corrupted position")
    }

    try {
        await save_position(<models.Position>{
            user_id: res.locals.user_id,
            torrent_id: sanitized_torrent_id,
            position: sanitized_position,
        })
    } catch (err: any) {
        return res.status(500).send("internal sever error")
    }

    return res.send()
}