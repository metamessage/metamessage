import {
  encodeFromValue,
  decodeToValue,
  mm,
  ValueType,
} from "../../../mm-ts/src/metamessage";

@mm({ desc: "用户" })
class User {
  @mm({ type: ValueType.Int64, desc: "用户ID", nullable: false })
  id: bigint = 0n;
  @mm({ desc: "昵称" })
  name: string = "";
  @mm({ type: ValueType.Uint8 })
  age: number = 0;
}

const u = new User();
u.id = 666n;
u.name = "abc";
u.age = 20;

// const node = ValueToNode(u);
// console.log('ValueToNode', node);
const wire = encodeFromValue(u);
console.log("wire", wire);
const u2 = decodeToValue(wire, User);

console.log(u2);

function bytesToHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}
