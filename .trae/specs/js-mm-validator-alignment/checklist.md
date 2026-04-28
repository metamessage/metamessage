# mm-js Validator Alignment Checklist

## Decimal Type
- [x] Decimal validation method added
- [x] Decimal validation returns error for invalid format
- [x] Decimal allows empty when allowEmpty is true

## IP Type
- [x] IP validation method added
- [x] IPv4 validation works correctly
- [x] IPv6 validation works correctly
- [x] IP version checking works
- [x] IP allows empty when allowEmpty is true

## URL Type
- [x] URL validation method added
- [x] URL validates scheme (http/https only)
- [x] URL validates host is not empty
- [x] URL allows empty when allowEmpty is true

## Slice Type
- [x] Slice validation method added
- [x] Slice validation similar to array
- [x] Slice childUnique checking works

## BigInt Enhancement
- [x] BigInt validation with min constraint
- [x] BigInt validation with max constraint
- [x] BigInt validation with allowEmpty

## Date/Time Enhancement
- [x] Date-only validation method added
- [x] Time-only validation method added
- [x] Date validation with location support
- [x] Time validation with location support

## Integration
- [x] validate() method dispatches to new validation methods
- [x] Validator syntax check passed
- [x] Validation results match Golang format
