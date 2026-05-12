import { something } from "module";

//ff:func feature=validate type=rule control=sequence
//ff:what annotation after import — A6 violation
export function tsA6Bad(): void {
    console.log(something);
}
