import { writable } from 'svelte/store';
import { get_me } from '../utils/cookies';

interface Me {
    id: number
    username: string
    firstname: string
    lastname: string
    email: string
    external: string
}

function create_me() {
    const { subscribe, set, update } = writable(<Me>{
        id: 0,
        username: '',
        firstname: '',
        lastname: '',
        email: '',
        external: '',
    });

    return {
        subscribe,
        refresh_from_cookies: () => {
            let me_from_cookie = get_me()

            if (me_from_cookie != undefined) {
                set(<Me>{
                    id: me_from_cookie.id,
                    username: me_from_cookie.username,
                    firstname: me_from_cookie.firstname,
                    lastname: me_from_cookie.lastname,
                    email: me_from_cookie.email,
                    external: me_from_cookie.external,
                })
            }
        },
        reset: () => set(<Me>{})
    };
}

export const me_store = create_me();