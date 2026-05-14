import { Node, MMValue, MMObject, MMArray } from '../ast/ast';
import { Tag } from '../ast/tag';
import { typeToString, ValueType } from '../ast/value-type';
import { META_KEY } from './mm';

const maxDepth = 32;

export function NilToNode(valueType: ValueType): MMValue {
  const tag = new Tag();
  tag.type = valueType;
  tag.isNull = true;
  return new MMValue(null, tag);
}

export function ValueToNode(v: any, tag?: Tag): Node {
  if (!tag) {
    tag = new Tag();
  }
  return valueToNode(v, tag, 0, '');
}

function valueToNode(v: any, tag: Tag, depth: number, path: string): Node {
  const meta = v?.constructor?.[META_KEY] ?? {};
  // console.log('meta', meta, tag);
  if (meta) {
    Object.assign(tag, meta.__class);
  }

  let data: any = null;
  let text: string = '';

  if (v === null || v === undefined) {
    if (tag.type === ValueType.Unknown) {
      throw new Error(
        `invalid input: v is untyped nil (no concrete type/value)`,
      );
    }
    tag.isNull = true;
    return createValueNode(data, text, tag, path);
  }

  if (
    typeof v === 'object' &&
    'value' in v &&
    'tag' in v &&
    typeof v.tag === 'object'
  ) {
    const mmValue = v;
    const mergedTag = new Tag();
    Object.assign(mergedTag, mmValue.tag);
    if (
      tag.type !== ValueType.Unknown &&
      mergedTag.type === ValueType.Unknown
    ) {
      mergedTag.type = tag.type;
    }
    return valueToNode(mmValue.value, mergedTag, depth, path);
  }

  switch (typeof v) {
    case 'boolean':
      if (tag.type === ValueType.Unknown) {
        tag.type = ValueType.Bool;
      }
      switch (tag.type) {
        case ValueType.Bool: {
          const result = tag.validateBool(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        default:
          throw new Error(`${tag.type} unsupported type: ${typeof v}`);
      }
      break;

    case 'number':
      if (tag.type === ValueType.Unknown) {
        tag.type = Number.isInteger(v) ? ValueType.Int : ValueType.Float64;
      }
      switch (tag.type) {
        case ValueType.Int: {
          const result = tag.validateInt8(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }

        case ValueType.Int8: {
          const result = tag.validateInt8(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Int16: {
          const result = tag.validateInt16(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Int32: {
          const result = tag.validateInt32(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Int64: {
          const result = tag.validateInt64(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Uint: {
          const result = tag.validateUint(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Uint8: {
          const result = tag.validateUint8(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Uint16: {
          const result = tag.validateUint16(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Uint32: {
          const result = tag.validateUint32(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Uint64: {
          const result = tag.validateUint64(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Float32: {
          if (Number.isNaN(v)) {
            throw new Error(`${tag.type} unsupported value: NaN`);
          }
          if (!isFinite(v)) {
            throw new Error(
              `${tag.type} unsupported value: ${v > 0 ? '+Inf' : '-Inf'}`,
            );
          }
          const result = tag.validateFloat32(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Float64: {
          if (Number.isNaN(v)) {
            throw new Error(`${tag.type} unsupported value: NaN`);
          }
          if (!isFinite(v)) {
            throw new Error(
              `${tag.type} unsupported value: ${v > 0 ? '+Inf' : '-Inf'}`,
            );
          }
          const result = tag.validateFloat64(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        default:
          throw new Error(
            `${typeToString(tag.type)} unsupported type: ${typeof v}`,
          );
      }
      break;

    case 'bigint':
      if (tag.type === ValueType.Unknown) {
        tag.type = ValueType.BigInt;
      }
      switch (tag.type) {
        case ValueType.BigInt:
        case ValueType.Int: {
          const result = tag.validateInt(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Int64:
        case ValueType.Uint64: {
          const result = tag.validateBigInt(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        default:
          throw new Error(
            `${typeToString(tag.type)} unsupported type: ${typeof v}`,
          );
      }
      break;

    case 'string':
      if (tag.type === ValueType.Unknown) {
        tag.type = ValueType.String;
      }
      switch (tag.type) {
        case ValueType.String: {
          const result = tag.validateString(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Decimal: {
          const result = tag.validateDecimal(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Email: {
          const result = tag.validateEmail(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.Enum: {
          if (!tag.enum) {
            throw new Error('enum empty');
          }
          const result = tag.validateEnum(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.UUID: {
          const result = tag.validateUUID(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.URL: {
          const result = tag.validateURL(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        case ValueType.IP: {
          const result = tag.validateIP(v);
          if (!result.valid) {
            throw new Error(`validate failed: ${result.error}`);
          }
          data = result.data;
          text = result.text || '';
          break;
        }
        default:
          throw new Error(`${tag.type} unsupported type: ${typeof v}`);
      }
      break;

    default:
      if (v instanceof Uint8Array) {
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.Bytes;
        }
        switch (tag.type) {
          case ValueType.Bytes: {
            const result = tag.validateBytes(v);
            if (!result.valid) {
              throw new Error(`validate failed: ${result.error}`);
            }
            data = result.data;
            text = result.text || '';
            break;
          }
          case ValueType.Image: {
            const result = tag.validateImage(v);
            if (!result.valid) {
              throw new Error(`validate failed: ${result.error}`);
            }
            data = result.data;
            text = result.text || '';
            break;
          }
          default:
            throw new Error(`${tag.type} unsupported type: Uint8Array`);
        }
      } else if (v instanceof Date) {
        if (tag.type === ValueType.Unknown) {
          tag.type = ValueType.DateTime;
        }
        switch (tag.type) {
          case ValueType.DateTime: {
            const result = tag.validateDateTime(v);
            if (!result.valid) {
              throw new Error(`validate failed: ${result.error}`);
            }
            data = result.data;
            text = result.text || '';
            break;
          }
          case ValueType.Date: {
            const result = tag.validateDate(v);
            if (!result.valid) {
              throw new Error(`validate failed: ${result.error}`);
            }
            data = result.data;
            text = result.text || '';
            break;
          }
          case ValueType.Time: {
            const result = tag.validateTime(v);
            if (!result.valid) {
              throw new Error(`validate failed: ${result.error}`);
            }
            data = result.data;
            text = result.text || '';
            break;
          }
          default:
            throw new Error(`${tag.type} unsupported type: Date`);
        }
      } else if (Array.isArray(v)) {
        return anyToArray(v, tag, depth, path);
      } else if (typeof v === 'object') {
        return anyToObject(v, tag, depth, path, meta);
      } else {
        throw new Error(`unsupported type: ${typeof v}`);
      }
  }

  return createValueNode(data, text, tag, path);
}

function createValueNode(
  data: any,
  text: string,
  tag: Tag,
  path: string,
): MMValue {
  const node = new MMValue(data, tag);
  node.setPath(path);
  return node;
}

function anyToObject(
  obj: Record<string, any>,
  tag: Tag,
  depth: number,
  path: string,
  meta: Record<string, Tag>,
): MMObject {
  depth++;
  if (depth > maxDepth) {
    throw new Error(`max depth: ${maxDepth}`);
  }

  tag.type = ValueType.Object;

  const objNode = new MMObject();
  objNode.setTag(tag);
  objNode.setPath(path);

  for (const [key, val] of Object.entries(obj)) {
    const fieldTag = new Tag();
    fieldTag.name = key;
    if (meta[key]) {
      Object.assign(fieldTag, meta[key]);
    }

    const fieldPath = path ? `${path}.${key}` : key;
    const fieldNode = valueToNode(val, fieldTag, depth, fieldPath);
    objNode.setProperty(key, fieldNode);
  }

  const result = tag.validateStruct();
  if (!result.valid) {
    throw new Error(`validate failed: ${result.error}`);
  }

  return objNode;
}

function anyToArray(
  arr: any[],
  tag: Tag,
  depth: number,
  path: string,
): MMArray {
  depth++;
  if (depth > maxDepth) {
    throw new Error(`max depth: ${maxDepth}`);
  }

  tag.type = ValueType.Slice;

  const arrNode = new MMArray();
  arrNode.setTag(tag);
  arrNode.setPath(path);

  for (let i = 0; i < arr.length; i++) {
    const itemTag = new Tag();
    const itemPath = `${path}[${i}]`;
    const itemNode = valueToNode(arr[i], itemTag, depth, itemPath);
    arrNode.addElement(itemNode);
  }

  const result = tag.validateSlice(arr);
  if (!result.valid) {
    throw new Error(`validate failed: ${result.error}`);
  }

  return arrNode;
}
