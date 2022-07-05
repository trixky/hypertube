/** Dispatch event on click outside of node */
// @source https://svelte.dev/repl/0ace7a508bd843b798ae599940a91783?version=3.16.7
export function clickOutside(node: HTMLElement) {
	const handleClick = (event: Event) => {
		if (
			node &&
			(!event.target || !node.contains(event.target as HTMLElement)) &&
			!event.defaultPrevented
		) {
			node.dispatchEvent(new CustomEvent<HTMLElement>('clickOutside', { detail: node }));
		}
	};

	document.addEventListener('click', handleClick, true);

	return {
		destroy() {
			document.removeEventListener('click', handleClick, true);
		}
	};
}
