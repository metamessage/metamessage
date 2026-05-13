import { Tag } from './tag';

type NodeType = 'unknown' | 'object' | 'array' | 'value' | 'doc';

export interface Node {
  getType(): string;
  getTag(): Tag;
  getPath(): string;
  setPath(path: string): void;
  setTag(tag: Tag): void;
}

export class MMValue implements Node {
  private value: any;
  private tag: Tag;
  private path: string;

  constructor(value: any, tag?: Tag) {
    this.value = value;
    this.tag = tag || new Tag();
    this.path = '';
  }

  getType(): NodeType {
    return 'value';
  }

  getTag(): Tag {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: Tag): void {
    this.tag = tag;
  }

  getValue(): any {
    return this.value;
  }

  setValue(value: any): void {
    this.value = value;
  }
}

export class MMObject implements Node {
  private properties: Record<string, Node>;
  private tag: Tag;
  private path: string;

  constructor(tag?: Tag) {
    this.properties = {};
    this.tag = tag || new Tag();
    this.path = '';
  }

  getType(): NodeType {
    return 'object';
  }

  getTag(): Tag {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: Tag): void {
    this.tag = tag;
  }

  getProperties(): Record<string, Node> {
    return this.properties;
  }

  setProperty(key: string, value: Node): void {
    this.properties[key] = value;
  }

  getProperty(key: string): Node | undefined {
    return this.properties[key];
  }

  hasProperty(key: string): boolean {
    return this.properties[key] !== undefined;
  }
}

export class MMArray implements Node {
  private elements: Node[];
  private tag: Tag;
  private path: string;

  constructor() {
    this.elements = [];
    this.tag = new Tag();
    this.path = '';
  }

  getType(): NodeType {
    return 'array';
  }

  getTag(): Tag {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: Tag): void {
    this.tag = tag;
  }

  getElements(): Node[] {
    return this.elements;
  }

  addElement(element: Node): void {
    this.elements.push(element);
  }

  getElement(index: number): Node | undefined {
    return this.elements[index];
  }

  size(): number {
    return this.elements.length;
  }
}

export class MMDoc implements Node {
  private root: Node;
  private tag: Tag;
  private path: string;

  constructor(root: Node, tag?: Tag) {
    this.root = root;
    this.tag = tag || new Tag();
    this.path = '';
  }

  getType(): NodeType {
    return 'doc';
  }

  getTag(): Tag {
    return this.tag;
  }

  getPath(): string {
    return this.path;
  }

  setPath(path: string): void {
    this.path = path;
  }

  setTag(tag: Tag): void {
    this.tag = tag;
  }

  getRoot(): Node {
    return this.root;
  }

  setRoot(root: Node): void {
    this.root = root;
  }
}
