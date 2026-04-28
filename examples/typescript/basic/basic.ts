import { encode, bind, mm } from '../../../mm-ts/src/mm/index.js';

// 使用 mm() 函数包装值，支持 tag 和自动类型推断
const person = {
  name: mm("Ed", { desc: "" }),
  email: mm("Ed@gmail.com", { desc: "", type: "email" }),
  score: mm(90, { desc: "", type: "uint8" }),
  age: mm(30, { desc: "" })
};

console.log('Original:', person);

// 编码到 Wire 格式
const wire = encode(person);
console.log('Encoded:', bytesToHex(wire));

// 从 Wire 解码并绑定到模板
const decoded = bind(wire, person);
console.log('Decoded:', decoded);

function bytesToHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');
}
