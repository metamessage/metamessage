import re

with open('/Users/lizongying/IdeaProjects/meta-message/mm-cs/src/MetaMessage/Jsonc/JsoncParser.cs', 'r') as f:
    content = f.read()

# Fix ParseObject signature and tag line
old_obj = (
    "private INode ParseObject(int openLine, string path)\n"
    "    {\n"
    "        _depth++;\n"
    "        if (_depth > MaxDepth)\n"
    "            throw new Exception(\"max depth: {MaxDepth}\");\n"
    "\n"
    "        Tag tag = ConsumeCommentsFor(openLine) ?? Tag.NewTag();"
)
new_obj = (
    "private INode ParseObject(int openLine, string path, Tag? existingTag = null)\n"
    "    {\n"
    "        _depth++;\n"
    "        if (_depth > MaxDepth)\n"
    "            throw new Exception(\"max depth: {MaxDepth}\");\n"
    "\n"
    "        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();"
)
content = content.replace(old_obj, new_obj)

# Fix ParseArray signature and tag line
old_arr = (
    "private INode ParseArray(int openLine, string path)\n"
    "    {\n"
    "        _depth++;\n"
    "        if (_depth > MaxDepth)\n"
    "            throw new Exception(\"max depth: {MaxDepth}\");\n"
    "\n"
    "        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();"
)
new_arr = (
    "private INode ParseArray(int openLine, string path, Tag? existingTag = null)\n"
    "    {\n"
    "        _depth++;\n"
    "        if (_depth > MaxDepth)\n"
    "            throw new Exception(\"max depth: {MaxDepth}\");\n"
    "\n"
    "        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();"
)
content = content.replace(old_arr, new_arr)

with open('/Users/lizongying/IdeaProjects/meta-message/mm-cs/src/MetaMessage/Jsonc/JsoncParser.cs', 'w') as f:
    f.write(content)

print('Done')