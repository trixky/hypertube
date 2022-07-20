/// <reference types="@sveltejs/kit" />

interface User {
	id: number;
	username: string;
	firstname: string;
	lastname: string;
	email: string;
	external: string;
}

// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare namespace App {
	interface Locals {
		locale: string;
		token?: string;
		user?: User;
	}
	// interface Platform {}
	interface Session {
		locale?: string;
		token?: string;
		user?: User;
	}
	// interface Stuff {}
}
declare namespace NodeJS {
	type Timeout = number;
	type Timer = number;
}

declare namespace svelte.JSX {
	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	interface HTMLProps<T> {
		onclickOutside?: () => void;
	}
}
