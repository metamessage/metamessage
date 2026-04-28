import { MMEncoder } from './encoder.js';
import { MMDecoder } from './decoder.js';
import { mm } from './types.js';
import * as constants from './constants.js';
import { MMBuffer, MMError } from './buffer.js';
import { MmValidator, validator, ValidationResult } from './validator.js';

export { mm } from './types.js';
export { MMEncoder, MMDecoder, MMBuffer, MMError };
export { constants };
export { MmValidator, validator, ValidationResult };

export function encode(value) {
  const encoder = new MMEncoder();
  encoder.encode(value);
  return encoder.result;
}

export function decode(data) {
  const decoder = new MMDecoder(data);
  return decoder.decode();
}

export function toJSONC(data) {
  const decoded = decode(data);
  return jsoncFromDecoded(decoded);
}

function jsoncFromDecoded(decoded) {
  switch (decoded.type) {
    case 'null':
      return 'null';
    case 'bool':
      return decoded.value ? 'true' : 'false';
    case 'int':
      return String(decoded.value);
    case 'float':
      return String(decoded.value);
    case 'string':
      return JSON.stringify(decoded.value);
    case 'bytes':
      return JSON.stringify(Array.from(decoded.value));
    case 'array':
      return '[' + decoded.value.map(jsoncFromDecoded).join(', ') + ']';
    case 'object':
      const entries = Object.entries(decoded.value);
      return '{' + entries.map(([k, v]) => JSON.stringify(k) + ': ' + jsoncFromDecoded(v)).join(', ') + '}';
    default:
      return 'null';
  }
}

export function fromJSONC(jsoncString) {
  return encode(JSON.parse(jsoncString));
}

export { parseJSONC, jsoncScanner } from '../jsonc/scanner.js';
export { printJSONC, printJSONCCompact } from '../jsonc/printer.js';