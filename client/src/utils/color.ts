import quantize from 'quantize';

export function extractPalette(image: HTMLImageElement) {
	const context = document.createElement('canvas').getContext('2d');
	if (!context) {
		return undefined;
	}
	context.imageSmoothingEnabled = true;
	context.drawImage(image, 0, 0, image.width, image.height);

	// Extract pixels RGB channels as an array of pixel data
	const pixels = context.getImageData(0, 0, image.width, image.height).data;
	const imageData: [number, number, number][] = [];
	for (let i = 0; i < pixels.length; i += 4) {
		const rgb: [number, number, number] = [pixels[i], pixels[i + 1], pixels[i + 2]];
		const a = pixels[i + 3];
		// If pixel is mostly opaque and not white
		if (typeof a === 'undefined' || a >= 125) {
			if (
				!(rgb[0] > 250 && rgb[1] > 250 && rgb[2] > 250) &&
				!(rgb[0] < 30 && rgb[1] < 30 && rgb[2] < 30)
			) {
				imageData.push(rgb);
			}
		}
	}

	// Extract a color palette
	const rawPalette: [number, number, number][] = quantize(imageData, 5, 10).palette();
	const palette = rawPalette.map((color) => {
		return `rgb(${color[0]}, ${color[1]}, ${color[2]})`;
	});

	// Clamp each channels to 150 to avoid bright colors
	let color = [rawPalette[0][0], rawPalette[0][1], rawPalette[0][2]];
	const difference = color.reduce((carry, value) => {
		if (value > 70) {
			return Math.max(carry, value - 70);
		}
		return carry;
	}, 0);
	color = color.map((color) => Math.max(0, color - difference)) as [number, number, number];
	return { palette, color: `rgb(${color[0]}, ${color[1]}, ${color[2]})` };
}
