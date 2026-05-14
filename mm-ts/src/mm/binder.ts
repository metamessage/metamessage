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
    // console.log('result', result);
    const { error } = this.bindNode(node, result);
    if (error) {
      throw new Error(error);
    }
    return result;
  }

  bindNode(node: Node, target: any): { value: any; error: string | null } {
    if (node instanceof MMValue) {
      return this.convertValue(node, target);
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

      // 🔥 拿到子节点的值，直接赋值！
      const { value, error } = this.bindNode(valueNode, target[camelKey]);
      if (error) {
        return {
          value: target,
          error: `failed to bind field ${camelKey}: ${error}`,
        };
      }

      // ✅ 赋值生效！解决指针问题！
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

  private convertValue(
    val: MMValue,
    target: any,
  ): { value: any; error: string | null } {
    const tag = val.getTag();
    if (!tag) return { value: target, error: 'tag is null' };

    const data = val.getValue();

    // 直接返回新值，不赋值！
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
        return typeof data === 'bigint'
          ? { value: Number(data), error: null }
          : { value: data, error: null };

      case ValueType.Float32:
      case ValueType.Float64:
        return { value: Number(data), error: null };

      case ValueType.String:
      case ValueType.Email:
      case ValueType.URL:
      case ValueType.IP:
      case ValueType.Decimal:
      case ValueType.Enum:
        return { value: String(data), error: null };

      case ValueType.Bool:
        return { value: Boolean(data), error: null };

      case ValueType.Bytes:
        return { value: new Uint8Array(data), error: null };

      case ValueType.BigInt:
        return { value: BigInt(data), error: null };

      case ValueType.DateTime:
      case ValueType.Date:
      case ValueType.Time:
        return { value: new Date(data), error: null };

      case ValueType.UUID:
        if (data instanceof Uint8Array) {
          const s = Array.from(data)
            .map((b) => b.toString(16).padStart(2, '0'))
            .join('');
          const uuid = `${s.slice(0, 8)}-${s.slice(8, 12)}-${s.slice(12, 16)}-${s.slice(16, 20)}-${s.slice(20)}`;
          return { value: uuid, error: null };
        }
        return { value: String(data), error: null };
    }

    return { value: target, error: `unsupported value type: ${tag.type}` };
  }
}

export function bindJSONC<T>(node: Node, Type: new () => T): T {
  const binder = new Binder();
  return binder.bind(node, Type);
}
