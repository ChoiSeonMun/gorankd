#!/usr/bin/env bash
set -euo pipefail

PROTO_DIR="api/proto"
GEN_DIR="api/proto/gen"

mkdir -p "$GEN_DIR"

protoc \
  --proto_path="$PROTO_DIR" \
  --go_out="$GEN_DIR" \
  --go_opt=paths=source_relative \
  --go-grpc_out="$GEN_DIR" \
  --go-grpc_opt=paths=source_relative \
  "$PROTO_DIR/ranking.proto"

echo "Proto generation complete."
