import { browser } from '$app/env';

export const labels = {
	password_token: 'token',
	token: 'token',
	user_info: 'userInfo'
};

export function add_a_cookie(name: string, value: string, days?: number) {
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

export function del_a_cookie(name: string) {
	if (browser) document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT; path=/';
}

// @source https://stackoverflow.com/questions/5639346/what-is-the-shortest-function-for-reading-a-cookie-by-name-in-javascript
export function extract_cookie(cookies: string, name: string): string | undefined {
	const cookie_value = cookies.match('(^|;)\\s*' + name + '\\s*=\\s*([^;]+)')?.pop();
	return cookie_value;
}

export function get_a_cookie(name: string): string | undefined {
	if (browser) {
		return extract_cookie(document.cookie, name);
	}
}

export function get_token(): string | undefined {
	return get_a_cookie(labels.token);
}

export function get_user(cookies: string) {
	const encodedUser = extract_cookie(cookies, labels.user_info);
	if (encodedUser != undefined) {
		const user = atob(encodedUser);
		const me = JSON.parse(user);
		if (me) {
			return <User>{
				id: me.id,
				username: me.username,
				firstname: me.firstname,
				lastname: me.lastname,
				email: me.email,
				external: me.external
			};
		}
	}
	return undefined;
}

export function get_me_from_cookie(): User | undefined {
	if (browser) {
		return get_user(document.cookie);
	}
	return undefined;
}
