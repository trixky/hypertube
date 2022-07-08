import { browser } from '$app/env';

interface Me {
    id: number;
    username: string;
    firstname: string;
    lastname: string;
    email: string;
    external: string;
}

const labels = {
    token: "token",
    me: "me",
}

function add_a_cookie(name: string, value: string) {
    if (browser) document.cookie = name + '=' + value + '; path=/'
}

function del_a_cookie(name: string) {
    if (browser) document.cookie = name + '=0 ; expires = Thu, 01 Jan 1970 00:00:00 GMT; path=/'
}


function get_a_cookie(name: string): string | undefined {
    // https://stackoverflow.com/questions/5639346/what-is-the-shortest-function-for-reading-a-cookie-by-name-in-javascript


    if (browser) {
        let cookie_value = document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)')?.pop()
        return cookie_value
    }
}

function get_token(): string | undefined {
    return get_a_cookie(labels.token)
}

function get_me(): Me | undefined {
    let me_64 = get_a_cookie(labels.me)

    if (me_64) {
        let me_json = atob(me_64)
        let me = JSON.parse(me_json)
        if (me) {
            let ret = <Me>{
                id: me.id,
                username: me.username,
                firstname: me.firstname,
                lastname: me.lastname,
                email: me.email,
                external: me.external,
            }
            return ret
        }
    }
    return undefined
}

export {
    labels,
    add_a_cookie,
    del_a_cookie,
    get_a_cookie,
    get_token,
    get_me,
}