#ifndef MMCPP_JSONC_PRINTER_HPP
#define MMCPP_JSONC_PRINTER_HPP

#include "../ir/ast.hpp"
#include <sstream>
#include <string>

namespace mmc {
namespace jsonc {

inline void printIndent(std::ostringstream &os, int indent) {
  for (int i = 0; i < indent; ++i)
    os << "\t";
}

inline void printValue(std::ostringstream &os, std::shared_ptr<ir::Value> val,
                       int indent);
inline void printArray(std::ostringstream &os, std::shared_ptr<ir::Array> arr,
                       int indent);
inline void printObject(std::ostringstream &os, std::shared_ptr<ir::Object> obj,
                        int indent);

inline void printTag(std::ostringstream &os, const ir::Tag *tag) {
  if (tag == nullptr)
    return;
  std::string ts = tag->toString();
  if (!ts.empty()) {
    os << " // mm: " << ts;
  }
}

inline void printLeadingComment(std::ostringstream &os, const ir::Tag *tag,
                                int indent) {
  if (tag == nullptr)
    return;
  std::string ts = tag->toString();
  if (!ts.empty()) {
    os << "\n";
    printIndent(os, indent);
    os << "// mm: " << ts << "\n";
  }
}

inline void printValue(std::ostringstream &os, std::shared_ptr<ir::Value> val,
                       int indent) {

  switch (val->getTag()->type) {
  case ir::ValueType::Str:
  case ir::ValueType::Email:
  case ir::ValueType::Url:
  case ir::ValueType::Enums: {
    os << "\"";
    for (char c : val->text) {
      switch (c) {
      case '"':
        os << "\\\"";
        break;
      case '\\':
        os << "\\\\";
        break;
      case '\n':
        os << "\\n";
        break;
      case '\r':
        os << "\\r";
        break;
      case '\t':
        os << "\\t";
        break;
      default:
        os << c;
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
  case ir::ValueType::I:
  case ir::ValueType::I8:
  case ir::ValueType::I16:
  case ir::ValueType::I32:
  case ir::ValueType::I64:
  case ir::ValueType::U:
  case ir::ValueType::U8:
  case ir::ValueType::U16:
  case ir::ValueType::U32:
  case ir::ValueType::U64:
  case ir::ValueType::F32:
  case ir::ValueType::F64:
    os << val->text;
    break;
  case ir::ValueType::Bigint:
    os << val->text;
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
    if (val->text.empty())
      os << "null";
    else
      os << "\"" << val->text << "\"";
    break;
  }
}

inline void printArray(std::ostringstream &os, std::shared_ptr<ir::Array> arr,
                       int indent) {
  os << "[\n";

  if (arr->items.empty()) {
    printIndent(os, indent);
    os << "]";
    return;
  }

  for (size_t i = 0; i < arr->items.size(); ++i) {
    printLeadingComment(os, arr->items[i]->getTag(), indent + 1);
    printIndent(os, indent + 1);

    switch (arr->items[i]->getType()) {
    case ir::NodeType::Object:
      printObject(os, std::static_pointer_cast<ir::Object>(arr->items[i]),
                  indent + 1);
      break;
    case ir::NodeType::Array:
      printArray(os, std::static_pointer_cast<ir::Array>(arr->items[i]),
                 indent + 1);
      break;
    case ir::NodeType::Value:
      printValue(os, std::static_pointer_cast<ir::Value>(arr->items[i]),
                 indent + 1);
      break;
    default:
      os << "null";
      break;
    }

    os << ",\n";
  }

  printIndent(os, indent);
  os << "]";
}

inline void printObject(std::ostringstream &os, std::shared_ptr<ir::Object> obj,
                        int indent) {
  os << "{\n";

  if (obj->fields.empty()) {
    printIndent(os, indent);
    os << "}";
    return;
  }

  for (size_t i = 0; i < obj->fields.size(); ++i) {
    auto &field = obj->fields[i];
    printLeadingComment(os, field.value->getTag(), indent + 1);
    printIndent(os, indent + 1);

    os << "\"";
    for (char c : field.key) {
      switch (c) {
      case '"':
        os << "\\\"";
        break;
      case '\\':
        os << "\\\\";
        break;
      default:
        os << c;
      }
    }
    os << "\": ";

    switch (field.value->getType()) {
    case ir::NodeType::Object:
      printObject(os, std::static_pointer_cast<ir::Object>(field.value),
                  indent + 1);
      break;
    case ir::NodeType::Array:
      printArray(os, std::static_pointer_cast<ir::Array>(field.value),
                 indent + 1);
      break;
    case ir::NodeType::Value:
      printValue(os, std::static_pointer_cast<ir::Value>(field.value),
                 indent + 1);
      break;
    default:
      os << "null";
      break;
    }

    os << ",\n";
  }

  printIndent(os, indent);
  os << "}";
}

inline std::string toJSONC(std::shared_ptr<ir::Node> node) {
  std::ostringstream os;
  printLeadingComment(os, node->getTag(), 0);
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