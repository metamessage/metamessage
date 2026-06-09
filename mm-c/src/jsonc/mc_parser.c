#include "mc_parser.h"
#include <ctype.h>
#include <inttypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_DEPTH 32

static void parser_error(MMC_Parser *parser, const char *message) {
  if (parser->error_message != NULL) {
    free(parser->error_message);
  }
  parser->error_message = (char *)malloc(256);
  if (parser->error_message != NULL) {
    snprintf(parser->error_message, 256, "Error at line %zu, column %zu: %s",
             parser->current_token ? parser->current_token->line : 0,
             parser->current_token ? parser->current_token->column : 0,
             message);
  }
}

static MMC_Node *parse_value(MMC_Parser *parser, MMC_Tag *inherit_tag);
static MMC_Node *parse_primitive(MMC_Parser *parser, MMC_Token *token,
                                 MMC_Tag *pending_tag);
static MMC_Node *parse_object(MMC_Parser *parser, MMC_Tag *tag,
                              MMC_Tag *inherit_tag);
static MMC_Node *parse_array(MMC_Parser *parser, MMC_Tag *tag,
                             MMC_Tag *inherit_tag);
static void next_token(MMC_Parser *parser);

static void add_pending_comment(MMC_Parser *parser, MMC_Token *token) {
  if (parser->pending_count >= parser->pending_capacity) {
    size_t new_cap =
        parser->pending_capacity == 0 ? 8 : parser->pending_capacity * 2;
    MMC_Token **new_arr = (MMC_Token **)realloc(parser->pending_comments,
                                                new_cap * sizeof(MMC_Token *));
    if (new_arr == NULL)
      return;
    parser->pending_comments = new_arr;
    parser->pending_capacity = new_cap;
  }
  MMC_Token *copy = (MMC_Token *)malloc(sizeof(MMC_Token));
  if (copy == NULL)
    return;
  copy->type = token->type;
  copy->line = token->line;
  copy->column = token->column;
  copy->literal = token->literal ? strdup(token->literal) : NULL;
  parser->pending_comments[parser->pending_count++] = copy;
}

static void clear_pending_comments(MMC_Parser *parser) {
  for (size_t i = 0; i < parser->pending_count; i++) {
    mmc_token_free(parser->pending_comments[i]);
    free(parser->pending_comments[i]);
  }
  parser->pending_count = 0;
}

static void camel_to_snake(const char *input, char *output,
                           size_t output_size) {
  if (input == NULL || output == NULL || output_size == 0)
    return;
  size_t j = 0;
  for (size_t i = 0; input[i] != '\0' && j + 1 < output_size; i++) {
    char c = input[i];
    if (isupper((unsigned char)c)) {
      if (i > 0) {
        char prev = input[i - 1];
        int prev_upper = isupper((unsigned char)prev) != 0;
        int next_upper =
            (input[i + 1] != '\0' && isupper((unsigned char)input[i + 1]) != 0);
        if (!prev_upper || (input[i + 1] != '\0' && !next_upper)) {
          if (j + 1 < output_size)
            output[j++] = '_';
        }
      }
      output[j++] = (char)tolower((unsigned char)c);
    } else {
      output[j++] = c;
    }
  }
  output[j] = '\0';
}

static void tag_from_comment(const char *literal, MMC_Tag *tag) {
  if (literal == NULL || tag == NULL)
    return;
  // Check if literal starts with "mm:" (prefix-only, like Go's CutPrefix)
  const char *mm = literal;
  while (*mm == ' ' || *mm == '\t' || *mm == '\r')
    mm++;
  if (mm[0] != 'm' || mm[1] != 'm' || mm[2] != ':')
    return;
  mm += 3;
  while (*mm == ' ' || *mm == '\t')
    mm++;
  // Simple tag parsing: key=value;key=value
  char buf[256];
  size_t i = 0;
  while (*mm != '\0' && i < sizeof(buf) - 1) {
    if (*mm == '\n')
      break;
    buf[i++] = *mm;
    mm++;
  }
  buf[i] = '\0';

  // Parse key=value pairs
  char *saveptr;
  char *token = strtok_r(buf, ";", &saveptr);
  while (token != NULL) {
    while (*token == ' ' || *token == '\t')
      token++;
    char *eq = strchr(token, '=');
    if (eq != NULL) {
      *eq = '\0';
      char *key = token;
      char *value = eq + 1;
      while (*value == ' ' || *value == '\t')
        value++;
      // Trim trailing space from key
      char *kend = key + strlen(key) - 1;
      while (kend >= key && (*kend == ' ' || *kend == '\t'))
        *kend-- = '\0';

      if (strcmp(key, "type") == 0) {
        if (strcmp(value, "str") == 0)
          tag->type = MMC_VALUE_TYPE_STR;
        else if (strcmp(value, "bool") == 0)
          tag->type = MMC_VALUE_TYPE_BOOL;
        else if (strcmp(value, "i") == 0)
          tag->type = MMC_VALUE_TYPE_I;
        else if (strcmp(value, "i8") == 0)
          tag->type = MMC_VALUE_TYPE_I8;
        else if (strcmp(value, "i16") == 0)
          tag->type = MMC_VALUE_TYPE_I16;
        else if (strcmp(value, "i32") == 0)
          tag->type = MMC_VALUE_TYPE_I32;
        else if (strcmp(value, "i64") == 0)
          tag->type = MMC_VALUE_TYPE_I64;
        else if (strcmp(value, "u") == 0)
          tag->type = MMC_VALUE_TYPE_U;
        else if (strcmp(value, "u8") == 0)
          tag->type = MMC_VALUE_TYPE_U8;
        else if (strcmp(value, "u16") == 0)
          tag->type = MMC_VALUE_TYPE_U16;
        else if (strcmp(value, "u32") == 0)
          tag->type = MMC_VALUE_TYPE_U32;
        else if (strcmp(value, "u64") == 0)
          tag->type = MMC_VALUE_TYPE_U64;
        else if (strcmp(value, "f32") == 0)
          tag->type = MMC_VALUE_TYPE_F32;
        else if (strcmp(value, "f64") == 0)
          tag->type = MMC_VALUE_TYPE_F64;
        else if (strcmp(value, "bytes") == 0)
          tag->type = MMC_VALUE_TYPE_BYTES;
        else if (strcmp(value, "bigint") == 0)
          tag->type = MMC_VALUE_TYPE_BIGINT;
        else if (strcmp(value, "arr") == 0)
          tag->type = MMC_VALUE_TYPE_ARR;
        else if (strcmp(value, "vec") == 0)
          tag->type = MMC_VALUE_TYPE_VEC;
        else if (strcmp(value, "obj") == 0)
          tag->type = MMC_VALUE_TYPE_OBJ;
        else if (strcmp(value, "map") == 0)
          tag->type = MMC_VALUE_TYPE_MAP;
        else if (strcmp(value, "datetime") == 0)
          tag->type = MMC_VALUE_TYPE_DATETIME;
        else if (strcmp(value, "date") == 0)
          tag->type = MMC_VALUE_TYPE_DATE;
        else if (strcmp(value, "time") == 0)
          tag->type = MMC_VALUE_TYPE_TIME;
        else if (strcmp(value, "uuid") == 0)
          tag->type = MMC_VALUE_TYPE_UUID;
        else if (strcmp(value, "decimal") == 0)
          tag->type = MMC_VALUE_TYPE_DECIMAL;
        else if (strcmp(value, "ip") == 0)
          tag->type = MMC_VALUE_TYPE_IP;
        else if (strcmp(value, "url") == 0)
          tag->type = MMC_VALUE_TYPE_URL;
        else if (strcmp(value, "email") == 0)
          tag->type = MMC_VALUE_TYPE_EMAIL;
        else if (strcmp(value, "enums") == 0)
          tag->type = MMC_VALUE_TYPE_ENUMS;
        else if (strcmp(value, "media") == 0)
          tag->type = MMC_VALUE_TYPE_MEDIA;
      } else if (strcmp(key, "size") == 0) {
        tag->size = atoi(value);
      } else if (strcmp(key, "deprecated") == 0) {
        tag->deprecated = 1;
      } else if (strcmp(key, "is_null") == 0) {
        tag->is_null = 1;
      } else if (strcmp(key, "example") == 0) {
        tag->example = 1;
      } else if (strcmp(key, "enums") == 0) {
        if (tag->enums != NULL)
          free(tag->enums);
        tag->enums = strdup(value);
      }
    }
    token = strtok_r(NULL, ";", &saveptr);
  }
}

static MMC_Tag *consume_comments_for(MMC_Parser *parser, size_t anchor_line) {
  if (parser->pending_count == 0)
    return NULL;

  MMC_Token *last = parser->pending_comments[parser->pending_count - 1];
  if (anchor_line - last->line > 1) {
    clear_pending_comments(parser);
    return NULL;
  }

  MMC_Tag *tag = (MMC_Tag *)malloc(sizeof(MMC_Tag));
  if (tag == NULL)
    return NULL;
  memset(tag, 0, sizeof(MMC_Tag));

  for (size_t i = 0; i < parser->pending_count; i++) {
    MMC_Token *ct = parser->pending_comments[i];
    tag_from_comment(ct->literal, tag);
  }

  clear_pending_comments(parser);
  return tag;
}

static void next_token(MMC_Parser *parser) {
  if (parser->current_token != NULL) {
    mmc_token_free(parser->current_token);
    free(parser->current_token);
  }
  parser->current_token = mmc_scanner_next_token(parser->scanner);
}

static MMC_Node *parse_value(MMC_Parser *parser, MMC_Tag *inherit_tag) {
  if (parser->current_token == NULL) {
    parser_error(parser, "No token available");
    return NULL;
  }

  // Accumulate comments
  if (parser->current_token->type == MMC_TOKEN_COMMENT) {
    if (parser->pending_count > 0) {
      MMC_Token *last = parser->pending_comments[parser->pending_count - 1];
      if (parser->current_token->line - last->line > 1) {
        clear_pending_comments(parser);
      }
    }
    add_pending_comment(parser, parser->current_token);
    next_token(parser);
    return parse_value(parser, inherit_tag);
  }

  // Consume pending comments for this value
  size_t anchor_line = parser->current_token->line;
  MMC_Tag *pending_tag = consume_comments_for(parser, anchor_line);

  switch (parser->current_token->type) {
  case MMC_TOKEN_LBRACE: {
    // Pass pending_tag as the object's own tag, and inherit_tag for children
    MMC_Node *obj = parse_object(parser, pending_tag, inherit_tag);
    if (pending_tag)
      free(pending_tag);
    return obj;
  }
  case MMC_TOKEN_LBRACKET: {
    // Pass pending_tag as the array's own tag, and inherit_tag for children
    MMC_Node *arr = parse_array(parser, pending_tag, inherit_tag);
    if (pending_tag)
      free(pending_tag);
    return arr;
  }
  case MMC_TOKEN_STRING:
  case MMC_TOKEN_NUMBER:
  case MMC_TOKEN_TRUE:
  case MMC_TOKEN_FALSE:
  case MMC_TOKEN_NULL: {
    // Merge pending tag before parsing primitive
    MMC_Node *node =
        parse_primitive(parser, parser->current_token, pending_tag);
    if (node && node->type == MMC_NODE_VALUE && pending_tag) {
      // Merge additional tag fields not handled by parse_primitive
      if (pending_tag->type != MMC_VALUE_TYPE_UNKNOWN)
        node->data.value.tag.type = pending_tag->type;
      if (pending_tag->size != 0)
        node->data.value.tag.size = pending_tag->size;
      if (pending_tag->deprecated)
        node->data.value.tag.deprecated = 1;
      if (pending_tag->is_null)
        node->data.value.tag.is_null = 1;
      if (pending_tag->example)
        node->data.value.tag.example = 1;
      if (pending_tag->enums) {
        if (node->data.value.tag.enums)
          free(node->data.value.tag.enums);
        node->data.value.tag.enums = strdup(pending_tag->enums);
      }
    }
    if (pending_tag)
      free(pending_tag);
    return node;
  }
  case MMC_TOKEN_EOF:
    if (pending_tag)
      free(pending_tag);
    parser_error(parser, "Unexpected EOF");
    return NULL;
  default:
    if (pending_tag)
      free(pending_tag);
    parser_error(parser, "Unexpected token");
    return NULL;
  }
}

static MMC_Node *parse_primitive(MMC_Parser *parser, MMC_Token *token,
                                 MMC_Tag *pending_tag) {
  // null is not supported
  if (token->type == MMC_TOKEN_NULL) {
    parser_error(parser, "null is not supported");
    next_token(parser);
    return NULL;
  }

  MMC_Node *node = mmc_node_new_value();
  if (node == NULL)
    return NULL;

  node->data.value.text = token->literal ? strdup(token->literal) : NULL;

  // Determine effective tag type: use pending_tag type if set, otherwise infer
  // from token
  MMC_ValueType effective_type = MMC_VALUE_TYPE_UNKNOWN;
  if (pending_tag && pending_tag->type != MMC_VALUE_TYPE_UNKNOWN) {
    effective_type = pending_tag->type;
  }

  switch (token->type) {
  case MMC_TOKEN_STRING:
    if (effective_type != MMC_VALUE_TYPE_UNKNOWN) {
      // Use the tag type for dispatch (like Go parser)
      node->data.value.tag.type = effective_type;
    } else {
      node->data.value.tag.type = MMC_VALUE_TYPE_STR;
    }
    node->data.value.data_type = MMC_DATA_STRING;
    node->data.value.data.string_value =
        token->literal ? strdup(token->literal) : NULL;
    break;
  case MMC_TOKEN_NUMBER:
    if (effective_type != MMC_VALUE_TYPE_UNKNOWN) {
      // Use tag type for dispatch (like Go parser)
      node->data.value.tag.type = effective_type;
      if (strchr(token->literal, '.') != NULL) {
        node->data.value.data_type = MMC_DATA_FLOAT;
        node->data.value.data.float_value = atof(token->literal);
      } else {
        node->data.value.data_type = MMC_DATA_INT;
        node->data.value.data.int_value = atoll(token->literal);
      }
    } else if (strchr(token->literal, '.') != NULL) {
      node->data.value.data_type = MMC_DATA_FLOAT;
      node->data.value.data.float_value = atof(token->literal);
      node->data.value.tag.type = MMC_VALUE_TYPE_F64;
    } else if (token->literal[0] == '-') {
      node->data.value.data_type = MMC_DATA_INT;
      node->data.value.data.int_value = atoll(token->literal);
      node->data.value.tag.type = MMC_VALUE_TYPE_I64;
    } else {
      node->data.value.data_type = MMC_DATA_UINT;
      node->data.value.data.uint_value = strtoull(token->literal, NULL, 10);
      node->data.value.tag.type = MMC_VALUE_TYPE_U64;
    }
    break;
  case MMC_TOKEN_TRUE:
    node->data.value.data_type = MMC_DATA_BOOL;
    node->data.value.data.bool_value = 1;
    node->data.value.tag.type = MMC_VALUE_TYPE_BOOL;
    break;
  case MMC_TOKEN_FALSE:
    node->data.value.data_type = MMC_DATA_BOOL;
    node->data.value.data.bool_value = 0;
    node->data.value.tag.type = MMC_VALUE_TYPE_BOOL;
    break;
  default:
    mmc_node_free(node);
    parser_error(parser, "Invalid primitive type");
    return NULL;
  }

  next_token(parser);
  return node;
}

static MMC_Node *parse_object(MMC_Parser *parser, MMC_Tag *tag,
                              MMC_Tag *inherit_tag) {
  parser->depth++;
  if (parser->depth > MAX_DEPTH) {
    parser_error(parser, "max depth: 32");
    return NULL;
  }

  MMC_Node *node = mmc_node_new_object();
  if (node == NULL) {
    parser->depth--;
    return NULL;
  }

  // Apply tag from comments before the object
  if (tag && tag->type != MMC_VALUE_TYPE_UNKNOWN) {
    node->data.object.tag.type = tag->type;
    if (tag->size != 0)
      node->data.object.tag.size = tag->size;
    if (tag->deprecated)
      node->data.object.tag.deprecated = 1;
    if (tag->is_null)
      node->data.object.tag.is_null = 1;
    if (tag->example)
      node->data.object.tag.example = 1;
    if (tag->enums) {
      node->data.object.tag.enums = strdup(tag->enums);
    }
  } else {
    node->data.object.tag.type = MMC_VALUE_TYPE_OBJ;
  }

  // For map types, also inherit from parent context
  if (inherit_tag && inherit_tag->type != MMC_VALUE_TYPE_UNKNOWN &&
      node->data.object.tag.type == MMC_VALUE_TYPE_UNKNOWN) {
    node->data.object.tag.type = inherit_tag->type;
  }

  // If type is MAP, propagate parent tag to children via inherit_tag
  MMC_Tag *child_inherit_tag = NULL;
  if (node->data.object.tag.type == MMC_VALUE_TYPE_MAP) {
    child_inherit_tag = inherit_tag ? inherit_tag : &node->data.object.tag;
  } else {
    child_inherit_tag = inherit_tag;
  }

  next_token(parser);

  size_t capacity = 4;
  node->data.object.fields = (MMC_Field *)malloc(capacity * sizeof(MMC_Field));
  if (node->data.object.fields == NULL) {
    mmc_node_free(node);
    parser->depth--;
    return NULL;
  }
  node->data.object.field_count = 0;

  while (parser->current_token != NULL &&
         parser->current_token->type != MMC_TOKEN_RBRACE) {
    if (parser->current_token->type == MMC_TOKEN_COMMA) {
      next_token(parser);
      continue;
    }

    // Accumulate comments
    if (parser->current_token->type == MMC_TOKEN_COMMENT) {
      if (parser->pending_count > 0) {
        MMC_Token *last = parser->pending_comments[parser->pending_count - 1];
        if (parser->current_token->line - last->line > 1) {
          clear_pending_comments(parser);
        }
      }
      add_pending_comment(parser, parser->current_token);
      next_token(parser);
      continue;
    }

    if (parser->current_token->type != MMC_TOKEN_STRING) {
      if (parser->current_token->type != MMC_TOKEN_EOF) {
        next_token(parser);
      }
      continue;
    }

    // Convert key to snake_case
    char snake_key[256];
    camel_to_snake(parser->current_token->literal, snake_key,
                   sizeof(snake_key));
    char *key = strdup(snake_key);
    next_token(parser);

    if (parser->current_token == NULL ||
        parser->current_token->type != MMC_TOKEN_COLON) {
      free(key);
      continue;
    }
    next_token(parser);

    MMC_Node *value = parse_value(parser, child_inherit_tag);
    if (value == NULL) {
      free(key);
      continue;
    }

    if (node->data.object.field_count >= capacity) {
      capacity *= 2;
      MMC_Field *new_fields = (MMC_Field *)realloc(
          node->data.object.fields, capacity * sizeof(MMC_Field));
      if (new_fields == NULL) {
        free(key);
        mmc_node_free(value);
        continue;
      }
      node->data.object.fields = new_fields;
    }

    node->data.object.fields[node->data.object.field_count].key = key;
    node->data.object.fields[node->data.object.field_count].value = value;
    node->data.object.field_count++;
  }

  if (parser->current_token != NULL &&
      parser->current_token->type == MMC_TOKEN_RBRACE) {
    next_token(parser);
  }

  parser->depth--;
  return node;
}

static MMC_Node *parse_array(MMC_Parser *parser, MMC_Tag *tag,
                             MMC_Tag *inherit_tag) {
  parser->depth++;
  if (parser->depth > MAX_DEPTH) {
    parser_error(parser, "max depth: 32");
    return NULL;
  }

  MMC_Node *node = mmc_node_new_array();
  if (node == NULL) {
    parser->depth--;
    return NULL;
  }

  // Apply tag from comments before the array
  if (tag && tag->type != MMC_VALUE_TYPE_UNKNOWN) {
    node->data.array.tag.type = tag->type;
    if (tag->size != 0)
      node->data.array.tag.size = tag->size;
    if (tag->deprecated)
      node->data.array.tag.deprecated = 1;
  } else {
    node->data.array.tag.type = MMC_VALUE_TYPE_ARR;
  }

  // Inherit from parent for child items
  MMC_Tag *child_inherit_tag = inherit_tag;

  next_token(parser);

  size_t capacity = 4;
  node->data.array.items = (MMC_Node **)malloc(capacity * sizeof(MMC_Node *));
  if (node->data.array.items == NULL) {
    mmc_node_free(node);
    parser->depth--;
    return NULL;
  }
  node->data.array.item_count = 0;

  while (parser->current_token != NULL &&
         parser->current_token->type != MMC_TOKEN_RBRACKET) {
    if (parser->current_token->type == MMC_TOKEN_COMMA) {
      next_token(parser);
      continue;
    }

    // Accumulate comments
    if (parser->current_token->type == MMC_TOKEN_COMMENT) {
      if (parser->pending_count > 0) {
        MMC_Token *last = parser->pending_comments[parser->pending_count - 1];
        if (parser->current_token->line - last->line > 1) {
          clear_pending_comments(parser);
        }
      }
      add_pending_comment(parser, parser->current_token);
      next_token(parser);
      continue;
    }

    MMC_Node *value = parse_value(parser, child_inherit_tag);
    if (value != NULL) {
      if (node->data.array.item_count >= capacity) {
        capacity *= 2;
        MMC_Node **new_items = (MMC_Node **)realloc(
            node->data.array.items, capacity * sizeof(MMC_Node *));
        if (new_items == NULL) {
          mmc_node_free(value);
          continue;
        }
        node->data.array.items = new_items;
      }
      node->data.array.items[node->data.array.item_count] = value;
      node->data.array.item_count++;
    } else {
      break;
    }
  }

  if (parser->current_token != NULL &&
      parser->current_token->type == MMC_TOKEN_RBRACKET) {
    next_token(parser);
  }

  parser->depth--;
  return node;
}

MMC_Parser *mmc_parser_new(const char *input) {
  MMC_Parser *parser = (MMC_Parser *)malloc(sizeof(MMC_Parser));
  if (parser == NULL)
    return NULL;

  parser->scanner = mmc_scanner_new(input);
  if (parser->scanner == NULL) {
    free(parser);
    return NULL;
  }

  parser->current_token = NULL;
  parser->error_message = NULL;
  parser->depth = 0;
  parser->pending_comments = NULL;
  parser->pending_count = 0;
  parser->pending_capacity = 0;

  next_token(parser);
  return parser;
}

void mmc_parser_free(MMC_Parser *parser) {
  if (parser == NULL)
    return;
  if (parser->scanner != NULL) {
    mmc_scanner_free(parser->scanner);
  }
  if (parser->current_token != NULL) {
    mmc_token_free(parser->current_token);
    free(parser->current_token);
  }
  clear_pending_comments(parser);
  if (parser->pending_comments != NULL) {
    free(parser->pending_comments);
  }
  if (parser->error_message != NULL) {
    free(parser->error_message);
  }
  free(parser);
}

MMC_Node *mmc_parser_parse(MMC_Parser *parser) {
  if (parser == NULL || parser->current_token == NULL) {
    return NULL;
  }

  MMC_Node *result = parse_value(parser, NULL);

  if (parser->error_message != NULL) {
    mmc_node_free(result);
    return NULL;
  }

  return result;
}

const char *mmc_parser_get_error(MMC_Parser *parser) {
  if (parser == NULL)
    return "Parser is NULL";
  return parser->error_message ? parser->error_message : NULL;
}