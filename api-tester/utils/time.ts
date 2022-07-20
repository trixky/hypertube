export async function sleep(duration: number): Promise<undefined> {
    return new Promise(resolve => {
        setTimeout(() => {
            resolve(undefined)
        }, duration);
    })
}