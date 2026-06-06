const { Tag } = require('../mm-ts/dist/ir/tag.js');
const { ValueType } = require('../mm-ts/dist/ir/value-type.js');

// Simulate the tag for the "names" array field
const tag = new Tag();
tag.childDesc = 'item name';
tag.type = ValueType.Vec;

console.log('tag with child_desc="item name"');
const tagBytes = tag.toBytes();
console.log('tag.toBytes() hex:', Buffer.from(tagBytes).toString('hex'));
console.log('tag.toBytes() length:', tagBytes.length);
console.log('tag.toString():', tag.toString());

// Also test with child_enums
const tag2 = new Tag();
tag2.childEnums = 'red|green|blue';
tag2.childType = ValueType.Enums;
tag2.type = ValueType.Vec;

console.log('\ntag with child_enums="red|green|blue"');
const tagBytes2 = tag2.toBytes();
console.log('tag.toBytes() hex:', Buffer.from(tagBytes2).toString('hex'));
console.log('tag.toBytes() length:', tagBytes2.length);
console.log('tag.toString():', tag2.toString());