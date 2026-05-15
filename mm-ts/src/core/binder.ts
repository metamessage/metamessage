import { MMValue, MMObject, MMArray } from '../ir/ast';
import { Tag } from '../ir/tag';
import { Node } from '../ir/ast';
import { ValueType } from '../ir/value-type';

export type Constructor<T> = new () => T;

export class Binder {
  bind<T>(node: Node, type: Constructor<T> | T): T {
    const result =
      typeof type === 'function' ? new (type as Constructor<T>)() : type;
    const { error } = this.bindNode(node, result);
    if (error) {
      throw new Error(error);
    }
    return result;
  }

  bindNode(node: Node, target: any): { value: any; error: string | null } {
    if (node instanceof MMValue) {
      return { value: node.getValue(), error: null };
    }

    if (node instanceof MMObject) {
      const tag = node.getTag();
      if (tag?.type === ValueType.Object) {
        return this.convertStruct(node, target);
      } else {
        return this.convertMap(node, target);
      }
    }

    if (node instanceof MMArray) {
      const tag = node.getTag();
      if (tag?.type === ValueType.Array) {
        return this.convertArray(node, target);
      } else {
        return this.convertSlice(node, target);
      }
    }

    return { value: target, error: `unsupported node type: ${typeof node}` };
  }

  private convertStruct(
    obj: MMObject,
    target: any,
  ): { value: any; error: string | null } {
    const tag = obj.getTag();
    if (tag?.nullable && target === null) {
      return {
        value: target,
        error: 'convertStruct requires object type, got null',
      };
    }
    if (typeof target !== 'object' || target === null) {
      return {
        value: target,
        error: `convertStruct requires object type, got ${typeof target}`,
      };
    }

    const properties = obj.getProperties();
    for (const [key, valueNode] of Object.entries(properties)) {
      const camelKey = key.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
      if (!(camelKey in target)) {
        return { value: target, error: `struct has no field '${camelKey}'` };
      }

      const { value, error } = this.bindNode(valueNode, target[camelKey]);
      if (error) {
        return {
          value: target,
          error: `failed to bind field ${camelKey}: ${error}`,
        };
      }

      target[camelKey] = value;
    }

    return { value: target, error: null };
  }

  private convertMap(
    obj: MMObject,
    target: any,
  ): { value: any; error: string | null } {
    if (
      typeof target !== 'object' ||
      target === null ||
      Array.isArray(target)
    ) {
      return { value: target, error: `convertMap requires object` };
    }

    const properties = obj.getProperties();
    for (const [key, valueNode] of Object.entries(properties)) {
      const { value, error } = this.bindNode(valueNode, target[key]);
      if (error) {
        return { value: target, error: `field ${key}: ${error}` };
      }
      target[key] = value;
    }
    return { value: target, error: null };
  }

  private convertArray(
    arr: MMArray,
    target: any,
  ): { value: any; error: string | null } {
    if (!Array.isArray(target)) {
      return { value: target, error: `convertArray requires array` };
    }

    const elements = arr.getElements();
    const size = arr.getTag()?.size || elements.length;

    if (target.length !== size) {
      return {
        value: target,
        error: `array length mismatch: expected ${size}, got ${target.length}`,
      };
    }

    for (let i = 0; i < elements.length; i++) {
      const el = elements[i];
      if (!el) continue;

      const { value, error } = this.bindNode(el, target[i]);
      if (error) {
        return { value: target, error: `array[${i}]: ${error}` };
      }
      target[i] = value;
    }

    return { value: target, error: null };
  }

  private convertSlice(
    arr: MMArray,
    target: any[],
  ): { value: any; error: string | null } {
    if (!Array.isArray(target)) {
      return { value: target, error: `convertSlice requires array` };
    }

    const elements = arr.getElements();
    target.length = 0;

    for (const el of elements) {
      const defaultValue = this.createDefaultValue(el.getTag());
      const { value, error } = this.bindNode(el, defaultValue);
      if (error) {
        return { value: target, error };
      }
      target.push(value);
    }

    return { value: target, error: null };
  }

  private createDefaultValue(tag: Tag | null): any {
    if (!tag) return null;

    switch (tag.type) {
      case ValueType.Int:
      case ValueType.Uint:
      case ValueType.Int64:
      case ValueType.Uint64:
      case ValueType.BigInt:
        return 0n;

      case ValueType.Int8:
      case ValueType.Int16:
      case ValueType.Int32:
      case ValueType.Uint8:
      case ValueType.Uint16:
      case ValueType.Uint32:
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
      case ValueType.Enum:
        return '';

      case ValueType.Decimal:
        return '0.0';

      case ValueType.Bool:
        return false;

      case ValueType.Bytes:
      case ValueType.Image:
      case ValueType.Video:
        return new Uint8Array();

      case ValueType.DateTime:
      case ValueType.Date:
      case ValueType.Time:
        return new Date(0);

      default:
        return null;
    }
  }
}

export function bindJSONC<T>(node: Node, Type: new () => T): T {
  const binder = new Binder();
  return binder.bind(node, Type);
}
