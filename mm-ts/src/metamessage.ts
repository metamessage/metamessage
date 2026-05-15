import { fromValue, fromJSONC, decode } from './core';
import { Tag } from './ir/tag';
import { Binder, Constructor } from './core/binder.js';

export { ValueType } from './ir/value-type';
export { mm } from './core/mm';

export function encodeFromValue(value: any, tag?: Tag): Uint8Array {
  return fromValue(value, tag);
}

export function encodeFromJSONC(jsonc: string): Uint8Array {
  return fromJSONC(jsonc);
}

export function decodeToValue<T>(
  wire: Uint8Array,
  type: Constructor<T> | T,
): any {
  const node = decode(wire);
  return new Binder().bind(node, type);
}

export function decodeToJSONC(jsonc: string): Uint8Array {
  return decodeToJSONC(jsonc);
}
