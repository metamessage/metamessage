#!/usr/bin/env node
const { encode, decode } = require('../../../mm-js/src/mm/index.js');

// 创建对象
const person = {
  name: "Ed",
  age: 30
};

console.log('Original:', person);

// 编码到 Wire 格式
const wire = encode(person);
console.log('Encoded:', bytesToHex(wire));

// 从 Wire 解码
const decoded = decode(wire);
console.log('Decoded:', decoded.value);

function bytesToHex(bytes) {
  return Array.from(bytes)
    .map(b => b.toString(16).padStart(2, '0'))
    .join('');
}
