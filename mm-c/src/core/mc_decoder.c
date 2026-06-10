#include "mc_decoder.h"
#include "../ir/mc_tag.h"
#include "../ir/mc_value_type.h"
#include "mc_constants.h"
#include <errno.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

static uint8_t dec_read_byte(mm_decoder_t *d) {
  if (d->offset >= d->size) {
    errno = EINVAL;
    return 0;
  }
  return d->data[d->offset++];
}

static const uint8_t *dec_read_bytes(mm_decoder_t *d, size_t len) {
  if (d->offset + len > d->size) {
    errno = EINVAL;
    return NULL;
  }
  const uint8_t *p = d->data + d->offset;
  d->offset += len;
  return p;
}

static uint64_t dec_read_uint64_be(mm_decoder_t *d, int byte_len) {
  uint64_t v = 0;
  for (int i = 0; i < byte_len; i++) {
    v = (v << 8) | dec_read_byte(d);
  }
  return v;
}

static const char *mime_id_to_string(uint8_t id) {
  switch (id) {
  case 1:
    return "image/jpeg";
  case 2:
    return "image/png";
  case 3:
    return "image/gif";
  case 4:
    return "image/webp";
  case 5:
    return "image/svg+xml";
  case 6:
    return "image/avif";
  case 7:
    return "image/bmp";
  case 8:
    return "image/x-icon";
  case 9:
    return "image/tiff";
  case 10:
    return "image/heic";
  case 11:
    return "image/heif";
  case 12:
    return "text/plain";
  case 13:
    return "text/html";
  case 14:
    return "text/css";
  case 15:
    return "text/javascript";
  case 16:
    return "application/json";
  case 17:
    return "text/csv";
  case 18:
    return "text/markdown";
  case 19:
    return "application/pdf";
  case 20:
    return "application/zip";
  case 21:
    return "application/gzip";
  case 22:
    return "application/x-tar";
  case 23:
    return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";
  case 24:
    return "application/"
           "vnd.openxmlformats-officedocument.wordprocessingml.document";
  case 25:
    return "application/octet-stream";
  case 26:
    return "video/mp4";
  case 27:
    return "video/webm";
  case 28:
    return "video/mov";
  case 29:
    return "audio/mpeg";
  case 30:
    return "audio/wav";
  case 31:
    return "audio/flac";
  case 32:
    return "font/woff2";
  case 33:
    return "font/ttf";
  default:
    return "";
  }
}

static const char base64_table[] =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

static int base64_decode_char(char c) {
  if (c >= 'A' && c <= 'Z')
    return c - 'A';
  if (c >= 'a' && c <= 'z')
    return c - 'a' + 26;
  if (c >= '0' && c <= '9')
    return c - '0' + 52;
  if (c == '+')
    return 62;
  if (c == '/')
    return 63;
  return 0xFF;
}

static uint8_t *mc_base64_decode(const char *text, size_t text_len,
                                 size_t *decoded_len) {
  if (!text || text_len == 0) {
    *decoded_len = 0;
    return NULL;
  }
  size_t padding = 0;
  if (text_len > 0 && text[text_len - 1] == '=')
    padding++;
  if (text_len > 1 && text[text_len - 2] == '=')
    padding++;
  size_t out_len = text_len / 4 * 3 - padding;
  uint8_t *out = (uint8_t *)malloc(out_len + 1);
  if (!out) {
    *decoded_len = 0;
    return NULL;
  }
  size_t o = 0;
  for (size_t i = 0; i < text_len; i += 4) {
    uint8_t a = (uint8_t)base64_decode_char(text[i]);
    uint8_t b =
        (i + 1 < text_len) ? (uint8_t)base64_decode_char(text[i + 1]) : 0;
    uint8_t c =
        (i + 2 < text_len) ? (uint8_t)base64_decode_char(text[i + 2]) : 0;
    uint8_t d =
        (i + 3 < text_len) ? (uint8_t)base64_decode_char(text[i + 3]) : 0;
    if (o < out_len)
      out[o++] = (uint8_t)((a << 2) | (b >> 4));
    if (o < out_len)
      out[o++] = (uint8_t)((b << 4) | (c >> 2));
    if (o < out_len)
      out[o++] = (uint8_t)((c << 6) | d);
  }
  *decoded_len = out_len;
  return out;
}

static char *mc_base64_encode(const uint8_t *data, size_t len) {
  if (!data || len == 0) {
    char *empty = (char *)malloc(1);
    if (empty)
      empty[0] = '\0';
    return empty;
  }
  size_t out_len = (len + 2) / 3 * 4;
  char *out = (char *)malloc(out_len + 1);
  if (!out)
    return NULL;
  size_t i = 0, o = 0;
  while (i < len) {
    uint32_t n = (uint32_t)data[i] << 16;
    if (i + 1 < len)
      n |= (uint32_t)data[i + 1] << 8;
    if (i + 2 < len)
      n |= data[i + 2];
    out[o++] = base64_table[(n >> 18) & 0x3F];
    out[o++] = base64_table[(n >> 12) & 0x3F];
    out[o++] = (i + 1 < len) ? base64_table[(n >> 6) & 0x3F] : '=';
    out[o++] = (i + 2 < len) ? base64_table[n & 0x3F] : '=';
    i += 3;
  }
  out[out_len] = '\0';
  return out;
}

static node_t *dec_decode_simple(mm_decoder_t *d, uint8_t b) {
  (void)d;
  int suffix = mm_suffix_of(b);
  node_t *node = node_new_scalar();
  node_scalar_t *val = &node->data.value;
  mm_tag_init(&val->tag);

  switch (suffix) {
  case MM_SIMPLE_NULL:
    node_free(node);
    node = node_new_null();
    break;
  case MM_SIMPLE_NULLBOOL:
    val->tag.type = MM_VALUE_BOOL;
    val->text = strdup("false");
    break;
  case MM_SIMPLE_NULLINT:
    val->tag.type = MM_VALUE_I;
    val->text = strdup("0");
    break;
  case MM_SIMPLE_NULLFLOAT:
    val->tag.type = MM_VALUE_F64;
    val->text = strdup("0.0");
    break;
  case MM_SIMPLE_NULLSTRING:
    val->tag.type = MM_VALUE_STR;
    val->text = strdup("");
    break;
  case MM_SIMPLE_NULLBYTES:
    val->tag.type = MM_VALUE_BYTES;
    val->text = strdup("");
    break;
  case MM_SIMPLE_FALSE:
    val->tag.type = MM_VALUE_BOOL;
    val->text = strdup("false");
    break;
  case MM_SIMPLE_TRUE:
    val->tag.type = MM_VALUE_BOOL;
    val->text = strdup("true");
    break;
  default:
    val->tag.type = MM_VALUE_STR;
    val->text = strdup("");
    break;
  }
  return node;
}

static node_t *dec_decode_int(mm_decoder_t *d, uint8_t b, int is_positive) {
  int byte_len = mm_int_len(b);
  uint64_t uv;

  if (byte_len == 0) {
    uv = (uint64_t)mm_suffix_of(b);
  } else {
    uv = dec_read_uint64_be(d, byte_len);
  }

  node_t *node = node_new_scalar();
  node_scalar_t *val = &node->data.value;
  mm_tag_init(&val->tag);

  if (is_positive) {
    val->tag.type = MM_VALUE_I;
    char buf[24];
    snprintf(buf, sizeof(buf), "%llu", (unsigned long long)uv);
    val->text = strdup(buf);
  } else {
    val->tag.type = MM_VALUE_I;
    char buf[24];
    snprintf(buf, sizeof(buf), "-%llu", (unsigned long long)uv);
    val->text = strdup(buf);
  }

  return node;
}

static node_t *dec_decode_float(mm_decoder_t *d, uint8_t b) {
  int byte_len = mm_float_len(b);
  int suffix = mm_suffix_of(b);
  int is_negative = (b & 0x10) != 0;

  double v;

  if (byte_len == 0) {
    int val = suffix & 0x0f;
    v = (double)val / 10.0;
    if (is_negative) {
      v = -v;
    }
  } else {
    int8_t exp = (int8_t)dec_read_byte(d);
    uint64_t mantissa;
    if (byte_len >= 1 && byte_len <= 8) {
      mantissa = dec_read_uint64_be(d, byte_len);
    } else {
      mantissa = 0;
    }
    char buf[128];
    char mantissa_str[32];
    snprintf(mantissa_str, sizeof(mantissa_str), "%llu",
             (unsigned long long)mantissa);
    int decimal_pos = (int)strlen(mantissa_str) + (int)exp;
    if (decimal_pos <= 0) {
      char pad[64];
      int pad_len = -decimal_pos;
      memset(pad, '0', (size_t)pad_len);
      pad[pad_len] = '\0';
      snprintf(buf, sizeof(buf), "0.%s%s", pad, mantissa_str);
    } else if (decimal_pos > 0 && decimal_pos < (int)strlen(mantissa_str)) {
      size_t m_len = strlen(mantissa_str);
      size_t int_part = (size_t)decimal_pos;
      char int_buf[64], frac_buf[64];
      memcpy(int_buf, mantissa_str, int_part);
      int_buf[int_part] = '\0';
      memcpy(frac_buf, mantissa_str + int_part, m_len - int_part + 1);
      snprintf(buf, sizeof(buf), "%s.%s", int_buf, frac_buf);
    } else {
      int trailing = decimal_pos - (int)strlen(mantissa_str);
      char pad[64];
      memset(pad, '0', (size_t)trailing);
      pad[trailing] = '\0';
      snprintf(buf, sizeof(buf), "%s%s", mantissa_str, pad);
    }
    if (is_negative) {
      char neg_buf[132];
      snprintf(neg_buf, sizeof(neg_buf), "-%s", buf);
      v = -strtod(buf, NULL);
    } else {
      v = strtod(buf, NULL);
    }
  }

  node_t *node = node_new_scalar();
  node_scalar_t *val = &node->data.value;
  mm_tag_init(&val->tag);
  val->tag.type = MM_VALUE_F64;
  char text_buf[64];
  snprintf(text_buf, sizeof(text_buf), "%g", v);
  val->text = strdup(text_buf);

  return node;
}

static node_t *dec_decode_string(mm_decoder_t *d, uint8_t b) {
  int extra_len = mm_string_extra_len(b);
  int inline_len = mm_string_inline_len(b);
  size_t str_len;

  if (extra_len == 0) {
    str_len = (size_t)inline_len;
  } else if (extra_len == 1) {
    str_len = (size_t)dec_read_byte(d);
  } else {
    uint8_t hi = dec_read_byte(d);
    uint8_t lo = dec_read_byte(d);
    str_len = ((size_t)hi << 8) | (size_t)lo;
  }

  char *text = NULL;
  if (str_len > 0) {
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (!raw) {
      text = strdup("");
    } else {
      text = (char *)malloc(str_len + 1);
      memcpy(text, raw, str_len);
      text[str_len] = '\0';
    }
  } else {
    text = strdup("");
  }

  node_t *node = node_new_scalar();
  node_scalar_t *val = &node->data.value;
  mm_tag_init(&val->tag);
  val->tag.type = MM_VALUE_STR;
  val->text = text;

  return node;
}

static char *dec_decode_bigint(const uint8_t *data, size_t len) {
  if (!data || len <= 1) {
    return strdup("0");
  }

  // data[0] is digit_count
  size_t digit_len = (size_t)data[0];
  if (digit_len == 0) {
    return strdup("0");
  }

  // Bit-packed data starts at data[1]
  const uint8_t *bits = data + 1;
  size_t bits_len = len - 1;
  size_t bit_pos = 0;

  // Read sign bit
  int byte_idx = (int)(bit_pos / 8);
  int bit_in_byte = 7 - (int)(bit_pos % 8);
  int sign = (bits[byte_idx] >> bit_in_byte) & 1;
  bit_pos++;

  // Buffer for result digits
  char *result = (char *)malloc(digit_len + 2);
  size_t pos = 0;

  size_t remaining = digit_len;
  while (remaining > 0) {
    int val = 0;
    int num_bits;
    if (remaining >= 3) {
      num_bits = 10;
    } else if (remaining == 2) {
      num_bits = 7;
    } else {
      num_bits = 4;
    }

    // Read num_bits bits
    for (int i = 0; i < num_bits; i++) {
      if (bit_pos >= bits_len * 8)
        break;
      byte_idx = (int)(bit_pos / 8);
      bit_in_byte = 7 - (int)(bit_pos % 8);
      int b = (bits[byte_idx] >> bit_in_byte) & 1;
      val = (val << 1) | b;
      bit_pos++;
    }

    // Format with leading zeros
    if (num_bits == 10) {
      pos += sprintf(result + pos, "%03d", val);
      remaining -= 3;
    } else if (num_bits == 7) {
      pos += sprintf(result + pos, "%02d", val);
      remaining -= 2;
    } else {
      pos += sprintf(result + pos, "%d", val);
      remaining -= 1;
    }
  }

  result[pos] = '\0';

  // Trim leading zeros
  char *trimmed = result;
  while (*trimmed == '0' && *(trimmed + 1) != '\0') {
    trimmed++;
  }

  char *final_str;
  if (sign) {
    final_str = (char *)malloc(strlen(trimmed) + 2);
    final_str[0] = '-';
    strcpy(final_str + 1, trimmed);
  } else {
    final_str = strdup(trimmed);
  }

  if (trimmed != result) {
    // trimmed points into result, free result instead
    free(result);
  } else {
    free(result);
  }

  return final_str;
}

static node_t *dec_decode_bytes(mm_decoder_t *d, uint8_t b, mm_tag_t *tag) {
  int extra_len = mm_bytes_extra_len(b);
  int inline_len = mm_bytes_inline_len(b);
  size_t bytes_len;

  if (extra_len == 0) {
    bytes_len = (size_t)inline_len;
  } else if (extra_len == 1) {
    bytes_len = (size_t)dec_read_byte(d);
  } else {
    uint8_t hi = dec_read_byte(d);
    uint8_t lo = dec_read_byte(d);
    bytes_len = ((size_t)hi << 8) | (size_t)lo;
  }

  node_t *node = node_new_scalar();
  node_scalar_t *val = &node->data.value;
  mm_tag_init(&val->tag);
  if (tag) {
    val->tag.type = tag->type;
    val->tag.is_inherit = tag->is_inherit;
  } else {
    val->tag.type = MM_VALUE_BYTES;
  }
  val->text = NULL;

  if (bytes_len > 0) {
    const uint8_t *raw = dec_read_bytes(d, bytes_len);
    if (raw) {
      if (tag && tag->type == MM_VALUE_BIGINT) {
        val->text = dec_decode_bigint(raw, bytes_len);
      } else if (tag && tag->type == MM_VALUE_UUID && bytes_len == 16) {
        char uuid_buf[37];
        snprintf(uuid_buf, sizeof(uuid_buf),
                 "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%"
                 "02x%02x%02x%02x",
                 raw[0], raw[1], raw[2], raw[3], raw[4], raw[5], raw[6], raw[7],
                 raw[8], raw[9], raw[10], raw[11], raw[12], raw[13], raw[14],
                 raw[15]);
        val->text = strdup(uuid_buf);
      } else {
        val->text = mc_base64_encode(raw, bytes_len);
      }
    }
  }
  if (!val->text) {
    val->text = strdup("");
  }

  return node;
}

static node_t *dec_decode_node(mm_decoder_t *d, mm_tag_t *parent_tag);

static node_t *dec_decode_array(mm_decoder_t *d, size_t total_len,
                                mm_tag_t *parent_tag) {
  node_t *node = node_new_array();
  node_array_t *arr = &node->data.array;
  mm_tag_init(&arr->tag);

  if (parent_tag) {
    // Propagate child_* attributes to arr->tag's child_* for children to
    // inherit
    if (parent_tag->child_desc) {
      arr->tag.child_desc = strdup(parent_tag->child_desc);
    }
    if (parent_tag->child_type != MM_VALUE_UNKNOWN) {
      arr->tag.child_type = parent_tag->child_type;
    }
    if (parent_tag->child_nullable) {
      arr->tag.child_nullable = true;
    }
    if (parent_tag->child_allow_empty) {
      arr->tag.child_allow_empty = true;
    }
    if (parent_tag->child_unique) {
      arr->tag.child_unique = true;
    }
    if (parent_tag->child_default_val) {
      arr->tag.child_default_val = strdup(parent_tag->child_default_val);
    }
    if (parent_tag->child_min) {
      arr->tag.child_min = strdup(parent_tag->child_min);
    }
    if (parent_tag->child_max) {
      arr->tag.child_max = strdup(parent_tag->child_max);
    }
    if (parent_tag->child_size != 0) {
      arr->tag.child_size = parent_tag->child_size;
    }
    if (parent_tag->child_enums) {
      arr->tag.child_enums = strdup(parent_tag->child_enums);
      arr->tag.child_type = MM_VALUE_ENUMS;
    }
    if (parent_tag->child_pattern) {
      arr->tag.child_pattern = strdup(parent_tag->child_pattern);
    }
    if (parent_tag->child_location_offset != 0) {
      arr->tag.child_location_offset = parent_tag->child_location_offset;
    }
    if (parent_tag->child_version != MM_TAG_DEFAULT_VERSION) {
      arr->tag.child_version = parent_tag->child_version;
    }
    if (parent_tag->child_mime) {
      arr->tag.child_mime = strdup(parent_tag->child_mime);
      arr->tag.child_type = MM_VALUE_MEDIA;
    }

    // Copy parent's own type/size to the array tag
    if (parent_tag->type != MM_VALUE_UNKNOWN) {
      arr->tag.type = parent_tag->type;
    }
    if (parent_tag->size != 0) {
      arr->tag.size = parent_tag->size;
    }
    if (!arr->tag.enums && parent_tag->enums) {
      arr->tag.enums = strdup(parent_tag->enums);
    }
    if (!arr->tag.mime && parent_tag->mime) {
      arr->tag.mime = strdup(parent_tag->mime);
    }
  }
  // Infer type from size if still unknown
  if (arr->tag.type == MM_VALUE_UNKNOWN) {
    if (arr->tag.size > 0) {
      arr->tag.type = MM_VALUE_ARR;
    } else {
      arr->tag.type = MM_VALUE_VEC;
    }
  }
  arr->tag.is_inherit = false;

  size_t end_offset = d->offset + total_len;
  while (d->offset < end_offset) {
    mm_tag_t item_tag;
    mm_tag_init(&item_tag);
    if (parent_tag) {
      mm_tag_inherit(&item_tag, parent_tag);
      if (parent_tag->child_type == MM_VALUE_UNKNOWN) {
        if (item_tag.type == MM_VALUE_UNKNOWN) {
          item_tag.type = parent_tag->type;
        }
        if (!item_tag.enums && parent_tag->enums) {
          item_tag.enums = strdup(parent_tag->enums);
        }
        if (!item_tag.mime && parent_tag->mime) {
          item_tag.mime = strdup(parent_tag->mime);
        }
        if (item_tag.version == MM_TAG_DEFAULT_VERSION &&
            parent_tag->version != MM_TAG_DEFAULT_VERSION) {
          item_tag.version = parent_tag->version;
        }
        if (item_tag.location_offset == 0 && parent_tag->location_offset != 0) {
          item_tag.location_offset = parent_tag->location_offset;
        }
      }
    }
    node_t *item = dec_decode_node(d, &item_tag);
    mm_tag_cleanup(&item_tag);
    if (!item)
      break;
    node_array_add_item(node, item);
  }

  if (d->offset < end_offset) {
    d->offset = end_offset;
  }

  return node;
}

static node_t *dec_decode_object(mm_decoder_t *d, size_t total_len,
                                 mm_tag_t *parent_tag) {
  node_t *node = node_new_object();
  node_object_t *obj = &node->data.object;
  mm_tag_init(&obj->tag);

  if (parent_tag) {
    // Propagate child_* attributes to obj->tag's child_* for children to
    // inherit
    if (parent_tag->child_desc) {
      obj->tag.child_desc = strdup(parent_tag->child_desc);
    }
    if (parent_tag->child_type != MM_VALUE_UNKNOWN) {
      obj->tag.child_type = parent_tag->child_type;
    }
    if (parent_tag->child_nullable) {
      obj->tag.child_nullable = true;
    }
    if (parent_tag->child_allow_empty) {
      obj->tag.child_allow_empty = true;
    }
    if (parent_tag->child_unique) {
      obj->tag.child_unique = true;
    }
    if (parent_tag->child_default_val) {
      obj->tag.child_default_val = strdup(parent_tag->child_default_val);
    }
    if (parent_tag->child_min) {
      obj->tag.child_min = strdup(parent_tag->child_min);
    }
    if (parent_tag->child_max) {
      obj->tag.child_max = strdup(parent_tag->child_max);
    }
    if (parent_tag->child_size != 0) {
      obj->tag.child_size = parent_tag->child_size;
    }
    if (parent_tag->child_enums) {
      obj->tag.child_enums = strdup(parent_tag->child_enums);
      obj->tag.child_type = MM_VALUE_ENUMS;
    }
    if (parent_tag->child_pattern) {
      obj->tag.child_pattern = strdup(parent_tag->child_pattern);
    }
    if (parent_tag->child_location_offset != 0) {
      obj->tag.child_location_offset = parent_tag->child_location_offset;
    }
    if (parent_tag->child_version != MM_TAG_DEFAULT_VERSION) {
      obj->tag.child_version = parent_tag->child_version;
    }
    if (parent_tag->child_mime) {
      obj->tag.child_mime = strdup(parent_tag->child_mime);
      obj->tag.child_type = MM_VALUE_MEDIA;
    }

    // Copy parent's own type/size to the object tag
    if (parent_tag->type != MM_VALUE_UNKNOWN) {
      obj->tag.type = parent_tag->type;
    }
    if (parent_tag->size != 0) {
      obj->tag.size = parent_tag->size;
    }
    if (!obj->tag.enums && parent_tag->enums) {
      obj->tag.enums = strdup(parent_tag->enums);
    }
    if (!obj->tag.mime && parent_tag->mime) {
      obj->tag.mime = strdup(parent_tag->mime);
    }
  }
  if (obj->tag.type == MM_VALUE_UNKNOWN) {
    obj->tag.type = MM_VALUE_OBJ;
  }
  obj->tag.is_inherit = false;

  size_t start_offset = d->offset;
  size_t end_offset = start_offset + total_len;

  node_t *key_array = mm_decoder_decode(d);
  if (!key_array || key_array->type != MM_NODE_ARRAY) {
    if (key_array)
      node_free(key_array);
    d->offset = end_offset;
    return node;
  }

  for (size_t i = 0;
       i < key_array->data.array.item_count && d->offset < end_offset; i++) {
    node_t *key_item = key_array->data.array.items[i];
    if (key_item->type != MM_NODE_VALUE || !key_item->data.value.text) {
      mm_decoder_decode(d);
      continue;
    }
    mm_tag_t field_tag;
    mm_tag_init(&field_tag);
    if (parent_tag) {
      mm_tag_inherit(&field_tag, parent_tag);
      if (parent_tag->child_type == MM_VALUE_UNKNOWN) {
        if (field_tag.type == MM_VALUE_UNKNOWN) {
          field_tag.type = parent_tag->type;
        }
        if (!field_tag.enums && parent_tag->enums) {
          field_tag.enums = strdup(parent_tag->enums);
        }
        if (!field_tag.mime && parent_tag->mime) {
          field_tag.mime = strdup(parent_tag->mime);
        }
        if (field_tag.version == MM_TAG_DEFAULT_VERSION &&
            parent_tag->version != MM_TAG_DEFAULT_VERSION) {
          field_tag.version = parent_tag->version;
        }
        if (field_tag.location_offset == 0 &&
            parent_tag->location_offset != 0) {
          field_tag.location_offset = parent_tag->location_offset;
        }
      }
    }
    node_t *val_node = dec_decode_node(d, &field_tag);
    mm_tag_cleanup(&field_tag);
    if (val_node) {
      node_object_add_field(node, key_item->data.value.text, val_node);
    }
  }

  node_free(key_array);
  d->offset = end_offset;

  return node;
}

static node_t *dec_decode_container(mm_decoder_t *d, uint8_t b,
                                    mm_tag_t *parent_tag) {
  int is_array = mm_is_array_container(b);
  int extra_len = mm_container_extra_len(b);
  int inline_len = mm_container_inline_len(b);
  size_t container_len;

  if (extra_len == 0) {
    container_len = (size_t)inline_len;
  } else if (extra_len == 1) {
    container_len = (size_t)dec_read_byte(d);
  } else {
    uint8_t hi = dec_read_byte(d);
    uint8_t lo = dec_read_byte(d);
    container_len = ((size_t)hi << 8) | (size_t)lo;
  }

  if (is_array) {
    return dec_decode_array(d, container_len, parent_tag);
  } else {
    return dec_decode_object(d, container_len, parent_tag);
  }
}

static bool dec_parse_one_tag_entry(mm_decoder_t *d, mm_tag_t *tag) {
  uint8_t b = dec_read_byte(d);
  if (errno)
    return false;

  int key = b & 0xF8;
  int payload = b & 0x07;

  switch (key) {
  case MM_TAG_KISNULL:
    tag->is_null = (payload & 0x01) != 0;
    if (tag->is_null) {
      tag->nullable = true;
    }
    return true;

  case MM_TAG_KEXAMPLE:
    tag->example = (payload & 0x01) != 0;
    return true;

  case MM_TAG_KDESC: {
    size_t str_len;
    if (payload <= 5) {
      str_len = (size_t)payload;
    } else if (payload == 6) {
      str_len = (size_t)dec_read_byte(d);
    } else {
      uint8_t hi = dec_read_byte(d);
      uint8_t lo = dec_read_byte(d);
      str_len = ((size_t)hi << 8) | (size_t)lo;
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->desc);
      tag->desc = (char *)malloc(str_len + 1);
      memcpy(tag->desc, raw, str_len);
      tag->desc[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KTYPE:
    tag->type = (mm_value_type_t)dec_read_byte(d);
    return true;

  case MM_TAG_KDEPRECATED:
    tag->deprecated = (payload & 0x01) != 0;
    return true;

  case MM_TAG_KNULLABLE:
    tag->nullable = (payload & 0x01) != 0;
    return true;

  case MM_TAG_KALLOWEMPTY:
    tag->allow_empty = (payload & 0x01) != 0;
    return true;

  case MM_TAG_KUNIQUE:
    tag->unique = (payload & 0x01) != 0;
    return true;

  case MM_TAG_KDEFAULTVAL: {
    size_t str_len;
    if (payload < 7) {
      str_len = (size_t)payload;
    } else {
      str_len = (size_t)dec_read_byte(d);
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->default_val);
      tag->default_val = (char *)malloc(str_len + 1);
      memcpy(tag->default_val, raw, str_len);
      tag->default_val[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KMIN: {
    size_t str_len;
    if (payload < 7) {
      str_len = (size_t)payload;
    } else {
      str_len = (size_t)dec_read_byte(d);
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->min);
      tag->min = (char *)malloc(str_len + 1);
      memcpy(tag->min, raw, str_len);
      tag->min[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KMAX: {
    size_t str_len;
    if (payload < 7) {
      str_len = (size_t)payload;
    } else {
      str_len = (size_t)dec_read_byte(d);
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->max);
      tag->max = (char *)malloc(str_len + 1);
      memcpy(tag->max, raw, str_len);
      tag->max[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KSIZE: {
    int nbytes = payload + 1;
    uint64_t val = 0;
    for (int i = 0; i < nbytes; i++) {
      val = (val << 8) | dec_read_byte(d);
    }
    tag->size = (int)val;
    return true;
  }

  case MM_TAG_KENUMS: {
    tag->type = MM_VALUE_ENUMS;
    size_t str_len;
    if (payload <= 5) {
      str_len = (size_t)payload;
    } else if (payload == 6) {
      str_len = (size_t)dec_read_byte(d);
    } else {
      uint8_t hi = dec_read_byte(d);
      uint8_t lo = dec_read_byte(d);
      str_len = ((size_t)hi << 8) | (size_t)lo;
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->enums);
      tag->enums = (char *)malloc(str_len + 1);
      memcpy(tag->enums, raw, str_len);
      tag->enums[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KPATTERN: {
    size_t str_len;
    if (payload < 7) {
      str_len = (size_t)payload;
    } else {
      str_len = (size_t)dec_read_byte(d);
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->pattern);
      tag->pattern = (char *)malloc(str_len + 1);
      memcpy(tag->pattern, raw, str_len);
      tag->pattern[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KLOCATION: {
    size_t str_len = (size_t)payload;
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      char buf[16];
      memcpy(buf, raw, str_len);
      buf[str_len] = '\0';
      tag->location_offset = (int)strtol(buf, NULL, 10);
    }
    return true;
  }

  case MM_TAG_KVERSION: {
    int nbytes = payload + 1;
    uint64_t val = 0;
    for (int i = 0; i < nbytes; i++) {
      val = (val << 8) | dec_read_byte(d);
    }
    tag->version = (int)val;
    return true;
  }

  case MM_TAG_KMIME: {
    int nbytes = payload + 1;
    uint8_t mime_id = 0;
    for (int i = 0; i < nbytes; i++) {
      mime_id = dec_read_byte(d);
    }
    free(tag->mime);
    tag->mime = strdup(mime_id_to_string(mime_id));
    tag->type = MM_VALUE_MEDIA;
    return true;
  }

  case MM_TAG_KCHILDDESC: {
    size_t str_len;
    if (payload <= 5) {
      str_len = (size_t)payload;
    } else if (payload == 6) {
      str_len = (size_t)dec_read_byte(d);
    } else {
      uint8_t hi = dec_read_byte(d);
      uint8_t lo = dec_read_byte(d);
      str_len = ((size_t)hi << 8) | (size_t)lo;
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->child_desc);
      tag->child_desc = (char *)malloc(str_len + 1);
      memcpy(tag->child_desc, raw, str_len);
      tag->child_desc[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KCHILDTYPE:
    tag->child_type = (mm_value_type_t)dec_read_byte(d);
    return true;

  case MM_TAG_KCHILDNULLABLE:
    tag->child_nullable = (payload & 0x01) != 0;
    return true;

  case MM_TAG_KCHILDALLOWEMPTY:
    tag->child_allow_empty = (payload & 0x01) != 0;
    return true;

  case MM_TAG_KCHILDUNIQUE:
    tag->child_unique = (payload & 0x01) != 0;
    return true;

  case MM_TAG_KCHILDDEFAULTVAL: {
    size_t str_len;
    if (payload < 7) {
      str_len = (size_t)payload;
    } else {
      str_len = (size_t)dec_read_byte(d);
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->child_default_val);
      tag->child_default_val = (char *)malloc(str_len + 1);
      memcpy(tag->child_default_val, raw, str_len);
      tag->child_default_val[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KCHILDMIN: {
    size_t str_len;
    if (payload < 7) {
      str_len = (size_t)payload;
    } else {
      str_len = (size_t)dec_read_byte(d);
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->child_min);
      tag->child_min = (char *)malloc(str_len + 1);
      memcpy(tag->child_min, raw, str_len);
      tag->child_min[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KCHILDMAX: {
    size_t str_len;
    if (payload < 7) {
      str_len = (size_t)payload;
    } else {
      str_len = (size_t)dec_read_byte(d);
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->child_max);
      tag->child_max = (char *)malloc(str_len + 1);
      memcpy(tag->child_max, raw, str_len);
      tag->child_max[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KCHILDSIZE: {
    int nbytes = payload + 1;
    uint64_t val = 0;
    for (int i = 0; i < nbytes; i++) {
      val = (val << 8) | dec_read_byte(d);
    }
    tag->child_size = (int)val;
    return true;
  }

  case MM_TAG_KCHILDENUMS: {
    tag->child_type = MM_VALUE_ENUMS;
    size_t str_len;
    if (payload <= 5) {
      str_len = (size_t)payload;
    } else if (payload == 6) {
      str_len = (size_t)dec_read_byte(d);
    } else {
      uint8_t hi = dec_read_byte(d);
      uint8_t lo = dec_read_byte(d);
      str_len = ((size_t)hi << 8) | (size_t)lo;
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->child_enums);
      tag->child_enums = (char *)malloc(str_len + 1);
      memcpy(tag->child_enums, raw, str_len);
      tag->child_enums[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KCHILDPATTERN: {
    size_t str_len;
    if (payload < 7) {
      str_len = (size_t)payload;
    } else {
      str_len = (size_t)dec_read_byte(d);
    }
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      free(tag->child_pattern);
      tag->child_pattern = (char *)malloc(str_len + 1);
      memcpy(tag->child_pattern, raw, str_len);
      tag->child_pattern[str_len] = '\0';
    }
    return true;
  }

  case MM_TAG_KCHILDLOCATION: {
    size_t str_len = (size_t)payload;
    const uint8_t *raw = dec_read_bytes(d, str_len);
    if (raw) {
      char buf[16];
      memcpy(buf, raw, str_len);
      buf[str_len] = '\0';
      tag->child_location_offset = (int)strtol(buf, NULL, 10);
    }
    return true;
  }

  case MM_TAG_KCHILDVERSION: {
    int nbytes = payload + 1;
    uint64_t val = 0;
    for (int i = 0; i < nbytes; i++) {
      val = (val << 8) | dec_read_byte(d);
    }
    tag->child_version = (int)val;
    return true;
  }

  case MM_TAG_KCHILDMIME: {
    int nbytes = payload + 1;
    uint8_t mime_id = 0;
    for (int i = 0; i < nbytes; i++) {
      mime_id = dec_read_byte(d);
    }
    free(tag->child_mime);
    tag->child_mime = strdup(mime_id_to_string(mime_id));
    tag->child_type = MM_VALUE_MEDIA;
    return true;
  }

  default:
    return true;
  }
}

static mm_tag_t dec_read_tag_bytes(mm_decoder_t *d) {
  mm_tag_t tag;
  mm_tag_init(&tag);

  uint8_t first = dec_read_byte(d);
  if (errno)
    return tag;

  size_t tag_bytes_len;
  if (first < 254) {
    tag_bytes_len = (size_t)first;
  } else if (first == 254) {
    tag_bytes_len = (size_t)dec_read_byte(d);
  } else {
    uint8_t hi = dec_read_byte(d);
    uint8_t lo = dec_read_byte(d);
    tag_bytes_len = ((size_t)hi << 8) | (size_t)lo;
  }

  size_t end_offset = d->offset + tag_bytes_len;

  while (d->offset < end_offset) {
    if (!dec_parse_one_tag_entry(d, &tag)) {
      break;
    }
    if (errno)
      break;
  }

  if (d->offset < end_offset) {
    d->offset = end_offset;
  }

  return tag;
}

static node_t *dec_decode_tag(mm_decoder_t *d, uint8_t b,
                              mm_tag_t *parent_tag) {
  int extra_len = mm_tag_extra_len(b);
  int inline_len = mm_tag_inline_len(b);
  size_t total_len;

  if (extra_len == 0) {
    total_len = (size_t)inline_len;
  } else if (extra_len == 1) {
    total_len = (size_t)dec_read_byte(d);
  } else {
    uint8_t hi = dec_read_byte(d);
    uint8_t lo = dec_read_byte(d);
    total_len = ((size_t)hi << 8) | (size_t)lo;
  }

  size_t tag_data_start = d->offset;

  mm_tag_t tag = dec_read_tag_bytes(d);

  size_t tag_bytes_consumed = d->offset - tag_data_start;

  node_t *inner = NULL;

  if (tag.is_null) {
    inner = node_new_scalar();
    node_scalar_t *val = &inner->data.value;
    mm_tag_cleanup(&val->tag);
    val->tag = tag;

    switch (tag.type) {
    case MM_VALUE_I:
    case MM_VALUE_I8:
    case MM_VALUE_I16:
    case MM_VALUE_I32:
    case MM_VALUE_I64:
    case MM_VALUE_U:
    case MM_VALUE_U8:
    case MM_VALUE_U16:
    case MM_VALUE_U32:
    case MM_VALUE_U64:
      val->text = strdup("0");
      break;
    case MM_VALUE_F32:
    case MM_VALUE_F64:
      val->text = strdup("0.0");
      break;
    case MM_VALUE_STR:
    case MM_VALUE_BYTES:
    case MM_VALUE_EMAIL:
    case MM_VALUE_UUID:
    case MM_VALUE_DECIMAL:
    case MM_VALUE_URL:
    case MM_VALUE_BIGINT:
      val->text = strdup("");
      break;
    case MM_VALUE_DATETIME:
    case MM_VALUE_DATE:
    case MM_VALUE_TIME:
      val->text = strdup("");
      break;
    case MM_VALUE_IP:
      val->text = strdup("");
      break;
    default:
      val->text = strdup("");
      break;
    }

    // If type is unknown, determine from remaining payload
    if (tag.is_null && val->tag.type == MM_VALUE_UNKNOWN) {
      size_t payload_remaining = total_len - tag_bytes_consumed;
      if (payload_remaining > 0) {
        uint8_t null_type_byte = dec_read_byte(d);
        switch (null_type_byte) {
        case MM_SIMPLE_NULL:
          // SimpleNull means no type - keep as unknown
          break;
        case MM_SIMPLE_NULLBOOL:
          val->tag.type = MM_VALUE_BOOL;
          free(val->text);
          val->text = strdup("false");
          break;
        case MM_SIMPLE_NULLINT:
          val->tag.type = MM_VALUE_I;
          free(val->text);
          val->text = strdup("0");
          break;
        case MM_SIMPLE_NULLFLOAT:
          val->tag.type = MM_VALUE_F64;
          free(val->text);
          val->text = strdup("0.0");
          break;
        case MM_SIMPLE_NULLSTRING:
          val->tag.type = MM_VALUE_STR;
          free(val->text);
          val->text = strdup("");
          break;
        case MM_SIMPLE_NULLBYTES:
          val->tag.type = MM_VALUE_BYTES;
          free(val->text);
          val->text = strdup("");
          break;
        default:
          val->tag.type = MM_VALUE_STR;
          free(val->text);
          val->text = strdup("");
          break;
        }
        tag_bytes_consumed += 1;
      }
    }

    size_t remaining = total_len - tag_bytes_consumed;
    if (d->offset + remaining <= d->size) {
      d->offset += remaining;
    }
  } else {
    inner = dec_decode_node(d, &tag);
    if (inner) {
      if (inner->type == MM_NODE_VALUE) {
        mm_tag_merge(&inner->data.value.tag, &tag);
        mm_tag_cleanup(&tag);
        if (inner->data.value.tag.type == MM_VALUE_UNKNOWN) {
          inner->data.value.tag.type = MM_VALUE_STR;
        }

        node_scalar_t *val = &inner->data.value;
        if (val->tag.type == MM_VALUE_DATETIME && val->text &&
            val->text[0] != '\0') {
          char *end = NULL;
          long long ts = strtoll(val->text, &end, 10);
          if (end && *end == '\0') {
            time_t t = (time_t)ts + val->tag.location_offset * 3600;
            struct tm tm;
            gmtime_r(&t, &tm);
            char buf[32];
            strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", &tm);
            free(val->text);
            val->text = strdup(buf);
          }
        } else if (val->tag.type == MM_VALUE_DATE && val->text &&
                   val->text[0] != '\0') {
          char *end = NULL;
          long long days = strtoll(val->text, &end, 10);
          if (end && *end == '\0') {
            time_t t = (time_t)days * 86400 + val->tag.location_offset * 3600;
            struct tm tm;
            gmtime_r(&t, &tm);
            char buf[16];
            strftime(buf, sizeof(buf), "%Y-%m-%d", &tm);
            free(val->text);
            val->text = strdup(buf);
          }
        } else if (val->tag.type == MM_VALUE_TIME && val->text &&
                   val->text[0] != '\0') {
          char *end = NULL;
          long long secs = strtoll(val->text, &end, 10);
          if (end && *end == '\0') {
            if (secs > 86399)
              secs = 86399;
            if (secs < 0)
              secs = 0;
            int hour = (int)(secs / 3600);
            int min = (int)((secs % 3600) / 60);
            int sec = (int)(secs % 60);
            char buf[16];
            snprintf(buf, sizeof(buf), "%02d:%02d:%02d", hour, min, sec);
            free(val->text);
            val->text = strdup(buf);
          }
        } else if (val->tag.type == MM_VALUE_UUID && val->text &&
                   val->text[0] != '\0') {
          size_t len = strlen(val->text);
          if (len == 24) {
            size_t decoded_len = 0;
            unsigned char *decoded =
                mc_base64_decode(val->text, len, &decoded_len);
            if (decoded && decoded_len == 16) {
              char uuid_buf[37];
              snprintf(uuid_buf, sizeof(uuid_buf),
                       "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%"
                       "02x%02x%02x%02x",
                       decoded[0], decoded[1], decoded[2], decoded[3],
                       decoded[4], decoded[5], decoded[6], decoded[7],
                       decoded[8], decoded[9], decoded[10], decoded[11],
                       decoded[12], decoded[13], decoded[14], decoded[15]);
              free(val->text);
              val->text = strdup(uuid_buf);
            }
            free(decoded);
          }
        }
      } else if (inner->type == MM_NODE_ARRAY) {
        mm_tag_merge(&inner->data.array.tag, &tag);
        mm_tag_cleanup(&tag);
        if (inner->data.array.tag.type == MM_VALUE_UNKNOWN) {
          if (inner->data.array.tag.size > 0) {
            inner->data.array.tag.type = MM_VALUE_ARR;
          } else {
            inner->data.array.tag.type = MM_VALUE_VEC;
          }
        }
      } else if (inner->type == MM_NODE_OBJECT) {
        mm_tag_merge(&inner->data.object.tag, &tag);
        mm_tag_cleanup(&tag);
        if (inner->data.object.tag.type == MM_VALUE_UNKNOWN) {
          inner->data.object.tag.type = MM_VALUE_OBJ;
        }
      }
    }
  }

  return inner;
}

static void dec_apply_tag_conversion(node_scalar_t *val) {
  if (!val || !val->text || val->text[0] == '\0')
    return;

  if (val->tag.type == MM_VALUE_TIME) {
    char *end = NULL;
    long long secs = strtoll(val->text, &end, 10);
    if (end && *end == '\0') {
      if (secs > 86399)
        secs = 86399;
      if (secs < 0)
        secs = 0;
      int hour = (int)(secs / 3600);
      int min = (int)((secs % 3600) / 60);
      int sec = (int)(secs % 60);
      char buf[16];
      snprintf(buf, sizeof(buf), "%02d:%02d:%02d", hour, min, sec);
      free(val->text);
      val->text = strdup(buf);
    }
  } else if (val->tag.type == MM_VALUE_DATE) {
    char *end = NULL;
    long long days = strtoll(val->text, &end, 10);
    if (end && *end == '\0') {
      time_t t = (time_t)days * 86400 + val->tag.location_offset * 3600;
      struct tm tm;
      gmtime_r(&t, &tm);
      char buf[16];
      strftime(buf, sizeof(buf), "%Y-%m-%d", &tm);
      free(val->text);
      val->text = strdup(buf);
    }
  } else if (val->tag.type == MM_VALUE_DATETIME) {
    char *end = NULL;
    long long ts = strtoll(val->text, &end, 10);
    if (end && *end == '\0') {
      time_t t = (time_t)ts + val->tag.location_offset * 3600;
      struct tm tm;
      gmtime_r(&t, &tm);
      char buf[32];
      strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", &tm);
      free(val->text);
      val->text = strdup(buf);
    }
  } else if (val->tag.type == MM_VALUE_UUID) {
    size_t len = strlen(val->text);
    if (len == 24) {
      size_t decoded_len = 0;
      unsigned char *decoded = mc_base64_decode(val->text, len, &decoded_len);
      if (decoded && decoded_len == 16) {
        char uuid_buf[37];
        snprintf(uuid_buf, sizeof(uuid_buf),
                 "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%"
                 "02x%02x%02x%02x",
                 decoded[0], decoded[1], decoded[2], decoded[3], decoded[4],
                 decoded[5], decoded[6], decoded[7], decoded[8], decoded[9],
                 decoded[10], decoded[11], decoded[12], decoded[13],
                 decoded[14], decoded[15]);
        free(val->text);
        val->text = strdup(uuid_buf);
      }
      free(decoded);
    }
  } else if (val->tag.type == MM_VALUE_ENUMS && val->tag.enums && val->text &&
             val->text[0] != '\0') {
    // Only convert if text is a pure numeric string (enum index),
    // not if it's already been converted to an enum value string
    char *endptr = NULL;
    long long idx = strtoll(val->text, &endptr, 10);
    if (endptr == val->text || *endptr != '\0') {
      return;
    }
    const char *enums = val->tag.enums;
    int current = 0;
    const char *segment = enums;

    // Find the idx-th '|'-separated segment
    while (*segment) {
      if (current == idx)
        break;
      if (*segment == '|')
        current++;
      segment++;
    }

    if (current == idx && *segment) {
      const char *end = segment;
      while (*end && *end != '|')
        end++;

      size_t len = (size_t)(end - segment);
      char *enum_str = (char *)malloc(len + 1);
      if (enum_str) {
        memcpy(enum_str, segment, len);
        enum_str[len] = '\0';
        free(val->text);
        val->text = enum_str;
      }
    }
  }
}

static node_t *dec_decode_node(mm_decoder_t *d, mm_tag_t *parent_tag) {
  if (!d || d->offset >= d->size) {
    return NULL;
  }

  errno = 0;
  uint8_t b = dec_read_byte(d);
  if (errno)
    return NULL;

  int prefix = mm_prefix_of(b);
  node_t *node = NULL;

  switch (prefix) {
  case MM_PREFIX_SIMPLE:
    node = dec_decode_simple(d, b);
    break;
  case MM_PREFIX_POSITIVEINT:
    node = dec_decode_int(d, b, 1);
    break;
  case MM_PREFIX_NEGATIVEINT:
    node = dec_decode_int(d, b, 0);
    break;
  case MM_PREFIX_FLOAT:
    node = dec_decode_float(d, b);
    break;
  case MM_PREFIX_STRING:
    node = dec_decode_string(d, b);
    break;
  case MM_PREFIX_BYTES:
    return dec_decode_bytes(d, b, parent_tag);
  case MM_PREFIX_CONTAINER:
    return dec_decode_container(d, b, parent_tag);
  case MM_PREFIX_TAG:
    return dec_decode_tag(d, b, parent_tag);
  default:
    return NULL;
  }

  // For VALUE nodes, apply parent tag type and conversion
  if (node && parent_tag) {
    if (node->type == MM_NODE_VALUE) {
      node_scalar_t *val = &node->data.value;
      mm_tag_merge(&val->tag, parent_tag);
      // Preserve inherited semantics: if parent's attributes were inherited,
      // the merged attributes should remain inherited (not shown in output)
      if (parent_tag->is_inherit) {
        val->tag.is_inherit = true;
      }
      dec_apply_tag_conversion(val);
    } else if (node->type == MM_NODE_NULL) {
      mm_tag_merge(&node->tag, parent_tag);
    }
  }

  return node;
}

mm_decoder_t *mm_decoder_new(const uint8_t *data, size_t size) {
  mm_decoder_t *d = (mm_decoder_t *)malloc(sizeof(mm_decoder_t));
  if (!d)
    return NULL;
  d->data = data;
  d->size = size;
  d->offset = 0;
  return d;
}

void mm_decoder_free(mm_decoder_t *d) { free(d); }

node_t *mm_decoder_decode(mm_decoder_t *d) { return dec_decode_node(d, NULL); }