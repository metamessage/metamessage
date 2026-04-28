use crate::jsonc::tag::Tag;
use crate::jsonc::value_type::ValueType;
use num_bigint::BigInt;
use std::any::Any;
use std::fmt::Debug;
use std::sync::LazyLock;
use regex::Regex;

#[derive(Debug, Clone)]
pub struct ValidationResult {
    pub is_valid: bool,
    pub errors: Vec<String>,
}

impl ValidationResult {
    pub fn new() -> Self {
        Self {
            is_valid: true,
            errors: Vec::new(),
        }
    }

    pub fn add_error(&mut self, error: String) {
        self.is_valid = false;
        self.errors.push(error);
    }
}

pub struct MmValidator {
    email_regex: &'static Regex,
}

impl Default for MmValidator {
    fn default() -> Self {
        Self {
            email_regex: &EMAIL_REGEX,
        }
    }
}

impl MmValidator {
    pub fn new() -> Self {
        Default::default()
    }

    pub fn validate(&self, value: &dyn Any, tag: &Tag) -> ValidationResult {
        let mut result = ValidationResult::new();

        match tag.value_type {
            ValueType::Bool => self.validate_bool(value, tag, &mut result),
            ValueType::Int => self.validate_int(value, tag, &mut result),
            ValueType::Int8 => self.validate_int8(value, tag, &mut result),
            ValueType::Int16 => self.validate_int16(value, tag, &mut result),
            ValueType::Int32 => self.validate_int32(value, tag, &mut result),
            ValueType::Int64 => self.validate_int64(value, tag, &mut result),
            ValueType::Uint => self.validate_uint(value, tag, &mut result),
            ValueType::Uint8 => self.validate_uint8(value, tag, &mut result),
            ValueType::Uint16 => self.validate_uint16(value, tag, &mut result),
            ValueType::Uint32 => self.validate_uint32(value, tag, &mut result),
            ValueType::Uint64 => self.validate_uint64(value, tag, &mut result),
            ValueType::BigInt => self.validate_big_int(value, tag, &mut result),
            ValueType::Float32 => self.validate_float32(value, tag, &mut result),
            ValueType::Float64 => self.validate_float64(value, tag, &mut result),
            ValueType::Decimal => self.validate_float(value, tag, &mut result),
            ValueType::String => self.validate_string(value, tag, &mut result),
            ValueType::Email => self.validate_email(value, tag, &mut result),
            ValueType::URL => self.validate_url(value, tag, &mut result),
            ValueType::Bytes => self.validate_bytes(value, tag, &mut result),
            ValueType::UUID => self.validate_uuid(value, tag, &mut result),
            ValueType::DateTime | ValueType::Date | ValueType::Time => self.validate_datetime(value, tag, &mut result),
            ValueType::Enum => self.validate_enum(value, tag, &mut result),
            ValueType::Array | ValueType::Slice => self.validate_array(value, tag, &mut result),
            ValueType::Struct => self.validate_struct(value, tag, &mut result),
            _ => {}
        }

        result
    }

    fn validate_bool(&self, value: &dyn Any, _tag: &Tag, result: &mut ValidationResult) {
        if value.downcast_ref::<bool>().is_some() {
        } else {
            result.add_error("value must be a boolean".to_string());
        }
    }

    fn validate_int(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<i64>() {
            if *int_value == 0 {
                result.add_error("type int not allow empty value 0".to_string());
                return;
            }
            let val64 = *int_value;
            if let Some(ref min_str) = tag.min {
                if let Ok(mini) = min_str.parse::<i64>() {
                    if val64 < mini {
                        result.add_error(format!("value {} is less than the minimum limit {}", val64, mini));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                if let Ok(maxi) = max_str.parse::<i64>() {
                    if val64 > maxi {
                        result.add_error(format!("value {} exceeds the maximum limit {}", val64, maxi));
                    }
                }
            }
        } else {
            result.add_error("value must be a number".to_string());
        }
    }

    fn validate_int8(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<i8>() {
            if *int_value == 0 {
                result.add_error("type int8 not allow empty value 0".to_string());
                return;
            }
            let val64 = *int_value as i64;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<i64>() {
                    Ok(mini) => {
                        if mini < i8::MIN as i64 || mini > i8::MAX as i64 {
                            result.add_error(format!("failed to parse tag.min as int8: {}", min_str));
                        } else if val64 < mini {
                            result.add_error(format!("value {} is less than the minimum limit {}", val64, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as int8: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<i64>() {
                    Ok(maxi) => {
                        if maxi < i8::MIN as i64 || maxi > i8::MAX as i64 {
                            result.add_error(format!("failed to parse tag.max as int8: {}", max_str));
                        } else if val64 > maxi {
                            result.add_error(format!("value {} exceeds the maximum limit {}", val64, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as int8: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be an int8".to_string());
        }
    }

    fn validate_int16(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<i16>() {
            if *int_value == 0 {
                result.add_error("type int16 not allow empty value 0".to_string());
                return;
            }
            let val64 = *int_value as i64;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<i64>() {
                    Ok(mini) => {
                        if mini < i16::MIN as i64 || mini > i16::MAX as i64 {
                            result.add_error(format!("failed to parse tag.min as int16: {}", min_str));
                        } else if val64 < mini {
                            result.add_error(format!("value {} is less than the minimum limit {}", val64, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as int16: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<i64>() {
                    Ok(maxi) => {
                        if maxi < i16::MIN as i64 || maxi > i16::MAX as i64 {
                            result.add_error(format!("failed to parse tag.max as int16: {}", max_str));
                        } else if val64 > maxi {
                            result.add_error(format!("value {} exceeds the maximum limit {}", val64, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as int16: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be an int16".to_string());
        }
    }

    fn validate_int32(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<i32>() {
            if *int_value == 0 {
                result.add_error("type int32 not allow empty value 0".to_string());
                return;
            }
            let val64 = *int_value as i64;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<i64>() {
                    Ok(mini) => {
                        if mini < i32::MIN as i64 || mini > i32::MAX as i64 {
                            result.add_error(format!("failed to parse tag.min as int32: {}", min_str));
                        } else if val64 < mini {
                            result.add_error(format!("value {} is less than the minimum limit {}", val64, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as int32: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<i64>() {
                    Ok(maxi) => {
                        if maxi < i32::MIN as i64 || maxi > i32::MAX as i64 {
                            result.add_error(format!("failed to parse tag.max as int32: {}", max_str));
                        } else if val64 > maxi {
                            result.add_error(format!("value {} exceeds the maximum limit {}", val64, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as int32: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be an int32".to_string());
        }
    }

    fn validate_int64(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<i64>() {
            if *int_value == 0 {
                result.add_error("type int64 not allow empty value 0".to_string());
                return;
            }
            let val = *int_value;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<i64>() {
                    Ok(mini) => {
                        if val < mini {
                            result.add_error(format!("value {} is less than the minimum limit {}", val, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as int64: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<i64>() {
                    Ok(maxi) => {
                        if val > maxi {
                            result.add_error(format!("value {} exceeds the maximum limit {}", val, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as int64: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be an int64".to_string());
        }
    }

    fn validate_uint(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<u32>() {
            if *int_value == 0 {
                result.add_error("type uint not allow empty value 0".to_string());
                return;
            }
            let val64 = *int_value as u64;
            if let Some(ref min_str) = tag.min {
                if let Ok(mini) = min_str.parse::<u64>() {
                    if val64 < mini {
                        result.add_error(format!("value {} is less than the minimum limit {}", val64, mini));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                if let Ok(maxi) = max_str.parse::<u64>() {
                    if val64 > maxi {
                        result.add_error(format!("value {} exceeds the maximum limit {}", val64, maxi));
                    }
                }
            }
        } else {
            result.add_error("value must be a uint".to_string());
        }
    }

    fn validate_uint8(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<u8>() {
            if *int_value == 0 {
                result.add_error("type uint8 not allow empty value 0".to_string());
                return;
            }
            let val64 = *int_value as u64;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<u64>() {
                    Ok(mini) => {
                        if mini < u8::MIN as u64 || mini > u8::MAX as u64 {
                            result.add_error(format!("failed to parse tag.min as uint8: {}", min_str));
                        } else if val64 < mini {
                            result.add_error(format!("value {} is less than the minimum limit {}", val64, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as uint8: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<u64>() {
                    Ok(maxi) => {
                        if maxi < u8::MIN as u64 || maxi > u8::MAX as u64 {
                            result.add_error(format!("failed to parse tag.max as uint8: {}", max_str));
                        } else if val64 > maxi {
                            result.add_error(format!("value {} exceeds the maximum limit {}", val64, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as uint8: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be a uint8".to_string());
        }
    }

    fn validate_uint16(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<u16>() {
            if *int_value == 0 {
                result.add_error("type uint16 not allow empty value 0".to_string());
                return;
            }
            let val64 = *int_value as u64;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<u64>() {
                    Ok(mini) => {
                        if mini < u16::MIN as u64 || mini > u16::MAX as u64 {
                            result.add_error(format!("failed to parse tag.min as uint16: {}", min_str));
                        } else if val64 < mini {
                            result.add_error(format!("value {} is less than the minimum limit {}", val64, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as uint16: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<u64>() {
                    Ok(maxi) => {
                        if maxi < u16::MIN as u64 || maxi > u16::MAX as u64 {
                            result.add_error(format!("failed to parse tag.max as uint16: {}", max_str));
                        } else if val64 > maxi {
                            result.add_error(format!("value {} exceeds the maximum limit {}", val64, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as uint16: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be a uint16".to_string());
        }
    }

    fn validate_uint32(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<u32>() {
            if *int_value == 0 {
                result.add_error("type uint32 not allow empty value 0".to_string());
                return;
            }
            let val64 = *int_value as u64;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<u64>() {
                    Ok(mini) => {
                        if mini < u32::MIN as u64 || mini > u32::MAX as u64 {
                            result.add_error(format!("failed to parse tag.min as uint32: {}", min_str));
                        } else if val64 < mini {
                            result.add_error(format!("value {} is less than the minimum limit {}", val64, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as uint32: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<u64>() {
                    Ok(maxi) => {
                        if maxi < u32::MIN as u64 || maxi > u32::MAX as u64 {
                            result.add_error(format!("failed to parse tag.max as uint32: {}", max_str));
                        } else if val64 > maxi {
                            result.add_error(format!("value {} exceeds the maximum limit {}", val64, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as uint32: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be a uint32".to_string());
        }
    }

    fn validate_uint64(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(int_value) = value.downcast_ref::<u64>() {
            if *int_value == 0 {
                result.add_error("type uint64 not allow empty value 0".to_string());
                return;
            }
            let val = *int_value;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<u64>() {
                    Ok(mini) => {
                        if val < mini {
                            result.add_error(format!("value {} is less than the minimum limit {}", val, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as uint64: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<u64>() {
                    Ok(maxi) => {
                        if val > maxi {
                            result.add_error(format!("value {} exceeds the maximum limit {}", val, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as uint64: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be a uint64".to_string());
        }
    }

    fn validate_big_int(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(big_int) = value.downcast_ref::<BigInt>() {
            if big_int.sign() == num_bigint::Sign::NoSign {
                result.add_error("type big.Int not allow empty value 0".to_string());
                return;
            }
            if let Some(ref min_str) = tag.min {
                if let Ok(mini) = min_str.parse::<BigInt>() {
                    if big_int < &mini {
                        result.add_error(format!("big.Int length {} < min {}", big_int, mini));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                if let Ok(maxi) = max_str.parse::<BigInt>() {
                    if big_int > &maxi {
                        result.add_error(format!("big.Int length {} > max {}", big_int, maxi));
                    }
                }
            }
        } else {
            result.add_error("value must be a big.Int".to_string());
        }
    }

    fn validate_float32(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(float_value) = value.downcast_ref::<f32>() {
            if *float_value == 0.0 {
                result.add_error("type float32 not allow empty value 0.0".to_string());
                return;
            }
            let val64 = *float_value as f64;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<f64>() {
                    Ok(mini) => {
                        if mini < f32::MIN as f64 || mini > f32::MAX as f64 {
                            result.add_error(format!("failed to parse tag.min as float32: {}", min_str));
                        } else if val64 < mini {
                            result.add_error(format!("{} < min {}", val64, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as float32: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<f64>() {
                    Ok(maxi) => {
                        if maxi < f32::MIN as f64 || maxi > f32::MAX as f64 {
                            result.add_error(format!("failed to parse tag.max as float32: {}", max_str));
                        } else if val64 > maxi {
                            result.add_error(format!("{} > max {}", val64, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as float32: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be a float32".to_string());
        }
    }

    fn validate_float64(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(float_value) = value.downcast_ref::<f64>() {
            if *float_value == 0.0 {
                result.add_error("type float64 not allow empty value 0.0".to_string());
                return;
            }
            let val = *float_value;
            if let Some(ref min_str) = tag.min {
                match min_str.parse::<f64>() {
                    Ok(mini) => {
                        if mini < f64::MIN || mini > f64::MAX {
                            result.add_error(format!("failed to parse tag.min as float64: {}", min_str));
                        } else if val < mini {
                            result.add_error(format!("{} < min {}", val, mini));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.min as float64: {}", min_str));
                    }
                }
            }
            if let Some(ref max_str) = tag.max {
                match max_str.parse::<f64>() {
                    Ok(maxi) => {
                        if maxi < f64::MIN || maxi > f64::MAX {
                            result.add_error(format!("failed to parse tag.max as float64: {}", max_str));
                        } else if val > maxi {
                            result.add_error(format!("{} > max {}", val, maxi));
                        }
                    }
                    Err(_) => {
                        result.add_error(format!("failed to parse tag.max as float64: {}", max_str));
                    }
                }
            }
        } else {
            result.add_error("value must be a float64".to_string());
        }
    }

    fn validate_float(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(float_value) = value.downcast_ref::<f64>() {
            match tag.value_type {
                ValueType::Float32 => {
                    if *float_value < f32::MIN as f64 || *float_value > f32::MAX as f64 {
                        result.add_error("value out of range for float32".to_string());
                    }
                }
                ValueType::Float64 => {
                    if *float_value < f64::MIN || *float_value > f64::MAX {
                        result.add_error("value out of range for float64".to_string());
                    }
                }
                ValueType::Decimal => {
                    if let Some(ref min_str) = tag.min {
                        if let Ok(mini) = min_str.parse::<f64>() {
                            if *float_value < mini {
                                result.add_error(format!("{} < min {}", float_value, mini));
                            }
                        }
                    }
                    if let Some(ref max_str) = tag.max {
                        if let Ok(maxi) = max_str.parse::<f64>() {
                            if *float_value > maxi {
                                result.add_error(format!("{} > max {}", float_value, maxi));
                            }
                        }
                    }
                }
                _ => {}
            }
        } else {
            result.add_error("value must be a float".to_string());
        }
    }

    fn validate_string(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(str_value) = value.downcast_ref::<String>() {
            if str_value.is_empty() {
                result.add_error("value is empty".to_string());
            }
        } else {
            result.add_error("value must be a string".to_string());
        }
    }

    fn validate_email(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(email) = value.downcast_ref::<String>() {
            if email.is_empty() {
                result.add_error("value is empty".to_string());
            } else if !self.email_regex.is_match(email) {
                result.add_error("value is not a valid email".to_string());
            }
        } else {
            result.add_error("value must be a string".to_string());
        }
    }

    fn validate_url(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(url) = value.downcast_ref::<String>() {
            if url.is_empty() {
                result.add_error("value is empty".to_string());
            } else if !url.starts_with("http://") && !url.starts_with("https://") {
                result.add_error("value is not a valid url".to_string());
            }
        } else {
            result.add_error("value must be a string".to_string());
        }
    }

    fn validate_bytes(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(bytes) = value.downcast_ref::<Vec<u8>>() {
            if bytes.is_empty() {
                result.add_error("value is empty".to_string());
            }
        } else {
            result.add_error("value must be a byte array".to_string());
        }
    }

    fn validate_uuid(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(uuid) = value.downcast_ref::<String>() {
            if uuid.is_empty() {
                result.add_error("value is empty".to_string());
            } else if uuid.len() != 36 {
                result.add_error("value is not a valid uuid".to_string());
            }
        } else {
            result.add_error("value must be a string".to_string());
        }
    }

    fn validate_datetime(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(datetime) = value.downcast_ref::<String>() {
            if datetime.is_empty() {
                result.add_error("value is empty".to_string());
            }
        } else {
            result.add_error("value must be a string".to_string());
        }
    }

    fn validate_enum(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(enum_value) = value.downcast_ref::<String>() {
            if enum_value.is_empty() {
                result.add_error("value is empty".to_string());
            } else if let Some(enum_values) = &tag.enum_values {
                if !enum_values.contains(enum_value) {
                    result.add_error("value is not in enum".to_string());
                }
            }
        } else {
            result.add_error("value must be a string".to_string());
        }
    }

    fn validate_array(&self, value: &dyn Any, tag: &Tag, result: &mut ValidationResult) {
        if let Some(array) = value.downcast_ref::<Vec<Box<dyn Any>>>() {
            if array.is_empty() {
                result.add_error("value is empty".to_string());
            }
            for (index, item) in array.iter().enumerate() {
                let child_tag = Tag {
                    value_type: tag.child_type.unwrap_or(ValueType::Unknown),
                    ..tag.clone()
                };
                let item_result = self.validate(item.as_ref(), &child_tag);
                if !item_result.is_valid {
                    for error in item_result.errors {
                        result.add_error(format!("item {}: {}", index, error));
                    }
                }
            }
        } else {
            result.add_error("value must be an array".to_string());
        }
    }

    fn validate_struct(&self, _value: &dyn Any, _tag: &Tag, _result: &mut ValidationResult) {
    }
}

static EMAIL_REGEX: LazyLock<Regex> = LazyLock::new(|| {
    Regex::new(r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$").unwrap()
});

pub static VALIDATOR: LazyLock<MmValidator> = LazyLock::new(|| MmValidator::new());

pub fn validate(value: &dyn Any, tag: &Tag) -> ValidationResult {
    VALIDATOR.validate(value, tag)
}
