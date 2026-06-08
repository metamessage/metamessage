import subprocess, sys
sys.path.insert(0, '.')
from metamessage.core.decoder import Decoder
from metamessage.ir.tag import Tag, ValueType

result = subprocess.run(['go', 'run', 'tests/harness/go/harness.go', '--encode', 'tests/fixtures/03_tags/constraints_tag.jsonc'], capture_output=True, text=True)
hex_str = result.stdout.strip()
data = bytes.fromhex(hex_str)
decoder = Decoder(data)
node = decoder.decode()

# Find the Arr node and check its tag
def find_arr_tag(node, depth=0):
    if hasattr(node, 'tag') and node.tag:
        t = node.tag
        if t.type == ValueType.Arr:
            print(f'Arr tag found at depth {depth}')
            print(f'  type={t.type} ({ValueType(t.type).name})')
            print(f'  size={t.size}')
            print(f'  is_inherit={t.is_inherit}')
            print(f'  str(tag)={repr(str(t))}')
            # Also check what __str__ evaluates
            parts = []
            if t.type != ValueType.Unknown and not t.is_inherit:
                if t.type in (ValueType.Str, ValueType.I, ValueType.F64, ValueType.Bool,
                             ValueType.Obj, ValueType.Vec):
                    print('  type is in suppression list')
                else:
                    print(f'  type is NOT in suppression list (type={t.type})')
                    if not (t.type == ValueType.Enums and t.enums) and not (t.type == ValueType.Media and t.mime):
                        parts.append("type=%s" % str(t.type))
                        print(f'  would add type={str(t.type)}')
                    else:
                        print(f'  skipped due to enums/media check')
            print(f'  parts={parts}')
            return
    if hasattr(node, 'fields'):
        for f in node.fields:
            find_arr_tag(f.value, depth+1)
    if hasattr(node, 'items') and not callable(node.items):
        for item in node.items:
            find_arr_tag(item, depth+2)

find_arr_tag(node)