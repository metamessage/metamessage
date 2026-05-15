import { MMEncoder } from './encoder';
import { MMDecoder } from './decoder';
import { MMBuffer, MMError } from './buffer';
import { mm } from './mm';
import * as constants from './constants';
import { Tag } from '../ir/tag';
import { ValueType } from '../ir/value-type';
import { Node, MMValue, MMObject, MMArray } from '../ir/ast';
import { ValueToNode } from './value-to-node';

import { parseJSONC, toJSONC } from '../jsonc/index';

export function encodeNode(node: Node): Uint8Array {
  const encoder = new MMEncoder();
  return encoder.encodeNode(node);
}

export function encodeValue(value: any, tag?: Tag): Uint8Array {
  const encoder = new MMEncoder();
  return encoder.encodeValue(value, tag);
}

export const fromValue = encodeValue;

export function fromJSONC(jsonc: string): Uint8Array {
  const node = parseJSONC(jsonc);
  const encoder = new MMEncoder();
  return encoder.encodeNode(node);
}

function nodeToDecodedValue(node: Node): DecodedValue {
  if (node instanceof MMValue) {
    return {
      type: node.getTag().type,
      value: node.getValue(),
    };
  } else if (node instanceof MMObject) {
    const props = node.getProperties();
    const value: Record<string, any> = {};
    for (const [key, valNode] of Object.entries(props)) {
      value[key] = nodeToDecodedValue(valNode).value;
    }
    return { type: ValueType.Object, value };
  } else if (node instanceof MMArray) {
    const elements = node.getElements();
    const value = elements.map((e) => nodeToDecodedValue(e).value);
    return { type: ValueType.Array, value };
  }
  return { type: ValueType.Unknown, value: null };
}

export interface DecodedValue {
  type: ValueType;
  value: any;
}

export function decode(data: Uint8Array): Node {
  const decoder = new MMDecoder(data);
  return decoder.decode();
}

export function decodeToValue(data: Uint8Array): DecodedValue {
  const node = decode(data);
  return nodeToDecodedValue(node);
}

export function decodeToJSONC(data: Uint8Array): string {
  const node = decode(data);
  return toJSONC(node);
}

export { mm } from './mm';
export { MMEncoder, MMDecoder, MMBuffer, MMError };
export * as constants from './constants';
export { ValueToNode } from './value-to-node';
export default {
  encodeNode,
  decode,
  decodeToJSONC,
  fromJSONC,
  mm,
  MMEncoder,
  MMDecoder,
  MMBuffer,
  MMError,
  constants,
  ValueToNode,
};
