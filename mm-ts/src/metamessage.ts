import { fromValue, fromJSONC, decode } from './mm';
import { Tag } from './ast/tag';
import { Binder, Constructor } from './mm/binder.js';

export { ValueType } from './ast/value-type';
export { mm } from './mm/mm';

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
