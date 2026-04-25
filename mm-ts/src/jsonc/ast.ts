export interface JSONCNode {
  getType(): string;
  getTag(): JSONCTag | null;
  getPath(): string;
  setPath(path: string): void;
  setTag(tag: JSONCTag): void;
}

export class JSONCValue implements JSONCNode {
  private value: any;
  private tag: JSONCTag | null;
  private path: string;

  constructor(value: any) {
    this.value = value;
    this.tag = null;
    this.path = '';
  }

  getType(): string {
    if (this.value === null) return 'null';
    if (typeof this.value === 'boolean') return 'boolean';
    if (typeof this.value === 'number') return 'number';
    if (typeof this.value === 'string') return 'string';
    return 'unknown';
  }

  getTag(): JSONCTag | null {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: JSONCTag): void {
    this.tag = tag;
  }

  getValue(): any {
    return this.value;
  }

  setValue(value: any): void {
    this.value = value;
  }
}

export class JSONCObject implements JSONCNode {
  private properties: Map<string, JSONCNode>;
  private tag: JSONCTag | null;
  private path: string;

  constructor() {
    this.properties = new Map();
    this.tag = null;
    this.path = '';
  }

  getType(): string {
    return 'object';
  }

  getTag(): JSONCTag | null {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: JSONCTag): void {
    this.tag = tag;
  }

  getProperties(): Map<string, JSONCNode> {
    return this.properties;
  }

  setProperty(key: string, value: JSONCNode): void {
    this.properties.set(key, value);
  }

  getProperty(key: string): JSONCNode | undefined {
    return this.properties.get(key);
  }

  hasProperty(key: string): boolean {
    return this.properties.has(key);
  }

  size(): number {
    return this.properties.size;
  }
}

export class JSONCArray implements JSONCNode {
  private elements: JSONCNode[];
  private tag: JSONCTag | null;
  private path: string;

  constructor() {
    this.elements = [];
    this.tag = null;
    this.path = '';
  }

  getType(): string {
    return 'array';
  }

  getTag(): JSONCTag | null {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: JSONCTag): void {
    this.tag = tag;
  }

  getElements(): JSONCNode[] {
    return this.elements;
  }

  addElement(element: JSONCNode): void {
    this.elements.push(element);
  }

  getElement(index: number): JSONCNode | undefined {
    return this.elements[index];
  }

  size(): number {
    return this.elements.length;
  }
}

export class JSONCDoc implements JSONCNode {
  private root: JSONCNode;
  private tag: JSONCTag | null;
  private path: string;

  constructor(root: JSONCNode) {
    this.root = root;
    this.tag = null;
    this.path = '';
  }

  getType(): string {
    return 'document';
  }

  getTag(): JSONCTag | null {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: JSONCTag): void {
    this.tag = tag;
  }

  getRoot(): JSONCNode {
    return this.root;
  }

  setRoot(root: JSONCNode): void {
    this.root = root;
  }
}

export class JSONCTag {
  private properties: Map<string, string>;

  constructor() {
    this.properties = new Map();
  }

  set(key: string, value: string): void {
    this.properties.set(key, value);
  }

  get(key: string): string | undefined {
    return this.properties.get(key);
  }

  has(key: string): boolean {
    return this.properties.has(key);
  }

  delete(key: string): boolean {
    return this.properties.delete(key);
  }

  clear(): void {
    this.properties.clear();
  }

  size(): number {
    return this.properties.size;
  }

  entries(): IterableIterator<[string, string]> {
    return this.properties.entries();
  }

  inherit(tag: JSONCTag): void {
    for (const [key, value] of tag.entries()) {
      if (!this.properties.has(key)) {
        this.properties.set(key, value);
      }
    }
  }

  toString(): string {
    const parts: string[] = [];
    for (const [key, value] of this.properties.entries()) {
      parts.push(`${key}=${value}`);
    }
    return parts.join(';');
  }
}

export function parseMMTag(tagStr: string): JSONCTag {
  const tag = new JSONCTag();
  const parts = tagStr.split(';');
  
  for (const part of parts) {
    const trimmed = part.trim();
    if (trimmed) {
      const equalsIndex = trimmed.indexOf('=');
      if (equalsIndex > 0) {
        const key = trimmed.substring(0, equalsIndex).trim();
        const value = trimmed.substring(equalsIndex + 1).trim();
        tag.set(key, value);
      }
    }
  }
  
  return tag;
}