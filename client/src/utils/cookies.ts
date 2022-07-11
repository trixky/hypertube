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
	token: 'token',
	user_info: 'userInfo'
};

function add_a_cookie(name: string, value: string, days?: number) {
	if (browser) {
		let expires = '';
		if (days) {
			const date = new Date();
			date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
			expires = '; expires=' + date.toUTCString();
		}
		document.cookie = name + '=' + (value || '') + expires + '; path=/';
	}
}

function del_a_cookie(name: string) {
	if (browser) document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/';
}

function get_a_cookie(name: string): string | undefined {
	// https://stackoverflow.com/questions/5639346/what-is-the-shortest-function-for-reading-a-cookie-by-name-in-javascript

	if (browser) {
		const cookie_value = document.cookie.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)')?.pop();
		return cookie_value;
	}
}

function get_token(): string | undefined {
	return get_a_cookie(labels.token);
}

function get_me_from_cookie(): Me | undefined {
	if (browser) {
		const me_64 = get_a_cookie(labels.user_info);

		if (me_64 != undefined) {
			const me_json = atob(me_64);
			const me = JSON.parse(me_json);
			if (me) {
				return <Me>{
					id: me.id,
					username: me.username,
					firstname: me.firstname,
					lastname: me.lastname,
					email: me.email,
					external: me.external
				};
			}
		}
	}
	return undefined;
}

export { labels, add_a_cookie, del_a_cookie, get_a_cookie, get_token, get_me_from_cookie };
