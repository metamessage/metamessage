package io.github.metamessage.core

// 与 Go 基线 internal/core 的 maxDepth 对齐：JSONC 解析与 valueToNode 共享同一深度上限，
// 错误消息格式统一为 "max depth: $MAX_DEPTH"。
object Limits {
    const val MAX_DEPTH = 1024
}
