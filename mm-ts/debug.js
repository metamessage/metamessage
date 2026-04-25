import { encode, decode, mm } from './dist/mm/index.js';

// 调试对象编码
const obj = { name: 'test', age: 25, active: true };
console.log('Original object:', obj);

const encoded = encode(obj);
console.log('Encoded bytes:', Array.from(encoded));

const decoded = decode(encoded);
console.log('Decoded object:', decoded);
console.log('Decoded age:', decoded.value.age);
console.log('Decoded age type:', typeof decoded.value.age);
console.log('Decoded age toString():', decoded.value.age.toString());

// 调试整数编码
const num = 25;
const encodedNum = encode(num);
console.log('\nEncoded 25:', Array.from(encodedNum));
const decodedNum = decode(encodedNum);
console.log('Decoded 25:', decodedNum);

// 调试字符串编码
const str = 'test';
const encodedStr = encode(str);
console.log('\nEncoded "test":', Array.from(encodedStr));
const decodedStr = decode(encodedStr);
console.log('Decoded "test":', decodedStr);

// 调试布尔编码
const bool = true;
const encodedBool = encode(bool);
console.log('\nEncoded true:', Array.from(encodedBool));
const decodedBool = decode(encodedBool);
console.log('Decoded true:', decodedBool);