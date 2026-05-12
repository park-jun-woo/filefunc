//ff:func feature=validate type=rule control=sequence
//ff:what deeply nested function for Q1 test
export function tsDeepNesting(): number {
    if (true) {
        if (true) {
            if (true) {
                return 1;
            }
        }
    }
    return 0;
}
