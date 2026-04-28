import Foundation

public class ValidationResult {
    public var isValid: Bool
    public var errors: [String]
    
    public init() {
        self.isValid = true
        self.errors = []
    }
    
    public func addError(_ error: String) {
        self.isValid = false
        self.errors.append(error)
    }
}

public class MmValidator {
    public init() {}
    
    public func validate(_ value: Any?, tag: JSONCTag) -> ValidationResult {
        let result = ValidationResult()
        
        if let value = value {
            switch tag.type {
            case .bool:
                validateBool(value, tag: tag, result: result)
            case .int:
                validateInt(value, tag: tag, result: result)
            case .int8:
                validateInt8(value, tag: tag, result: result)
            case .int16:
                validateInt16(value, tag: tag, result: result)
            case .int32:
                validateInt32(value, tag: tag, result: result)
            case .int64:
                validateInt64(value, tag: tag, result: result)
            case .uint:
                validateUint(value, tag: tag, result: result)
            case .uint8:
                validateUint8(value, tag: tag, result: result)
            case .uint16:
                validateUint16(value, tag: tag, result: result)
            case .uint32:
                validateUint32(value, tag: tag, result: result)
            case .uint64:
                validateUint64(value, tag: tag, result: result)
            case .bigInt:
                validateBigInt(value, tag: tag, result: result)
            case .float32:
                validateFloat32(value, tag: tag, result: result)
            case .float64, .decimal:
                validateFloat64(value, tag: tag, result: result)
            case .string:
                validateString(value, tag: tag, result: result)
            case .email:
                validateEmail(value, tag: tag, result: result)
            case .url:
                validateUrl(value, tag: tag, result: result)
            case .bytes:
                validateBytes(value, tag: tag, result: result)
            case .uuid:
                validateUUID(value, tag: tag, result: result)
            case .dateTime, .date, .time:
                validateDateTime(value, tag: tag, result: result)
            case .enumValue:
                validateEnum(value, tag: tag, result: result)
            case .array, .slice:
                validateArray(value, tag: tag, result: result)
            case .structType:
                validateStruct(value, tag: tag, result: result)
            default:
                break
            }
        } else {
            if !tag.nullable {
                result.addError("value is required")
            }
        }
        
        return result
    }
    
    private func validateBool(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if !(value is Bool) {
            result.addError("value must be a boolean")
        }
    }

    private func validateInt(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < Int(Int.min) || intValue > Int.max {
            result.addError("value out of range for int")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as int: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as int: \(maxStr)")
            }
        }
    }

    private func validateInt8(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < -128 || intValue > 127 {
            result.addError("value out of range for int8")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < -128 || mini > 127 {
                    result.addError("tag.min \(mini) is out of int8 range [-128, 127]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as int8: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < -128 || maxi > 127 {
                    result.addError("tag.max \(maxi) is out of int8 range [-128, 127]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as int8: \(maxStr)")
            }
        }
    }

    private func validateInt16(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < -32768 || intValue > 32767 {
            result.addError("value out of range for int16")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < -32768 || mini > 32767 {
                    result.addError("tag.min \(mini) is out of int16 range [-32768, 32767]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as int16: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < -32768 || maxi > 32767 {
                    result.addError("tag.max \(maxi) is out of int16 range [-32768, 32767]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as int16: \(maxStr)")
            }
        }
    }

    private func validateInt32(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < -2147483648 || intValue > 2147483647 {
            result.addError("value out of range for int32")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < -2147483648 || mini > 2147483647 {
                    result.addError("tag.min \(mini) is out of int32 range [-2147483648, 2147483647]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as int32: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < -2147483648 || maxi > 2147483647 {
                    result.addError("tag.max \(maxi) is out of int32 range [-2147483648, 2147483647]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as int32: \(maxStr)")
            }
        }
    }

    private func validateInt64(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < -9223372036854775808 || intValue > 9223372036854775807 {
            result.addError("value out of range for int64")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < -9223372036854775808 || mini > 9223372036854775807 {
                    result.addError("tag.min \(mini) is out of int64 range [-9223372036854775808, 9223372036854775807]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as int64: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < -9223372036854775808 || maxi > 9223372036854775807 {
                    result.addError("tag.max \(maxi) is out of int64 range [-9223372036854775808, 9223372036854775807]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as int64: \(maxStr)")
            }
        }
    }

    private func validateUint(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < 0 || intValue > 4294967295 {
            result.addError("value out of range for uint")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < 0 || mini > 4294967295 {
                    result.addError("tag.min \(mini) is out of uint range [0, 4294967295]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as uint: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < 0 || maxi > 4294967295 {
                    result.addError("tag.max \(maxi) is out of uint range [0, 4294967295]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as uint: \(maxStr)")
            }
        }
    }

    private func validateUint8(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < 0 || intValue > 255 {
            result.addError("value out of range for uint8")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < 0 || mini > 255 {
                    result.addError("tag.min \(mini) is out of uint8 range [0, 255]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as uint8: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < 0 || maxi > 255 {
                    result.addError("tag.max \(maxi) is out of uint8 range [0, 255]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as uint8: \(maxStr)")
            }
        }
    }

    private func validateUint16(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < 0 || intValue > 65535 {
            result.addError("value out of range for uint16")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < 0 || mini > 65535 {
                    result.addError("tag.min \(mini) is out of uint16 range [0, 65535]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as uint16: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < 0 || maxi > 65535 {
                    result.addError("tag.max \(maxi) is out of uint16 range [0, 65535]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as uint16: \(maxStr)")
            }
        }
    }

    private func validateUint32(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < 0 || intValue > 4294967295 {
            result.addError("value out of range for uint32")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < 0 || mini > 4294967295 {
                    result.addError("tag.min \(mini) is out of uint32 range [0, 4294967295]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as uint32: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < 0 || maxi > 4294967295 {
                    result.addError("tag.max \(maxi) is out of uint32 range [0, 4294967295]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as uint32: \(maxStr)")
            }
        }
    }

    private func validateUint64(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let intValue = value as? Int else {
            result.addError("value must be a number")
            return
        }
        if intValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        if intValue < 0 {
            result.addError("value out of range for uint64")
        }
        
        if let minStr = tag.min {
            if let mini = Int(minStr) {
                if mini < 0 {
                    result.addError("tag.min \(mini) is out of uint64 range [0, 18446744073709551615]")
                } else if intValue < mini {
                    result.addError("value \(intValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as uint64: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Int(maxStr) {
                if maxi < 0 {
                    result.addError("tag.max \(maxi) is out of uint64 range [0, 18446744073709551615]")
                } else if intValue > maxi {
                    result.addError("value \(intValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as uint64: \(maxStr)")
            }
        }
    }

    private func validateBigInt(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let stringValue = value as? String else {
            result.addError("value must be a string")
            return
        }
        if stringValue.isEmpty && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        let bigIntRegex = "^-?[0-9]+$"
        let bigIntPredicate = NSPredicate(format: "SELF MATCHES %@", bigIntRegex)
        if !bigIntPredicate.evaluate(with: stringValue) {
            result.addError("value is not a valid bigInt")
        }
    }

    private func validateFloat32(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let floatValue = value as? Float else {
            result.addError("value must be a float")
            return
        }
        if floatValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        
        if let minStr = tag.min {
            if let mini = Float(minStr) {
                if mini < -3.402823466e+38 || mini > 3.402823466e+38 {
                    result.addError("tag.min \(mini) is out of float32 range [-3.402823466e+38, 3.402823466e+38]")
                } else if floatValue < mini {
                    result.addError("value \(floatValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as float32: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Float(maxStr) {
                if maxi < -3.402823466e+38 || maxi > 3.402823466e+38 {
                    result.addError("tag.max \(maxi) is out of float32 range [-3.402823466e+38, 3.402823466e+38]")
                } else if floatValue > maxi {
                    result.addError("value \(floatValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as float32: \(maxStr)")
            }
        }
    }

    private func validateFloat64(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        guard let doubleValue = value as? Double else {
            result.addError("value must be a double")
            return
        }
        if doubleValue == 0 && !tag.allowEmpty {
            result.addError("value is required")
            return
        }
        
        if let minStr = tag.min {
            if let mini = Double(minStr) {
                if mini < -1.7976931348623157e+308 || mini > 1.7976931348623157e+308 {
                    result.addError("tag.min \(mini) is out of float64 range [-1.7976931348623157e+308, 1.7976931348623157e+308]")
                } else if doubleValue < mini {
                    result.addError("value \(doubleValue) is less than the minimum limit \(mini)")
                }
            } else {
                result.addError("failed to parse tag.min as float64: \(minStr)")
            }
        }
        
        if let maxStr = tag.max {
            if let maxi = Double(maxStr) {
                if maxi < -1.7976931348623157e+308 || maxi > 1.7976931348623157e+308 {
                    result.addError("tag.max \(maxi) is out of float64 range [-1.7976931348623157e+308, 1.7976931348623157e+308]")
                } else if doubleValue > maxi {
                    result.addError("value \(doubleValue) exceeds the maximum limit \(maxi)")
                }
            } else {
                result.addError("failed to parse tag.max as float64: \(maxStr)")
            }
        }
    }
    
    private func validateString(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if let strValue = value as? String {
            if strValue.isEmpty && !tag.allowEmpty {
                result.addError("value is empty")
            }
        } else {
            result.addError("value must be a string")
        }
    }
    
    private func validateEmail(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if let email = value as? String {
            if email.isEmpty && !tag.allowEmpty {
                result.addError("value is empty")
                return
            }
            
            if !email.isEmpty {
                let emailRegex = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
                let emailPredicate = NSPredicate(format: "SELF MATCHES %@", emailRegex)
                if !emailPredicate.evaluate(with: email) {
                    result.addError("value is not a valid email")
                }
            }
        } else {
            result.addError("value must be a string")
        }
    }
    
    private func validateUrl(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if let url = value as? String {
            if url.isEmpty && !tag.allowEmpty {
                result.addError("value is empty")
                return
            }
            
            if !url.isEmpty {
                if let urlObj = URL(string: url) {
                    if !urlObj.scheme!.starts(with: "http") {
                        result.addError("value is not a valid url")
                    }
                } else {
                    result.addError("value is not a valid url")
                }
            }
        } else {
            result.addError("value must be a string")
        }
    }
    
    private func validateBytes(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if let bytes = value as? Data {
            if bytes.isEmpty && !tag.allowEmpty {
                result.addError("value is empty")
            }
        } else {
            result.addError("value must be a byte array")
        }
    }
    
    private func validateUUID(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if let uuid = value as? String {
            if uuid.isEmpty && !tag.allowEmpty {
                result.addError("value is empty")
                return
            }
            
            if !uuid.isEmpty {
                if UUID(uuidString: uuid) == nil {
                    result.addError("value is not a valid uuid")
                }
            }
        } else {
            result.addError("value must be a string")
        }
    }
    
    private func validateDateTime(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if !(value is String) && !(value is Date) {
            result.addError("value must be a datetime")
        }
    }
    
    private func validateEnum(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if let enumValue = value as? String {
            if enumValue.isEmpty && !tag.allowEmpty {
                result.addError("value is empty")
                return
            }
            
            if !enumValue.isEmpty && !tag.enumValues.isEmpty {
                let enumValues = tag.enumValues.split(separator: "|").map { $0.trimmingCharacters(in: .whitespaces) }
                if !enumValues.contains(enumValue) {
                    result.addError("value is not in enum")
                }
            }
        } else {
            result.addError("value must be a string")
        }
    }
    
    private func validateArray(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if let array = value as? [Any] {
            if array.isEmpty && !tag.allowEmpty {
                result.addError("value is empty")
            }
            
            // 验证数组元素
            for (index, item) in array.enumerated() {
                let childTag = JSONCTag()
                childTag.inherit(from: tag)
                let itemResult = validate(item, tag: childTag)
                if !itemResult.isValid {
                    for error in itemResult.errors {
                        result.addError("item \(index): \(error)")
                    }
                }
            }
        } else {
            result.addError("value must be an array")
        }
    }
    
    private func validateStruct(_ value: Any, tag: JSONCTag, result: ValidationResult) {
        if !(value is [String: Any]) {
            result.addError("value must be an object")
        }
    }
}

public let validator = MmValidator()
