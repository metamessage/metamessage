#ifndef MMCPP_AST_HPP
#define MMCPP_AST_HPP

#include <string>
#include <vector>
#include <variant>
#include <optional>
#include <memory>

namespace mmc {

enum class ValueType {
    Unknown,
    String,
    Bool,
    Int,
    Int8,
    Int16,
    Int32,
    Int64,
    Uint,
    Uint8,
    Uint16,
    Uint32,
    Uint64,
    Float32,
    Float64,
    Bytes,
    BigInt,
    Array,
    Slice,
    Object
};

struct Tag {
    ValueType type = ValueType::Unknown;
    int size = -1;
    int required = -1;
    int minValue = 0;
    int maxValue = 0;
    int minLength = 0;
    int maxLength = 0;
    int deprecated = 0;
};

using ValueData = std::variant<
    bool,
    std::string,
    int64_t,
    uint64_t,
    double,
    std::vector<uint8_t>,
    std::monostate
>;

class Node;
class Value;
class Object;
class Array;

class Value {
public:
    ValueData data;
    std::string text;
    Tag tag;

    Value() : data(std::monostate{}) {}
};

class Object {
public:
    std::vector<std::pair<std::string, std::unique_ptr<Node>>> fields;
    Tag tag;
};

class Array {
public:
    std::vector<std::unique_ptr<Node>> items;
    Tag tag;
};

class Node {
public:
    enum class Type { Value, Object, Array } type;
    std::variant<Value, Object, Array> data;

    explicit Node(Type t) : type(t) {}
};

std::unique_ptr<Node> createValueNode();
std::unique_ptr<Node> createObjectNode();
std::unique_ptr<Node> createArrayNode();

}

#endif
