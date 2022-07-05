import * as cookies from './cookies'

function disconnect() {
    cookies.del_a_cookie('token')

    window.location.href = window.location.origin + '/login';
}

function already_connected(browser: boolean): boolean {
    if (browser && cookies.get_a_cookie(cookies.labels.token)) {
		window.location.href = window.location.origin + '/';
        
        return true
	}

    return false
}

function not_connected(browser: boolean): boolean {
    if (browser && !cookies.get_a_cookie(cookies.labels.token)) {
		window.location.href = window.location.origin + '/login';

        return true
	}

    return false
}

export {
    disconnect,
    already_connected,
    not_connected
}