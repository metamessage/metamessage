import Foundation

public enum MMPrefix: UInt8 {
    case simple
    case positiveInt
    case negativeInt
    case prefixFloat
    case prefixString
    case prefixBytes
    case container
    case prefixTag

    public var rawValue: UInt8 {
        switch self {
        case .simple: return 0b000 << 5
        case .positiveInt: return 0b001 << 5
        case .negativeInt: return 0b010 << 5
        case .prefixFloat: return 0b011 << 5
        case .prefixString: return 0b100 << 5
        case .prefixBytes: return 0b101 << 5
        case .container: return 0b110 << 5
        case .prefixTag: return 0b111 << 5
        }
    }

    public init?(rawValue: UInt8) {
        switch rawValue {
        case 0b000 << 5: self = .simple
        case 0b001 << 5: self = .positiveInt
        case 0b010 << 5: self = .negativeInt
        case 0b011 << 5: self = .prefixFloat
        case 0b100 << 5: self = .prefixString
        case 0b101 << 5: self = .prefixBytes
        case 0b110 << 5: self = .container
        case 0b111 << 5: self = .prefixTag
        default: return nil
        }
    }
}

public enum MMSimpleValue: UInt8 {
    case null = 0
    case nullBool = 1
    case nullInt = 2
    case nullFloat = 3
    case nullString = 4
    case nullBytes = 5
    case falseValue = 6
    case trueValue = 7
    case code = 8
    case message = 9
    case data = 10
    case success = 11
    case error = 12
    case unknown = 13
    case page = 14
    case limit = 15
    case offset = 16
    case total = 17
    case id = 18
    case name = 19
    case description = 20
    case typeValue = 21
    case version = 22
    case status = 23
    case url = 24
    case createTime = 25
    case updateTime = 26
    case deleteTime = 27
    case account = 28
    case token = 29
    case expireTime = 30
    case key = 31
}

public enum MMConstants {
    public static let max1Byte: UInt64 = 0xFF
    public static let max2Byte: UInt64 = 0xFFFF
    public static let max3Byte: UInt64 = 0xFFFFFF
    public static let max4Byte: UInt64 = 0xFFFFFFFF
    public static let max5Byte: UInt64 = 0xFFFFFFFFFF
    public static let max6Byte: UInt64 = 0xFFFFFFFFFFFF
    public static let max7Byte: UInt64 = 0xFFFFFFFFFFFFFF
    public static let max8Byte: UInt64 = 0xFFFFFFFFFFFFFFFF

    public static let intLenMask: UInt8 = 0b11111
    public static let intLen1Byte: UInt8 = 0b11000
    public static let intLen2Byte: UInt8 = 0b11001
    public static let intLen3Byte: UInt8 = 0b11010
    public static let intLen4Byte: UInt8 = 0b11011
    public static let intLen5Byte: UInt8 = 0b11100
    public static let intLen6Byte: UInt8 = 0b11101
    public static let intLen7Byte: UInt8 = 0b11110
    public static let intLen8Byte: UInt8 = 0b11111

    public static let floatPositiveNegativeMask: UInt8 = 0b10000
    public static let floatLenMask: UInt8 = 0b01111
    public static let floatLen1Byte: UInt8 = 0b01000
    public static let floatLen2Byte: UInt8 = 0b01001
    public static let floatLen3Byte: UInt8 = 0b01010
    public static let floatLen4Byte: UInt8 = 0b01011
    public static let floatLen5Byte: UInt8 = 0b01100
    public static let floatLen6Byte: UInt8 = 0b01101
    public static let floatLen7Byte: UInt8 = 0b01110
    public static let floatLen8Byte: UInt8 = 0b01111

    public static let stringLenMask: UInt8 = 0b11111
    public static let stringLen1Byte: UInt8 = 0b11110
    public static let stringLen2Byte: UInt8 = 0b11111

    public static let bytesLenMask: UInt8 = 0b11111
    public static let bytesLen1Byte: UInt8 = 0b11110
    public static let bytesLen2Byte: UInt8 = 0b11111

    public static let containerMask: UInt8 = 0b10000
    public static let containerObject: UInt8 = 0b00000
    public static let containerArray: UInt8 = 0b10000
    public static let containerLenMask: UInt8 = 0b01111
    public static let containerLen1Byte: UInt8 = 0b01110
    public static let containerLen2Byte: UInt8 = 0b01111

    public static let tagLenMask: UInt8 = 0b11111
    public static let tagLen1Byte: UInt8 = 0b11110
    public static let tagLen2Byte: UInt8 = 0b11111

    public static let prefixMask: UInt8 = 0b11100000
    public static let suffixMask: UInt8 = 0b00011111
}

public func getPrefix(_ b: UInt8) -> MMPrefix? {
    return MMPrefix(rawValue: b & MMConstants.prefixMask)
}

public func getSuffix(_ b: UInt8) -> UInt8 {
    return b & MMConstants.suffixMask
}

public func intLen(_ b: UInt8) -> (extraBytes: Int, len: Int) {
    let l = Int(b & MMConstants.intLenMask)
    if l < Int(MMConstants.intLen1Byte) {
        return (0, l)
    } else {
        return (l - Int(MMConstants.intLen1Byte) + 1, 0)
    }
}

public func floatLen(_ b: UInt8) -> (extraBytes: Int, len: Int) {
    let l = Int(b & MMConstants.floatLenMask)
    if l < Int(MMConstants.floatLen1Byte) {
        return (0, l)
    } else {
        return (l - Int(MMConstants.floatLen1Byte) + 1, 0)
    }
}

public func stringLen(_ b: UInt8) -> (extraBytes: Int, len: Int) {
    let l = Int(b & MMConstants.stringLenMask)
    if l < Int(MMConstants.stringLen1Byte) {
        return (0, l)
    } else if l == Int(MMConstants.stringLen1Byte) {
        return (1, l)
    } else {
        return (2, l)
    }
}

public func bytesLen(_ b: UInt8) -> (extraBytes: Int, len: Int) {
    let l = Int(b & MMConstants.bytesLenMask)
    if l < Int(MMConstants.bytesLen1Byte) {
        return (0, l)
    } else if l == Int(MMConstants.bytesLen1Byte) {
        return (1, l)
    } else {
        return (2, l)
    }
}

public func containerLen(_ b: UInt8) -> (extraBytes: Int, len: Int) {
    let l = Int(b & MMConstants.containerLenMask)
    if l < Int(MMConstants.containerLen1Byte) {
        return (0, l)
    } else if l == Int(MMConstants.containerLen1Byte) {
        return (1, l)
    } else {
        return (2, l)
    }
}

public func tagLen(_ b: UInt8) -> (extraBytes: Int, len: Int) {
    let l = Int(b & MMConstants.tagLenMask)
    if l < Int(MMConstants.tagLen1Byte) {
        return (0, l)
    } else if l == Int(MMConstants.tagLen1Byte) {
        return (1, l)
    } else {
        return (2, l)
    }
}

public func isArray(_ b: UInt8) -> Bool {
    return (b & MMConstants.containerMask) == MMConstants.containerArray
}