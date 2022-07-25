export const Origin = import.meta.env.VITE_DOMAIN ?? 'localhost';

export function cleanRoute(route: string) {
	if (route.startsWith('/')) {
		return route.slice(1);
	}
	return route;
}

export const ApiMediaPort = import.meta.env.VITE_API_MEDIA_PORT ?? '7072';
export function apiMedia(route: string) {
	return `http://${Origin}:${ApiMediaPort}/${cleanRoute(route)}`;
}

export const TmdbProxyPort = import.meta.env.VITE_TMDB_PROXY_PORT ?? '7260';
export function imageProxy(route: string): string {
	return `http://${Origin}:${TmdbProxyPort}/${cleanRoute(route)}`;
}
export const ApiStreamingPort = import.meta.env.VITE_API_STREAMING_PORT ?? '3030';
export function apiStreaming(route: string): string {
	return `http://${Origin}:${ApiStreamingPort}/${cleanRoute(route)}`;
}

export const ApiUserPort = import.meta.env.VITE_API_USER_PORT ?? '7170';
export function apiUser(route: string): string {
	return `http://${Origin}:${ApiUserPort}/${cleanRoute(route)}`;
}

export const ApiAuthPort = import.meta.env.VITE_TMDB_PROXY_PORT ?? '7070';
export function apiAuth(route: string): string {
	return `http://${Origin}:${ApiAuthPort}/${cleanRoute(route)}`;
}

export const ApiPositionPort = import.meta.env.VITE_API_POSITION_PORT ?? '3040';
export function apiPosition(route: string): string {
	return `http://${Origin}:${ApiPositionPort}/${cleanRoute(route)}`;
}

console.log(
	Origin,
	ApiMediaPort,
	TmdbProxyPort,
	ApiStreamingPort,
	ApiUserPort,
	ApiAuthPort,
	ApiPositionPort
);
