export function imageUrl(path: string): string {
	// TODO replace by env value
	if (path.startsWith('/')) {
		path = path.slice(1);
	}
	return `http://localhost:7260/${path}`;
}
