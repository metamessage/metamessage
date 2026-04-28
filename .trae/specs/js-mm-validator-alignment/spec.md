# mm-js Validator Reference mm-ts - Product Requirement Document

## Why
mm-js validator needs to align with mm-ts implementation and match Golang's validation functionality completely. Currently mm-js is missing some types and validation features that exist in Golang.

## What Changes
- Add missing validation types: decimal, ip, url, slice
- Implement all validation methods that mm-ts has
- Ensure mm-js validation behavior matches Golang exactly

## Impact
- Affected Specs: `per-type-validation`
- Affected Code:
  - `mm-js/src/mm/validator.js`

## ADDED Requirements

### Requirement: Decimal Type Validation
mm-js validator shall support decimal type validation like Golang.

#### Scenario: Decimal Validation Success
- **WHEN**: Validating a decimal string like "123.45"
- **THEN**: Returns valid result

#### Scenario: Decimal Validation Failure
- **WHEN**: Validating an invalid decimal like "abc"
- **THEN**: Returns error "invalid decimal..."

### Requirement: IP Type Validation
mm-js validator shall support IP address validation with version checking.

#### Scenario: IPv4 Validation
- **WHEN**: Validating "192.168.1.1" with tag.version = 4
- **THEN**: Returns valid result

#### Scenario: IPv6 Validation
- **WHEN**: Validating "::1" with tag.version = 6
- **THEN**: Returns valid result

### Requirement: URL Type Validation
mm-js validator shall support URL validation with scheme and host checking.

#### Scenario: Valid URL
- **WHEN**: Validating "https://example.com"
- **THEN**: Returns valid result

#### Scenario: Invalid URL Scheme
- **WHEN**: Validating "ftp://example.com"
- **THEN**: Returns error "invalid url..."

### Requirement: Slice Type Validation
mm-js validator shall support slice type validation similar to array but for slice semantics.

## MODIFIED Requirements

### Requirement: BigInt Validation Enhancement
mm-js shall handle big.Int validation with proper big.Int type checking and min/max validation.

### Requirement: DateTime Validation Enhancement
mm-js shall support date and time types separately, matching Golang's date-only and time-only validation.

## Acceptance Criteria

### AC-1: All Golang Types Supported
- **Given**: Any type that Golang validates
- **WHEN**: mm-js validates the same type
- **THEN**: Returns the same validation result

### AC-2: Type Coverage
- [x] array validation
- [x] struct validation
- [x] string validation with pattern
- [x] bytes validation with base64 encoding
- [x] bool validation
- [x] int/int8/int16/int32/int64 validation with range
- [x] uint/uint8/uint16/uint32/uint64 validation with range
- [x] float32/float64 validation with range
- [x] bigInt validation
- [x] datetime/date/time validation
- [x] uuid validation with version
- [x] email validation
- [x] enum validation
- [x] image validation
- [ ] decimal validation
- [ ] ip validation
- [ ] url validation
- [ ] slice validation
