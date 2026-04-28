import { MMEncoder } from './encoder';
import { MMDecoder, DecodedValue } from './decoder';
import { MMBuffer, MMError } from './buffer';
import { mm, MMValue, Tag, ValueType } from './types';
import * as constants from './constants';
import { MmValidator, validator, ValidationResult } from './validator';

export function encode(value: any): Uint8Array {
  const encoder = new MMEncoder();
  encoder.encode(value);
  return encoder.result;
}

export function decode(data: Uint8Array | ArrayBuffer | number[]): DecodedValue {
  const decoder = new MMDecoder(data);
  return decoder.decode();
}

export function toJSONC(data: Uint8Array | ArrayBuffer | number[]): string {
  const decoded = decode(data);
  return jsoncFromDecoded(decoded);
}

function jsoncFromDecoded(decoded: { type: string; value: any }): string {
  switch (decoded.type) {
    case 'null':
      return 'null';
    case 'bool':
      return decoded.value ? 'true' : 'false';
    case 'int':
      return decoded.value.toString();
    case 'float':
      return decoded.value.toString();
    case 'string':
      return JSON.stringify(decoded.value);
    case 'bytes':
      return JSON.stringify(Array.from(decoded.value));
    case 'array':
      return '[' + decoded.value.map((v: any) => {
        if (typeof v === 'object' && v !== null && 'type' in v && 'value' in v) {
          return jsoncFromDecoded(v as { type: string; value: any });
        }
        return jsoncFromDecoded({ type: getTypeFromValue(v), value: v });
      }).join(', ') + ']';
    case 'object':
      const entries = Object.entries(decoded.value);
      return '{' + entries.map(([k, v]) => {
        const keyStr = JSON.stringify(k);
        if (typeof v === 'object' && v !== null && 'type' in v && 'value' in v) {
          return keyStr + ': ' + jsoncFromDecoded(v as { type: string; value: any });
        }
        return keyStr + ': ' + jsoncFromDecoded({ type: getTypeFromValue(v), value: v });
      }).join(', ') + '}';
    case 'simple':
      return decoded.value.toString();
    case 'tag':
      return JSON.stringify(decoded.value);
    default:
      return 'null';
  }
}

function getTypeFromValue(value: any): string {
  if (value === null) return 'null';
  if (typeof value === 'boolean') return 'bool';
  if (typeof value === 'number') return 'float';
  if (typeof value === 'bigint') return 'int';
  if (typeof value === 'string') return 'string';
  if (value instanceof Uint8Array) return 'bytes';
  if (Array.isArray(value)) return 'array';
  if (typeof value === 'object') return 'object';
  return 'null';
}

export function fromJSONC(jsoncString: string): Uint8Array {
  const parsed = JSON.parse(jsoncString);
  return encode(parsed);
}

export function bind<T>(data: Uint8Array | ArrayBuffer | number[], template: T): T {
  const decoded = decode(data);
  return bindDecodedToTemplate(decoded, template) as T;
}

function bindDecodedToTemplate(decoded: any, template: any): any {
  if (template === null || template === undefined) {
    return decoded;
  }
  
  // 处理解码结果的结构 { type: string, value: any }
  const decodedValue = decoded && typeof decoded === 'object' && 'value' in decoded ? decoded.value : decoded;
  
  if (typeof template === 'object' && 'value' in template && 'options' in template) {
    // MMValue 结构
    const mmValue = template;
    return {
      value: bindDecodedToTemplate(decodedValue, mmValue.value),
      options: mmValue.options
    };
  } else if (Array.isArray(template)) {
    // 数组
    if (Array.isArray(decodedValue)) {
      return decodedValue.map((item, index) => {
        const templateItem = template[index];
        return bindDecodedToTemplate(item, templateItem);
      });
    } else {
      return [];
    }
  } else if (typeof template === 'object') {
    // 对象
    const result: any = {};
    for (const key in template) {
      if (template.hasOwnProperty(key)) {
        const templateValue = template[key];
        const itemValue = decodedValue && typeof decodedValue === 'object' ? decodedValue[key] : undefined;
        result[key] = bindDecodedToTemplate(itemValue, templateValue);
      }
    }
    return result;
  } else {
    // 基本类型
    return decodedValue;
  }
}

export { mm, MMValue, Tag, ValueType } from './types';
export { MMEncoder, MMDecoder, MMBuffer, MMError };
export type { DecodedValue } from './decoder';
export * as constants from './constants';
export { MmValidator, validator, ValidationResult } from './validator';
export default {
  encode,
  decode,
  toJSONC,
  fromJSONC,
  bind,
  mm,
  MMEncoder,
  MMDecoder,
  MMBuffer,
  MMError,
  constants,
  MmValidator,
  validator,
};