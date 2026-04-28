import { Tag, ValueType, MMValue } from './types';

// Regex patterns
const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
const decimalRegex = /^-?\d+\.\d+$/;
const uuidRegex = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/;

// Validation result interface
export interface ValidationResult {
  valid: boolean;
  error?: string;
  data?: any;
  text?: string;
}

// Validator class
export class MmValidator {
  // Validate array
  static validateArray(value: any[], tag: Tag): ValidationResult {
    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type array not support location UTC${tag.location}` };
    }

    const length = value.length;

    if (length === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type array not allow empty' };
      }
      return { valid: true, data: value, text: JSON.stringify(value) };
    }

    if (tag.size && length > tag.size) {
      return { valid: false, error: 'type array over size' };
    }

    if (tag.childUnique) {
      const seen = new Set();
      for (let i = 0; i < value.length; i++) {
        const data = value[i];
        const key = typeof data === 'object' ? JSON.stringify(data) : data;
        if (seen.has(key)) {
          return { valid: false, error: `array duplicate value found: ${data}, index: ${i}` };
        }
        seen.add(key);
      }
    }

    return { valid: true, data: value, text: JSON.stringify(value) };
  }

  // Validate struct
  static validateStruct(tag: Tag): ValidationResult {
    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type struct not support location UTC${tag.location}` };
    }

    return { valid: true };
  }

  // Validate string
  static validateString(value: string, tag: Tag): ValidationResult {
    if (value === '') {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type string not allow empty value "${value}"` };
      }
      return { valid: true, data: value, text: value };
    }

    if (tag.pattern) {
      try {
        const regex = new RegExp(tag.pattern);
        if (!regex.test(value)) {
          return { valid: false, error: `value "${value}" does not match pattern ${tag.pattern}` };
        }
      } catch (error) {
        return { valid: false, error: `pattern "${tag.pattern}" compile err: ${error}` };
      }
    }

    const length = value.length;

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as int: ${tag.min}` };
      }
      if (length < mini) {
        return { valid: false, error: `string length ${length} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as int: ${tag.max}` };
      }
      if (length > maxi) {
        return { valid: false, error: `string length ${length} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.size && length !== tag.size) {
      return { valid: false, error: `string length ${length} != size ${tag.size}` };
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type string not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value };
  }

  // Validate bytes
  static validateBytes(value: Uint8Array | number[], tag: Tag): ValidationResult {
    const length = Array.isArray(value) ? value.length : value.byteLength;

    if (length === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type []byte not allow empty value []byte{}' };
      }
      return { valid: true, data: value, text: '' };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as int: ${tag.min}` };
      }
      if (length < mini) {
        return { valid: false, error: `[]byte length ${length} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as int: ${tag.max}` };
      }
      if (length > maxi) {
        return { valid: false, error: `[]byte length ${length} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.size && length !== tag.size) {
      return { valid: false, error: `[]byte length ${length} != size ${tag.size}` };
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type []byte not support location UTC${tag.location}` };
    }

    // Convert to base64 string
    let text = '';
    if (length > 0) {
      if (Array.isArray(value)) {
        const uint8Array = new Uint8Array(value);
        text = btoa(String.fromCharCode(...uint8Array));
      } else {
        text = btoa(String.fromCharCode(...value));
      }
    }

    return { valid: true, data: value, text };
  }

  // Validate bool
  static validateBool(value: boolean, tag: Tag): ValidationResult {
    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.allowEmpty) {
      return { valid: false, error: 'type bool not support allow empty' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type bool not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate int8
  static validateInt8(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type int8 not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < -128 || value > 127) {
      return { valid: false, error: `value ${value} is out of int8 range [-128, 127]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as int8: ${tag.min}` };
      }
      if (mini < -128 || mini > 127) {
        return { valid: false, error: `tag.min ${mini} is out of int8 range [-128, 127]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as int8: ${tag.max}` };
      }
      if (maxi < -128 || maxi > 127) {
        return { valid: false, error: `tag.max ${maxi} is out of int8 range [-128, 127]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type int8 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate int16
  static validateInt16(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type int16 not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < -32768 || value > 32767) {
      return { valid: false, error: `value ${value} is out of int16 range [-32768, 32767]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as int16: ${tag.min}` };
      }
      if (mini < -32768 || mini > 32767) {
        return { valid: false, error: `tag.min ${mini} is out of int16 range [-32768, 32767]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as int16: ${tag.max}` };
      }
      if (maxi < -32768 || maxi > 32767) {
        return { valid: false, error: `tag.max ${maxi} is out of int16 range [-32768, 32767]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type int16 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate int32
  static validateInt32(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type int32 not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < -2147483648 || value > 2147483647) {
      return { valid: false, error: `value ${value} is out of int32 range [-2147483648, 2147483647]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as int32: ${tag.min}` };
      }
      if (mini < -2147483648 || mini > 2147483647) {
        return { valid: false, error: `tag.min ${mini} is out of int32 range [-2147483648, 2147483647]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as int32: ${tag.max}` };
      }
      if (maxi < -2147483648 || maxi > 2147483647) {
        return { valid: false, error: `tag.max ${maxi} is out of int32 range [-2147483648, 2147483647]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type int32 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate int64
  static validateInt64(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type int64 not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < -9223372036854775808 || value > 9223372036854775807) {
      return { valid: false, error: `value ${value} is out of int64 range [-9223372036854775808, 9223372036854775807]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as int64: ${tag.min}` };
      }
      if (mini < -9223372036854775808 || mini > 9223372036854775807) {
        return { valid: false, error: `tag.min ${mini} is out of int64 range [-9223372036854775808, 9223372036854775807]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as int64: ${tag.max}` };
      }
      if (maxi < -9223372036854775808 || maxi > 9223372036854775807) {
        return { valid: false, error: `tag.max ${maxi} is out of int64 range [-9223372036854775808, 9223372036854775807]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type int64 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate uint
  static validateUint(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type uint not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < 0 || value > 4294967295) {
      return { valid: false, error: `value ${value} is out of uint range [0, 4294967295]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as uint: ${tag.min}` };
      }
      if (mini < 0 || mini > 4294967295) {
        return { valid: false, error: `tag.min ${mini} is out of uint range [0, 4294967295]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as uint: ${tag.max}` };
      }
      if (maxi < 0 || maxi > 4294967295) {
        return { valid: false, error: `tag.max ${maxi} is out of uint range [0, 4294967295]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type uint not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate uint8
  static validateUint8(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type uint8 not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < 0 || value > 255) {
      return { valid: false, error: `value ${value} is out of uint8 range [0, 255]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as uint8: ${tag.min}` };
      }
      if (mini < 0 || mini > 255) {
        return { valid: false, error: `tag.min ${mini} is out of uint8 range [0, 255]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as uint8: ${tag.max}` };
      }
      if (maxi < 0 || maxi > 255) {
        return { valid: false, error: `tag.max ${maxi} is out of uint8 range [0, 255]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type uint8 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate uint16
  static validateUint16(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type uint16 not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < 0 || value > 65535) {
      return { valid: false, error: `value ${value} is out of uint16 range [0, 65535]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as uint16: ${tag.min}` };
      }
      if (mini < 0 || mini > 65535) {
        return { valid: false, error: `tag.min ${mini} is out of uint16 range [0, 65535]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as uint16: ${tag.max}` };
      }
      if (maxi < 0 || maxi > 65535) {
        return { valid: false, error: `tag.max ${maxi} is out of uint16 range [0, 65535]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type uint16 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate uint32
  static validateUint32(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type uint32 not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < 0 || value > 4294967295) {
      return { valid: false, error: `value ${value} is out of uint32 range [0, 4294967295]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as uint32: ${tag.min}` };
      }
      if (mini < 0 || mini > 4294967295) {
        return { valid: false, error: `tag.min ${mini} is out of uint32 range [0, 4294967295]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as uint32: ${tag.max}` };
      }
      if (maxi < 0 || maxi > 4294967295) {
        return { valid: false, error: `tag.max ${maxi} is out of uint32 range [0, 4294967295]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type uint32 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate uint64
  static validateUint64(value: number, tag: Tag): ValidationResult {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type uint64 not allow empty value ${value}` };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (value < 0 || value > 18446744073709551615) {
      return { valid: false, error: `value ${value} is out of uint64 range [0, 18446744073709551615]` };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as uint64: ${tag.min}` };
      }
      if (mini < 0 || mini > 18446744073709551615) {
        return { valid: false, error: `tag.min ${mini} is out of uint64 range [0, 18446744073709551615]` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as uint64: ${tag.max}` };
      }
      if (maxi < 0 || maxi > 18446744073709551615) {
        return { valid: false, error: `tag.max ${maxi} is out of uint64 range [0, 18446744073709551615]` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type uint64 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate float32
  static validateFloat32(value: number, tag: Tag): ValidationResult {
    if (value === 0.0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type float32 not allow empty value 0.0' };
      }
      return { valid: true, data: value, text: '0.0' };
    }

    if (tag.min) {
      const mini = parseFloat(tag.min);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as float32: ${tag.min}` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseFloat(tag.max);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as float32: ${tag.max}` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type float32 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate float64
  static validateFloat64(value: number, tag: Tag): ValidationResult {
    if (value === 0.0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type float64 not allow empty value 0.0' };
      }
      return { valid: true, data: value, text: '0.0' };
    }

    if (tag.min) {
      const mini = parseFloat(tag.min);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as float64: ${tag.min}` };
      }
      if (value < mini) {
        return { valid: false, error: `value ${value} is less than the minimum limit ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseFloat(tag.max);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as float64: ${tag.max}` };
      }
      if (value > maxi) {
        return { valid: false, error: `value ${value} exceeds the maximum limit ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type float64 not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate bigInt
  static validateBigInt(value: bigint, tag: Tag): ValidationResult {
    if (value === BigInt(0)) {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type bigInt not allow empty value 0' };
      }
      return { valid: true, data: value, text: '0' };
    }

    if (tag.min) {
      const mini = BigInt(tag.min);
      if (value < mini) {
        return { valid: false, error: `bigInt value ${value} < min ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = BigInt(tag.max);
      if (value > maxi) {
        return { valid: false, error: `bigInt value ${value} > max ${maxi}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type bigInt not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value.toString() };
  }

  // Validate datetime
  static validateDateTime(value: Date, tag: Tag): ValidationResult {
    if (value.getTime() === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: `type datetime not allow empty ${value.toISOString()}` };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    // Format as datetime string
    const format = value.toISOString().replace('T', ' ').substring(0, 19);

    return { valid: true, data: value, text: format };
  }

  // Validate UUID
  static validateUUID(value: string, tag: Tag): ValidationResult {
    if (value === '') {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type uuid not allow empty value ""' };
      }
      return { valid: true, data: new Uint8Array(16), text: value };
    }

    if (!uuidRegex.test(value)) {
      return { valid: false, error: `value '${value}' does not match UUID pattern` };
    }

    // Parse UUID to bytes (simplified)
    const uuidBytes = new Uint8Array(16);
    const parts = value.replace(/-/g, '').match(/.{2}/g);
    if (parts) {
      parts.forEach((part, index) => {
        if (index < 16) {
          uuidBytes[index] = parseInt(part, 16);
        }
      });
    }

    if (tag.version) {
      // Extract version from UUID
      const version = parseInt(value.substring(14, 15), 16);
      if (tag.version !== version) {
        return { valid: false, error: 'invalid uuid version' };
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type uuid not support location UTC${tag.location}` };
    }

    return { valid: true, data: uuidBytes, text: value };
  }

  // Validate email
  static validateEmail(value: string, tag: Tag): ValidationResult {
    if (value === '') {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type email not allow empty value ""' };
      }
      return { valid: true, data: value, text: value };
    }

    if (!emailRegex.test(value)) {
      return { valid: false, error: `value '${value}' does not match email pattern` };
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type email not support location UTC${tag.location}` };
    }

    return { valid: true, data: value, text: value };
  }

  // Validate enum
  static validateEnum(value: string, tag: Tag): ValidationResult {
    if (value === '') {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type enum not allow empty value ""' };
      }
      return { valid: true, data: -1, text: value };
    }

    if (!tag.enum) {
      return { valid: false, error: 'enum not defined' };
    }

    const enums = tag.enum?.split('|') || [];
    let idx = -1;
    for (let i = 0; i < enums.length; i++) {
      const enumValue = enums[i];
      if (enumValue && enumValue.trim() === value) {
        idx = i;
        break;
      }
    }

    if (idx === -1) {
      return { valid: false, error: `value '${value}' not found in enum: ${enums}` };
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type enum not support location UTC${tag.location}` };
    }

    return { valid: true, data: idx, text: value };
  }

  // Validate image
  static validateImage(value: Uint8Array | number[], tag: Tag): ValidationResult {
    const length = Array.isArray(value) ? value.length : value.byteLength;

    if (length === 0) {
      if (!tag.allowEmpty) {
        return { valid: false, error: 'type image not allow empty value []byte{}' };
      }
      return { valid: true, data: value, text: '' };
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return { valid: false, error: `failed to parse tag.min as int: ${tag.min}` };
      }
      if (length < mini) {
        return { valid: false, error: `[]byte length ${length} < min ${mini}` };
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return { valid: false, error: `failed to parse tag.max as int: ${tag.max}` };
      }
      if (length > maxi) {
        return { valid: false, error: `[]byte length ${length} > max ${maxi}` };
      }
    }

    if (tag.size && length !== tag.size) {
      return { valid: false, error: `[]byte length ${length} != size ${tag.size}` };
    }

    if (tag.desc && tag.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (tag.location && tag.location !== 0) {
      return { valid: false, error: `type image not support location UTC${tag.location}` };
    }

    // Convert to base64 string
    let text = '';
    if (length > 0) {
      if (Array.isArray(value)) {
        const uint8Array = new Uint8Array(value);
        text = btoa(String.fromCharCode(...uint8Array));
      } else {
        text = btoa(String.fromCharCode(...value));
      }
    }

    return { valid: true, data: value, text };
  }

  // Validate any value based on type
  static validate(value: any, tag: Tag): ValidationResult {
    const type = tag.type || 'unknown';

    switch (type) {
      case 'array':
        if (Array.isArray(value)) {
          return this.validateArray(value, tag);
        }
        return { valid: false, error: `expected array, got ${typeof value}` };
      
      case 'struct':
        return this.validateStruct(tag);
      
      case 'string':
        if (typeof value === 'string') {
          return this.validateString(value, tag);
        }
        return { valid: false, error: `expected string, got ${typeof value}` };
      
      case 'bytes':
        if (Array.isArray(value) || value instanceof Uint8Array) {
          return this.validateBytes(value, tag);
        }
        return { valid: false, error: `expected bytes, got ${typeof value}` };
      
      case 'bool':
        if (typeof value === 'boolean') {
          return this.validateBool(value, tag);
        }
        return { valid: false, error: `expected bool, got ${typeof value}` };
      
      case 'int':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt8(value, tag);
        }
        return { valid: false, error: `expected integer, got ${typeof value}` };

      case 'int8':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt8(value, tag);
        }
        return { valid: false, error: `expected int8, got ${typeof value}` };

      case 'int16':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt16(value, tag);
        }
        return { valid: false, error: `expected int16, got ${typeof value}` };

      case 'int32':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt32(value, tag);
        }
        return { valid: false, error: `expected int32, got ${typeof value}` };

      case 'int64':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt64(value, tag);
        }
        return { valid: false, error: `expected int64, got ${typeof value}` };

      case 'uint':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint(value, tag);
        }
        return { valid: false, error: `expected uint, got ${typeof value}` };

      case 'uint8':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint8(value, tag);
        }
        return { valid: false, error: `expected uint8, got ${typeof value}` };

      case 'uint16':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint16(value, tag);
        }
        return { valid: false, error: `expected uint16, got ${typeof value}` };

      case 'uint32':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint32(value, tag);
        }
        return { valid: false, error: `expected uint32, got ${typeof value}` };

      case 'uint64':
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint64(value, tag);
        }
        return { valid: false, error: `expected uint64, got ${typeof value}` };

      case 'float32':
        if (typeof value === 'number' && !Number.isInteger(value)) {
          return this.validateFloat32(value, tag);
        }
        return { valid: false, error: `expected float32, got ${typeof value}` };

      case 'float64':
        if (typeof value === 'number' && !Number.isInteger(value)) {
          return this.validateFloat64(value, tag);
        }
        return { valid: false, error: `expected float64, got ${typeof value}` };

      case 'bigint':
        if (typeof value === 'bigint') {
          return this.validateBigInt(value, tag);
        }
        return { valid: false, error: `expected bigint, got ${typeof value}` };
      
      case 'datetime':
      case 'date':
      case 'time':
        if (value instanceof Date) {
          return this.validateDateTime(value, tag);
        }
        return { valid: false, error: `expected Date, got ${typeof value}` };
      
      case 'uuid':
        if (typeof value === 'string') {
          return this.validateUUID(value, tag);
        }
        return { valid: false, error: `expected string, got ${typeof value}` };
      
      case 'email':
        if (typeof value === 'string') {
          return this.validateEmail(value, tag);
        }
        return { valid: false, error: `expected string, got ${typeof value}` };
      
      case 'enum':
        if (typeof value === 'string') {
          return this.validateEnum(value, tag);
        }
        return { valid: false, error: `expected string, got ${typeof value}` };
      
      case 'image':
        if (Array.isArray(value) || value instanceof Uint8Array) {
          return this.validateImage(value, tag);
        }
        return { valid: false, error: `expected bytes, got ${typeof value}` };
      
      default:
        return { valid: true, data: value, text: String(value) };
    }
  }

  // Validate MMValue
  static validateMMValue<T>(mmValue: MMValue<T>): ValidationResult {
    return this.validate(mmValue.value, mmValue.options);
  }
}

// Export the validator
export const validator = MmValidator;