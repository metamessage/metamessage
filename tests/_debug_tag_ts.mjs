import { Tag, ValueType } from '../mm-ts/dist/ir/tag.js';

// Simulate the tag for the "names" array field
const tag = new Tag();
tag.childDesc = 'item name';
tag.type = ValueType.Vec;

console.log('tag with child_desc="item name"');
const tagBytes = tag.toBytes();
console.log('tag.toBytes() hex:', Buffer.from(tagBytes).toString('hex'));
console.log('tag.toBytes() length:', tagBytes.length);
console.log('tag.toString():', tag.toString());