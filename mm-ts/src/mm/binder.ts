import { MMValue, MMObject, MMArray } from '../ast/ast';
import { Tag } from '../ast/tag';
import { Node } from '../ast/ast';
import { ValueType } from '../ast/value-type';

export type Constructor<T> = new () => T;

export class Binder {
  bind<T>(node: Node, type: Constructor<T> | T): T {
    // const result = new Type();
    const result =
      typeof type === 'function' ? new (type as Constructor<T>)() : type;
    const err = this.bindNode(node, result);
    if (err) {
      throw new Error(err);
    }
    return result;
  }

  bindNode(node: Node, target: any): string | null {
    if (node instanceof MMValue) {
      return this.convertValue(node, target);
    } else if (node instanceof MMObject) {
      const tag = node.getTag();
      if (tag && tag.type === ValueType.Object) {
        return this.convertStruct(node, target);
      } else {
        return this.convertMap(node, target);
      }
    } else if (node instanceof MMArray) {
      const tag = node.getTag();
      if (tag && tag.type === ValueType.Array) {
        return this.convertArray(node, target);
      } else {
        return this.convertSlice(node, target);
      }
    }
    return `unsupported node type: ${typeof node}`;
  }

  private convertStruct(obj: MMObject, target: any): string | null {
    const tag = obj.getTag();
    if (tag && tag.nullable) {
      if (target === null || typeof target !== 'object') {
        return 'convertStruct requires object type, got null or non-object';
      }
    }

    const targetType = typeof target;
    if (targetType !== 'object' || target === null) {
      return `convertStruct requires object type, got ${targetType}`;
    }

    const properties = obj.getProperties();
    for (const [key, valueNode] of Object.entries(properties)) {
      const fieldKey = key;
      const upperKey = fieldKey.charAt(0).toUpperCase() + fieldKey.slice(1);

      let found = false;
      for (const prop in target) {
        if (prop === upperKey || prop === fieldKey) {
          found = true;
          if (!target.hasOwnProperty(prop)) {
            continue;
          }

          const fieldVal = target[prop];
          const err = this.bindNode(valueNode, fieldVal);
          if (err) {
            return `failed to bind struct field ${fieldKey}: ${err}`;
          }
          break;
        }
      }

      if (!found) {
        return `struct has no field '${fieldKey}'`;
      }
    }

    return null;
  }

  private convertMap(obj: MMObject, target: any): string | null {
    const tag = obj.getTag();
    if (tag && tag.nullable) {
      if (target === null || typeof target !== 'object') {
        return 'convertMap requires object type, got null or non-object';
      }
    }

    if (typeof target !== 'object' || target === null) {
      return `convertMap requires object type, got ${typeof target}`;
    }

    if (Array.isArray(target)) {
      return 'convertMap requires object type, got array';
    }

    const properties = obj.getProperties();
    for (const [key, valueNode] of Object.entries(properties)) {
      const fieldKey = key;
      const err = this.bindNode(valueNode, target[fieldKey]);
      if (err) {
        return `failed to convert field ${fieldKey}: ${err}`;
      }
    }

    return null;
  }

  private convertArray(arr: MMArray, target: any): string | null {
    const tag = arr.getTag();
    if (tag && tag.nullable) {
      if (target === null || !Array.isArray(target)) {
        return 'convertArray requires array type, got null or non-array';
      }
    }

    if (!Array.isArray(target)) {
      return `convertArray requires array type, got ${typeof target}`;
    }

    const elements = arr.getElements();
    const arrayLen = target.length;
    const size = tag?.size || elements.length;

    if (size !== arrayLen) {
      return `array length mismatch: target array length ${arrayLen}, got ${size} items`;
    }

    for (let i = 0; i < elements.length; i++) {
      const element = elements[i];
      if (!element) {
        continue;
      }
      const err = this.bindNode(element, target[i]);
      if (err) {
        return `failed to convert array item ${i}: ${err}`;
      }
    }

    return null;
  }

  private convertSlice(arr: MMArray, target: any[]): string | null {
    const tag = arr.getTag();
    if (tag && tag.nullable) {
      if (!Array.isArray(target)) {
        return 'convertSlice requires array type, got non-array';
      }
    }

    if (!Array.isArray(target)) {
      return `convertSlice requires array type, got ${typeof target}`;
    }

    const elements = arr.getElements();
    target.length = 0;

    for (let i = 0; i < elements.length; i++) {
      const item = elements[i];
      if (!item) {
        continue;
      }
      let result: any;

      if (item instanceof MMValue) {
        result = this.createDefaultValue(item.getTag());
      } else if (item instanceof MMObject) {
        result = {};
      } else if (item instanceof MMArray) {
        result = [];
      } else {
        result = null;
      }

      const err = this.bindNode(item, result);
      if (err) {
        return `failed to convert array item ${i}: ${err}`;
      }

      target.push(result);
    }

    return null;
  }

  private createDefaultValue(tag: Tag | null): any {
    if (!tag) return null;

    switch (tag.type) {
      case ValueType.Int:
      case ValueType.Int8:
      case ValueType.Int16:
      case ValueType.Int32:
      case ValueType.Int64:
        return 0;
      case ValueType.Uint:
      case ValueType.Uint8:
      case ValueType.Uint16:
      case ValueType.Uint32:
      case ValueType.Uint64:
        return 0;
      case ValueType.Float32:
      case ValueType.Float64:
        return 0.0;
      case ValueType.String:
      case ValueType.Email:
      case ValueType.URL:
      case ValueType.IP:
      case ValueType.UUID:
      case ValueType.Decimal:
        return '';
      case ValueType.Bool:
        return false;
      case ValueType.Bytes:
        return new Uint8Array();
      case ValueType.BigInt:
        return BigInt(0);
      case ValueType.DateTime:
      case ValueType.Date:
      case ValueType.Time:
        return new Date(0);
      case ValueType.Enum:
        return '';
      default:
        return null;
    }
  }

  private convertValue(val: MMValue, target: any): string | null {
    const tag = val.getTag();
    if (!tag) {
      return 'tag is null';
    }

    const data = val.getValue();

    if (tag.nullable) {
      if (target === null) {
        return null;
      }
    }

    switch (tag.type) {
      case ValueType.Int:
      case ValueType.Int8:
      case ValueType.Int16:
      case ValueType.Int32:
      case ValueType.Int64:
      case ValueType.Uint:
      case ValueType.Uint8:
      case ValueType.Uint16:
      case ValueType.Uint32:
      case ValueType.Uint64:
        if (typeof target !== 'number' && typeof target !== 'bigint') {
          return `target type must be number or bigint, got ${typeof target}`;
        }
        if (typeof data === 'bigint') {
          if (typeof target === 'number') {
            target = Number(data);
          } else {
            target = data;
          }
        } else {
          target = data;
        }
        break;

      case ValueType.Float32:
      case ValueType.Float64:
        if (typeof target !== 'number') {
          return `target type must be number, got ${typeof target}`;
        }
        target = typeof data === 'number' ? data : Number(data);
        break;

      case ValueType.String:
      case ValueType.Email:
      case ValueType.URL:
      case ValueType.IP:
      case ValueType.Decimal:
        if (typeof target !== 'string') {
          return `target type must be string, got ${typeof target}`;
        }
        target = String(data);
        break;

      case ValueType.Bool:
        if (typeof target !== 'boolean') {
          return `target type must be boolean, got ${typeof target}`;
        }
        target = Boolean(data);
        break;

      case ValueType.Bytes:
        if (!(target instanceof Uint8Array)) {
          return `target type must be Uint8Array, got ${typeof target}`;
        }
        target = data instanceof Uint8Array ? data : new Uint8Array(data);
        break;

      case ValueType.BigInt:
        if (typeof target !== 'bigint') {
          return `target type must be bigint, got ${typeof target}`;
        }
        target = typeof data === 'bigint' ? data : BigInt(data);
        break;

      case ValueType.DateTime:
      case ValueType.Date:
      case ValueType.Time:
        if (!(target instanceof Date)) {
          return `target type must be Date, got ${typeof target}`;
        }
        target = data instanceof Date ? data : new Date(data);
        break;

      case ValueType.UUID:
        if (typeof target !== 'string') {
          return `target type must be string, got ${typeof target}`;
        }
        if (data instanceof Uint8Array) {
          const uuidStr = Array.from(data)
            .map((b) => b.toString(16).padStart(2, '0'))
            .join('');
          target = `${uuidStr.slice(0, 8)}-${uuidStr.slice(8, 12)}-${uuidStr.slice(12, 16)}-${uuidStr.slice(16, 20)}-${uuidStr.slice(20)}`;
        } else {
          target = String(data);
        }
        break;

      case ValueType.Enum:
        if (typeof target !== 'string') {
          return `target type must be string, got ${typeof target}`;
        }
        target = String(data);
        break;

      case ValueType.Image:
      case ValueType.Video:
      case ValueType.Doc:
        return `unsupported type: ${tag.type}`;

      case ValueType.Map:
      case ValueType.Object:
      case ValueType.Array:
      case ValueType.Slice:
      case ValueType.Unknown:
      default:
        return `unsupported type: ${tag.type}`;
    }

    return null;
  }
}

export function bindJSONC<T>(node: Node, Type: new () => T): T {
  const binder = new Binder();
  return binder.bind(node, Type);
}
