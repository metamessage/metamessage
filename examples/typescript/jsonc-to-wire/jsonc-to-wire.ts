import { fromJSONC, toJSONC, encode, decode } from '../../../mm-ts/src/mm/index.js';
import { parseFromString, toString } from '../../../mm-ts/src/jsonc/index.js';

// JSONC 字符串
const jsonc = `{
    // mm: type=datetime; desc=创建时间
    "create_time": "2026-01-01 00:00:00",
    // mm: type=str; desc=用户名称
    "user_name": "Alice",
    // mm: type=bool; desc=是否激活
    "is_active": true,
    // mm: type=array; child_type=i
    "scores": [95, 87, 92]
}`;

console.log('Input JSONC:');
console.log(jsonc);

// 解析 JSONC
const parsed = parseFromString(jsonc);
console.log('\nParsed:');
console.log(toString(parsed));

// 编码到 Wire 格式
const wire = encode(parsed);
console.log('\nEncoded Wire:');
console.log(bytesToHex(wire));

// 解码回 JSONC
const jsoncOut = toJSONC(wire);
console.log('\nDecoded to JSONC:');
console.log(jsoncOut);

function bytesToHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');
}
