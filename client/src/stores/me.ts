import { writable } from 'svelte/store';
import { get_me_from_cookie } from '$utils/cookies';

function create_me() {
	const { subscribe, set, update } = writable(<User>{
		id: 0,
		username: '',
		firstname: '',
		lastname: '',
		email: '',
		external: ''
	});

	return {
		subscribe,
		refresh_from_cookies: () => {
			const me_from_cookie = get_me_from_cookie();
			if (me_from_cookie != undefined) {
				set(<User>{
					id: me_from_cookie.id,
					username: me_from_cookie.username,
					firstname: me_from_cookie.firstname,
					lastname: me_from_cookie.lastname,
					email: me_from_cookie.email,
					external: me_from_cookie.external
				});
			}
		},
		reset: () => set(<User>{})
	};
}

export const me_store = create_me();
