// ValueType aligns with Go's ValueType enum
export type ValueType =
  | 'unknown'
  | 'doc'
  | 'array'
  | 'slice'
  | 'struct'
  | 'map'
  | 'string'
  | 'bytes'
  | 'bool'
  | 'int'
  | 'int8'
  | 'int16'
  | 'int32'
  | 'int64'
  | 'uint'
  | 'uint8'
  | 'uint16'
  | 'uint32'
  | 'uint64'
  | 'float32'
  | 'float64'
  | 'bigint'
  | 'datetime'
  | 'date'
  | 'time'
  | 'uuid'
  | 'decimal'
  | 'ip'
  | 'url'
  | 'email'
  | 'enum'
  | 'image'
  | 'video';

export interface Tag {
  name?: string;
  isNull?: boolean;
  example?: boolean;
  desc?: string;
  type?: ValueType;
  raw?: boolean;
  nullable?: boolean;
  allowEmpty?: boolean;
  unique?: boolean;
  default?: string;
  min?: string;
  max?: string;
  size?: number;
  enum?: string;
  pattern?: string;
  location?: number;
  version?: number;
  mime?: string;
  childDesc?: string;
  childType?: ValueType;
  childRaw?: boolean;
  childNullable?: boolean;
  childAllowEmpty?: boolean;
  childUnique?: boolean;
  childDefault?: string;
  childMin?: string;
  childMax?: string;
  childSize?: number;
  childEnum?: string;
  childPattern?: string;
  childLocation?: number;
  childVersion?: number;
  childMime?: string;
  isInherit?: boolean;
}

export interface MMValue<T = any> {
  value: T;
  options: Tag;
}

export interface MM {
  <T>(value: T, options?: Tag): MMValue<T>;
  int: (v: number, options?: Tag) => MMValue<number>;
  str: (v: string, options?: Tag) => MMValue<string>;
  bool: (v: boolean, options?: Tag) => MMValue<boolean>;
  float: (v: number, options?: Tag) => MMValue<number>;
  bytes: (v: Uint8Array | number[], options?: Tag) => MMValue<Uint8Array | number[]>;
  array: <T>(v: T[], options?: Tag) => MMValue<T[]>;
  struct: <T>(v: T, options?: Tag) => MMValue<T>;
  bigint: (v: bigint, options?: Tag) => MMValue<bigint>;
  uuid: (v: string, options?: Tag) => MMValue<string>;
  datetime: (v: Date, options?: Tag) => MMValue<Date>;
  date: (v: Date, options?: Tag) => MMValue<Date>;
  time: (v: Date, options?: Tag) => MMValue<Date>;
  email: (v: string, options?: Tag) => MMValue<string>;
  url: (v: string, options?: Tag) => MMValue<string>;
  ip: (v: string, options?: Tag) => MMValue<string>;
  decimal: (v: string, options?: Tag) => MMValue<string>;
  i8: (v: number, options?: Tag) => MMValue<number>;
  i16: (v: number, options?: Tag) => MMValue<number>;
  i32: (v: number, options?: Tag) => MMValue<number>;
  i64: (v: bigint, options?: Tag) => MMValue<bigint>;
  u: (v: number, options?: Tag) => MMValue<number>;
  u8: (v: number, options?: Tag) => MMValue<number>;
  u16: (v: number, options?: Tag) => MMValue<number>;
  u32: (v: number, options?: Tag) => MMValue<number>;
  u64: (v: bigint, options?: Tag) => MMValue<bigint>;
  f32: (v: number, options?: Tag) => MMValue<number>;
  f64: (v: number, options?: Tag) => MMValue<number>;
}

function mmFunction<T>(value: T, options: Tag = {}): MMValue<T> {
  // Auto type inference
  if (!options.type) {
    if (typeof value === 'string') {
      options.type = 'string';
    } else if (typeof value === 'number' && Number.isInteger(value)) {
      options.type = 'int';
    } else if (typeof value === 'boolean') {
      options.type = 'bool';
    } else if (typeof value === 'number' && !Number.isInteger(value)) {
      options.type = 'float64';
    } else if (value === null) {
      // For null values, don't set a specific type
    } else if (Array.isArray(value)) {
      options.type = 'array';
    } else if (value instanceof Uint8Array) {
      options.type = 'bytes';
    } else if (typeof value === 'object' && value !== null) {
      options.type = 'struct';
    }
  }
  return { value, options };
}

const mmInstance: MM = Object.assign(mmFunction, {
  int: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'int' }),
  str: (v: string, options?: Tag) => mmFunction(v, { ...options, type: 'string' }),
  bool: (v: boolean, options?: Tag) => mmFunction(v, { ...options, type: 'bool' }),
  float: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'float64' }),
  bytes: (v: Uint8Array | number[], options?: Tag) => mmFunction(v, { ...options, type: 'bytes' }),
  array: <T>(v: T[], options?: Tag) => mmFunction(v, { ...options, type: 'array' }),
  struct: <T>(v: T, options?: Tag) => mmFunction(v, { ...options, type: 'struct' }),
  bigint: (v: bigint, options?: Tag) => mmFunction(v, { ...options, type: 'bigint' }),
  uuid: (v: string, options?: Tag) => mmFunction(v, { ...options, type: 'uuid' }),
  datetime: (v: Date, options?: Tag) => mmFunction(v, { ...options, type: 'datetime' }),
  date: (v: Date, options?: Tag) => mmFunction(v, { ...options, type: 'date' }),
  time: (v: Date, options?: Tag) => mmFunction(v, { ...options, type: 'time' }),
  email: (v: string, options?: Tag) => mmFunction(v, { ...options, type: 'email' }),
  url: (v: string, options?: Tag) => mmFunction(v, { ...options, type: 'url' }),
  ip: (v: string, options?: Tag) => mmFunction(v, { ...options, type: 'ip' }),
  decimal: (v: string, options?: Tag) => mmFunction(v, { ...options, type: 'decimal' }),
  i8: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'int8' }),
  i16: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'int16' }),
  i32: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'int32' }),
  i64: (v: bigint, options?: Tag) => mmFunction(v, { ...options, type: 'int64' }),
  u: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'uint' }),
  u8: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'uint8' }),
  u16: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'uint16' }),
  u32: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'uint32' }),
  u64: (v: bigint, options?: Tag) => mmFunction(v, { ...options, type: 'uint64' }),
  f32: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'float32' }),
  f64: (v: number, options?: Tag) => mmFunction(v, { ...options, type: 'float64' }),
});

export const mm = mmInstance;