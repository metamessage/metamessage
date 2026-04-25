import { encode, decode, mm } from './src/mm/index.ts';

// 调试对象编码
const obj = { name: 'test', age: 25, active: true };
console.log('Original object:', obj);

const encoded = encode(obj);
console.log('Encoded bytes:', Array.from(encoded));

console.log('\nDecoding process:');
const decoded = decode(encoded);
console.log('Decoded object:', decoded);
console.log('Decoded age:', decoded.value.age);
console.log('Decoded age type:', typeof decoded.value.age);
console.log('Decoded age toString():', decoded.value.age.toString());

// 调试单个字段编码
console.log('\nDebugging individual fields:');

// 调试 name 字段
const nameEncoded = encode('name');
console.log('Encoded "name":', Array.from(nameEncoded));
const nameDecoded = decode(nameEncoded);
console.log('Decoded "name":', nameDecoded);

// 调试 test 字段
const testEncoded = encode('test');
console.log('Encoded "test":', Array.from(testEncoded));
const testDecoded = decode(testEncoded);
console.log('Decoded "test":', testDecoded);

// 调试 age 字段
const ageEncoded = encode(25);
console.log('Encoded 25:', Array.from(ageEncoded));
const ageDecoded = decode(ageEncoded);
console.log('Decoded 25:', ageDecoded);

// 调试 active 字段
const activeEncoded = encode(true);
console.log('Encoded true:', Array.from(activeEncoded));
const activeDecoded = decode(activeEncoded);
console.log('Decoded true:', activeDecoded);