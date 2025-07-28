#!/bin/bash
# Generate Go code from all proto files
PROTO_DIR=./proto
OUT_DIR=./proto/generated
mkdir -p $OUT_DIR
protoc -I=$PROTO_DIR --go_out=$OUT_DIR --go-grpc_out=$OUT_DIR $PROTO_DIR/user/*.proto $PROTO_DIR/product/*.proto $PROTO_DIR/order/*.proto 