function uppercase_first_character(str: string): string {
	return str.charAt(0).toUpperCase() + str.slice(1);
}

export { uppercase_first_character };
