import { ValueType, stringToType, typeToString } from './value-type';

const KIsNull = 0 << 3;
const KExample = 1 << 3;
const KDesc = 2 << 3;
const KType = 3 << 3;
const KRaw = 4 << 3;
const KNullable = 5 << 3;
const KAllowEmpty = 6 << 3;
const KUnique = 7 << 3;
const KDefault = 8 << 3;
const KMin = 9 << 3;
const KMax = 10 << 3;
const KSize = 11 << 3;
const KEnum = 12 << 3;
const KPattern = 13 << 3;
const KLocation = 14 << 3;
const KVersion = 15 << 3;
const KMime = 16 << 3;
const KChildDesc = 17 << 3;
const KChildType = 18 << 3;
const KChildRaw = 19 << 3;
const KChildNullable = 20 << 3;
const KChildAllowEmpty = 21 << 3;
const KChildUnique = 22 << 3;
const KChildDefault = 23 << 3;
const KChildMin = 24 << 3;
const KChildMax = 25 << 3;
const KChildSize = 26 << 3;
const KChildEnum = 27 << 3;
const KChildPattern = 28 << 3;
const KChildLocation = 29 << 3;
const KChildVersion = 30 << 3;
const KChildMime = 31 << 3;

const Max1Byte = 0xFF;
const Max2Byte = 0xFFFF;
const Max3Byte = 0xFFFFFF;
const Max4Byte = 0xFFFFFFFF;
const Max5Byte = 0xFFFFFFFFFFn;
const Max6Byte = 0xFFFFFFFFFFFFn;
const Max7Byte = 0xFFFFFFFFFFFFFFn;

export interface ValidationResult {
  valid: boolean;
  error?: string;
  data?: any;
  text?: string;
}

const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
const decimalRegex = /^-?\d+\.\d+$/;
const uuidRegex = /^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/;

export class Tag {
  name?: string;
  isNull: boolean = false;
  example?: boolean;
  desc: string = '';
  type: ValueType = ValueType.Unknown;
  raw: boolean = false;
  nullable: boolean = false;
  allowEmpty: boolean = false;
  unique: boolean = false;
  default: string = '';
  min: string = '';
  max: string = '';
  size: bigint = 0n;
  enum: string = '';
  pattern: string = '';
  location: number = 0;
  version: number = 0;
  mime: string = '';
  childDesc: string= '';
  childType: ValueType = ValueType.Unknown;
  childRaw: boolean = false;
  childNullable: boolean=false;
  childAllowEmpty: boolean=false;
  childUnique: boolean=false;
  childDefault: string='';
  childMin: string='';
  childMax: string='';
  childSize: number=0;
  childEnum: string='';
  childPattern: string='';
  childLocation: number=0;
  childVersion: number=0;
  childMime: string='';
  isInherit: boolean=false;

  constructor(data: Partial<Tag> = {}) {
    Object.assign(this, data);
  }

  inherit(tag: Tag): void {
    if (!tag) return;
    
    if (this.type === ValueType.Unknown) {
      this.type = tag.type;
    }
    if (!this.desc) {
      this.desc = tag.desc;
    }
    if (!this.nullable) {
      this.nullable = tag.nullable;
    }
    if (!this.allowEmpty) {
      this.allowEmpty = tag.allowEmpty;
    }
    if (!this.unique) {
      this.unique = tag.unique;
    }
    if (!this.min) {
      this.min = tag.min;
    }
    if (!this.max) {
      this.max = tag.max;
    }
    if (this.size === 0n) {
      this.size = tag.size;
    }
    if (!this.enum) {
      this.enum = tag.enum;
    }
    if (!this.pattern) {
      this.pattern = tag.pattern;
    }
    if (!this.location) {
      this.location = tag.location;
    }
    if (this.version === 0) {
      this.version = tag.version;
    }
    if (!this.mime) {
      this.mime = tag.mime;
    }
    
    if (!this.childType) {
      this.childType = tag.childType;
    }
    if (!this.childDesc) {
      this.childDesc = tag.childDesc;
    }
    if (!this.childNullable) {
      this.childNullable = tag.childNullable;
    }
    if (!this.childAllowEmpty) {
      this.childAllowEmpty = tag.childAllowEmpty;
    }
    if (!this.childUnique) {
      this.childUnique = tag.childUnique;
    }
    if (!this.childMin) {
      this.childMin = tag.childMin;
    }
    if (!this.childMax) {
      this.childMax = tag.childMax;
    }
    if (!this.childSize) {
      this.childSize = tag.childSize;
    }
    if (!this.childEnum) {
      this.childEnum = tag.childEnum;
    }
    if (!this.childPattern) {
      this.childPattern = tag.childPattern;
    }
    if (!this.childLocation) {
      this.childLocation = tag.childLocation;
    }
    if (!this.childVersion) {
      this.childVersion = tag.childVersion;
    }
    if (!this.childMime) {
      this.childMime = tag.childMime;
    }
  }

  toString(): string {
    const parts: string[] = [];
    if (this.type !== ValueType.Unknown) {
      if (this.type === ValueType.Int ||
          this.type === ValueType.Float64 ||
          this.type === ValueType.Bool ||
          this.type === ValueType.Slice) {
      } else {
        if (this.type === ValueType.Array && this.size > 0 ||
            this.type === ValueType.Enum && this.enum !== "") {
        } else {
          parts.push(`type=${typeToString(this.type)}`);
        }
      }
    }

    if (this.example) {
      parts.push('example');
    }

    if (this.isNull) {
      parts.push('is_null');
    }

    if (this.nullable) {
      if (!this.isNull) {
        parts.push('nullable');
      }
    }

    if (this.desc) {
      parts.push(`desc=${this.desc}`);
    }

    if (this.raw) {
      parts.push('raw');
    }

    if (this.allowEmpty) {
      parts.push('allow_empty');
    }

    if (this.unique) {
      parts.push('unique');
    }

    if (this.default) {
      parts.push(`default=${this.default}`);
    }

    if (this.min) {
      parts.push(`min=${this.min}`);
    }

    if (this.max) {
      parts.push(`max=${this.max}`);
    }

    if (this.size !== 0n) {
      parts.push(`size=${this.size}`);
    }

    if (this.enum) {
      parts.push(`enum=${this.enum}`);
    }

    if (this.pattern) {
      parts.push(`pattern=${this.pattern}`);
    }

    if (this.location) {
      parts.push(`location=${this.location}`);
    }

    if (this.version !== 0) {
      parts.push(`version=${this.version}`);
    }

    if (this.mime) {
      parts.push(`mime=${this.mime}`);
    }

    if (this.childDesc) {
      parts.push(`child_desc="${this.childDesc}"`);
    }

    if (this.childType) {
      if (this.childType === ValueType.String ||
          this.childType === ValueType.Int ||
          this.childType === ValueType.Float64 ||
          this.childType === ValueType.Bool ||
          this.childType === ValueType.Object ||
          this.childType === ValueType.Slice) {
      } else {
        if (this.childType === ValueType.Array && this.childSize && this.childSize > 0 ||
            this.childType === ValueType.Enum && this.childEnum) {
        } else {
          parts.push(`child_type=${this.childType}`);
        }
      }
    }

    if (this.childRaw) {
      parts.push('child_raw');
    }

    if (this.childNullable) {
      parts.push('child_nullable');
    }

    if (this.childAllowEmpty) {
      parts.push('child_allow_empty');
    }

    if (this.childUnique) {
      parts.push('child_unique');
    }

    if (this.childDefault) {
      parts.push(`child_default=${this.childDefault}`);
    }

    if (this.childMin) {
      parts.push(`child_min=${this.childMin}`);
    }

    if (this.childMax) {
      parts.push(`child_max=${this.childMax}`);
    }

    if (this.childSize) {
      parts.push(`child_size=${this.childSize}`);
    }

    if (this.childEnum) {
      parts.push(`child_enum=${this.childEnum}`);
    }

    if (this.childPattern) {
      parts.push(`child_pattern=${this.childPattern}`);
    }

    if (this.childLocation) {
      parts.push(`child_location=${this.childLocation}`);
    }

    if (this.childVersion) {
      parts.push(`child_version=${this.childVersion}`);
    }

    if (this.childMime) {
      parts.push(`child_mime=${this.childMime}`);
    }

    return parts.join('; ');
  }

  private getLocationOffsetHour(location: number): number {
    return location;
  }

  private parseMIME(mime: string): number {
    return mime.length;
  }

  private encodeUint64(buf: number[], sign: number, uv: bigint): void {
    const uvNum = Number(uv);
    if (uv <= Max1Byte) {
      sign |= 0;
      buf.push(sign);
      buf.push(uvNum & 0xFF);
    } else if (uv <= Max2Byte) {
      sign |= 1;
      buf.push(sign);
      buf.push((uvNum >> 8) & 0xFF);
      buf.push(uvNum & 0xFF);
    } else if (uv <= Max3Byte) {
      sign |= 2;
      buf.push(sign);
      buf.push((uvNum >> 16) & 0xFF);
      buf.push((uvNum >> 8) & 0xFF);
      buf.push(uvNum & 0xFF);
    } else if (uv <= Max4Byte) {
      sign |= 3;
      buf.push(sign);
      buf.push((uvNum >> 24) & 0xFF);
      buf.push((uvNum >> 16) & 0xFF);
      buf.push((uvNum >> 8) & 0xFF);
      buf.push(uvNum & 0xFF);
    } else if (uv <= Max5Byte) {
      sign |= 4;
      buf.push(sign);
      const uv32 = uvNum;
      const uv32High = Math.floor(uv32 / 0x100000000);
      buf.push(uv32High & 0xFF);
      buf.push((uv32 >> 24) & 0xFF);
      buf.push((uv32 >> 16) & 0xFF);
      buf.push((uv32 >> 8) & 0xFF);
      buf.push(uv32 & 0xFF);
    } else if (uv <= Max6Byte) {
      sign |= 5;
      buf.push(sign);
      const uv32 = uvNum;
      const uv32High = Math.floor(uv32 / 0x100000000);
      const uv40 = Math.floor(uv32High / 0x100);
      buf.push(uv40 & 0xFF);
      buf.push(uv32High & 0xFF);
      buf.push((uv32 >> 24) & 0xFF);
      buf.push((uv32 >> 16) & 0xFF);
      buf.push((uv32 >> 8) & 0xFF);
      buf.push(uv32 & 0xFF);
    } else if (uv <= Max7Byte) {
      sign |= 6;
      buf.push(sign);
      const uv32 = uvNum;
      const uv32High = Math.floor(uv32 / 0x100000000);
      const uv48 = Math.floor(uv32High / 0x10000);
      buf.push(uv48 & 0xFF);
      buf.push((uv48 >> 8) & 0xFF);
      buf.push(uv32High & 0xFF);
      buf.push((uv32 >> 24) & 0xFF);
      buf.push((uv32 >> 16) & 0xFF);
      buf.push((uv32 >> 8) & 0xFF);
      buf.push(uv32 & 0xFF);
    } else {
      sign |= 7;
      buf.push(sign);
      const uv32 = uvNum;
      const uv32High = Math.floor(uv32 / 0x100000000);
      const uv56 = Math.floor(uv32High / 0x1000000);
      buf.push(uv56 & 0xFF);
      buf.push((uv56 >> 8) & 0xFF);
      buf.push((uv56 >> 16) & 0xFF);
      buf.push(uv32High & 0xFF);
      buf.push((uv32 >> 24) & 0xFF);
      buf.push((uv32 >> 16) & 0xFF);
      buf.push((uv32 >> 8) & 0xFF);
      buf.push(uv32 & 0xFF);
    }
  }

  toBytes(): Uint8Array {
    const buf: number[] = [];
    
    if (this.example) {
      buf.push(KExample | 1);
    }

    if (this.isNull) {
      buf.push(KIsNull | 1);
    }

    if (this.nullable && !this.isInherit) {
      if (!this.isNull) {
        buf.push(KNullable | 1);
      }
    }

    if (this.desc !== "" && !this.isInherit) {
      const l = this.desc.length;
      const encoder = new TextEncoder();
      const descBytes = encoder.encode(this.desc);
      
      if (l <= 5) {
        buf.push(KDesc | l);
        buf.push(...descBytes);
      } else if (l <= 0xFF) {
        buf.push(KDesc | 6);
        buf.push(l);
        buf.push(...descBytes);
      } else if (l <= 0xFFFF) {
        buf.push(KDesc | 7);
        buf.push((l >> 8) & 0xFF);
        buf.push(l & 0xFF);
        buf.push(...descBytes);
      }
    }

    if (this.type !== ValueType.Unknown && !this.isInherit) {
      if (this.type === ValueType.String ||
          this.type === ValueType.Bytes ||
          this.type === ValueType.Int ||
          this.type === ValueType.Float64 ||
          this.type === ValueType.Bool ||
          this.type === ValueType.Object ||
          this.type === ValueType.Slice) {
      } else {
        if (this.type === ValueType.Array && this.size > 0 ||
            this.type === ValueType.Enum && this.enum !== "") {
        } else {
          buf.push(KType);
          buf.push(this.type);
        }
      }
    }

    if (this.raw && !this.isInherit) {
      buf.push(KRaw | 1);
    }

    if (this.allowEmpty && !this.isInherit) {
      buf.push(KAllowEmpty | 1);
    }

    if (this.unique && !this.isInherit) {
      buf.push(KUnique | 1);
    }

    if (this.default !== "" && !this.isInherit) {
      const l = this.default.length;
      const encoder = new TextEncoder();
      const defaultBytes = encoder.encode(this.default);
      
      if (l < 7) {
        buf.push(KDefault | l);
        buf.push(...defaultBytes);
      } else {
        buf.push(KDefault | 7);
        buf.push(l);
        buf.push(...defaultBytes);
      }
    }

    if (this.min !== "" && !this.isInherit) {
      const l = this.min.length;
      const encoder = new TextEncoder();
      const minBytes = encoder.encode(this.min);
      
      if (l < 7) {
        buf.push(KMin | l);
        buf.push(...minBytes);
      } else {
        buf.push(KMin | 7);
        buf.push(l);
        buf.push(...minBytes);
      }
    }

    if (this.max !== "" && !this.isInherit) {
      const l = this.max.length;
      const encoder = new TextEncoder();
      const maxBytes = encoder.encode(this.max);
      
      if (l < 7) {
        buf.push(KMax | l);
        buf.push(...maxBytes);
      } else {
        buf.push(KMax | 7);
        buf.push(l);
        buf.push(...maxBytes);
      }
    }

    if (this.size !== 0n && !this.isInherit) {
      this.encodeUint64(buf, KSize, this.size);
    }

    if (this.enum !== "" && !this.isInherit) {
      const l = this.enum.length;
      const encoder = new TextEncoder();
      const enumBytes = encoder.encode(this.enum);
      
      if (l <= 5) {
        buf.push(KEnum | l);
        buf.push(...enumBytes);
      } else if (l <= 0xFF) {
        buf.push(KEnum | 6);
        buf.push(l);
        buf.push(...enumBytes);
      } else if (l <= 0xFFFF) {
        buf.push(KEnum | 7);
        buf.push((l >> 8) & 0xFF);
        buf.push(l & 0xFF);
        buf.push(...enumBytes);
      }
    }

    if (this.pattern !== "" && !this.isInherit) {
      const l = this.pattern.length;
      const encoder = new TextEncoder();
      const patternBytes = encoder.encode(this.pattern);
      
      if (l < 7) {
        buf.push(KPattern | l);
        buf.push(...patternBytes);
      } else {
        buf.push(KPattern | 7);
        buf.push(l);
        buf.push(...patternBytes);
      }
    }

    const locationOffsetHour = this.getLocationOffsetHour(this.location);
    if (locationOffsetHour !== 0 && !this.isInherit) {
      const v = String(locationOffsetHour);
      const encoder = new TextEncoder();
      const locationBytes = encoder.encode(v);
      buf.push(KLocation | locationBytes.length);
      buf.push(...locationBytes);
    }

    if (this.version !== 0 && !this.isInherit) {
      this.encodeUint64(buf, KVersion, BigInt(this.version));
    }

    if (this.mime !== "" && !this.isInherit) {
      const l = this.parseMIME(this.mime);
      if (l < 7) {
        buf.push(KMime | l);
      } else {
        buf.push(KMime | 7);
        buf.push(l);
      }
    }

    if (this.childDesc !== "") {
      const l = this.childDesc.length;
      const encoder = new TextEncoder();
      const childDescBytes = encoder.encode(this.childDesc);
      
      if (l <= 5) {
        buf.push(KChildDesc | l);
        buf.push(...childDescBytes);
      } else if (l <= 0xFF) {
        buf.push(KChildDesc | 6);
        buf.push(l);
        buf.push(...childDescBytes);
      } else if (l <= 0xFFFF) {
        buf.push(KChildDesc | 7);
        buf.push((l >> 8) & 0xFF);
        buf.push(l & 0xFF);
        buf.push(...childDescBytes);
      }
    }

    if (this.childType !== ValueType.Unknown) {
      if (this.childType === ValueType.String ||
          this.childType === ValueType.Int ||
          this.childType === ValueType.Float64 ||
          this.childType === ValueType.Bool ||
          this.childType === ValueType.Object ||
          this.childType === ValueType.Slice) {
      } else {
        if (this.childType === ValueType.Array && this.childSize > 0 ||
            this.childType === ValueType.Enum && this.childEnum !== "") {
        } else {
          buf.push(KChildType);
          buf.push(this.childType);
        }
      }
    }

    if (this.childRaw) {
      buf.push(KChildRaw | 1);
    }

    if (this.childNullable) {
      buf.push(KChildNullable | 1);
    }

    if (this.childAllowEmpty) {
      buf.push(KChildAllowEmpty | 1);
    }

    if (this.childUnique) {
      buf.push(KChildUnique | 1);
    }

    if (this.childDefault !== "") {
      const l = this.childDefault.length;
      const encoder = new TextEncoder();
      const childDefaultBytes = encoder.encode(this.childDefault);
      
      if (l < 7) {
        buf.push(KChildDefault | l);
        buf.push(...childDefaultBytes);
      } else {
        buf.push(KChildDefault | 7);
        buf.push(l);
        buf.push(...childDefaultBytes);
      }
    }

    if (this.childMin !== "") {
      const l = this.childMin.length;
      const encoder = new TextEncoder();
      const childMinBytes = encoder.encode(this.childMin);
      
      if (l < 7) {
        buf.push(KChildMin | l);
        buf.push(...childMinBytes);
      } else {
        buf.push(KChildMin | 7);
        buf.push(l);
        buf.push(...childMinBytes);
      }
    }

    if (this.childMax !== "") {
      const l = this.childMax.length;
      const encoder = new TextEncoder();
      const childMaxBytes = encoder.encode(this.childMax);
      
      if (l < 7) {
        buf.push(KChildMax | l);
        buf.push(...childMaxBytes);
      } else {
        buf.push(KChildMax | 7);
        buf.push(l);
        buf.push(...childMaxBytes);
      }
    }

    if (this.childSize !== 0) {
      this.encodeUint64(buf, KChildSize, BigInt(this.childSize));
    }

    if (this.childEnum !== "") {
      const l = this.childEnum.length;
      const encoder = new TextEncoder();
      const childEnumBytes = encoder.encode(this.childEnum);
      
      if (l <= 5) {
        buf.push(KChildEnum | l);
        buf.push(...childEnumBytes);
      } else if (l <= 0xFF) {
        buf.push(KChildEnum | 6);
        buf.push(l);
        buf.push(...childEnumBytes);
      } else if (l <= 0xFFFF) {
        buf.push(KChildEnum | 7);
        buf.push((l >> 8) & 0xFF);
        buf.push(l & 0xFF);
        buf.push(...childEnumBytes);
      }
    }

    if (this.childPattern !== "") {
      const l = this.childPattern.length;
      const encoder = new TextEncoder();
      const childPatternBytes = encoder.encode(this.childPattern);
      
      if (l < 7) {
        buf.push(KChildPattern | l);
        buf.push(...childPatternBytes);
      } else {
        buf.push(KChildPattern | 7);
        buf.push(l);
        buf.push(...childPatternBytes);
      }
    }

    const childLocationOffsetHour = this.getLocationOffsetHour(this.childLocation);
    if (childLocationOffsetHour !== 0) {
      const v = String(childLocationOffsetHour);
      const encoder = new TextEncoder();
      const childLocationBytes = encoder.encode(v);
      buf.push(KChildLocation | childLocationBytes.length);
      buf.push(...childLocationBytes);
    }

    if (this.childVersion !== 0) {
      this.encodeUint64(buf, KChildVersion, BigInt(this.childVersion));
    }

    if (this.childMime !== "") {
      const l = this.parseMIME(this.childMime);
      if (l < 7) {
        buf.push(KChildMime | l);
      } else {
        buf.push(KChildMime | 7);
        buf.push(l);
      }
    }

    return new Uint8Array(buf);
  }

  validateString(val: string): ValidationResult {
    if (val === '') {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: val };
      }
      return { valid: false, error: `type string not allow empty value "${val}"` };
    }

    if (this.pattern) {
      try {
        const regex = new RegExp(this.pattern);
        if (!regex.test(val)) {
          return { valid: false, error: `value "${val}" does not match pattern ${this.pattern}` };
        }
      } catch (error) {
        return { valid: false, error: `pattern "${this.pattern}" compile err: ${error}` };
      }
    }

    const length = val.length;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (length < mini) {
          return { valid: false, error: `string length ${length} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (length > maxi) {
          return { valid: false, error: `string length ${length} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.size !== 0n) {
      if (BigInt(length) !== this.size) {
        return { valid: false, error: `string length ${length} != size ${this.size}` };
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val };
  }

  validateBytes(val: Uint8Array | number[]): ValidationResult {
    const arr = val instanceof Uint8Array ? val : new Uint8Array(val);
    const length = arr.length;

    if (length === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: arr, text: '' };
      }
      return { valid: false, error: 'type []byte not allow empty value []byte{}' };
    }

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (length < mini) {
          return { valid: false, error: `[]byte length ${length} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (length > maxi) {
          return { valid: false, error: `[]byte length ${length} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.size !== 0n) {
      if (BigInt(length) !== this.size) {
        return { valid: false, error: `[]byte length ${length} != size ${this.size}` };
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: arr, text: Buffer.from(arr).toString('base64') };
  }

  validateBool(val: boolean): ValidationResult {
    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (this.allowEmpty) {
      return { valid: false, error: 'type bool not support allow empty' };
    }

    return { valid: true, data: val, text: val ? 'true' : 'false' };
  }

  validateInt(val: bigint): ValidationResult {
    if (val === 0n) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type int not allow empty value ${val}` };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `value ${val64} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `value ${val64} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateInt8(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type int8 not allow empty value ${val}` };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `value ${val64} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `value ${val64} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateInt16(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type int16 not allow empty value ${val}` };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `value ${val64} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `value ${val64} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateInt32(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type int32 not allow empty value ${val}` };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `value ${val64} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `value ${val64} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateInt64(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type int64 not allow empty value ${val}` };
    }

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (val < mini) {
        return { valid: false, error: `value ${val} is less than the minimum limit ${mini}` };
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (val > maxi) {
        return { valid: false, error: `value ${val} exceeds the maximum limit ${maxi}` };
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateUint(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type uint not allow empty value ${val}` };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `value ${val} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `value ${val} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateUint8(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type uint8 not allow empty value ${val}` };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `value ${val} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `value ${val} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateUint16(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type uint16 not allow empty value ${val}` };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `value ${val} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `value ${val} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateUint32(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type uint32 not allow empty value ${val}` };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `value ${val} is less than the minimum limit ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `value ${val} exceeds the maximum limit ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateUint64(val: number): ValidationResult {
    if (val === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: `type uint64 not allow empty value ${val}` };
    }

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (val < mini) {
        return { valid: false, error: `value ${val} is less than the minimum limit ${mini}` };
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (val > maxi) {
        return { valid: false, error: `value ${val} exceeds the maximum limit ${maxi}` };
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateFloat32(val: number): ValidationResult {
    if (val === 0.0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0.0' };
      }
      return { valid: false, error: 'type float32 not allow empty value 0.0' };
    }

    const val64 = val;

    if (this.min) {
      const mini = parseFloat(this.min);
      if (!isNaN(mini)) {
        if (val64 < mini) {
          return { valid: false, error: `${val64} < min ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseFloat(this.max);
      if (!isNaN(maxi)) {
        if (val64 > maxi) {
          return { valid: false, error: `${val64} > max ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateFloat64(val: number): ValidationResult {
    if (val === 0.0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0.0' };
      }
      return { valid: false, error: 'type float64 not allow empty value 0.0' };
    }

    if (this.min) {
      const mini = parseFloat(this.min);
      if (!isNaN(mini)) {
        if (val < mini) {
          return { valid: false, error: `${val} < min ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseFloat(this.max);
      if (!isNaN(maxi)) {
        if (val > maxi) {
          return { valid: false, error: `${val} > max ${maxi}` };
        }
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateBigInt(val: bigint): ValidationResult {
    if (val === BigInt(0)) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '0' };
      }
      return { valid: false, error: 'type big.Int not allow empty value 0' };
    }

    if (this.min) {
      const mini = BigInt(this.min);
      if (val < mini) {
        return { valid: false, error: `big.Int ${val} < min ${mini}` };
      }
    }

    if (this.max) {
      const maxi = BigInt(this.max);
      if (val > maxi) {
        return { valid: false, error: `big.Int ${val} > max ${maxi}` };
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val.toString() };
  }

  validateDateTime(val: Date): ValidationResult {
    val = new Date(Math.floor(val.getTime() / 1000) * 1000);
    const format = val.toISOString().replace('T', ' ').substring(0, 19);

    if (val.getTime() === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: format };
      }
      return { valid: false, error: `type datetime not allow empty ${format}` };
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: format };
  }

  validateDate(val: Date): ValidationResult {
    val = new Date(Math.floor(val.getTime() / 1000) * 1000);
    const format = val.toISOString().substring(0, 10);

    if (val.getTime() === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: format };
      }
      return { valid: false, error: `type date not allow empty ${format}` };
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: format };
  }

  validateTime(val: Date): ValidationResult {
    val = new Date(Math.floor(val.getTime() / 1000) * 1000);
    const format = val.toISOString().substring(11, 19);

    if (val.getTime() === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: format };
      }
      return { valid: false, error: `type time not allow empty ${format}` };
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: format };
  }

  validateUUID(val: string): ValidationResult {
    if (val === '') {
      if (this.allowEmpty) {
        return { valid: true, data: new Uint8Array(16), text: val };
      }
      return { valid: false, error: 'type uuid not allow empty value ""' };
    }

    if (!uuidRegex.test(val)) {
      return { valid: false, error: `value '${val}' does not match UUID pattern` };
    }

    const uuidBytes = new Uint8Array(16);
    const parts = val.replace(/-/g, '').match(/.{2}/g);
    if (parts) {
      parts.forEach((part, index) => {
        if (index < 16) {
          uuidBytes[index] = parseInt(part, 16);
        }
      });
    }

    if (this.version) {
      const version = parseInt(val.substring(14, 15), 16);
      if (this.version !== version) {
        return { valid: false, error: 'invalid uuid version' };
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: uuidBytes, text: val };
  }

  validateDecimal(val: string): ValidationResult {
    if (val === '') {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: val };
      }
      return { valid: false, error: 'type decimal not allow empty value ""' };
    }

    if (!decimalRegex.test(val)) {
      return { valid: false, error: `invalid decimal "${val}", must be like "0.0"` };
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val };
  }

  validateIP(val: string): ValidationResult {
    if (val === '') {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '' };
      }
      return { valid: false, error: 'type ip not allow empty value ""' };
    }

    if (this.version === 4) {
      const ipv4Regex = /^(\d{1,3}\.){3}\d{1,3}$/;
      if (!ipv4Regex.test(val)) {
        return { valid: false, error: `invalid ipv4: ${val}` };
      }
    }

    if (this.version === 6) {
      const ipv6Regex = /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::$/;
      if (!ipv6Regex.test(val)) {
        return { valid: false, error: `invalid ipv6: ${val}` };
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val };
  }

  validateURL(val: string): ValidationResult {
    if (val === '') {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: '' };
      }
      return { valid: false, error: 'type url not allow empty value ""' };
    }

    try {
      const url = new URL(val);
      if (url.protocol !== 'http:' && url.protocol !== 'https:') {
        return { valid: false, error: `invalid url: ${val}` };
      }
      if (!url.host) {
        return { valid: false, error: `invalid url: ${val}` };
      }
    } catch {
      return { valid: false, error: `invalid url: ${val}` };
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val };
  }

  validateEmail(val: string): ValidationResult {
    if (val === '') {
      if (this.allowEmpty) {
        return { valid: true, data: val, text: val };
      }
      return { valid: false, error: 'type email not allow empty value ""' };
    }

    if (!emailRegex.test(val)) {
      return { valid: false, error: `value '${val}' does not match email pattern` };
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: val, text: val };
  }

  validateEnum(val: string): ValidationResult {
    if (val === '') {
      if (this.allowEmpty) {
        return { valid: true, data: -1, text: val };
      }
      return { valid: false, error: 'type enum not allow empty value ""' };
    }

    const enums = this.enum.split('|');
    let idx = -1;
    for (let i = 0; i < enums.length; i++) {
      const e = enums[i];
      if (e && e.trim() === val) {
        idx = i;
        break;
      }
    }

    if (idx === -1) {
      return { valid: false, error: `value '${val}' not found in enum: ${enums}` };
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: idx, text: val };
  }

  validateImage(val: Uint8Array | number[]): ValidationResult {
    const arr = val instanceof Uint8Array ? val : new Uint8Array(val);
    const length = arr.length;

    if (length === 0) {
      if (this.allowEmpty) {
        return { valid: true, data: arr, text: '' };
      }
      return { valid: false, error: 'type image not allow empty value []byte{}' };
    }

    if (this.min) {
      const mini = parseInt(this.min, 10);
      if (!isNaN(mini)) {
        if (length < mini) {
          return { valid: false, error: `[]byte length ${length} < min ${mini}` };
        }
      }
    }

    if (this.max) {
      const maxi = parseInt(this.max, 10);
      if (!isNaN(maxi)) {
        if (length > maxi) {
          return { valid: false, error: `[]byte length ${length} > max ${maxi}` };
        }
      }
    }

    if (this.size !== 0n) {
      if (BigInt(length) !== this.size) {
        return { valid: false, error: `[]byte length ${length} != size ${this.size}` };
      }
    }

    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    return { valid: true, data: arr, text: Buffer.from(arr).toString('base64') };
  }

  validateStruct(): ValidationResult {
    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (this.location && this.location !== 0) {
      return { valid: false, error: `type struct not support location UTC${this.location}` };
    }

    return { valid: true };
  }

  validateMap(): ValidationResult {
    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (this.location && this.location !== 0) {
      return { valid: false, error: `type map not support location UTC${this.location}` };
    }

    return { valid: true };
  }

  validateSlice(value: any[]): ValidationResult {
    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (this.location && this.location !== 0) {
      return { valid: false, error: `type slice not support location UTC${this.location}` };
    }

    const length = value.length;

    if (length === 0) {
      if (this.allowEmpty) {
        return { valid: true };
      }
      return { valid: false, error: 'type slice not allow empty' };
    }

    if (this.childUnique) {
      const seen = new Set();
      for (let i = 0; i < value.length; i++) {
        const data = value[i];
        const key = typeof data === 'object' ? JSON.stringify(data) : data;
        if (seen.has(key)) {
          return { valid: false, error: `slice duplicate value found: ${data}, index: ${i}` };
        }
        seen.add(key);
      }
    }

    return { valid: true };
  }

  validateArray(value: any[]): ValidationResult {
    if (this.desc && this.desc.length > 65535) {
      return { valid: false, error: 'desc length exceeds 65535 bytes' };
    }

    if (this.location && this.location !== 0) {
      return { valid: false, error: `type array not support location UTC${this.location}` };
    }

    const length = value.length;

    if (length === 0) {
      if (this.allowEmpty) {
        return { valid: true };
      }
      return { valid: false, error: 'type array not allow empty' };
    }

    if (this.size && length > this.size) {
      return { valid: false, error: 'type array over size' };
    }

    if (this.childUnique) {
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

    return { valid: true };
  }
}

export function parseMMTag(tagStr: string): Tag {
  const tag = new Tag();
  const parts = tagStr.split(';');

  for (const part of parts) {
    const trimmed = part.trim();
    if (!trimmed) continue;
    const kv = trimmed.split('=', 2);
    if (kv.length === 0 || !kv[0]) continue;
    const key = kv[0].toLowerCase();
    const value = (kv.length > 1 && kv[1]) ? kv[1].trim() : '';

    switch (key) {
      case 'is_null':
        tag.isNull = true;
        tag.nullable = true;
        break;
      case 'example':
        tag.example = true;
        break;
      case 'desc':
        tag.desc = value.replace(/^"|"$/g, '');
        break;
      case 'type':
        if (value === 'struct') {
          tag.type = ValueType.Object;
        } else {
          tag.type = stringToType(value);
        }
        break;
      case 'nullable':
        tag.nullable = true;
        break;
      case 'raw':
        tag.raw = true;
        break;
      case 'allow_empty':
        tag.allowEmpty = true;
        break;
      case 'unique':
        tag.unique = true;
        break;
      case 'default':
        tag.default = value;
        break;
      case 'min':
        tag.min = value;
        break;
      case 'max':
        tag.max = value;
        break;
      case 'size':
        tag.size = BigInt(value) || 0n;
        break;
      case 'enum':
        tag.type = ValueType.Enum;
        tag.enum = value;
        break;
      case 'pattern':
        tag.pattern = value;
        break;
      case 'location':
        tag.location = parseInt(value, 10) || 0;
        break;
      case 'version':
        tag.version = parseInt(value, 10) || 0;
        break;
      case 'mime':
        tag.mime = value;
        break;
      case 'child_desc':
        tag.childDesc = value.replace(/^"|"$/g, '');
        break;
      case 'child_type':
        tag.childType = stringToType(value);
        break;
      case 'child_raw':
        tag.childRaw = true;
        break;
      case 'child_nullable':
        tag.childNullable = true;
        break;
      case 'child_allow_empty':
        tag.childAllowEmpty = true;
        break;
      case 'child_unique':
        tag.childUnique = true;
        break;
      case 'child_default':
        tag.childDefault = value;
        break;
      case 'child_min':
        tag.childMin = value;
        break;
      case 'child_max':
        tag.childMax = value;
        break;
      case 'child_size':
        tag.childSize = parseInt(value, 10) || 0;
        break;
      case 'child_enum':
        tag.childEnum = value;
        break;
      case 'child_pattern':
        tag.childPattern = value;
        break;
      case 'child_location':
        tag.childLocation = parseInt(value, 10) || 0;
        break;
      case 'child_version':
        tag.childVersion = parseInt(value, 10) || 0;
        break;
      case 'child_mime':
        tag.childMime = value;
        break;
    }
  }

  return tag;
}
