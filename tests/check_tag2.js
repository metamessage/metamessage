#!/usr/bin/env node
// Compare Go and TS tag encoding for child_tags scenario
const { ValueType } = require('../mm-ts/dist/metamessage');
const path = require('path');

// We need access to Tag class - let's check what's exported
const mm = require('../mm-ts/dist/metamessage');
console.log('Exports:', Object.keys(mm));

// Check Tag class via internal path
const { Tag } = require('../mm-ts/dist/ir/tag');
console.log('Tag available:', typeof Tag);

// Simulate the "names" array tag: child_desc=item name
const namesTag = new Tag();
namesTag.type = ValueType.Vec;
namesTag.childDesc = 'item name';
console.log('\nnames array tag:');
console.log('  toBytes():', Buffer.from(namesTag.toBytes()).toString('hex'));

// Simulate the "i8_arr" array tag: child_type=i8
const i8ArrTag = new Tag();
i8ArrTag.type = ValueType.Vec;
i8ArrTag.childType = ValueType.I8;
console.log('\ni8_arr tag:');
console.log('  toBytes():', Buffer.from(i8ArrTag.toBytes()).toString('hex'));

// Simulate the "enum_arr" array tag: child_enums=red|green|blue
const enumArrTag = new Tag();
enumArrTag.type = ValueType.Vec;
enumArrTag.childEnums = 'red|green|blue';
console.log('\nenum_arr tag:');
console.log('  toBytes():', Buffer.from(enumArrTag.toBytes()).toString('hex'));

// Simulate child element tag after inherit
console.log('\n--- Element tags after inherit ---');

// For "names" array element "alice": child_desc inherited
const parentTag1 = new Tag();
parentTag1.childDesc = 'item name';
const child1 = new Tag();
child1.type = ValueType.Str;
child1.inherit(parentTag1);
console.log('names element tag (after inherit):');
console.log('  type:', child1.type, 'isInherit:', child1.isInherit, 'desc:', child1.desc);
console.log('  toBytes():', Buffer.from(child1.toBytes()).toString('hex'));

// For "i8_arr" element 1: child_type inherited
const parentTag2 = new Tag();
parentTag2.childType = ValueType.I8;
const child2 = new Tag();
child2.type = ValueType.I;
child2.inherit(parentTag2);
console.log('\ni8_arr element tag (after inherit):');
console.log('  type:', child2.type, 'isInherit:', child2.isInherit);
console.log('  toBytes():', Buffer.from(child2.toBytes()).toString('hex'));

// For "scores" element 80: child_min=1; child_max=100; child_nullable
const parentTag3 = new Tag();
parentTag3.childMin = '1';
parentTag3.childMax = '100';
parentTag3.childNullable = true;
const child3 = new Tag();
child3.type = ValueType.I;
child3.inherit(parentTag3);
console.log('\nscores element tag (after inherit):');
console.log('  type:', child3.type, 'isInherit:', child3.isInherit, 'min:', child3.min, 'max:', child3.max, 'nullable:', child3.nullable);
console.log('  toBytes():', Buffer.from(child3.toBytes()).toString('hex'));