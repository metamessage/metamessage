// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "MetaMessage",
    platforms: [
        .macOS(.v10_15),
        .iOS(.v13),
    ],
    products: [
        .library(
            name: "MetaMessage",
            targets: ["MetaMessage"]),
    ],
    dependencies: [],
    targets: [
        .target(
            name: "MetaMessage",
            dependencies: [],
            path: "Sources/MetaMessage"),
        .testTarget(
            name: "MetaMessageTests",
            dependencies: ["MetaMessage"],
            path: "Tests/MetaMessageTests"),
    ])