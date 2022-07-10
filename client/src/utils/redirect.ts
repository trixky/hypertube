import * as cookies from './cookies'
import { goto } from '$app/navigation';

function disconnect() {
    cookies.del_a_cookie(cookies.labels.token)
    cookies.del_a_cookie(cookies.labels.user_info)
    goto('/login')
}

function already_connected(browser: boolean): boolean {
    if (browser && cookies.get_token()) {
        goto('/')

        return true
    }

    return false
}

function not_connected(browser: boolean): boolean {
    if (browser && !cookies.get_token()) {
        goto('/login')

        return true
    }

    return false
}

export {
    disconnect,
    already_connected,
    not_connected
}