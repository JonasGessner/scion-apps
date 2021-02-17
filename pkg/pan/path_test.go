// Copyright 2021 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pan

import (
	"fmt"
	"testing"

	"github.com/scionproto/scion/go/lib/slayers/path/scion"
	"github.com/stretchr/testify/assert"
)

func TestInterfacesFromDecoded(t *testing.T) {
	// Not a great test case...
	rawPath := []byte("\x00\x00\x20\x80\x00\x00\x01\x11\x00\x00\x01\x00\x01\x00\x02\x22\x00\x00" +
		"\x01\x00\x00\x3f\x00\x01\x00\x00\x01\x02\x03\x04\x05\x06\x00\x3f\x00\x03\x00\x02\x01\x02\x03" +
		"\x04\x05\x06\x00\x3f\x00\x00\x00\x02\x01\x02\x03\x04\x05\x06\x00\x3f\x00\x01\x00\x00\x01\x02" +
		"\x03\x04\x05\x06")

	sp := scion.Decoded{}
	sp.DecodeFromBytes(rawPath)
	ifaces := interfacesFromDecoded(sp)
	fmt.Println(ifaces)
	expected := []IfID{1, 2, 2, 1}
	assert.Equal(t, ifaces, expected)
}
