export class JSONCBinder {
  bind(node, Type) {
    if (!node) return null;

    switch (node.getType()) {
      case 'object':
        return this.bindObject(node, Type);
      case 'array':
        return this.bindArray(node, Type);
      case 'value':
        return this.bindValue(node, Type);
      default:
        return null;
    }
  }

  bindObject(node, Type) {
    if (!node || node.getType() !== 'object') {
      throw new Error('Expected object');
    }

    if (Type === Object || Type === null) {
      const result = {};
      for (const field of node.fields) {
        result[field.key] = this.bind(field.value, null);
      }
      return result;
    }

    const result = {};
    const typeInfo = this.getTypeInfo(Type);

    for (const field of node.fields) {
      const key = field.key;
      const value = field.value;

      const targetKey = this.findMatchingKey(key, typeInfo.fields);

      if (targetKey) {
        const fieldType = typeInfo.fields[targetKey];
        result[targetKey] = this.bind(value, fieldType);
      } else {
        result[key] = this.bind(value, null);
      }
    }

    return result;
  }

  bindArray(node, Type) {
    if (!node || node.getType() !== 'array') {
      throw new Error('Expected array');
    }

    const result = [];
    const itemType = this.getArrayItemType(Type);

    for (const item of node.items) {
      result.push(this.bind(item, itemType));
    }

    return result;
  }

  bindValue(node, Type) {
    if (!node || node.getType() !== 'value') {
      throw new Error('Expected value');
    }

    if (node.tag && node.tag.isNull) {
      return null;
    }

    const data = node.data;

    if (Type === Boolean || Type === 'boolean') {
      if (typeof data === 'boolean') return data;
      if (data === 'true' || data === 'false') return data === 'true';
      return Boolean(data);
    }

    if (Type === Number || Type === 'number') {
      if (typeof data === 'number') return data;
      if (typeof data === 'bigint') return Number(data);
      return Number(data);
    }

    if (Type === BigInt || Type === 'bigint') {
      if (typeof data === 'bigint') return data;
      return BigInt(data);
    }

    if (Type === String || Type === 'string') {
      if (typeof data === 'string') return data;
      return String(data);
    }

    if (Type === Uint8Array || Type === 'bytes') {
      if (data instanceof Uint8Array) return data;
      if (Array.isArray(data)) return new Uint8Array(data);
      return data;
    }

    if (Type === null || Type === undefined || Type === Object) {
      if (data === null || data === undefined) return null;

      if (typeof data === 'boolean') return data;
      if (typeof data === 'number') return data;
      if (typeof data === 'bigint') return data;
      if (typeof data === 'string') return data;
      if (Array.isArray(data)) {
        return this.bindArray(node, null);
      }
      if (typeof data === 'object') {
        return this.bindObject(node, null);
      }
      return data;
    }

    return data;
  }

  getTypeInfo(Type) {
    if (!Type || Type === Object) {
      return { fields: {} };
    }

    const fields = {};

    if (typeof Type === 'function') {
      const params = Type.toString().match(/\((.*?)\)/);
      if (params) {
        const args = params[1].split(',').map(p => p.trim());
        for (const arg of args) {
          if (arg) {
            fields[arg] = null;
          }
        }
      }
    }

    return { fields };
  }

  getArrayItemType(Type) {
    if (!Type) return null;

    if (Array.isArray(Type)) {
      return Type[0];
    }

    return null;
  }

  findMatchingKey(key, fields) {
    const lowerKey = key.toLowerCase();

    for (const fieldName of Object.keys(fields)) {
      if (fieldName.toLowerCase() === lowerKey) {
        return fieldName;
      }
    }

    const snakeToCamel = (str) => str.replace(/_([a-z])/g, (_, c) => c.toUpperCase());
    const camelToSnake = (str) => str.replace(/([A-Z])/g, '_$1').toLowerCase();

    for (const fieldName of Object.keys(fields)) {
      if (snakeToCamel(fieldName) === snakeToCamel(key)) {
        return fieldName;
      }
      if (camelToSnake(fieldName) === camelToSnake(key)) {
        return fieldName;
      }
    }

    return null;
  }
}