import Foundation

enum Mime {
    private static let strToMime: [String: UInt64] = {
        var dict: [String: UInt64] = [:]
        dict["image/jpeg"] = 1
        dict["image/jpg"] = 1
        dict["image/png"] = 2
        dict["image/gif"] = 3
        dict["image/webp"] = 4
        dict["image/svg+xml"] = 5
        dict["image/avif"] = 6
        dict["image/bmp"] = 7
        dict["image/x-icon"] = 8
        dict["image/tiff"] = 9
        dict["image/heic"] = 10
        dict["image/heif"] = 11
        dict["text/plain"] = 12
        dict["text/html"] = 13
        dict["text/css"] = 14
        dict["text/javascript"] = 15
        dict["application/json"] = 16
        dict["text/csv"] = 17
        dict["text/markdown"] = 18
        dict["application/pdf"] = 19
        dict["application/zip"] = 20
        dict["application/gzip"] = 21
        dict["application/x-tar"] = 22
        dict["application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"] = 23
        dict["application/vnd.openxmlformats-officedocument.wordprocessingml.document"] = 24
        dict["application/octet-stream"] = 25
        dict["video/mp4"] = 26
        dict["video/webm"] = 27
        dict["video/quicktime"] = 28
        dict["audio/mpeg"] = 29
        dict["audio/wav"] = 30
        dict["audio/flac"] = 31
        dict["font/woff2"] = 32
        dict["font/ttf"] = 33
        return dict
    }()

    static func parse(_ s: String) -> UInt64 {
        let key = s.trimmingCharacters(in: .whitespaces).lowercased()
        return strToMime[key] ?? 0
    }

    static func toString(_ code: UInt64) -> String {
        switch code {
        case 1: return "image/jpeg"
        case 2: return "image/png"
        case 3: return "image/gif"
        case 4: return "image/webp"
        case 5: return "image/svg+xml"
        case 6: return "image/avif"
        case 7: return "image/bmp"
        case 8: return "image/x-icon"
        case 9: return "image/tiff"
        case 10: return "image/heic"
        case 11: return "image/heif"
        case 12: return "text/plain"
        case 13: return "text/html"
        case 14: return "text/css"
        case 15: return "text/javascript"
        case 16: return "application/json"
        case 17: return "text/csv"
        case 18: return "text/markdown"
        case 19: return "application/pdf"
        case 20: return "application/zip"
        case 21: return "application/gzip"
        case 22: return "application/x-tar"
        case 23: return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        case 24: return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        case 25: return "application/octet-stream"
        case 26: return "video/mp4"
        case 27: return "video/webm"
        case 28: return "video/quicktime"
        case 29: return "audio/mpeg"
        case 30: return "audio/wav"
        case 31: return "audio/flac"
        case 32: return "font/woff2"
        case 33: return "font/ttf"
        default: return "unknown"
        }
    }
}