import { Tag } from '../ast/tag';
import { typeToString, ValueType } from '../ast/value-type';
import { MMArray, MMObject, MMValue, Node } from '../ast/ast';

export const META_KEY = Symbol('mm_meta');

export interface Options extends Partial<Tag> {}

export function mm(options: Options) {
  return function (target: any, context: any) {
    if (context.kind === 'class') {
      const cls = target;
      if (!cls[META_KEY]) cls[META_KEY] = {};
      cls[META_KEY].__class = options;
      return;
    }

    if (context.kind === 'field') {
      context.addInitializer(function (this: any) {
        const cls = this.constructor;
        cls[META_KEY] = cls[META_KEY] || {};
        cls[META_KEY][context.name] = options;
      });
    }
  };
}

export function toMM(inst: any) {
  const meta = inst.constructor[META_KEY] || {};

  const result: any = {};

  result.class = meta.__class || {};

  for (const key of Object.keys(meta).filter((k) => k !== '__class')) {
    result[key] = {
      ...meta[key],
      value: inst[key],
    };
  }

  return result;
}

mm.str = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.String } as Tag);
mm.bool = (v: boolean, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Bool } as Tag);
mm.bytes = (v: Uint8Array, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Bytes } as Tag);
mm.arr = <T extends Node>(v: T[], size?: bigint, tag?: Tag) => {
  const arr = new MMArray();
  arr.setTag(
    Object.assign({}, tag, {
      type: ValueType.Array,
      size: size ?? v.length,
    }),
  );

  if (v.length > 0) {
    const firstType = v[0]!.getTag().type;

    for (const item of v) {
      const currentType = item.getTag().type;
      if (currentType !== firstType) {
        throw new Error(
          `Array element type mismatch: expected ${typeToString(firstType)}, got ${typeToString(currentType)}`,
        );
      }
      arr.addElement(item);
    }
  }

  return arr;
};

mm.slice = <T extends Node>(v: T[], tag?: Tag) => {
  const arr = new MMArray();
  arr.setTag(
    Object.assign({}, tag, {
      type: ValueType.Slice,
    }),
  );

  for (const item of v) {
    arr.addElement(item);
  }

  return arr;
};

mm.obj = <T extends Node>(v: Record<string, T>, tag?: Tag) => {
  const obj = new MMObject();
  obj.setTag(
    Object.assign({}, tag, {
      type: ValueType.Object,
    }),
  );

  for (const [key, value] of Object.entries(v)) {
    obj.setProperty(key, value);
  }

  return obj;
};

mm.map = <T extends Node>(v: Record<string, T>, tag?: Tag) => {
  const obj = new MMObject();
  obj.setTag(
    Object.assign({}, tag, {
      type: ValueType.Map,
    }),
  );

  for (const [key, value] of Object.entries(v)) {
    obj.setProperty(key, value);
  }

  return obj;
};

mm.doc = <T>(v: T, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Doc } as Tag);
mm.bigint = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.BigInt } as Tag);
mm.uuid = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.UUID } as Tag);
mm.datetime = (v: Date, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.DateTime } as Tag);
mm.date = (v: Date, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Date } as Tag);
mm.time = (v: Date, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Time } as Tag);
mm.email = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Email } as Tag);
mm.url = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.URL } as Tag);
mm.ip = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.IP } as Tag);
mm.decimal = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Decimal } as Tag);
mm.enum = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Enum } as Tag);
mm.image = (v: Uint8Array, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Image } as Tag);
mm.video = (v: Uint8Array, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Video } as Tag);
mm.i = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Int } as Tag);
mm.i8 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Int8 } as Tag);
mm.i16 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Int16 } as Tag);
mm.i32 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Int32 } as Tag);
mm.i64 = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Int64 } as Tag);
mm.u = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Uint } as Tag);
mm.u8 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Uint8 } as Tag);
mm.u16 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Uint16 } as Tag);
mm.u32 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Uint32 } as Tag);
mm.u64 = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Uint64 } as Tag);
mm.f32 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Float32 } as Tag);
mm.f64 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Float64 } as Tag);
