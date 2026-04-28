import { ValueType } from './constants.js';

// Regex patterns
const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
const decimalRegex = /^-?\d+\.\d+$/;
const uuidRegex = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/;

// Validation result class
class ValidationResult {
  constructor(valid, error = null, data = null, text = null) {
    this.valid = valid;
    this.error = error;
    this.data = data;
    this.text = text;
  }
}

// Validator class
class MmValidator {
  // Validate array
  static validateArray(value, tag) {
    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type array not support location UTC${tag.location}`);
    }

    const length = value.length;

    if (length === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type array not allow empty');
      }
      return new ValidationResult(true, null, value, JSON.stringify(value));
    }

    if (tag.size && length > tag.size) {
      return new ValidationResult(false, 'type array over size');
    }

    if (tag.childUnique) {
      const seen = new Set();
      for (let i = 0; i < value.length; i++) {
        const data = value[i];
        const key = typeof data === 'object' ? JSON.stringify(data) : data;
        if (seen.has(key)) {
          return new ValidationResult(false, `array duplicate value found: ${data}, index: ${i}`);
        }
        seen.add(key);
      }
    }

    return new ValidationResult(true, null, value, JSON.stringify(value));
  }

  // Validate struct
  static validateStruct(tag) {
    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type struct not support location UTC${tag.location}`);
    }

    return new ValidationResult(true);
  }

  // Validate string
  static validateString(value, tag) {
    if (value === '') {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type string not allow empty value "${value}"`);
      }
      return new ValidationResult(true, null, value, value);
    }

    if (tag.pattern) {
      try {
        const regex = new RegExp(tag.pattern);
        if (!regex.test(value)) {
          return new ValidationResult(false, `value "${value}" does not match pattern ${tag.pattern}`);
        }
      } catch (error) {
        return new ValidationResult(false, `pattern "${tag.pattern}" compile err: ${error}`);
      }
    }

    const length = value.length;

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as int: ${tag.min}`);
      }
      if (length < mini) {
        return new ValidationResult(false, `string length ${length} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as int: ${tag.max}`);
      }
      if (length > maxi) {
        return new ValidationResult(false, `string length ${length} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.size && length !== tag.size) {
      return new ValidationResult(false, `string length ${length} != size ${tag.size}`);
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type string not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value);
  }

  // Validate bytes
  static validateBytes(value, tag) {
    const length = Array.isArray(value) ? value.length : value.byteLength;

    if (length === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type []byte not allow empty value []byte{}');
      }
      return new ValidationResult(true, null, value, '');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as int: ${tag.min}`);
      }
      if (length < mini) {
        return new ValidationResult(false, `[]byte length ${length} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as int: ${tag.max}`);
      }
      if (length > maxi) {
        return new ValidationResult(false, `[]byte length ${length} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.size && length !== tag.size) {
      return new ValidationResult(false, `[]byte length ${length} != size ${tag.size}`);
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type []byte not support location UTC${tag.location}`);
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

    return new ValidationResult(true, null, value, text);
  }

  // Validate bool
  static validateBool(value, tag) {
    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.allowEmpty) {
      return new ValidationResult(false, 'type bool not support allow empty');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type bool not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate int
  static validateInt(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type int not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as int: ${tag.min}`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as int: ${tag.max}`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type int not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate int8
  static validateInt8(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type int8 not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as int8: ${tag.min}`);
      }
      if (mini < -128 || mini > 127) {
        return new ValidationResult(false, `tag.min ${mini} is out of int8 range [-128, 127]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as int8: ${tag.max}`);
      }
      if (maxi < -128 || maxi > 127) {
        return new ValidationResult(false, `tag.max ${maxi} is out of int8 range [-128, 127]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type int8 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate int16
  static validateInt16(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type int16 not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as int16: ${tag.min}`);
      }
      if (mini < -32768 || mini > 32767) {
        return new ValidationResult(false, `tag.min ${mini} is out of int16 range [-32768, 32767]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as int16: ${tag.max}`);
      }
      if (maxi < -32768 || maxi > 32767) {
        return new ValidationResult(false, `tag.max ${maxi} is out of int16 range [-32768, 32767]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type int16 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate int32
  static validateInt32(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type int32 not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as int32: ${tag.min}`);
      }
      if (mini < -2147483648 || mini > 2147483647) {
        return new ValidationResult(false, `tag.min ${mini} is out of int32 range [-2147483648, 2147483647]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as int32: ${tag.max}`);
      }
      if (maxi < -2147483648 || maxi > 2147483647) {
        return new ValidationResult(false, `tag.max ${maxi} is out of int32 range [-2147483648, 2147483647]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type int32 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate int64
  static validateInt64(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type int64 not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as int64: ${tag.min}`);
      }
      if (mini < -9223372036854775808 || mini > 9223372036854775807) {
        return new ValidationResult(false, `tag.min ${mini} is out of int64 range [-9223372036854775808, 9223372036854775807]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as int64: ${tag.max}`);
      }
      if (maxi < -9223372036854775808 || maxi > 9223372036854775807) {
        return new ValidationResult(false, `tag.max ${maxi} is out of int64 range [-9223372036854775808, 9223372036854775807]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type int64 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate uint
  static validateUint(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type uint not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as uint: ${tag.min}`);
      }
      if (mini < 0 || mini > 4294967295) {
        return new ValidationResult(false, `tag.min ${mini} is out of uint range [0, 4294967295]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as uint: ${tag.max}`);
      }
      if (maxi < 0 || maxi > 4294967295) {
        return new ValidationResult(false, `tag.max ${maxi} is out of uint range [0, 4294967295]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type uint not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate uint8
  static validateUint8(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type uint8 not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as uint8: ${tag.min}`);
      }
      if (mini < 0 || mini > 255) {
        return new ValidationResult(false, `tag.min ${mini} is out of uint8 range [0, 255]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as uint8: ${tag.max}`);
      }
      if (maxi < 0 || maxi > 255) {
        return new ValidationResult(false, `tag.max ${maxi} is out of uint8 range [0, 255]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type uint8 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate uint16
  static validateUint16(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type uint16 not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as uint16: ${tag.min}`);
      }
      if (mini < 0 || mini > 65535) {
        return new ValidationResult(false, `tag.min ${mini} is out of uint16 range [0, 65535]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as uint16: ${tag.max}`);
      }
      if (maxi < 0 || maxi > 65535) {
        return new ValidationResult(false, `tag.max ${maxi} is out of uint16 range [0, 65535]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type uint16 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate uint32
  static validateUint32(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type uint32 not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as uint32: ${tag.min}`);
      }
      if (mini < 0 || mini > 4294967295) {
        return new ValidationResult(false, `tag.min ${mini} is out of uint32 range [0, 4294967295]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as uint32: ${tag.max}`);
      }
      if (maxi < 0 || maxi > 4294967295) {
        return new ValidationResult(false, `tag.max ${maxi} is out of uint32 range [0, 4294967295]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type uint32 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate uint64
  static validateUint64(value, tag) {
    if (value === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type uint64 not allow empty value ${value}`);
      }
      return new ValidationResult(true, null, value, '0');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as uint64: ${tag.min}`);
      }
      if (mini < 0 || mini > 18446744073709551615) {
        return new ValidationResult(false, `tag.min ${mini} is out of uint64 range [0, 18446744073709551615]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `value ${value} is less than the minimum limit ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as uint64: ${tag.max}`);
      }
      if (maxi < 0 || maxi > 18446744073709551615) {
        return new ValidationResult(false, `tag.max ${maxi} is out of uint64 range [0, 18446744073709551615]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `value ${value} exceeds the maximum limit ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type uint64 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate float32
  static validateFloat32(value, tag) {
    if (value === 0.0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type float32 not allow empty value 0.0');
      }
      return new ValidationResult(true, null, value, '0.0');
    }

    if (tag.min) {
      const mini = parseFloat(tag.min);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as float32: ${tag.min}`);
      }
      if (mini < -3.402823466e+38 || mini > 3.402823466e+38) {
        return new ValidationResult(false, `tag.min ${mini} is out of float32 range [-3.402823466e+38, 3.402823466e+38]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `${value} < min ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseFloat(tag.max);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as float32: ${tag.max}`);
      }
      if (maxi < -3.402823466e+38 || maxi > 3.402823466e+38) {
        return new ValidationResult(false, `tag.max ${maxi} is out of float32 range [-3.402823466e+38, 3.402823466e+38]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `${value} > max ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type float32 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate float64
  static validateFloat64(value, tag) {
    if (value === 0.0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type float64 not allow empty value 0.0');
      }
      return new ValidationResult(true, null, value, '0.0');
    }

    if (tag.min) {
      const mini = parseFloat(tag.min);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as float64: ${tag.min}`);
      }
      if (mini < -1.7976931348623157e+308 || mini > 1.7976931348623157e+308) {
        return new ValidationResult(false, `tag.min ${mini} is out of float64 range [-1.7976931348623157e+308, 1.7976931348623157e+308]`);
      }
      if (value < mini) {
        return new ValidationResult(false, `${value} < min ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseFloat(tag.max);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as float64: ${tag.max}`);
      }
      if (maxi < -1.7976931348623157e+308 || maxi > 1.7976931348623157e+308) {
        return new ValidationResult(false, `tag.max ${maxi} is out of float64 range [-1.7976931348623157e+308, 1.7976931348623157e+308]`);
      }
      if (value > maxi) {
        return new ValidationResult(false, `${value} > max ${maxi}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type float64 not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value.toString());
  }

  // Validate BigInt
  static validateBigInt(value, tag) {
    const bigIntValue = BigInt(value);

    if (bigIntValue === 0n) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type big.Int not allow empty value 0');
      }
      return new ValidationResult(true, null, bigIntValue, '0');
    }

    if (tag.min) {
      const mini = BigInt(tag.min);
      if (bigIntValue < mini) {
        return new ValidationResult(false, `big.Int length ${bigIntValue.toString()} < min ${mini.toString()}`);
      }
    }

    if (tag.max) {
      const maxi = BigInt(tag.max);
      if (bigIntValue > maxi) {
        return new ValidationResult(false, `big.Int length ${bigIntValue.toString()} > max ${maxi.toString()}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type big.Int not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, bigIntValue, bigIntValue.toString());
  }

  // Validate datetime
  static validateDateTime(value, tag) {
    if (value.getTime() === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type datetime not allow empty ${value.toISOString()}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    // Format as datetime string
    const format = value.toISOString().replace('T', ' ').substring(0, 19);

    return new ValidationResult(true, null, value, format);
  }

  // Validate UUID
  static validateUUID(value, tag) {
    if (value === '') {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type uuid not allow empty value ""');
      }
      return new ValidationResult(true, null, new Uint8Array(16), value);
    }

    if (!uuidRegex.test(value)) {
      return new ValidationResult(false, `value '${value}' does not match UUID pattern`);
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
        return new ValidationResult(false, 'invalid uuid version');
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type uuid not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, uuidBytes, value);
  }

  // Validate email
  static validateEmail(value, tag) {
    if (value === '') {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type email not allow empty value ""');
      }
      return new ValidationResult(true, null, value, value);
    }

    if (!emailRegex.test(value)) {
      return new ValidationResult(false, `value '${value}' does not match email pattern`);
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type email not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value);
  }

  // Validate enum
  static validateEnum(value, tag) {
    if (value === '') {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type enum not allow empty value ""');
      }
      return new ValidationResult(true, null, -1, value);
    }

    if (!tag.enum) {
      return new ValidationResult(false, 'enum not defined');
    }

    const enums = tag.enum.split('|');
    let idx = -1;
    for (let i = 0; i < enums.length; i++) {
      const enumValue = enums[i];
      if (enumValue && enumValue.trim() === value) {
        idx = i;
        break;
      }
    }

    if (idx === -1) {
      return new ValidationResult(false, `value '${value}' not found in enum: ${enums}`);
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type enum not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, idx, value);
  }

  // Validate image
  static validateImage(value, tag) {
    const length = Array.isArray(value) ? value.length : value.byteLength;

    if (length === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type image not allow empty value []byte{}');
      }
      return new ValidationResult(true, null, value, '');
    }

    if (tag.min) {
      const mini = parseInt(tag.min, 10);
      if (isNaN(mini)) {
        return new ValidationResult(false, `failed to parse tag.min as int: ${tag.min}`);
      }
      if (length < mini) {
        return new ValidationResult(false, `[]byte length ${length} < min ${mini}`);
      }
    }

    if (tag.max) {
      const maxi = parseInt(tag.max, 10);
      if (isNaN(maxi)) {
        return new ValidationResult(false, `failed to parse tag.max as int: ${tag.max}`);
      }
      if (length > maxi) {
        return new ValidationResult(false, `[]byte length ${length} > max ${maxi}`);
      }
    }

    if (tag.size && length !== tag.size) {
      return new ValidationResult(false, `[]byte length ${length} != size ${tag.size}`);
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type image not support location UTC${tag.location}`);
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

    return new ValidationResult(true, null, value, text);
  }

  // Validate decimal
  static validateDecimal(value, tag) {
    if (value === '') {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type decimal not allow empty value ""');
      }
      return new ValidationResult(true, null, value, value);
    }

    if (!decimalRegex.test(value)) {
      return new ValidationResult(false, `invalid decimal "${value}", must be like "0.0"`);
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type decimal not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value);
  }

  // Validate IP
  static validateIP(value, tag) {
    if (value === '') {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type ip not allow empty value ""');
      }
      return new ValidationResult(true, null, null, '');
    }

    const ipv4Regex = /^(\d{1,3}\.){3}\d{1,3}$/;
    const ipv6Regex = /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::$|^([0-9a-fA-F]{1,4}:)*:$|^([0-9a-fA-F]{1,4}:){1,7}:$/;

    const isIPv4 = ipv4Regex.test(value);
    const isIPv6 = ipv6Regex.test(value);

    if (!isIPv4 && !isIPv6) {
      return new ValidationResult(false, `invalid ip "${value}"`);
    }

    if (isIPv4 && tag.version === 6) {
      return new ValidationResult(false, `invalid ipv4: ${value}`);
    }

    if (isIPv6 && tag.version === 4) {
      return new ValidationResult(false, `invalid ipv6: ${value}`);
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type ip not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value);
  }

  // Validate URL
  static validateURL(value, tag) {
    if (value === '') {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type url not allow empty value ""');
      }
      return new ValidationResult(true, null, null, '');
    }

    let url;
    try {
      url = new URL(value);
    } catch (e) {
      return new ValidationResult(false, `invalid url: ${value}`);
    }

    if (url.protocol !== 'http:' && url.protocol !== 'https:') {
      return new ValidationResult(false, `invalid url: ${value}`);
    }

    if (!url.host) {
      return new ValidationResult(false, `invalid url: ${value}`);
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type url not support location UTC${tag.location}`);
    }

    return new ValidationResult(true, null, value, value);
  }

  // Validate slice
  static validateSlice(value, tag) {
    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    if (tag.location && tag.location !== 0) {
      return new ValidationResult(false, `type slice not support location UTC${tag.location}`);
    }

    const length = value.length;

    if (length === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, 'type slice not allow empty');
      }
      return new ValidationResult(true, null, value, JSON.stringify(value));
    }

    if (tag.size && length > tag.size) {
      return new ValidationResult(false, 'type slice over size');
    }

    if (tag.childUnique) {
      const seen = new Set();
      for (let i = 0; i < value.length; i++) {
        const data = value[i];
        const key = typeof data === 'object' ? JSON.stringify(data) : data;
        if (seen.has(key)) {
          return new ValidationResult(false, `slice duplicate value found: ${data}, index: ${i}`);
        }
        seen.add(key);
      }
    }

    return new ValidationResult(true, null, value, JSON.stringify(value));
  }

  // Validate date
  static validateDate(value, tag) {
    if (value.getTime() === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type date not allow empty ${value.toISOString()}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    const format = value.toISOString().substring(0, 10);

    return new ValidationResult(true, null, value, format);
  }

  // Validate time
  static validateTime(value, tag) {
    if (value.getTime() === 0) {
      if (!tag.allowEmpty) {
        return new ValidationResult(false, `type time not allow empty ${value.toISOString()}`);
      }
    }

    if (tag.desc && tag.desc.length > 65535) {
      return new ValidationResult(false, 'desc length exceeds 65535 bytes');
    }

    const format = value.toISOString().substring(11, 19);

    return new ValidationResult(true, null, value, format);
  }

  // Validate any value based on type
  static validate(value, tag) {
    const type = tag.type || ValueType.Unknown;

    switch (type) {
      case ValueType.Array:
        if (Array.isArray(value)) {
          return this.validateArray(value, tag);
        }
        return new ValidationResult(false, `expected array, got ${typeof value}`);
      
      case ValueType.Struct:
        return this.validateStruct(tag);
      
      case ValueType.String:
        if (typeof value === 'string') {
          return this.validateString(value, tag);
        }
        return new ValidationResult(false, `expected string, got ${typeof value}`);
      
      case ValueType.Bytes:
        if (Array.isArray(value) || value instanceof Uint8Array) {
          return this.validateBytes(value, tag);
        }
        return new ValidationResult(false, `expected bytes, got ${typeof value}`);
      
      case ValueType.Bool:
        if (typeof value === 'boolean') {
          return this.validateBool(value, tag);
        }
        return new ValidationResult(false, `expected bool, got ${typeof value}`);
      
      case ValueType.Int:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt(value, tag);
        }
        return new ValidationResult(false, `expected int, got ${typeof value}`);

      case ValueType.Int8:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt8(value, tag);
        }
        return new ValidationResult(false, `expected int8, got ${typeof value}`);

      case ValueType.Int16:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt16(value, tag);
        }
        return new ValidationResult(false, `expected int16, got ${typeof value}`);

      case ValueType.Int32:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt32(value, tag);
        }
        return new ValidationResult(false, `expected int32, got ${typeof value}`);

      case ValueType.Int64:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateInt64(value, tag);
        }
        return new ValidationResult(false, `expected int64, got ${typeof value}`);

      case ValueType.Uint:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint(value, tag);
        }
        return new ValidationResult(false, `expected uint, got ${typeof value}`);

      case ValueType.Uint8:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint8(value, tag);
        }
        return new ValidationResult(false, `expected uint8, got ${typeof value}`);

      case ValueType.Uint16:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint16(value, tag);
        }
        return new ValidationResult(false, `expected uint16, got ${typeof value}`);

      case ValueType.Uint32:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint32(value, tag);
        }
        return new ValidationResult(false, `expected uint32, got ${typeof value}`);

      case ValueType.Uint64:
        if (typeof value === 'number' && Number.isInteger(value)) {
          return this.validateUint64(value, tag);
        }
        return new ValidationResult(false, `expected uint64, got ${typeof value}`);

      case ValueType.Float32:
        if (typeof value === 'number' && !Number.isInteger(value)) {
          return this.validateFloat32(value, tag);
        }
        return new ValidationResult(false, `expected float32, got ${typeof value}`);

      case ValueType.Float64:
        if (typeof value === 'number' && !Number.isInteger(value)) {
          return this.validateFloat64(value, tag);
        }
        return new ValidationResult(false, `expected float64, got ${typeof value}`);

      case ValueType.BigInt:
        if (typeof value === 'bigint') {
          return this.validateBigInt(value, tag);
        }
        return new ValidationResult(false, `expected bigint, got ${typeof value}`);
      
      case ValueType.DateTime:
        if (value instanceof Date) {
          return this.validateDateTime(value, tag);
        }
        return new ValidationResult(false, `expected Date, got ${typeof value}`);

      case ValueType.Date:
        if (value instanceof Date) {
          return this.validateDate(value, tag);
        }
        return new ValidationResult(false, `expected Date, got ${typeof value}`);

      case ValueType.Time:
        if (value instanceof Date) {
          return this.validateTime(value, tag);
        }
        return new ValidationResult(false, `expected Date, got ${typeof value}`);

      case ValueType.Uuid:
        if (typeof value === 'string') {
          return this.validateUUID(value, tag);
        }
        return new ValidationResult(false, `expected string, got ${typeof value}`);
      
      case ValueType.Email:
        if (typeof value === 'string') {
          return this.validateEmail(value, tag);
        }
        return new ValidationResult(false, `expected string, got ${typeof value}`);
      
      case ValueType.Enum:
        if (typeof value === 'string') {
          return this.validateEnum(value, tag);
        }
        return new ValidationResult(false, `expected string, got ${typeof value}`);
      
      case ValueType.Image:
        if (Array.isArray(value) || value instanceof Uint8Array) {
          return this.validateImage(value, tag);
        }
        return new ValidationResult(false, `expected bytes, got ${typeof value}`);

      case ValueType.Decimal:
        if (typeof value === 'string') {
          return this.validateDecimal(value, tag);
        }
        return new ValidationResult(false, `expected string, got ${typeof value}`);

      case ValueType.Ip:
        if (typeof value === 'string') {
          return this.validateIP(value, tag);
        }
        return new ValidationResult(false, `expected string, got ${typeof value}`);

      case ValueType.Url:
        if (typeof value === 'string') {
          return this.validateURL(value, tag);
        }
        return new ValidationResult(false, `expected string, got ${typeof value}`);

      case ValueType.Slice:
        if (Array.isArray(value)) {
          return this.validateSlice(value, tag);
        }
        return new ValidationResult(false, `expected array, got ${typeof value}`);

      default:
        return new ValidationResult(true, null, value, String(value));
    }
  }

  // Validate MM value
  static validateMMValue(mmValue) {
    if (!mmValue || typeof mmValue !== 'object' || !('value' in mmValue) || !('options' in mmValue)) {
      return new ValidationResult(false, 'invalid MM value');
    }
    return this.validate(mmValue.value, mmValue.options);
  }
}

// Export the validator
export const validator = MmValidator;
export { MmValidator, ValidationResult };
