// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json_test

import (
	"fmt"
	"json"
	"testing"
)

func TestKeep(t *testing.T) {
	eks := TestKeepStruct{}
	json.StructKeepType = json.KeepEmptyNumber | json.KeepEmptyArray
	bytes, _ := json.Marshal(&eks)
	fmt.Println(string(bytes))
}

type TestKeepStruct struct {
	Bo  bool
	It  uint
	St  string
	Ar  []string
	St2 string `json:"Data,keepEmpty"`
}
