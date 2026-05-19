import { fromValue, fromJSONC, decode, ValueToNode } from './core';
import { Tag } from './ir/tag';
import { Binder, Constructor } from './core/binder.js';
import { parseJSONC, toJSONC } from './jsonc';

export { ValueType } from './ir/value-type';
export { mm } from './core/mm';

export function encodeFromValue(value: any, tag?: Tag): Uint8Array {
  return fromValue(value, tag);
}

export function encodeFromJsonc(jsonc: string): Uint8Array {
  return fromJSONC(jsonc);
}

export function decodeToValue<T>(
  wire: Uint8Array,
  type: Constructor<T> | T,
): any {
  const node = decode(wire);
  return new Binder().bind(node, type);
}

export function decodeToJsonc(wire: Uint8Array): string {
  const node = decode(wire);
  return toJSONC(node);
}

export function valueToJsonc(value: any, tag?: Tag): string {
  const node = ValueToNode(value, tag);
  return toJSONC(node);
}

export function jsoncToValue<T>(jsonc: string, type: Constructor<T> | T): T {
  const node = parseJSONC(jsonc);
  return new Binder().bind(node, type);
}
