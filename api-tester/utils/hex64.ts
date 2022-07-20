import atob from 'atob'

export interface UserInfo {
    id: number,
    username: string,
    firstname: string,
    lastname: string,
    email: string,
    external: string | null,
}

export function extract_hex64(hex64: string): UserInfo | undefined {
    const user = atob(hex64);
    const me = JSON.parse(user);
    if (me) {
        return <UserInfo>{
            id: me.id,
            username: me.username,
            firstname: me.firstname,
            lastname: me.lastname,
            email: me.email,
            external: me.external
        };
    }

    return undefined
}