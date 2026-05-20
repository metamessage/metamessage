#ifndef MMCPP_JSONC_PRINTER_HPP
#define MMCPP_JSONC_PRINTER_HPP

#include "../ir/ast.hpp"
#include <string>
#include <sstream>

namespace mmc {
namespace jsonc {

inline void printIndent(std::ostringstream& os, int indent) {
    for (int i = 0; i < indent; ++i) os << "    ";
}

inline void printValue(std::ostringstream& os, std::shared_ptr<ir::Value> val, int indent);
inline void printArray(std::ostringstream& os, std::shared_ptr<ir::Array> arr, int indent);
inline void printObject(std::ostringstream& os, std::shared_ptr<ir::Object> obj, int indent);

inline void printTag(std::ostringstream& os, const ir::Tag* tag) {
    if (tag == nullptr) return;
    std::string ts = tag->toString();
    if (!ts.empty()) {
        os << " // mm: " << ts;
    }
}

inline void printValue(std::ostringstream& os, std::shared_ptr<ir::Value> val, int indent) {
    printTag(os, val->getTag());

    switch (val->getTag()->type) {
        case ir::ValueType::Str:
        case ir::ValueType::Email:
        case ir::ValueType::Url: {
            os << "\"";
            for (char c : val->text) {
                switch (c) {
                    case '"':  os << "\\\""; break;
                    case '\\': os << "\\\\"; break;
                    case '\n': os << "\\n"; break;
                    case '\r': os << "\\r"; break;
                    case '\t': os << "\\t"; break;
                    default:   os << c;
                }
            }
            os << "\"";
            break;
        }
        case ir::ValueType::Bytes:
            os << "\"" << val->text << "\"";
            break;
        case ir::ValueType::Bool:
            os << (val->text == "true" || val->text == "1" ? "true" : "false");
            break;
        case ir::ValueType::I: case ir::ValueType::I8:
        case ir::ValueType::I16: case ir::ValueType::I32: case ir::ValueType::I64:
        case ir::ValueType::U: case ir::ValueType::U8:
        case ir::ValueType::U16: case ir::ValueType::U32: case ir::ValueType::U64:
        case ir::ValueType::F32: case ir::ValueType::F64:
        case ir::ValueType::Enum:
            os << val->text;
            break;
        case ir::ValueType::Bigint:
            os << val->text << "n";
            break;
        case ir::ValueType::Datetime:
        case ir::ValueType::Date:
        case ir::ValueType::Time:
            os << "\"" << val->text << "\"";
            break;
        case ir::ValueType::Uuid:
            os << "\"" << val->text << "\"";
            break;
        case ir::ValueType::Decimal:
            os << val->text;
            break;
        case ir::ValueType::Ip:
            os << "\"" << val->text << "\"";
            break;
        default:
            if (val->text.empty()) os << "null";
            else os << "\"" << val->text << "\"";
            break;
    }
}

inline void printArray(std::ostringstream& os, std::shared_ptr<ir::Array> arr, int indent) {
    printTag(os, arr->getTag());
    os << "[";

    if (arr->items.empty()) {
        os << "]";
        return;
    }

    bool multiline = false;
    size_t objCount = 0;
    for (auto& item : arr->items) {
        if (item->getType() == ir::NodeType::Object) {
            objCount++;
            if (!multiline && item->getType() == ir::NodeType::Object) {
                auto obj = std::static_pointer_cast<ir::Object>(item);
                if (obj->fields.size() > 1) multiline = true;
            }
        }
    }
    if (objCount > 1) multiline = true;

    if (multiline) os << "\n";

    for (size_t i = 0; i < arr->items.size(); ++i) {
        if (multiline) printIndent(os, indent + 1);

        switch (arr->items[i]->getType()) {
            case ir::NodeType::Object:
                printObject(os, std::static_pointer_cast<ir::Object>(arr->items[i]),
                           multiline ? indent + 1 : indent);
                break;
            case ir::NodeType::Array:
                printArray(os, std::static_pointer_cast<ir::Array>(arr->items[i]),
                          multiline ? indent + 1 : indent);
                break;
            case ir::NodeType::Value:
                printValue(os, std::static_pointer_cast<ir::Value>(arr->items[i]),
                          multiline ? indent + 1 : indent);
                break;
            default:
                os << "null";
                break;
        }

        if (i < arr->items.size() - 1) os << ",";
        if (multiline) os << "\n";
    }

    if (multiline) printIndent(os, indent);
    os << "]";
}

inline void printObject(std::ostringstream& os, std::shared_ptr<ir::Object> obj, int indent) {
    printTag(os, obj->getTag());
    os << "{";

    if (obj->fields.empty()) {
        os << "}";
        return;
    }

    bool multiline = false;
    for (auto& f : obj->fields) {
        if (f.value->getType() == ir::NodeType::Object) multiline = true;
    }
    if (obj->fields.size() > 1) multiline = true;

    if (multiline) os << "\n";

    for (size_t i = 0; i < obj->fields.size(); ++i) {
        auto& field = obj->fields[i];
        if (multiline) printIndent(os, indent + 1);

        os << "\"";
        for (char c : field.key) {
            switch (c) {
                case '"':  os << "\\\""; break;
                case '\\': os << "\\\\"; break;
                default:   os << c;
            }
        }
        os << "\": ";

        switch (field.value->getType()) {
            case ir::NodeType::Object:
                printObject(os, std::static_pointer_cast<ir::Object>(field.value),
                           multiline ? indent + 1 : indent);
                break;
            case ir::NodeType::Array:
                printArray(os, std::static_pointer_cast<ir::Array>(field.value),
                          multiline ? indent + 1 : indent);
                break;
            case ir::NodeType::Value:
                printValue(os, std::static_pointer_cast<ir::Value>(field.value),
                          multiline ? indent + 1 : indent);
                break;
            default:
                os << "null";
                break;
        }

        if (i < obj->fields.size() - 1) os << ",";
        if (multiline) os << "\n";
    }

    if (multiline) printIndent(os, indent);
    os << "}";
}

inline std::string toJSONC(std::shared_ptr<ir::Node> node) {
    std::ostringstream os;
    switch (node->getType()) {
        case ir::NodeType::Object:
            printObject(os, std::static_pointer_cast<ir::Object>(node), 0);
            break;
        case ir::NodeType::Array:
            printArray(os, std::static_pointer_cast<ir::Array>(node), 0);
            break;
        case ir::NodeType::Value:
            printValue(os, std::static_pointer_cast<ir::Value>(node), 0);
            break;
        default:
            os << "null";
            break;
    }
    return os.str();
}

} // namespace jsonc
} // namespace mmc

#endif