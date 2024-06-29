#!/bin/bash

CC="$(pwd)/vector-gobot/vic-toolchain/arm-linux-gnueabi/bin/arm-linux-gnueabi-gcc" \
CGO_LDFLAGS="-L$(pwd)/vector-gobot/build -L$(pwd)/vector-gobot/build/libjpeg-turbo/lib" \
GOARM=7 \
GOARCH=arm \
CGO_ENABLED=1 \
go build -o caprice $@
