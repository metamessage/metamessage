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

set(IR_SOURCES
    src/ir/mc_ast.c
    src/ir/mc_tag.c
)

set(CORE_SOURCES
    src/core/mc_encoder.c
    src/core/mc_decoder.c
)

set(MM_SOURCES
    src/jsonc/mc_printer.c
    src/jsonc/mc_ir_parser.c
    src/mm.c
)

set(ALL_SOURCES
    ${JSONC_SOURCES}
    ${IR_SOURCES}
    ${CORE_SOURCES}
    ${MM_SOURCES}
)

add_library(mmc_jsonc STATIC ${JSONC_SOURCES})

enable_testing()

add_executable(mmc_test
    tests/test_jsonc.c
    ${JSONC_SOURCES}
)

add_executable(mmc_comprehensive_test
    tests/test_comprehensive.c
    ${ALL_SOURCES}
)

target_include_directories(mmc_test PRIVATE src)
target_include_directories(mmc_comprehensive_test PRIVATE src)

add_test(NAME mmc_jsonc_test COMMAND mmc_test)
add_test(NAME mmc_comprehensive_test COMMAND mmc_comprehensive_test)