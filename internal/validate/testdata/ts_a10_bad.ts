//ff:func feature=validate type=rule control=selection
//ff:what selection declared but no switch
export function tsA10Bad(x: number): string {
    if (x > 0) {
        return "positive";
    }
    return "non-positive";
}
