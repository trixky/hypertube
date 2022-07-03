const labels = {
    token: "token"
}

function add_a_cookie(name: string, value: string) {
    document.cookie = name + '=' + value + ';'
}

function del_a_cookie(name: string) {
    document.cookie = name + '=0 ; max-age=0;'
}

function get_a_cookie(name: string): string | undefined {
    // https://stackoverflow.com/questions/5639346/what-is-the-shortest-function-for-reading-a-cookie-by-name-in-javascript
    
    return document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)')?.pop()
}

export {
    labels,
    add_a_cookie,
    del_a_cookie,
    get_a_cookie,
}