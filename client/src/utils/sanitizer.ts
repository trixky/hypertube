import { uppercase_first_character } from './str';


function name(name: string, can_be_null: boolean = false): string {
    if (!name.length) return can_be_null ? '' : 'is missing';

    if (name.length < 3) return 'is too short, needs at least 3 characters';
    else if (name.length > 20)
        return 'is too long, must contain a maximum of 20 characters';

    return '';
}

function email(email: string, can_be_null: boolean = false): string {
    // https://stackoverflow.com/questions/46155/how-can-i-validate-an-email-address-in-javascript
    const regex =
        /[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?/;

    if (!email.length) return can_be_null ? '' : 'is missing';
    if (!regex.test(email)) return 'is bad formatted';

    return '';
}

function password(password: string, can_be_null: boolean = false): string {
    let password_warnings = [];

    if (!password.length) return 'is missing';

    if (password.length < 8)
        password_warnings.push('is too short, needs at least 8 characters');
    else if (password.length > 30)
        password_warnings.push('is too long, must contain a maximum of 20 characters');
    if (!/[a-z]/.test(password)) {
        password_warnings.push('must contain at least one lowercase character (a-z)');
    }
    if (!/[A-Z]/.test(password)) {
        password_warnings.push('must contain at least one uppercase character (A-Z)');
    }
    if (!/\d/.test(password)) {
        password_warnings.push('must contain at least one numeric character (0-9)');
    }
    if (!/[ !@#$%^&*()-=_+[\]{}\\|'\";:/?.>,<`~]/.test(password)) {
        password_warnings.push('must contain at least one specific character (!@#...)');
    }

    let password_warning = password_warnings
        .map((str) => uppercase_first_character(str))
        .join('\n- ');

    if (password_warnings.length > 1)
        password_warning = '- ' + password_warning;

    return password_warning
}

function confirm_password(password: string, confirm_password: string, can_be_null: boolean = false): string {
    if (!confirm_password.length) return 'is missing';

    if (confirm_password != password)
        return 'passwords must be the same';

    return '';
}

export {
    name,
    email,
    password,
    confirm_password,
}