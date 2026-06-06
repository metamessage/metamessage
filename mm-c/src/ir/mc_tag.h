#ifndef MMC_IR_TAG_H
#define MMC_IR_TAG_H

#include "mc_value_type.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

enum {
  MM_TAG_KISNULL = 0 << 3,
  MM_TAG_KEXAMPLE = 1 << 3,
  MM_TAG_KDEPRECATED = 2 << 3,
  MM_TAG_KDESC = 3 << 3,
  MM_TAG_KTYPE = 4 << 3,
  MM_TAG_KNULLABLE = 5 << 3,
  MM_TAG_KALLOWEMPTY = 6 << 3,
  MM_TAG_KUNIQUE = 7 << 3,
  MM_TAG_KDEFAULTVAL = 8 << 3,
  MM_TAG_KMIN = 9 << 3,
  MM_TAG_KMAX = 10 << 3,
  MM_TAG_KSIZE = 11 << 3,
  MM_TAG_KENUMS = 12 << 3,
  MM_TAG_KPATTERN = 13 << 3,
  MM_TAG_KLOCATION = 14 << 3,
  MM_TAG_KVERSION = 15 << 3,
  MM_TAG_KMIME = 16 << 3,
  MM_TAG_KCHILDDESC = 17 << 3,
  MM_TAG_KCHILDTYPE = 18 << 3,
  MM_TAG_KCHILDNULLABLE = 19 << 3,
  MM_TAG_KCHILDALLOWEMPTY = 20 << 3,
  MM_TAG_KCHILDUNIQUE = 21 << 3,
  MM_TAG_KCHILDDEFAULTVAL = 22 << 3,
  MM_TAG_KCHILDMIN = 23 << 3,
  MM_TAG_KCHILDMAX = 24 << 3,
  MM_TAG_KCHILDSIZE = 25 << 3,
  MM_TAG_KCHILDENUMS = 26 << 3,
  MM_TAG_KCHILDPATTERN = 27 << 3,
  MM_TAG_KCHILDLOCATION = 28 << 3,
  MM_TAG_KCHILDVERSION = 29 << 3,
  MM_TAG_KCHILDMIME = 30 << 3,
  MM_TAG_KMORE = 31 << 3
};

#define MM_TAG_DEFAULT_VERSION 0
#define MM_TAG_DEFAULT_LOCATION 0

typedef struct {
  char *name;
  bool is_null;
  bool example;
  char *desc;
  mm_value_type_t type;
  bool deprecated;
  bool nullable;
  bool allow_empty;
  bool unique;
  char *default_val;
  char *min;
  char *max;
  int size;
  char *enums;
  char *pattern;
  int location_offset;
  int version;
  char *mime;
  char *child_desc;
  mm_value_type_t child_type;
  bool child_nullable;
  bool child_allow_empty;
  bool child_unique;
  char *child_default_val;
  char *child_min;
  char *child_max;
  int child_size;
  char *child_enums;
  char *child_pattern;
  int child_location_offset;
  int child_version;
  char *child_mime;
  int more;
  bool is_inherit;
} mm_tag_t;

void mm_tag_init(mm_tag_t *tag);
void mm_tag_cleanup(mm_tag_t *tag);
void mm_tag_inherit(mm_tag_t *tag, const mm_tag_t *parent);
char *mm_tag_to_string(const mm_tag_t *tag);
mm_tag_t mm_tag_parse(const char *tag_str);
uint8_t *mm_tag_bytes(const mm_tag_t *tag, size_t *out_len);
void mm_tag_merge(mm_tag_t *dst, const mm_tag_t *src);

// Validate a slice (vec) value
// Returns 0 on success, -1 on error (error message written to error_out if not NULL)
int mm_validate_vec(const mm_tag_t *tag, size_t item_count, char **error_out);
// Validate an array value
// Returns 0 on success, -1 on error (error message written to error_out if not NULL)
int mm_validate_arr(const mm_tag_t *tag, size_t item_count, char **error_out);

#endif