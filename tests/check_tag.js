#!/usr/bin/env node
const { Tag, ValueType } = require('../mm-ts/dist/metamessage');

// Simulate what happens for "names" array element "alice"
// The array has child_desc=item name
const parentTag = new Tag();
parentTag.childDesc = 'item name';
parentTag.childType = ValueType.Str;

// After inherit, the child tag should have desc="item name", type=Str, isInherit=true
const childTag = new Tag();
childTag.type = ValueType.Str; // Set by parser since it's a string value
console.log('Before inherit:');
console.log('  type:', childTag.type);
console.log('  desc:', childTag.desc);
console.log('  isInherit:', childTag.isInherit);
console.log('  toBytes():', Buffer.from(childTag.toBytes()).toString('hex'));

childTag.inherit(parentTag);
console.log('\nAfter inherit:');
console.log('  type:', childTag.type);
console.log('  desc:', childTag.desc);
console.log('  isInherit:', childTag.isInherit);
console.log('  toBytes():', Buffer.from(childTag.toBytes()).toString('hex'));
console.log('  toString():', childTag.toString());

// Also check for i8_arr element
console.log('\n--- i8_arr element ---');
const parentTag2 = new Tag();
parentTag2.childType = ValueType.I8;

const childTag2 = new Tag();
childTag2.type = ValueType.I;
console.log('Before inherit:');
console.log('  toBytes():', Buffer.from(childTag2.toBytes()).toString('hex'));
childTag2.inherit(parentTag2);
console.log('After inherit:');
console.log('  type:', childTag2.type);
console.log('  isInherit:', childTag2.isInherit);
console.log('  toBytes():', Buffer.from(childTag2.toBytes()).toString('hex'));

// Check enum_arr element
console.log('\n--- enum_arr element ---');
const parentTag3 = new Tag();
parentTag3.childEnums = 'red|green|blue';

const childTag3 = new Tag();
childTag3.type = ValueType.Str;
childTag3.inherit(parentTag3);
console.log('After inherit:');
console.log('  type:', childTag3.type);
console.log('  enums:', childTag3.enums);
console.log('  isInherit:', childTag3.isInherit);
console.log('  toBytes():', Buffer.from(childTag3.toBytes()).toString('hex'));