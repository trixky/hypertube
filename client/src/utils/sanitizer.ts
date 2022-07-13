import { uppercase_first_character } from './str';
import { $t } from './translate';

function name(name: string, can_be_null = false): string {
	if (!name.length) return can_be_null ? '' : $t('sanitizer.missing');

	if (name.length < 3) return $t('sanitizer.too_short', { values: { amount: 3 } });
	else if (name.length > 20) return $t('sanitizer.too_long', { values: { amount: 20 } });

	return '';
}

function email(email: string, can_be_null = false): string {
	// https://stackoverflow.com/questions/46155/how-can-i-validate-an-email-address-in-javascript
	const regex =
		/[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?/;

	if (!email.length) return can_be_null ? '' : $t('sanitizer.missing');
	if (!regex.test(email)) return $t('sanitizer.badly_formatted');

	return '';
}

function password(password: string, _ /* can_be_null */ = false): string {
	const password_warnings = [];

	if (!password.length) return $t('sanitizer.missing');

	if (password.length < 8)
		password_warnings.push($t('sanitizer.too_short', { values: { amount: 8 } }));
	else if (password.length > 30)
		password_warnings.push($t('sanitizer.too_long', { values: { amount: 20 } }));
	if (!/[a-z]/.test(password)) {
		password_warnings.push($t('sanitizer.missing_lowercase'));
	}
	if (!/[A-Z]/.test(password)) {
		password_warnings.push($t('sanitizer.missing_uppercase'));
	}
	if (!/\d/.test(password)) {
		password_warnings.push($t('sanitizer.missing_digit'));
	}
	if (!/[ !@#$%^&*()-=_+[\]{}\\|'";:/?.>,<`~]/.test(password)) {
		password_warnings.push($t('sanitizer.missing_special'));
	}

	let password_warning = password_warnings
		.map((str) => uppercase_first_character(str))
		.join('\n- ');

	if (password_warnings.length > 1) password_warning = '- ' + password_warning;

	return password_warning;
}

function confirm_password(
	password: string,
	confirm_password: string,
	_ /* can_be_null */ = false
): string {
	if (!confirm_password.length) return $t('sanitizer.missing');

	if (confirm_password != password) return $t('sanitizer.passwords_dont_match');

	return '';
}

export { name, email, password, confirm_password };
