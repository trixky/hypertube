export function sanitize_string_to_positive_integer(input: any): number {
    let sanitized_input: number = NaN
    const type_of_input = typeof input

    if (type_of_input != "number") {
        if (typeof input != "string")
            throw new Error("input corrupted");

        sanitized_input = parseInt(input);
    } else {
        sanitized_input = input
    }

    if (isNaN(sanitized_input)) throw new Error("input corrupted");
    if (sanitized_input < 0) throw new Error("input need to be positive");

    return Math.round(sanitized_input);
}