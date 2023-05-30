@_default:
    just --list

set dotenv-load := true

# build:
#     {{ go }} build

# test:
#     {{ app }} build test

run:
    go run ./cmd