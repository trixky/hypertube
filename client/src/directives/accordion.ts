// @source https://svelte.dev/repl/7729845536404efcaf1f6c65328df3f2?version=3.38.2
export function accordion(node: HTMLElement, isOpen: boolean) {
	const initialHeight = node.offsetHeight;
	if (isOpen) {
		node.style.height = 'auto';
		node.style.overflow = 'hidden';
	} else {
		node.style.height = '0';
		node.style.overflow = 'hidden';
	}
	return {
		update(isOpen: boolean) {
			const animation = node.animate(
				[
					{
						height: `${initialHeight}px`,
						overflow: 'hidden'
					},
					{
						height: 0,
						overflow: 'hidden'
					}
				],
				{ duration: 250, fill: 'both' }
			);
			animation.pause();
			if (!isOpen) {
				animation.play();
			} else {
				animation.reverse();
			}
		}
	};
}
