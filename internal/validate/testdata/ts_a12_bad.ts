//ff:func feature=validate type=rule control=sequence
//ff:what sequence declared but has for loop
export function tsA12Bad(items: string[]): void {
    for (const item of items) {
        console.log(item);
    }
}
