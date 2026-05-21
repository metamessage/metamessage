// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "mm-harness-swift",
    platforms: [.macOS(.v13)],
    dependencies: [
        .package(path: "../../../mm-swift"),
    ],
    targets: [
        .executableTarget(
            name: "mm-harness-swift",
            dependencies: [
                .product(name: "MetaMessage", package: "mm-swift"),
            ],
            path: ".",
            sources: ["main.swift"]
        ),
    ]
)