export function sanitize_string_to_positive_integer(input: string | number): number {
	let sanitized_input = NaN;

	if (typeof input === 'number') {
		sanitized_input = input;
	} else {
		if (typeof input != 'string') {
			throw new Error('input corrupted');
		}
		sanitized_input = parseInt(input);
	}

	if (isNaN(sanitized_input)) throw new Error('input corrupted');
	if (sanitized_input < 0) throw new Error('input need to be positive');

	return Math.round(sanitized_input);
}
