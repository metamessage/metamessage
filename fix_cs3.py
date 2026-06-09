content = open('/Users/lizongying/IdeaProjects/meta-message/mm-cs/src/MetaMessage/Jsonc/JsoncParser.cs', 'r').read()

# Fix ParseObject
old_obj = 'private INode ParseObject(int openLine, string path)\n    {\n        _depth++;\n        if (_depth > MaxDepth)\n            throw new Exception($"max depth: {MaxDepth}");\n\n        Tag tag = ConsumeCommentsFor(openLine) ?? Tag.NewTag();'
new_obj = 'private INode ParseObject(int openLine, string path, Tag? existingTag = null)\n    {\n        _depth++;\n        if (_depth > MaxDepth)\n            throw new Exception($"max depth: {MaxDepth}");\n\n        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();'
if old_obj in content:
    content = content.replace(old_obj, new_obj)
    print('ParseObject fixed')
else:
    print('ParseObject pattern NOT FOUND')

# Fix ParseArray
old_arr = 'private INode ParseArray(int openLine, string path)\n    {\n        _depth++;\n        if (_depth > MaxDepth)\n            throw new Exception($"max depth: {MaxDepth}");\n\n        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();'
new_arr = 'private INode ParseArray(int openLine, string path, Tag? existingTag = null)\n    {\n        _depth++;\n        if (_depth > MaxDepth)\n            throw new Exception($"max depth: {MaxDepth}");\n\n        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();'
if old_arr in content:
    content = content.replace(old_arr, new_arr)
    print('ParseArray fixed')
else:
    print('ParseArray pattern NOT FOUND')

open('/Users/lizongying/IdeaProjects/meta-message/mm-cs/src/MetaMessage/Jsonc/JsoncParser.cs', 'w').write(content)
print('Done')