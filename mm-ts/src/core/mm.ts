import { Tag } from '../ir/tag';
import { typeToString, ValueType } from '../ir/value-type';
import { MMArray, MMObject, MMValue, Node } from '../ir/ast';

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
  new MMValue(v, { ...tag, type: ValueType.Str } as Tag);
mm.bool = (v: boolean, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Bool } as Tag);
mm.bytes = (v: Uint8Array, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Bytes } as Tag);
mm.arr = <T extends Node>(v: T[], size?: bigint, tag?: Tag) => {
  const arr = new MMArray();
  arr.setTag(
    Object.assign({}, tag, {
      type: ValueType.Arr,
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
      type: ValueType.Vec,
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
      type: ValueType.Obj,
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
  new MMValue(v, { ...tag, type: ValueType.Bigint } as Tag);
mm.uuid = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Uuid } as Tag);
mm.datetime = (v: Date, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Datetime } as Tag);
mm.date = (v: Date, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Date } as Tag);
mm.time = (v: Date, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Time } as Tag);
mm.email = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Email } as Tag);
mm.url = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Url } as Tag);
mm.ip = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Ip } as Tag);
mm.decimal = (v: string, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Decimal } as Tag);
mm.enum = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Enums } as Tag);
mm.image = (v: Uint8Array, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Image } as Tag);
mm.video = (v: Uint8Array, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.Video } as Tag);
mm.i = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.I } as Tag);
mm.i8 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.I8 } as Tag);
mm.i16 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.I16 } as Tag);
mm.i32 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.I32 } as Tag);
mm.i64 = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.I64 } as Tag);
mm.u = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.U } as Tag);
mm.u8 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.U8 } as Tag);
mm.u16 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.U16 } as Tag);
mm.u32 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.U32 } as Tag);
mm.u64 = (v: bigint, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.U64 } as Tag);
mm.f32 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.F32 } as Tag);
mm.f64 = (v: number, tag?: Tag) =>
  new MMValue(v, { ...tag, type: ValueType.F64 } as Tag);
