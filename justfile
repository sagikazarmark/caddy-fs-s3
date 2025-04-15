[private]
default:
  @just --list

check:
    dagger call check

test:
    dagger call test

lint:
    dagger call lint

fmt:
    golangci-lint fmt
