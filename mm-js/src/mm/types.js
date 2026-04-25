import { ValueType } from './constants.js';

export function mm(value, options = {}) {
  return { value, options };
}

mm.int = (v) => mm(v, { type: ValueType.Int });
mm.str = (v) => mm(v, { type: ValueType.String });
mm.bool = (v) => mm(v, { type: ValueType.Bool });
mm.float = (v) => mm(v, { type: ValueType.Float64 });
mm.bytes = (v) => mm(v, { type: ValueType.Bytes });
mm.array = (v, options = {}) => mm(v, { ...options, type: ValueType.Array });
mm.struct = (v, options = {}) => mm(v, { ...options, type: ValueType.Struct });
mm.bigint = (v) => mm(v, { type: ValueType.BigInt });
mm.uuid = (v) => mm(v, { type: ValueType.Uuid });
mm.datetime = (v) => mm(v, { type: ValueType.DateTime });
mm.date = (v) => mm(v, { type: ValueType.Date });
mm.time = (v) => mm(v, { type: ValueType.Time });
mm.email = (v) => mm(v, { type: ValueType.Email });
mm.url = (v) => mm(v, { type: ValueType.Url });
mm.ip = (v) => mm(v, { type: ValueType.Ip });
mm.decimal = (v) => mm(v, { type: ValueType.Decimal });

mm.i8 = (v) => mm(v, { type: ValueType.Int8 });
mm.i16 = (v) => mm(v, { type: ValueType.Int16 });
mm.i32 = (v) => mm(v, { type: ValueType.Int32 });
mm.i64 = (v) => mm(v, { type: ValueType.Int64 });
mm.u = (v) => mm(v, { type: ValueType.Uint });
mm.u8 = (v) => mm(v, { type: ValueType.Uint8 });
mm.u16 = (v) => mm(v, { type: ValueType.Uint16 });
mm.u32 = (v) => mm(v, { type: ValueType.Uint32 });
mm.u64 = (v) => mm(v, { type: ValueType.Uint64 });
mm.f32 = (v) => mm(v, { type: ValueType.Float32 });
mm.f64 = (v) => mm(v, { type: ValueType.Float64 });

mm.null = () => mm(null, { type: ValueType.Unknown, isNull: true });
mm.nil = () => mm(null, { type: ValueType.Unknown, isNull: true });

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