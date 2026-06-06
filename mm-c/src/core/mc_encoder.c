#include "mc_encoder.h"
#include "../ir/mc_tag.h"
#include "../ir/mc_value_type.h"
#include "mc_constants.h"
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#define DEBUG_ENCODE 0

static uint8_t base64_decode_char(char c) {
  if (c >= 'A' && c <= 'Z')
    return (uint8_t)(c - 'A');
  if (c >= 'a' && c <= 'z')
    return (uint8_t)(c - 'a' + 26);
  if (c >= '0' && c <= '9')
    return (uint8_t)(c - '0' + 52);
  if (c == '-' || c == '+')
    return 62;
  if (c == '_' || c == '/')
    return 63;
  return 0xFF;
}

static uint8_t *mc_base64_decode_url(const char *text, size_t text_len,
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
    uint8_t a = base64_decode_char(text[i]);
    uint8_t b = (i + 1 < text_len) ? base64_decode_char(text[i + 1]) : 0;
    uint8_t c = (i + 2 < text_len) ? base64_decode_char(text[i + 2]) : 0;
    uint8_t d = (i + 3 < text_len) ? base64_decode_char(text[i + 3]) : 0;
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

typedef struct {
  uint8_t *buf;
  size_t size;
  size_t cap;
} encoder_t;

static void enc_encode_node_object(encoder_t *e, mm_object_t *obj);
static void enc_encode_node_array(encoder_t *e, mm_array_t *arr);
static void enc_encode_node_value(encoder_t *e, mm_value_t *val);
static void enc_encode_node_doc(encoder_t *e, mm_doc_t *doc);

static void enc_write_byte(encoder_t *e, uint8_t b) {
  if (e->size >= e->cap) {
    e->cap = e->cap == 0 ? 64 : e->cap * 2;
    e->buf = realloc(e->buf, e->cap);
  }
  e->buf[e->size++] = b;
}

static void enc_write_bytes(encoder_t *e, const uint8_t *data, size_t len) {
  if (len == 0)
    return;
  while (e->size + len > e->cap) {
    e->cap = e->cap == 0 ? 64 : e->cap * 2;
    e->buf = realloc(e->buf, e->cap);
  }
  memcpy(e->buf + e->size, data, len);
  e->size += len;
}

static void enc_encode_i(encoder_t *e, uint8_t sign, const char *text) {
  uint64_t uv;
  if (sign == MM_PREFIX_NEGATIVEINT) {
    int64_t sv = (int64_t)strtoll(text, NULL, 10);
    if (sv < 0) {
      uv = (uint64_t)(-sv);
    } else {
      uv = (uint64_t)sv;
      sign = MM_PREFIX_POSITIVEINT;
    }
  } else {
    uv = strtoull(text, NULL, 10);
  }

  if (uv <= 23) {
    enc_write_byte(e, (uint8_t)(sign | uv));
  } else if (uv <= 0xFF) {
    enc_write_byte(e, (uint8_t)(sign | MM_INTLEN1BYTE));
    enc_write_byte(e, (uint8_t)uv);
  } else if (uv <= 0xFFFF) {
    enc_write_byte(e, (uint8_t)(sign | MM_INTLEN2BYTE));
    enc_write_byte(e, (uint8_t)(uv >> 8));
    enc_write_byte(e, (uint8_t)uv);
  } else if (uv <= 0xFFFFFF) {
    enc_write_byte(e, (uint8_t)(sign | MM_INTLEN3BYTE));
    enc_write_byte(e, (uint8_t)(uv >> 16));
    enc_write_byte(e, (uint8_t)(uv >> 8));
    enc_write_byte(e, (uint8_t)uv);
  } else if (uv <= 0xFFFFFFFF) {
    enc_write_byte(e, (uint8_t)(sign | MM_INTLEN4BYTE));
    enc_write_byte(e, (uint8_t)(uv >> 24));
    enc_write_byte(e, (uint8_t)(uv >> 16));
    enc_write_byte(e, (uint8_t)(uv >> 8));
    enc_write_byte(e, (uint8_t)uv);
  } else {
    enc_write_byte(e, (uint8_t)(sign | MM_INTLEN8BYTE));
    enc_write_byte(e, (uint8_t)(uv >> 56));
    enc_write_byte(e, (uint8_t)(uv >> 48));
    enc_write_byte(e, (uint8_t)(uv >> 40));
    enc_write_byte(e, (uint8_t)(uv >> 32));
    enc_write_byte(e, (uint8_t)(uv >> 24));
    enc_write_byte(e, (uint8_t)(uv >> 16));
    enc_write_byte(e, (uint8_t)(uv >> 8));
    enc_write_byte(e, (uint8_t)uv);
  }
}

static bool parse_float_string(const char *text, bool *is_negative,
                               int8_t *exponent, uint64_t *mantissa) {
  const char *p = text;
  *is_negative = false;
  *exponent = 0;
  *mantissa = 0;

  if (*p == '-') {
    *is_negative = true;
    p++;
  } else if (*p == '+') {
    p++;
  }

  if (*p == '\0')
    return false;

  // Split into integer and fractional parts
  const char *int_start = p;
  size_t int_len = 0;
  while (*p >= '0' && *p <= '9') {
    p++;
    int_len++;
  }

  size_t frac_len = 0;
  const char *frac_start = NULL;
  if (*p == '.') {
    p++;
    frac_start = p;
    while (*p >= '0' && *p <= '9') {
      p++;
      frac_len++;
    }
  }

  // Handle scientific notation
  int64_t sci_exp = 0;
  if (*p == 'e' || *p == 'E') {
    p++;
    int neg_exp = 0;
    if (*p == '-') {
      neg_exp = 1;
      p++;
    } else if (*p == '+') {
      p++;
    }
    char *end = NULL;
    long e = strtol(p, &end, 10);
    if (end == p)
      return false;
    sci_exp = neg_exp ? -e : e;
    p = end;
  }

  if (*p != '\0')
    return false;

  // No digits at all
  if (int_len == 0 && frac_len == 0)
    return false;

  // Build mantissa string: concatenate int and frac parts, stripping leading
  // zeros
  char buf[64];
  size_t buf_len = 0;
  bool leading = true;

  for (size_t i = 0; i < int_len; i++) {
    if (leading && int_start[i] == '0')
      continue;
    leading = false;
    buf[buf_len++] = int_start[i];
  }

  for (size_t i = 0; i < frac_len; i++) {
    if (leading && frac_start[i] == '0')
      continue;
    leading = false;
    buf[buf_len++] = frac_start[i];
  }

  if (buf_len == 0) {
    buf[buf_len++] = '0';
  }

  buf[buf_len] = '\0';

  // Parse mantissa
  char *end = NULL;
  unsigned long long uv = strtoull(buf, &end, 10);
  if (end != buf + buf_len)
    return false;
  *mantissa = (uint64_t)uv;

  // Compute exponent: -(frac_len) + sci_exp
  int64_t exp = -(int64_t)frac_len + sci_exp;
  if (exp < INT8_MIN || exp > INT8_MAX)
    return false;
  *exponent = (int8_t)exp;

  return true;
}

static void enc_encode_float(encoder_t *e, const char *text) {
  bool is_negative;
  int8_t exponent;
  uint64_t mantissa;

  if (!parse_float_string(text, &is_negative, &exponent, &mantissa)) {
    return;
  }

  uint8_t sign = MM_PREFIX_FLOAT;
  if (is_negative) {
    sign |= MM_FLOAT_NEG_MASK;
  }

  // Inline encoding: 0.0 through 0.7
  if (exponent == -1 && mantissa <= 7) {
    enc_write_byte(e, (uint8_t)(sign | mantissa));
    return;
  }

  int mantissa_bytes = 0;
  if (mantissa <= 0xFF) {
    mantissa_bytes = 1;
  } else if (mantissa <= 0xFFFF) {
    mantissa_bytes = 2;
  } else if (mantissa <= 0xFFFFFF) {
    mantissa_bytes = 3;
  } else if (mantissa <= 0xFFFFFFFF) {
    mantissa_bytes = 4;
  } else if (mantissa <= 0xFFFFFFFFFF) {
    mantissa_bytes = 5;
  } else if (mantissa <= 0xFFFFFFFFFFFF) {
    mantissa_bytes = 6;
  } else if (mantissa <= 0xFFFFFFFFFFFFFF) {
    mantissa_bytes = 7;
  } else {
    mantissa_bytes = 8;
  }

  int len_val = mantissa_bytes + 7;
  sign |= len_val;

  enc_write_byte(e, sign);
  enc_write_byte(e, (uint8_t)(int8_t)exponent);

  // Write mantissa bytes in big-endian
  switch (mantissa_bytes) {
  case 1:
    enc_write_byte(e, (uint8_t)mantissa);
    break;
  case 2:
    enc_write_byte(e, (uint8_t)(mantissa >> 8));
    enc_write_byte(e, (uint8_t)mantissa);
    break;
  case 3:
    enc_write_byte(e, (uint8_t)(mantissa >> 16));
    enc_write_byte(e, (uint8_t)(mantissa >> 8));
    enc_write_byte(e, (uint8_t)mantissa);
    break;
  case 4:
    enc_write_byte(e, (uint8_t)(mantissa >> 24));
    enc_write_byte(e, (uint8_t)(mantissa >> 16));
    enc_write_byte(e, (uint8_t)(mantissa >> 8));
    enc_write_byte(e, (uint8_t)mantissa);
    break;
  case 5:
    enc_write_byte(e, (uint8_t)(mantissa >> 32));
    enc_write_byte(e, (uint8_t)(mantissa >> 24));
    enc_write_byte(e, (uint8_t)(mantissa >> 16));
    enc_write_byte(e, (uint8_t)(mantissa >> 8));
    enc_write_byte(e, (uint8_t)mantissa);
    break;
  case 6:
    enc_write_byte(e, (uint8_t)(mantissa >> 40));
    enc_write_byte(e, (uint8_t)(mantissa >> 32));
    enc_write_byte(e, (uint8_t)(mantissa >> 24));
    enc_write_byte(e, (uint8_t)(mantissa >> 16));
    enc_write_byte(e, (uint8_t)(mantissa >> 8));
    enc_write_byte(e, (uint8_t)mantissa);
    break;
  case 7:
    enc_write_byte(e, (uint8_t)(mantissa >> 48));
    enc_write_byte(e, (uint8_t)(mantissa >> 40));
    enc_write_byte(e, (uint8_t)(mantissa >> 32));
    enc_write_byte(e, (uint8_t)(mantissa >> 24));
    enc_write_byte(e, (uint8_t)(mantissa >> 16));
    enc_write_byte(e, (uint8_t)(mantissa >> 8));
    enc_write_byte(e, (uint8_t)mantissa);
    break;
  case 8:
    enc_write_byte(e, (uint8_t)(mantissa >> 56));
    enc_write_byte(e, (uint8_t)(mantissa >> 48));
    enc_write_byte(e, (uint8_t)(mantissa >> 40));
    enc_write_byte(e, (uint8_t)(mantissa >> 32));
    enc_write_byte(e, (uint8_t)(mantissa >> 24));
    enc_write_byte(e, (uint8_t)(mantissa >> 16));
    enc_write_byte(e, (uint8_t)(mantissa >> 8));
    enc_write_byte(e, (uint8_t)mantissa);
    break;
  }
}

static void enc_encode_string(encoder_t *e, const char *s) {
  size_t len = strlen(s);
  uint8_t sign = MM_PREFIX_STRING;

  if (len <= 29) {
    enc_write_byte(e, (uint8_t)(sign | len));
  } else if (len <= 255) {
    enc_write_byte(e, (uint8_t)(sign | MM_STRINGLEN1BYTE));
    enc_write_byte(e, (uint8_t)len);
  } else if (len <= 65535) {
    enc_write_byte(e, (uint8_t)(sign | MM_STRINGLEN2BYTE));
    enc_write_byte(e, (uint8_t)(len >> 8));
    enc_write_byte(e, (uint8_t)(len & 0xFF));
  }

  enc_write_bytes(e, (const uint8_t *)s, len);
}

static void enc_encode_bytes(encoder_t *e, const uint8_t *data, size_t len) {
  uint8_t sign = MM_PREFIX_BYTES;

  if (len <= 29) {
    enc_write_byte(e, (uint8_t)(sign | len));
  } else if (len <= 255) {
    enc_write_byte(e, (uint8_t)(sign | MM_BYTESLEN1BYTE));
    enc_write_byte(e, (uint8_t)len);
  } else if (len <= 65535) {
    enc_write_byte(e, (uint8_t)(sign | MM_BYTESLEN2BYTE));
    enc_write_byte(e, (uint8_t)(len >> 8));
    enc_write_byte(e, (uint8_t)(len & 0xFF));
  }

  enc_write_bytes(e, data, len);
}

static void enc_encode_bigint(encoder_t *e, const char *text) {
  if (!text || !*text)
    return;

  bool neg = false;
  if (*text == '-') {
    neg = true;
    text++;
  }

  size_t len = strlen(text);
  if (len == 0)
    return;

  size_t groups = len / 3;
  size_t rem = len % 3;
  size_t total_bits = 1 + groups * 10;
  if (rem == 2)
    total_bits += 7;
  else if (rem == 1)
    total_bits += 4;

  size_t byte_count = (total_bits + 7) / 8;
  uint8_t *bits = calloc(1, byte_count);

  uint8_t current = 0;
  int bit_offset = 0;
  size_t byte_idx = 0;

#define WRITE_BIT(b)                                                           \
  do {                                                                         \
    current = (current << 1) | ((b) & 1);                                      \
    bit_offset++;                                                              \
    if (bit_offset == 8) {                                                     \
      bits[byte_idx++] = current;                                              \
      current = 0;                                                             \
      bit_offset = 0;                                                          \
    }                                                                          \
  } while (0)

  WRITE_BIT(neg ? 1 : 0);

  for (size_t i = 0; i < len;) {
    size_t rem_group = len - i;
    int val;
    int num_bits;
    if (rem_group >= 3) {
      val = (text[i] - '0') * 100 + (text[i + 1] - '0') * 10 +
            (text[i + 2] - '0');
      num_bits = 10;
      i += 3;
    } else if (rem_group == 2) {
      val = (text[i] - '0') * 10 + (text[i + 1] - '0');
      num_bits = 7;
      i += 2;
    } else {
      val = text[i] - '0';
      num_bits = 4;
      i += 1;
    }
    for (int b = num_bits - 1; b >= 0; b--) {
      WRITE_BIT((val >> b) & 1);
    }
  }

  if (bit_offset > 0) {
    current <<= (8 - bit_offset);
    bits[byte_idx++] = current;
  }

#undef WRITE_BIT

  uint8_t *payload = malloc(1 + byte_count);
  payload[0] = (uint8_t)len;
  memcpy(payload + 1, bits, byte_count);

  enc_encode_bytes(e, payload, 1 + byte_count);

  free(payload);
  free(bits);
}

static void enc_encode_simple(encoder_t *e, uint8_t value) {
  enc_write_byte(e, (uint8_t)(MM_PREFIX_SIMPLE | value));
}

static void enc_encode_bool(encoder_t *e, const char *text) {
  if (strcmp(text, "true") == 0) {
    enc_encode_simple(e, MM_SIMPLE_TRUE);
  } else {
    enc_encode_simple(e, MM_SIMPLE_FALSE);
  }
}

static void enc_encode_array(encoder_t *e, const uint8_t *data, size_t len) {
  uint8_t sign = (uint8_t)(MM_PREFIX_CONTAINER | MM_CONTAINER_ARRAY);

  if (len <= 13) {
    enc_write_byte(e, (uint8_t)(sign | len));
  } else if (len <= 255) {
    enc_write_byte(e, (uint8_t)(sign | MM_CONTAINERLEN1BYTE));
    enc_write_byte(e, (uint8_t)len);
  } else if (len <= 65535) {
    enc_write_byte(e, (uint8_t)(sign | MM_CONTAINERLEN2BYTE));
    enc_write_byte(e, (uint8_t)(len >> 8));
    enc_write_byte(e, (uint8_t)(len & 0xFF));
  }

  enc_write_bytes(e, data, len);
}

static void enc_encode_container(encoder_t *e, uint8_t container_type,
                                 const uint8_t *data, size_t len) {
  uint8_t sign = (uint8_t)(MM_PREFIX_CONTAINER | container_type);

  if (len <= 13) {
    enc_write_byte(e, (uint8_t)(sign | len));
  } else if (len <= 255) {
    enc_write_byte(e, (uint8_t)(sign | MM_CONTAINERLEN1BYTE));
    enc_write_byte(e, (uint8_t)len);
  } else if (len <= 65535) {
    enc_write_byte(e, (uint8_t)(sign | MM_CONTAINERLEN2BYTE));
    enc_write_byte(e, (uint8_t)(len >> 8));
    enc_write_byte(e, (uint8_t)(len & 0xFF));
  }

  enc_write_bytes(e, data, len);
}

static void enc_encode_tag(encoder_t *e, mm_tag_t *tag,
                           const uint8_t *inner_data, size_t inner_len) {
  size_t tag_len;
  uint8_t *tag_bytes = mm_tag_bytes(tag, &tag_len);

  if (tag_bytes == NULL || tag_len == 0) {
    free(tag_bytes);
    enc_write_bytes(e, inner_data, inner_len);
    return;
  }

  encoder_t tag_enc = {0};

  if (tag_len < 254) {
    enc_write_byte(&tag_enc, (uint8_t)tag_len);
  } else if (tag_len < 257) {
    enc_write_byte(&tag_enc, 254);
    enc_write_byte(&tag_enc, (uint8_t)tag_len);
  } else {
    enc_write_byte(&tag_enc, 255);
    enc_write_byte(&tag_enc, (uint8_t)(tag_len >> 8));
    enc_write_byte(&tag_enc, (uint8_t)(tag_len & 0xFF));
  }

  enc_write_bytes(&tag_enc, tag_bytes, tag_len);
  free(tag_bytes);

  size_t combined_len = tag_enc.size + inner_len;

  if (combined_len < MM_TAGLEN1BYTE) {
    enc_write_byte(e, (uint8_t)(MM_PREFIX_TAG | combined_len));
  } else if (combined_len < 256) {
    enc_write_byte(e, (uint8_t)(MM_PREFIX_TAG | MM_TAGLEN1BYTE));
    enc_write_byte(e, (uint8_t)combined_len);
  } else if (combined_len < 65536) {
    enc_write_byte(e, (uint8_t)(MM_PREFIX_TAG | MM_TAGLEN2BYTE));
    enc_write_byte(e, (uint8_t)(combined_len >> 8));
    enc_write_byte(e, (uint8_t)(combined_len & 0xFF));
  }

  enc_write_bytes(e, tag_enc.buf, tag_enc.size);
  enc_write_bytes(e, inner_data, inner_len);

  free(tag_enc.buf);
}

static void enc_encode_node_value(encoder_t *e, mm_value_t *val) {
  encoder_t tmp = {0};

  switch (val->tag.type) {
  case MM_VALUE_STR:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLSTRING);
    } else {
      enc_encode_string(&tmp, val->text);
    }
    break;

  case MM_VALUE_BOOL:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLBOOL);
    } else {
      enc_encode_bool(&tmp, val->text);
    }
    break;

  case MM_VALUE_I:
  case MM_VALUE_I8:
  case MM_VALUE_I16:
  case MM_VALUE_I32:
  case MM_VALUE_I64: {
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
    } else {
      int64_t sv = (int64_t)strtoll(val->text, NULL, 10);
      if (sv < 0) {
        enc_encode_i(&tmp, MM_PREFIX_NEGATIVEINT, val->text);
      } else {
        enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, val->text);
      }
    }
    break;
  }

  case MM_VALUE_U:
  case MM_VALUE_U8:
  case MM_VALUE_U16:
  case MM_VALUE_U32:
  case MM_VALUE_U64:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
    } else {
      enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, val->text);
    }
    break;

  case MM_VALUE_F32:
  case MM_VALUE_F64:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLFLOAT);
    } else {
      enc_encode_float(&tmp, val->text);
    }
    break;

  case MM_VALUE_BYTES:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLBYTES);
    } else {
      size_t decoded_len = 0;
      uint8_t *decoded =
          mc_base64_decode_url(val->text, strlen(val->text), &decoded_len);
      if (decoded) {
        enc_encode_bytes(&tmp, decoded, decoded_len);
        free(decoded);
      }
    }
    break;

  case MM_VALUE_MEDIA:
  case MM_VALUE_IMAGE:
  case MM_VALUE_VIDEO:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLBYTES);
    } else {
      size_t decoded_len = 0;
      uint8_t *decoded =
          mc_base64_decode_url(val->text, strlen(val->text), &decoded_len);
      if (decoded) {
        enc_encode_bytes(&tmp, decoded, decoded_len);
        free(decoded);
      }
    }
    break;

  case MM_VALUE_DECIMAL:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLFLOAT);
    } else {
      enc_encode_float(&tmp, val->text);
    }
    break;

  case MM_VALUE_EMAIL:
  case MM_VALUE_URL:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLSTRING);
    } else {
      enc_encode_string(&tmp, val->text);
    }
    break;

  case MM_VALUE_BIGINT:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLSTRING);
    } else {
      enc_encode_bigint(&tmp, val->text);
    }
    break;

  case MM_VALUE_UUID:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLBYTES);
    } else {
      uint8_t uuid_bytes[16];
      size_t uuid_len = 0;
      char hex[33];
      size_t j = 0;
      for (const char *p = val->text; *p && j < 32; p++) {
        if (*p != '-') {
          hex[j++] = *p;
        }
      }
      hex[j] = '\0';
      for (size_t i = 0; i < 16 && i * 2 < 32; i++) {
        char byte_str[3] = {hex[i * 2], hex[i * 2 + 1] ? hex[i * 2 + 1] : '\0',
                            '\0'};
        uuid_bytes[i] = (uint8_t)strtol(byte_str, NULL, 16);
      }
      uuid_len = 16;
      enc_encode_bytes(&tmp, uuid_bytes, uuid_len);
    }
    break;

  case MM_VALUE_DATE:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
    } else {
      struct tm tm = {0};
      if (sscanf(val->text, "%d-%d-%d", &tm.tm_year, &tm.tm_mon, &tm.tm_mday) >=
          3) {
        tm.tm_year -= 1900;
        tm.tm_mon -= 1;
        tm.tm_isdst = -1;
        int64_t epoch = (int64_t)timegm(&tm);
        int64_t days = epoch / 86400;
        char buf[24];
        snprintf(buf, sizeof(buf), "%lld", (long long)days);
        enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, buf);
      } else {
        enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, "0");
      }
    }
    break;

  case MM_VALUE_TIME:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
    } else {
      int hour = 0, minute = 0, sec = 0;
      if (sscanf(val->text, "%d:%d:%d", &hour, &minute, &sec) >= 2) {
        int total_secs = hour * 3600 + minute * 60 + sec;
        char buf[24];
        snprintf(buf, sizeof(buf), "%d", total_secs);
        enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, buf);
      } else {
        enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, "0");
      }
    }
    break;

  case MM_VALUE_DATETIME:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
    } else {
      struct tm tm = {0};
      if (sscanf(val->text, "%d-%d-%d %d:%d:%d", &tm.tm_year, &tm.tm_mon,
                 &tm.tm_mday, &tm.tm_hour, &tm.tm_min, &tm.tm_sec) >= 3) {
        tm.tm_year -= 1900;
        tm.tm_mon -= 1;
        tm.tm_isdst = -1;
        int64_t epoch = (int64_t)timegm(&tm);
        epoch -= val->tag.location_offset * 3600;
        char buf[24];
        snprintf(buf, sizeof(buf), "%lld", (long long)epoch);
        if (epoch < 0) {
          enc_encode_i(&tmp, MM_PREFIX_NEGATIVEINT, buf);
        } else {
          enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, buf);
        }
      } else {
        enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, "0");
      }
    }
    break;

  case MM_VALUE_IP:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLSTRING);
    } else {
      enc_encode_string(&tmp, val->text);
    }
    break;

  case MM_VALUE_ENUMS:
    if (val->tag.is_null) {
      enc_encode_simple(&tmp, MM_SIMPLE_NULLINT);
    } else {
      int enum_index = 0;
      if (val->tag.enums) {
        char *enums_copy = strdup(val->tag.enums);
        if (enums_copy) {
          char *token = strtok(enums_copy, "|");
          int idx = 0;
          while (token) {
            while (*token == ' ')
              token++;
            char *end = token + strlen(token);
            while (end > token && *(end - 1) == ' ')
              end--;
            *end = '\0';
            if (strcmp(token, val->text) == 0) {
              enum_index = idx;
              break;
            }
            token = strtok(NULL, "|");
            idx++;
          }
          free(enums_copy);
        }
      }
      char buf[24];
      snprintf(buf, sizeof(buf), "%d", enum_index);
      enc_encode_i(&tmp, MM_PREFIX_POSITIVEINT, buf);
    }
    break;

  default:
    enc_encode_string(&tmp, val->text);
    break;
  }

  enc_encode_tag(e, &val->tag, tmp.buf, tmp.size);
  free(tmp.buf);
}

static void enc_encode_node_array(encoder_t *e, mm_array_t *arr) {
  encoder_t items = {0};

  for (size_t i = 0; i < arr->item_count; i++) {
    encoder_t tmp = {0};
    mm_node_t *item = arr->items[i];

    switch (item->type) {
    case MM_NODE_OBJECT:
      enc_encode_node_object(&tmp, &item->data.object);
      break;
    case MM_NODE_ARRAY:
      enc_encode_node_array(&tmp, &item->data.array);
      break;
    case MM_NODE_VALUE:
      enc_encode_node_value(&tmp, &item->data.value);
      break;
    case MM_NODE_DOC:
      enc_encode_node_doc(&tmp, &item->data.doc);
      break;
    default:
      break;
    }

    enc_write_bytes(&items, tmp.buf, tmp.size);
    free(tmp.buf);
  }

  encoder_t container_enc = {0};
  enc_encode_array(&container_enc, items.buf, items.size);
  free(items.buf);

  enc_encode_tag(e, &arr->tag, container_enc.buf, container_enc.size);
  free(container_enc.buf);
}

static void enc_encode_node_object(encoder_t *e, mm_object_t *obj) {
  encoder_t keys = {0};
  encoder_t vals = {0};

  for (size_t i = 0; i < obj->field_count; i++) {
    mm_field_t *field = &obj->fields[i];

    encoder_t tmp_val = {0};
    switch (field->value->type) {
    case MM_NODE_OBJECT:
      enc_encode_node_object(&tmp_val, &field->value->data.object);
      break;
    case MM_NODE_ARRAY:
      enc_encode_node_array(&tmp_val, &field->value->data.array);
      break;
    case MM_NODE_VALUE:
      enc_encode_node_value(&tmp_val, &field->value->data.value);
      break;
    case MM_NODE_DOC:
      enc_encode_node_doc(&tmp_val, &field->value->data.doc);
      break;
    default:
      break;
    }
    enc_write_bytes(&vals, tmp_val.buf, tmp_val.size);
    free(tmp_val.buf);

    encoder_t tmp_key = {0};
    enc_encode_string(&tmp_key, field->key);
    enc_write_bytes(&keys, tmp_key.buf, tmp_key.size);
    free(tmp_key.buf);
  }

  encoder_t arr_enc = {0};
  enc_encode_array(&arr_enc, keys.buf, keys.size);
  free(keys.buf);

  encoder_t combined = {0};
  enc_write_bytes(&combined, arr_enc.buf, arr_enc.size);
  enc_write_bytes(&combined, vals.buf, vals.size);
  free(arr_enc.buf);
  free(vals.buf);

  encoder_t container_enc = {0};
  enc_encode_container(&container_enc, MM_CONTAINER_OBJECT, combined.buf,
                       combined.size);
  free(combined.buf);

  enc_encode_tag(e, &obj->tag, container_enc.buf, container_enc.size);
  free(container_enc.buf);
}

static void enc_encode_node_doc(encoder_t *e, mm_doc_t *doc) {
  encoder_t keys = {0};
  encoder_t vals = {0};

  for (size_t i = 0; i < doc->field_count; i++) {
    mm_field_t *field = &doc->fields[i];

    encoder_t tmp_val = {0};
    switch (field->value->type) {
    case MM_NODE_OBJECT:
      enc_encode_node_object(&tmp_val, &field->value->data.object);
      break;
    case MM_NODE_ARRAY:
      enc_encode_node_array(&tmp_val, &field->value->data.array);
      break;
    case MM_NODE_VALUE:
      enc_encode_node_value(&tmp_val, &field->value->data.value);
      break;
    case MM_NODE_DOC:
      enc_encode_node_doc(&tmp_val, &field->value->data.doc);
      break;
    default:
      break;
    }
    enc_write_bytes(&vals, tmp_val.buf, tmp_val.size);
    free(tmp_val.buf);

    encoder_t tmp_key = {0};
    enc_encode_string(&tmp_key, field->key);
    enc_write_bytes(&keys, tmp_key.buf, tmp_key.size);
    free(tmp_key.buf);
  }

  encoder_t arr_enc = {0};
  enc_encode_array(&arr_enc, keys.buf, keys.size);
  free(keys.buf);

  encoder_t combined = {0};
  enc_write_bytes(&combined, arr_enc.buf, arr_enc.size);
  enc_write_bytes(&combined, vals.buf, vals.size);
  free(arr_enc.buf);
  free(vals.buf);

  encoder_t container_enc = {0};
  enc_encode_container(&container_enc, MM_CONTAINER_OBJECT, combined.buf,
                       combined.size);
  free(combined.buf);

  enc_encode_tag(e, &doc->tag, container_enc.buf, container_enc.size);
  free(container_enc.buf);
}

static void enc_encode(encoder_t *e, mm_node_t *node) {
  if (node == NULL)
    return;

  switch (node->type) {
  case MM_NODE_OBJECT:
    enc_encode_node_object(e, &node->data.object);
    break;
  case MM_NODE_ARRAY:
    enc_encode_node_array(e, &node->data.array);
    break;
  case MM_NODE_VALUE:
    enc_encode_node_value(e, &node->data.value);
    break;
  case MM_NODE_DOC:
    enc_encode_node_doc(e, &node->data.doc);
    break;
  default:
    break;
  }
}

mm_encoder_buffer_t *mm_encoder_encode(mm_node_t *node) {
  encoder_t e = {0};
  enc_encode(&e, node);
  mm_encoder_buffer_t *buf = malloc(sizeof(mm_encoder_buffer_t));
  buf->data = e.buf;
  buf->size = e.size;
  buf->capacity = e.cap;
  return buf;
}

void mm_encoder_buffer_free(mm_encoder_buffer_t *buf) {
  if (buf) {
    free(buf->data);
    free(buf);
  }
}