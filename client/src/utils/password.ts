import 'crypto'

// https://stackoverflow.com/questions/18338890/are-there-any-sha-256-javascript-implementations-that-are-generally-considered-t/48161723#48161723

async function encrypt_password(password: string): Promise<string> {
    if (!password.length) return ''
        
    const customized_password = "hYpErTuBe." + password + ".hYpErTuBe"
    const password_buffer = new TextEncoder().encode(customized_password);
    const hash_bytes = await crypto.subtle.digest('SHA-256', password_buffer);
    const hash_array = Array.from(new Uint8Array(hash_bytes));
    const hash = hash_array.map(b => b.toString(16).padStart(2, '0')).join('');

    return hash
}

export {
    encrypt_password
}