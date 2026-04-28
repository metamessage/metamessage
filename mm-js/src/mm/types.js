import { ValueType } from './constants.js';

export function mm(value, options = {}) {
  return { value, options };
}

mm.int = (v, options = {}) => mm(v, { ...options, type: ValueType.Int });
mm.str = (v, options = {}) => mm(v, { ...options, type: ValueType.String });
mm.bool = (v, options = {}) => mm(v, { ...options, type: ValueType.Bool });
mm.float = (v, options = {}) => mm(v, { ...options, type: ValueType.Float64 });
mm.bytes = (v, options = {}) => mm(v, { ...options, type: ValueType.Bytes });
mm.array = (v, options = {}) => mm(v, { ...options, type: ValueType.Array });
mm.struct = (v, options = {}) => mm(v, { ...options, type: ValueType.Struct });
mm.bigint = (v, options = {}) => mm(v, { ...options, type: ValueType.BigInt });
mm.uuid = (v, options = {}) => mm(v, { ...options, type: ValueType.Uuid });
mm.datetime = (v, options = {}) => mm(v, { ...options, type: ValueType.DateTime });
mm.date = (v, options = {}) => mm(v, { ...options, type: ValueType.Date });
mm.time = (v, options = {}) => mm(v, { ...options, type: ValueType.Time });
mm.email = (v, options = {}) => mm(v, { ...options, type: ValueType.Email });
mm.url = (v, options = {}) => mm(v, { ...options, type: ValueType.Url });
mm.ip = (v, options = {}) => mm(v, { ...options, type: ValueType.Ip });
mm.decimal = (v, options = {}) => mm(v, { ...options, type: ValueType.Decimal });

mm.i8 = (v, options = {}) => mm(v, { ...options, type: ValueType.Int8 });
mm.i16 = (v, options = {}) => mm(v, { ...options, type: ValueType.Int16 });
mm.i32 = (v, options = {}) => mm(v, { ...options, type: ValueType.Int32 });
mm.i64 = (v, options = {}) => mm(v, { ...options, type: ValueType.Int64 });
mm.u = (v, options = {}) => mm(v, { ...options, type: ValueType.Uint });
mm.u8 = (v, options = {}) => mm(v, { ...options, type: ValueType.Uint8 });
mm.u16 = (v, options = {}) => mm(v, { ...options, type: ValueType.Uint16 });
mm.u32 = (v, options = {}) => mm(v, { ...options, type: ValueType.Uint32 });
mm.u64 = (v, options = {}) => mm(v, { ...options, type: ValueType.Uint64 });
mm.f32 = (v, options = {}) => mm(v, { ...options, type: ValueType.Float32 });
mm.f64 = (v, options = {}) => mm(v, { ...options, type: ValueType.Float64 });

export function isMM(obj) {
  return obj && typeof obj === 'object' && 'value' in obj && 'options' in obj;
}

export function getMMType(maybeMM) {
  if (isMM(maybeMM)) {
    return maybeMM.options?.type || ValueType.Unknown;
  }
  return ValueType.Unknown;
}

export function getMMValue(maybeMM) {
  if (isMM(maybeMM)) {
    return maybeMM.value;
  }
  return maybeMM;
}