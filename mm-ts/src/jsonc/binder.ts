import { JSONCValue, JSONCObject, JSONCArray } from './ast';

export interface JSONCNode {
  getType(): string;
  getTag(): any;
  getPath(): string;
  setPath(path: string): void;
  setTag(tag: any): void;
}

export class JSONCBinder {
  bind<T>(node: JSONCNode, Type: new () => T): T {
    const result = new Type();
    this.bindNode(node, result);
    return result;
  }

  private bindNode(node: JSONCNode, target: any): void {
    if (node instanceof JSONCValue) {
      this.bindValue(node, target);
    } else if (node instanceof JSONCObject) {
      this.bindObject(node, target);
    } else if (node instanceof JSONCArray) {
      this.bindArray(node, target);
    }
  }

  private bindValue(value: JSONCValue, target: any): void {
    // For direct value assignment
    if (target !== undefined) {
      target = value.getValue();
    }
  }

  private bindObject(obj: JSONCObject, target: any): void {
    for (const [key, valueNode] of obj.getProperties().entries()) {
      if (key in target) {
        const property = target[key];
        if (typeof property === 'object' && property !== null) {
          this.bindNode(valueNode, property);
        } else {
          if (valueNode instanceof JSONCValue) {
            target[key] = valueNode.getValue();
          } else if (valueNode instanceof JSONCObject) {
            target[key] = {};
            this.bindObject(valueNode, target[key]);
          } else if (valueNode instanceof JSONCArray) {
            target[key] = [];
            this.bindArray(valueNode, target[key]);
          }
        }
      }
    }
  }

  private bindArray(array: JSONCArray, target: any[]): void {
    for (const elementNode of array.getElements()) {
      if (elementNode instanceof JSONCValue) {
        target.push(elementNode.getValue());
      } else if (elementNode instanceof JSONCObject) {
        const obj: any = {};
        this.bindObject(elementNode, obj);
        target.push(obj);
      } else if (elementNode instanceof JSONCArray) {
        const arr: any[] = [];
        this.bindArray(elementNode, arr);
        target.push(arr);
      }
    }
  }
}

export function bindJSONC<T>(node: JSONCNode, Type: new () => T): T {
  const binder = new JSONCBinder();
  return binder.bind(node, Type);
}