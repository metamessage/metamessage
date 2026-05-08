cmake_minimum_required(VERSION 3.10)
project(mmc C)

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_STANDARD_REQUIRED ON)

include_directories(src)

set(JSONC_SOURCES
    src/jsonc/mc_token.c
    src/jsonc/mc_scanner.c
    src/jsonc/mc_ast.c
    src/jsonc/mc_parser.c
)

add_library(mmc_jsonc STATIC ${JSONC_SOURCES})

enable_testing()

add_executable(mmc_test
    tests/test_jsonc.c
    src/jsonc/mc_token.c
    src/jsonc/mc_scanner.c
    src/jsonc/mc_ast.c
    src/jsonc/mc_parser.c
)

target_include_directories(mmc_test PRIVATE src)

add_test(NAME mmc_jsonc_test COMMAND mmc_test)
