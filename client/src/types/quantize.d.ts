declare module 'quantize' {
	type RGB = [number, number, number];
	type RGBArray = RGB[];

	export default function quantize(
		image: RGBArray,
		colorLimit?: number,
		quality?: number
	): {
		palette: () => RGBArray;
		size: () => number;
	};
}
