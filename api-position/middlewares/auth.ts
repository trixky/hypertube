import express from "express";
import { client } from "../redis/db";

async function get_user_id_from_token(token: string): Promise<number | null> {
    const keys: Array<string> = await client.sendCommand(["KEYS", "token." + token + ".*"])
    if (keys.length == 1) {
        return parseInt(keys[0].split('.').pop()!);
    }

    return null
}

export default async function auth_middleware(req: express.Request, res: express.Response, next: express.NextFunction) {
    if (req.cookies.token != undefined) {
        const user_id = await get_user_id_from_token(req.cookies.token)
        
        if (user_id != null) {
            res.locals.user_id = user_id
            return next()
        }
    }
    
    res.status(401)
    res.send()
}