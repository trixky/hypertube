import { _ } from 'svelte-i18n';

// * > svelte-i18n types
type FormatXMLElementFn<T, R = string | T | Array<string | T>> = (parts: Array<string | T>) => R;
type InterpolationValues =
	| Record<
			string,
			string | number | boolean | Date | FormatXMLElementFn<unknown> | null | undefined
	  >
	| undefined;
interface MessageObject {
	id: string;
	locale?: string;
	format?: string;
	default?: string;
	values?: InterpolationValues;
}
type MessageFormatter = (id: string | MessageObject, options?: Omit<MessageObject, 'id'>) => string;
// * <

let translator: MessageFormatter;

_.subscribe((formater) => (translator = formater));

export function $t(id: string | MessageObject, options?: Omit<MessageObject, 'id'>) {
	return translator(id, options);
}
