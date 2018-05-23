#!/bin/bash
# Copyright 2017 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

set -e

go build -o doc-helper
./doc-helper -help >doc.go
gofmt -w doc.go
rm doc-helper
