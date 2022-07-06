/// <reference types="@sveltejs/kit" />

// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare namespace App {
	// interface Locals {}
	// interface Platform {}
	// interface Session {}
	// interface Stuff {}
}
declare namespace NodeJS {
	type Timeout = number;
	type Timer = number;
}

declare namespace svelte.JSX {
	interface HTMLProps<T> {
		onclickOutside?: () => void;
	}
}
