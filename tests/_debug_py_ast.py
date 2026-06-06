import sys
sys.path.insert(0, 'mm-py')
from metamessage.jsonc.jsonc import parse_jsonc
from metamessage.ir.ast import Obj, Arr, Val

source = open('tests/fixtures/03_tags/child_tags.jsonc').read()
node = parse_jsonc(source)

def walk_node(n, depth=0):
    if isinstance(n, Obj):
        for f in n.fields:
            print(f'{"  "*depth}Field "{f.key}":')
            walk_node(f.value, depth+1)
    elif isinstance(n, Arr):
        print(f'{"  "*depth}Array (items: {len(n.items)}):')
        if n.tag:
            print(f'{"  "*depth}  ArrayTag: IsInherit={n.tag.is_inherit} Type={n.tag.type} Desc="{n.tag.desc}" ChildDesc="{n.tag.child_desc}"')
            print(f'{"  "*depth}  ArrayTag: ChildType={n.tag.child_type} Enums="{n.tag.enums}" ChildEnums="{n.tag.child_enums}"')
            tb = n.tag.bytes()
            print(f'{"  "*depth}  Bytes(): {tb.hex()} (len={len(tb)})')
        for i, item in enumerate(n.items):
            if isinstance(item, Val):
                print(f'{"  "*depth}  [{i}] Value data={item.data} text="{item.text}":')
                if item.tag:
                    print(f'{"  "*depth}    Tag: IsInherit={item.tag.is_inherit} Type={item.tag.type} Desc="{item.tag.desc}" ChildDesc="{item.tag.child_desc}"')
                    print(f'{"  "*depth}    Tag: ChildType={item.tag.child_type} Enums="{item.tag.enums}" ChildEnums="{item.tag.child_enums}"')
                    tb = item.tag.bytes()
                    print(f'{"  "*depth}    Bytes(): {tb.hex()} (len={len(tb)})')
                else:
                    print(f'{"  "*depth}    Tag: None')
            elif isinstance(n, Arr):
                walk_node(item, depth+2)
            else:
                walk_node(item, depth+2)
    elif isinstance(n, Val):
        print(f'{"  "*depth}Value: data={n.data} text="{n.text}"')

walk_node(node)