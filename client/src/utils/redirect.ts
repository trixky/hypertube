import * as cookies from './cookies'

function already_connected(browser: boolean) {
    if (browser && cookies.get_a_cookie(cookies.labels.token)) {
		window.location.href = window.location.origin + '/';
	}
}

export {
    already_connected
}