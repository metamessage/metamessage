import re

with open('src/ir/tag.rs', 'r') as f:
    content = f.read()

# child_default -> child_default_val (only field accesses)
content = re.sub(r'parent\.child_default\b', 'parent.child_default_val', content)
content = re.sub(r'(?<![a-zA-Z"\'])tag\.child_default\b', 'tag.child_default_val', content)
content = re.sub(r'(?<![a-zA-Z"\'])self\.child_default\b', 'self.child_default_val', content)

# child_enum -> child_enums (only field accesses)
content = re.sub(r'parent\.child_enum\b', 'parent.child_enums', content)
content = re.sub(r'(?<![a-zA-Z"\'])tag\.child_enum\b', 'tag.child_enums', content)
content = re.sub(r'(?<![a-zA-Z"\'])self\.child_enum\b', 'self.child_enums', content)

# default -> default_val (only field accesses, not string literals)
content = re.sub(r'(?<![a-zA-Z"\'])self\.default\b(?!_)', 'self.default_val', content)
content = re.sub(r'(?<![a-zA-Z"\'])tag\.default\b(?!_)', 'tag.default_val', content)
content = re.sub(r'src\.default\b(?!_)', 'src.default_val', content)
content = re.sub(r'dst\.default\b(?!_)', 'dst.default_val', content)

# enum_values -> enums (only field accesses)
content = re.sub(r'(?<![a-zA-Z"\'])self\.enum_values\b', 'self.enums', content)
content = re.sub(r'(?<![a-zA-Z"\'])tag\.enum_values\b', 'tag.enums', content)
content = re.sub(r'src\.enum_values\b', 'src.enums', content)
content = re.sub(r'dst\.enum_values\b', 'dst.enums', content)

with open('src/ir/tag.rs', 'w') as f:
    f.write(content)

print("Done")