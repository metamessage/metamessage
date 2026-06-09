content = open('/Users/lizongying/IdeaProjects/meta-message/mm-cs/src/MetaMessage/Jsonc/JsoncParser.cs', 'r').read()

old = 'private INode ParseArray(int openLine, string path)\n    {\n        _depth++;\n        if (_depth > MaxDepth)\n            throw new Exception($"max depth: {MaxDepth}");\n\n        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();'
new = 'private INode ParseArray(int openLine, string path, Tag? existingTag = null)\n    {\n        _depth++;\n        if (_depth > MaxDepth)\n            throw new Exception($"max depth: {MaxDepth}");\n\n        Tag tag = existingTag ?? ConsumeCommentsFor(openLine) ?? Tag.NewTag();'

assert old in content, 'Pattern not found!'
content = content.replace(old, new)
open('/Users/lizongying/IdeaProjects/meta-message/mm-cs/src/MetaMessage/Jsonc/JsoncParser.cs', 'w').write(content)
print('Done - ParseArray signature fixed')