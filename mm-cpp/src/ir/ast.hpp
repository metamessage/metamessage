#ifndef MMCPP_IR_AST_HPP
#define MMCPP_IR_AST_HPP

#include "tag.hpp"
#include <string>
#include <vector>
#include <memory>
#include <variant>
#include <cstdint>

namespace mmc {
namespace ir {

enum class NodeType : uint8_t {
    Unknown = 0,
    NodeObject,
    NodeArray,
    Value,
    Doc
};

inline std::string nodeTypeToString(NodeType nt) {
    switch (nt) {
        case NodeType::NodeObject: return "object";
        case NodeType::NodeArray:  return "array";
        case NodeType::Value:  return "value";
        case NodeType::Doc:    return "doc";
        default:               return "unknown";
    }
}

class Node;
class NodeObject;
class NodeArray;
class NodeScalar;
class Doc;

struct Field {
    std::string key;
    std::shared_ptr<Node> value;

    Field(std::string k, std::shared_ptr<Node> v)
        : key(std::move(k)), value(std::move(v)) {}
};

class Node {
public:
    virtual ~Node() = default;
    virtual NodeType getType() const = 0;
    virtual Tag* getTag() = 0;
    virtual const Tag* getTag() const = 0;
    virtual const std::string& getPath() const { return path_; }
    virtual void setPath(const std::string& p) { path_ = p; }

protected:
    std::string path_;
};

class NodeObject : public Node {
public:
    std::vector<Field> fields;
    Tag tag;

    NodeType getType() const override { return NodeType::NodeObject; }
    Tag* getTag() override { return &tag; }
    const Tag* getTag() const override { return &tag; }
};

class NodeArray : public Node {
public:
    std::vector<std::shared_ptr<Node>> items;
    Tag tag;

    NodeType getType() const override { return NodeType::NodeArray; }
    Tag* getTag() override { return &tag; }
    const Tag* getTag() const override { return &tag; }
};

class NodeScalar : public Node {
public:
    std::string text;
    int64_t data = 0;
    Tag tag;

    NodeType getType() const override { return NodeType::Value; }
    Tag* getTag() override { return &tag; }
    const Tag* getTag() const override { return &tag; }
};

class Doc : public Node {
public:
    std::vector<Field> fields;
    Tag tag;

    NodeType getType() const override { return NodeType::Doc; }
    Tag* getTag() override { return &tag; }
    const Tag* getTag() const override { return &tag; }
};

inline std::shared_ptr<NodeObject> makeNodeObject() {
    return std::make_shared<NodeObject>();
}

inline std::shared_ptr<NodeArray> makeNodeArray() {
    return std::make_shared<NodeArray>();
}

inline std::shared_ptr<NodeScalar> makeNodeScalar() {
    return std::make_shared<NodeScalar>();
}

inline std::shared_ptr<Doc> makeDoc() {
    return std::make_shared<Doc>();
}

} // namespace ir
} // namespace mmc

#endif