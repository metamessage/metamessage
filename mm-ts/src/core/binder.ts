import { NodeScalar, NodeObject, NodeArray, NodeNull } from '../ir/ast';
import { Tag } from '../ir/tag';
import { Node } from '../ir/ast';
import { ValueType } from '../ir/value-type';

export type Constructor<T> = new () => T;

export class Binder {
  bind<T>(node: Node, type: Constructor<T> | T): T {
    const result =
      typeof type === 'function' ? new (type as Constructor<T>)() : type;
    const { value, error } = this.bindNode(node, result);
    if (error) {
      throw new Error(error);
    }
    return value;
  }

  /** 空容器注入的 example 哨兵项：只为保留元素类型，不参与绑定 */
  private isExample(node: Node | null | undefined): boolean {
    if (!node) return false;
    return node.getTag()?.example === true;
  }

  /** 编码侧真实值为 nil（is_null）：绑定成 null */
  private isNullNode(node: Node | null | undefined): boolean {
    if (!node) return false;
    return node.getTag()?.isNull === true;
  }

  bindNode(node: Node, target: any): { value: any; error: string | null } {
    // Handle NodeNull - return null value
    if (node instanceof NodeNull) {
      return { value: null, error: null };
    }

    // is_null：与 Go 对齐，绑定成 null（不展开 example 占位子节点）
    if (this.isNullNode(node)) {
      return { value: null, error: null };
    }

    if (node instanceof NodeScalar) {
      return { value: node.getValue(), error: null };
    }

    // When target is null/undefined, create dynamic containers
    // This mirrors Go's interface{} handling in Bind
    if (target === null || target === undefined) {
      if (node instanceof NodeObject) {
        const properties = node.getProperties();
        const result: Record<string, any> = {};
        for (const [key, valueNode] of Object.entries(properties)) {
          if (this.isExample(valueNode)) {
            continue; // 空容器占位项不绑定
          }
          const { value, error } = this.bindNode(valueNode, undefined);
          if (error) {
            return {
              value: result,
              error: `failed to bind field ${key}: ${error}`,
            };
          }
          result[key] = value;
        }
        return { value: result, error: null };
      }

      if (node instanceof NodeArray) {
        const elements = node.getElements();
        const result: any[] = [];
        for (let i = 0; i < elements.length; i++) {
          const el = elements[i]!;
          if (this.isExample(el)) {
            continue; // 空容器占位项不绑定
          }
          const { value, error } = this.bindNode(el, undefined);
          if (error) {
            return {
              value: result,
              error: `failed to bind array item ${i}: ${error}`,
            };
          }
          result.push(value);
        }
        return { value: result, error: null };
      }
    }

    if (node instanceof NodeObject) {
      const tag = node.getTag();
      if (tag?.type === ValueType.Obj) {
        return this.convertObj(node, target);
      } else {
        return this.convertMap(node, target);
      }
    }

    if (node instanceof NodeArray) {
      const tag = node.getTag();
      if (tag?.type === ValueType.Arr) {
        return this.convertArr(node, target);
      } else {
        return this.convertVec(node, target);
      }
    }

    return { value: target, error: `unsupported node type: ${typeof node}` };
  }

  private convertObj(
    obj: NodeObject,
    target: any,
  ): { value: any; error: string | null } {
    const tag = obj.getTag();
    if (tag?.nullable && target === null) {
      return {
        value: target,
        error: 'convertObj requires object type, got null',
      };
    }
    if (typeof target !== 'object' || target === null) {
      return {
        value: target,
        error: `convertObj requires object type, got ${typeof target}`,
      };
    }

    const properties = obj.getProperties();
    for (const [key, valueNode] of Object.entries(properties)) {
      if (this.isExample(valueNode)) {
        continue; // 空容器占位项不绑定
      }
      const camelKey = key.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
      if (!(camelKey in target)) {
        // For plain objects (not class instances), create fields dynamically
        // This mirrors Go's interface{} handling in Bind
        if (
          target.constructor === Object ||
          target.constructor === undefined
        ) {
          const { value, error } = this.bindNode(valueNode, undefined);
          if (error) {
            return {
              value: target,
              error: `failed to bind field ${camelKey}: ${error}`,
            };
          }
          target[camelKey] = value;
          continue;
        }
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
    obj: NodeObject,
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
      // 空 map 的 example 占位项：跳过后得到非 null 空对象
      if (this.isExample(valueNode)) {
        continue;
      }
      const { value, error } = this.bindNode(valueNode, target[key]);
      if (error) {
        return { value: target, error: `field ${key}: ${error}` };
      }
      target[key] = value;
    }
    return { value: target, error: null };
  }

  private convertArr(
    arr: NodeArray,
    target: any,
  ): { value: any; error: string | null } {
    if (!Array.isArray(target)) {
      return { value: target, error: `convertArr requires array` };
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
      if (this.isExample(el)) {
        continue; // 固定数组长度由 tag.size 决定，example 占位项不占下标
      }

      const { value, error } = this.bindNode(el, target[i]);
      if (error) {
        return { value: target, error: `array[${i}]: ${error}` };
      }
      target[i] = value;
    }

    return { value: target, error: null };
  }

  private convertVec(
    arr: NodeArray,
    target: any,
  ): { value: any; error: string | null } {
    if (!Array.isArray(target)) {
      // Create a new array for dynamic binding
      target = [];
    }

    const elements = arr.getElements();
    target.length = 0;

    for (const el of elements) {
      // example 哨兵项只为保留元素类型，不参与绑定；
      // 全部为 example（编码侧为空容器）时产出长度 0 的空数组。
      if (this.isExample(el)) {
        continue;
      }
      const defaultVal = this.createDefaultValue(el.getTag());
      const { value, error } = this.bindNode(el, defaultVal);
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
      case ValueType.I:
      case ValueType.U:
      case ValueType.I64:
      case ValueType.U64:
      case ValueType.Bigint:
        return 0n;

      case ValueType.I8:
      case ValueType.I16:
      case ValueType.I32:
      case ValueType.U8:
      case ValueType.U16:
      case ValueType.U32:
        return 0;

      case ValueType.F32:
      case ValueType.F64:
        return 0.0;

      case ValueType.Str:
      case ValueType.Email:
      case ValueType.Url:
      case ValueType.Ip:
      case ValueType.Uuid:
      case ValueType.Decimal:
      case ValueType.Enums:
        return '';

      case ValueType.Decimal:
        return '0.0';

      case ValueType.Bool:
        return false;

      case ValueType.Bytes:
      case ValueType.Media:
        return new Uint8Array();

      case ValueType.Datetime:
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
