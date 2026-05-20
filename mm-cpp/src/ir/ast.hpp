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
    Object,
    Array,
    Value,
    Doc
};

inline std::string nodeTypeToString(NodeType nt) {
    switch (nt) {
        case NodeType::Object: return "object";
        case NodeType::Array:  return "array";
        case NodeType::Value:  return "value";
        case NodeType::Doc:    return "doc";
        default:               return "unknown";
    }
}

class Node;
class Object;
class Array;
class Value;
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

class Object : public Node {
public:
    std::vector<Field> fields;
    Tag tag;

    NodeType getType() const override { return NodeType::Object; }
    Tag* getTag() override { return &tag; }
    const Tag* getTag() const override { return &tag; }
};

class Array : public Node {
public:
    std::vector<std::shared_ptr<Node>> items;
    Tag tag;

    NodeType getType() const override { return NodeType::Array; }
    Tag* getTag() override { return &tag; }
    const Tag* getTag() const override { return &tag; }
};

class Value : public Node {
public:
    std::string text;
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

inline std::shared_ptr<Object> makeObject() {
    return std::make_shared<Object>();
}

inline std::shared_ptr<Array> makeArray() {
    return std::make_shared<Array>();
}

inline std::shared_ptr<Value> makeValue() {
    return std::make_shared<Value>();
}

inline std::shared_ptr<Doc> makeDoc() {
    return std::make_shared<Doc>();
}

} // namespace ir
} // namespace mmc

#endif