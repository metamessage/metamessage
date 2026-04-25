export interface MMOptions {
  type?: string;
  desc?: string;
  nullable?: boolean;
  default?: any;
  min?: string;
  max?: string;
  size?: number;
  enum?: string;
  pattern?: string;
  location?: number;
  version?: number;
  mime?: string;
  [key: string]: any;
}

export interface MMValue<T = any> {
  value: T;
  options: MMOptions;
}

export interface MM {
  <T>(value: T, options?: MMOptions): MMValue<T>;
  int: (v: number | bigint) => MMValue<number | bigint>;
  str: (v: string) => MMValue<string>;
  bool: (v: boolean) => MMValue<boolean>;
  float: (v: number) => MMValue<number>;
  bytes: (v: Uint8Array | number[]) => MMValue<Uint8Array | number[]>;
  array: <T>(v: T[], options?: MMOptions) => MMValue<T[]>;
  struct: <T>(v: T, options?: MMOptions) => MMValue<T>;
  bigint: (v: bigint) => MMValue<bigint>;
  null: (type?: string) => MMValue<null>;
}

function mmFunction<T>(value: T, options: MMOptions = {}): MMValue<T> {
  return { value, options };
}

const mmInstance: MM = Object.assign(mmFunction, {
  int: (v: number | bigint) => mmFunction(v, { type: 'int' }),
  str: (v: string) => mmFunction(v, { type: 'str' }),
  bool: (v: boolean) => mmFunction(v, { type: 'bool' }),
  float: (v: number) => mmFunction(v, { type: 'float' }),
  bytes: (v: Uint8Array | number[]) => mmFunction(v, { type: 'bytes' }),
  array: <T>(v: T[], options?: MMOptions) => mmFunction(v, { ...options, type: 'array' }),
  struct: <T>(v: T, options?: MMOptions) => mmFunction(v, { ...options, type: 'struct' }),
  bigint: (v: bigint) => mmFunction(v, { type: 'bigint' }),
  null: (type?: string) => mmFunction(null, type ? { type } : {}),
});

export const mm = mmInstance;