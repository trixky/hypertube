import { writable } from 'svelte/store';

function random_image_option(): string {
    return Math.random().toString() + Math.random().toString();
}

function createProfilePicture(): any {
    const { subscribe, set } = writable(random_image_option());

    return {
        subscribe,
        refresh: () => {
            setTimeout(() => {
                set(random_image_option())
            }, 500);
        }
    };
}

export const profilePicture = createProfilePicture();