#!/bin/bash
# Generate Go code from all proto files
PROTO_DIR=./proto
OUT_DIR=./proto/generated
mkdir -p $OUT_DIR

# Generate from user service proto
protoc -I=$PROTO_DIR --go_out=$OUT_DIR --go-grpc_out=$OUT_DIR ./proto/user/user.proto

# Generate from product service proto
protoc -I=$PROTO_DIR --go_out=$OUT_DIR --go-grpc_out=$OUT_DIR ./proto/product/product.proto

# Generate from order service proto
protoc -I=$PROTO_DIR --go_out=$OUT_DIR --go-grpc_out=$OUT_DIR ./proto/order/order.proto

# Generate from other service protos (when ready)
# protoc -I=$PROTO_DIR --go_out=$OUT_DIR --go-grpc_out=$OUT_DIR $PROTO_DIR/payment/*.proto

echo "Proto files generated successfully!" 