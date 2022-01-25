#!/usr/bin/env bash

set -eo pipefail

protoc_gen_gocosmos() {
  if ! grep "github.com/gogo/protobuf => github.com/regen-network/protobuf" go.mod &>/dev/null; then
    echo -e "\tPlease run this command from somewhere inside the cosmos-sdk folder."
    return 1
  fi

  go get github.com/regen-network/cosmos-proto/protoc-gen-gocosmos@latest 2>/dev/null
}

protoc_gen_gocosmos

if [ -f ./build/cosmos-sdk/README.md ]; then
  if [ ! -f ./build/upgrade.zip ]; then
    wget -c https://github.com/functionx/cosmos-sdk/archive/refs/heads/feature/upgrade.zip -O ./build/upgrade.zip
  fi
  (cd build && unzip -o upgrade.zip && mv cosmos-sdk-feature-upgrade cosmos-sdk)
fi

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  buf protoc \
    -I "proto" \
    -I "build/cosmos-sdk/proto" \
    -I "build/cosmos-sdk/third_party/proto" \
    --gocosmos_out=plugins=interfacetype+grpc, Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
    --grpc-gateway_out=logtostderr=true:. \
    $(find "${dir}" -maxdepth 1 -name '*.proto')
done

# command to generate docs using protoc-gen-doc
buf protoc \
  -I "proto" \
  -I "build/cosmos-sdk/proto" \
  -I "build/cosmos-sdk/third_party/proto" \
  --doc_out=./docs/core \
  --doc_opt=./docs/protodoc-markdown.tmpl,proto-docs.md \
  $(find "$(pwd)/proto" -maxdepth 5 -name '*.proto')

go mod tidy
