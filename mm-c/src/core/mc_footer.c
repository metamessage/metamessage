#include "mc_footer.h"
#include "mc_constants.h"
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>

/* footer_slot_bytes：编码尾部长度槽。 */
size_t mc_footer_slot_bytes(size_t meta_len, uint8_t *out) {
  if (meta_len <= 253) {
    out[0] = (uint8_t)meta_len;
    return 1;
  } else if (meta_len <= 0xFF) {
    out[0] = (uint8_t)(meta_len & 0xFF);
    out[1] = 254;
    return 2;
  } else if (meta_len <= 0xFFFF) {
    out[0] = (uint8_t)((meta_len >> 8) & 0xFF);
    out[1] = (uint8_t)(meta_len & 0xFF);
    out[2] = 255;
    return 3;
  }
  return 0;
}

void mc_reverse_bytes(const uint8_t *src, uint8_t *dst, size_t len) {
  for (size_t i = 0; i < len; i++)
    dst[i] = src[len - 1 - i];
}

/* ---- 动态字节缓冲 ---- */
typedef struct {
  uint8_t *buf;
  size_t len;
  size_t cap;
} dbuf_t;

static void db_init(dbuf_t *b) {
  b->buf = NULL;
  b->len = 0;
  b->cap = 0;
}
static void db_free(dbuf_t *b) {
  free(b->buf);
  b->buf = NULL;
  b->len = b->cap = 0;
}
static int db_reserve(dbuf_t *b, size_t extra) {
  if (b->len + extra <= b->cap)
    return 0;
  size_t nc = b->cap ? b->cap : 16;
  while (b->len + extra > nc)
    nc *= 2;
  uint8_t *nb = realloc(b->buf, nc);
  if (!nb)
    return -1;
  b->buf = nb;
  b->cap = nc;
  return 0;
}
static int db_push(dbuf_t *b, uint8_t c) {
  if (db_reserve(b, 1) != 0)
    return -1;
  b->buf[b->len++] = c;
  return 0;
}
static int db_write(dbuf_t *b, const uint8_t *d, size_t n) {
  if (n == 0)
    return 0;
  if (db_reserve(b, n) != 0)
    return -1;
  memcpy(b->buf + b->len, d, n);
  b->len += n;
  return 0;
}

/* 各类型档位：返回前缀，回填 (w, inline)；未知前缀返回 -1。 */
static int f_tier(uint8_t sign, int *w, int *inline_len) {
  *w = 0;
  *inline_len = 0;
  int p = mm_prefix_of(sign);
  switch (p) {
  case MM_PREFIX_SIMPLE:
    return p;
  case MM_PREFIX_POSITIVEINT:
  case MM_PREFIX_NEGATIVEINT: {
    int l = sign & 0b11111;
    if (l < MM_INTLEN1BYTE)
      *inline_len = l;
    else
      *w = l - MM_INTLEN1BYTE + 1;
    return p;
  }
  case MM_PREFIX_FLOAT: {
    int l = sign & 0b01111;
    if (l < MM_FLOATLEN1BYTE)
      *inline_len = l;
    else
      *w = l - MM_FLOATLEN1BYTE + 1;
    return p;
  }
  case MM_PREFIX_STRING: {
    int l = sign & 0b11111;
    *inline_len = l;
    if (l < MM_STRINGLEN1BYTE)
      *w = 0;
    else if (l == MM_STRINGLEN1BYTE)
      *w = 1;
    else
      *w = 2;
    return p;
  }
  case MM_PREFIX_BYTES: {
    int l = sign & 0b11111;
    switch (l) {
    case MM_BYTESLEN1BYTE: *w = 1; break;
    case MM_BYTESLEN2BYTE: *w = 2; break;
    case MM_BYTESLEN3BYTE: *w = 3; break;
    case MM_BYTESLEN4BYTE: *w = 4; break;
    default: *w = 0; *inline_len = l;
    }
    return p;
  }
  case MM_PREFIX_CONTAINER: {
    int l = sign & 0b01111;
    switch (l) {
    case MM_CONTAINERLEN1BYTE: *w = 1; break;
    case MM_CONTAINERLEN2BYTE: *w = 2; break;
    case MM_CONTAINERLEN3BYTE: *w = 3; break;
    case MM_CONTAINERLEN4BYTE: *w = 4; break;
    default: *w = 0; *inline_len = l;
    }
    return p;
  }
  case MM_PREFIX_TAG: {
    int l = sign & 0b11111;
    switch (l) {
    case MM_TAGLEN1BYTE: *w = 1; break;
    case MM_TAGLEN2BYTE: *w = 2; break;
    case MM_TAGLEN3BYTE: *w = 3; break;
    case MM_TAGLEN4BYTE: *w = 4; break;
    default: *w = 0; *inline_len = l;
    }
    return p;
  }
  default:
    return -1;
  }
}

static int read_tier_be(const uint8_t *data, size_t end, int w, size_t *out) {
  if ((ssize_t)end - w < 0)
    return -1;
  size_t v = 0;
  for (size_t i = end - w; i < end; i++)
    v = (v << 8) | data[i];
  *out = v;
  return 0;
}

static int tier_bytes_be(dbuf_t *b, size_t length, int w) {
  switch (w) {
  case 0:
    return 0;
  case 1:
    return db_push(b, (uint8_t)(length & 0xFF));
  case 2:
    return db_push(b, (uint8_t)((length >> 8) & 0xFF)) ||
           db_push(b, (uint8_t)(length & 0xFF));
  case 3:
    return db_push(b, (uint8_t)((length >> 16) & 0xFF)) ||
           db_push(b, (uint8_t)((length >> 8) & 0xFF)) ||
           db_push(b, (uint8_t)(length & 0xFF));
  default:
    return db_push(b, (uint8_t)((length >> 24) & 0xFF)) ||
           db_push(b, (uint8_t)((length >> 16) & 0xFF)) ||
           db_push(b, (uint8_t)((length >> 8) & 0xFF)) ||
           db_push(b, (uint8_t)(length & 0xFF));
  }
}

static int parse_footer_slot(const uint8_t *data, size_t node_end, size_t *ml,
                             size_t *slot_w) {
  if (node_end < 1)
    return -1;
  uint8_t marker = data[node_end - 1];
  if (marker == 254) {
    if (node_end < 2)
      return -1;
    *ml = data[node_end - 2];
    *slot_w = 2;
    return 0;
  } else if (marker == 255) {
    if (node_end < 3)
      return -1;
    *ml = ((size_t)data[node_end - 3] << 8) | data[node_end - 2];
    *slot_w = 3;
    return 0;
  }
  *ml = marker;
  *slot_w = 1;
  return 0;
}

static int meta_prefix_bytes(dbuf_t *b, size_t ml) {
  if (ml < 254)
    return db_push(b, (uint8_t)ml);
  if (ml <= 255)
    return db_push(b, 254) || db_push(b, (uint8_t)(ml & 0xFF));
  return db_push(b, 255) || db_push(b, (uint8_t)((ml >> 8) & 0xFF)) ||
         db_push(b, (uint8_t)(ml & 0xFF));
}

static int node_tail_start(const uint8_t *wire, size_t wlen, size_t end,
                           size_t *start) {
  if (end < 1 || end > wlen)
    return -1;
  uint8_t sign = wire[end - 1];
  int w, inl;
  int p = f_tier(sign, &w, &inl);
  if (p < 0)
    return -1;
  switch (p) {
  case MM_PREFIX_SIMPLE:
    *start = end - 1;
    return 0;
  case MM_PREFIX_POSITIVEINT:
  case MM_PREFIX_NEGATIVEINT:
    *start = (w == 0) ? end - 1 : end - 1 - (size_t)w;
    return 0;
  case MM_PREFIX_FLOAT:
    *start = (w == 0) ? end - 1 : end - 2 - (size_t)w;
    return 0;
  case MM_PREFIX_STRING:
  case MM_PREFIX_BYTES:
  case MM_PREFIX_CONTAINER: {
    if (w == 0) {
      *start = end - 1 - (size_t)inl;
    } else {
      size_t l;
      if (read_tier_be(wire, end - 1, w, &l) != 0)
        return -1;
      *start = end - 1 - (size_t)w - l;
    }
    return 0;
  }
  case MM_PREFIX_TAG: {
    size_t total;
    if (w > 0) {
      if (read_tier_be(wire, end - 1, w, &total) != 0)
        return -1;
    } else {
      total = inl;
    }
    size_t slot_end = end - 1 - (size_t)w;
    size_t ml, slot_w;
    if (parse_footer_slot(wire, slot_end, &ml, &slot_w) != 0)
      return -1;
    /* 反推用实际尾槽宽 slot_w：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
     * 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。 */
    ssize_t inner = (ssize_t)total - (ssize_t)ml - (ssize_t)slot_w;
    if (inner < 0)
      return -1;
    *start = (size_t)((ssize_t)slot_end - (ssize_t)slot_w - (ssize_t)ml - inner);
    return 0;
  }
  default:
    return -1;
  }
}

static int untailify_children(const uint8_t *payload, size_t plen, dbuf_t *out);

static int untailify_node(const uint8_t *wire, size_t start, size_t end,
                          dbuf_t *out) {
  if (start >= end)
    return -1;
  uint8_t sign = wire[end - 1];
  int w, inl;
  int p = f_tier(sign, &w, &inl);
  if (p < 0)
    return -1;

  switch (p) {
  case MM_PREFIX_SIMPLE:
    return db_push(out, sign);
  case MM_PREFIX_POSITIVEINT:
  case MM_PREFIX_NEGATIVEINT:
  case MM_PREFIX_FLOAT: {
    if (w == 0)
      return db_push(out, sign);
    if (db_push(out, sign) != 0)
      return -1;
    return db_write(out, wire + start, end - 1 - start);
  }
  case MM_PREFIX_STRING:
  case MM_PREFIX_BYTES:
  case MM_PREFIX_CONTAINER: {
    size_t l;
    if (w > 0) {
      if (read_tier_be(wire, end - 1, w, &l) != 0)
        return -1;
    } else {
      l = inl;
    }
    size_t body_len = end - 1 - (size_t)w - start;
    if (body_len != l)
      return -1;
    if (db_push(out, sign) != 0)
      return -1;
    if (tier_bytes_be(out, l, w) != 0)
      return -1;
    if (p == MM_PREFIX_CONTAINER)
      return untailify_children(wire + start, body_len, out);
    return db_write(out, wire + start, body_len);
  }
  case MM_PREFIX_TAG: {
    size_t total;
    if (w > 0) {
      if (read_tier_be(wire, end - 1, w, &total) != 0)
        return -1;
    } else {
      total = inl;
    }
    size_t slot_end = end - 1 - (size_t)w;
    size_t ml, slot_w;
    if (parse_footer_slot(wire, slot_end, &ml, &slot_w) != 0)
      return -1;
    ssize_t inner_len_s = (ssize_t)total - (ssize_t)ml - (ssize_t)slot_w;
    if (inner_len_s < 0)
      return -1;
    size_t inner_len = (size_t)inner_len_s;
    size_t meta_end = slot_end - slot_w;
    if (meta_end < ml)
      return -1;
    size_t inner_start_abs = start;
    size_t inner_wire_len = meta_end - ml - inner_start_abs;
    if (inner_wire_len != inner_len)
      return -1;

    if (db_push(out, sign) != 0)
      return -1;
    if (tier_bytes_be(out, total, w) != 0)
      return -1;
    if (meta_prefix_bytes(out, ml) != 0)
      return -1;
    /* meta：wire[meta_end-ml .. meta_end) 是逆序区，还原成正序。 */
    uint8_t *meta = malloc(ml ? ml : 1);
    if (!meta)
      return -1;
    mc_reverse_bytes(wire + (meta_end - ml), meta, ml);
    int rc = db_write(out, meta, ml);
    free(meta);
    if (rc != 0)
      return -1;
    if (inner_len > 0) {
      dbuf_t child;
      db_init(&child);
      if (untailify_node(wire, inner_start_abs, inner_start_abs + inner_len, &child) != 0) {
        db_free(&child);
        return -1;
      }
      rc = db_write(out, child.buf, child.len);
      db_free(&child);
      if (rc != 0)
        return -1;
    }
    return 0;
  }
  default:
    return -1;
  }
}

static int untailify_children(const uint8_t *payload, size_t plen, dbuf_t *out) {
  dbuf_t *nodes = NULL;
  size_t ncount = 0, ncap = 0;
  size_t end = plen;
  while (end > 0) {
    size_t start;
    if (node_tail_start(payload, plen, end, &start) != 0)
      goto fail;
    if (start >= end)
      goto fail;
    if (ncount == ncap) {
      size_t nn = ncap ? ncap * 2 : 8;
      dbuf_t *tmp = realloc(nodes, nn * sizeof(dbuf_t));
      if (!tmp)
        goto fail;
      nodes = tmp;
      ncap = nn;
    }
    db_init(&nodes[ncount]);
    if (untailify_node(payload, start, end, &nodes[ncount]) != 0) {
      db_free(&nodes[ncount]);
      goto fail;
    }
    ncount++;
    end = start;
  }
  for (size_t i = ncount; i-- > 0;) {
    if (db_write(out, nodes[i].buf, nodes[i].len) != 0) {
      goto fail;
    }
    db_free(&nodes[i]);
  }
  free(nodes);
  return 0;
fail:
  for (size_t i = 0; i < ncount; i++)
    db_free(&nodes[i]);
  free(nodes);
  return -1;
}

uint8_t *mc_untailify(const uint8_t *wire, size_t len, size_t *out_len) {
  if (len == 0)
    return NULL;
  size_t start;
  if (node_tail_start(wire, len, len, &start) != 0)
    return NULL;
  if (start != 0)
    return NULL;
  dbuf_t out;
  db_init(&out);
  if (untailify_node(wire, start, len, &out) != 0) {
    db_free(&out);
    return NULL;
  }
  *out_len = out.len;
  return out.buf;
}
