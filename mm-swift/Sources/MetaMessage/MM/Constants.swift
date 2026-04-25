import Foundation

public enum MMPrefix: UInt8 {
    case simple = 0b000 << 5
    case positiveInt = 0b001 << 5
    case negativeInt = 0b010 << 5
    case prefixFloat = 0b011 << 5
    case prefixString = 0b100 << 5
    case prefixBytes = 0b101 << 5
    case container = 0b110 << 5
    case prefixTag = 0b111 << 5
}

public enum MMSimpleValue: UInt8 {
    case nullBool = 0
    case nullInt = 1
    case nullFloat = 2
    case nullString = 3
    case nullBytes = 4
    case falseValue = 5
    case trueValue = 6
    case code = 7
    case message = 8
    case data = 9
    case success = 10
    case error = 11
    case unknown = 12
    case page = 13
    case limit = 14
    case offset = 15
    case total = 16
    case id = 17
    case name = 18
    case description = 19
    case typeValue = 20
    case version = 21
    case status = 22
    case url = 23
    case createTime = 24
    case updateTime = 25
    case deleteTime = 26
    case account = 27
    case token = 28
    case expireTime = 29
    case key = 30
    case val = 31
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
    public static let floatLen1Byte: UInt8 = 0b00111
    public static let floatLen2Byte: UInt8 = 0b01000
    public static let floatLen3Byte: UInt8 = 0b01001
    public static let floatLen4Byte: UInt8 = 0b01010
    public static let floatLen5Byte: UInt8 = 0b01011
    public static let floatLen6Byte: UInt8 = 0b01100
    public static let floatLen7Byte: UInt8 = 0b01101
    public static let floatLen8Byte: UInt8 = 0b01110

    public static let stringLenMask: UInt8 = 0b11111
    public static let stringLen1Byte: UInt8 = 0b11110
    public static let stringLen2Byte: UInt8 = 0b11111

    public static let bytesLenMask: UInt8 = 0b11111
    public static let bytesLen1Byte: UInt8 = 0b11110
    public static let bytesLen2Byte: UInt8 = 0b11111

    public static let containerMask: UInt8 = 0b10000
    public static let containerMap: UInt8 = 0b00000
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