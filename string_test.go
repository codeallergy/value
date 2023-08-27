/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */


package value_test


import (
	"testing"
	"github.com/stretchr/testify/require"
	val "github.com/codeallergy/value"
	"bytes"
	"strconv"
	"encoding/json"
)

var testStrings = map[string]string {

	"": "a0",
	"test": "a474657374",
	"hello": "a568656c6c6f",
}

func TestUtf8String(t *testing.T) {

	for str, hex := range testStrings {

		b := val.Utf8(str)

		require.Equal(t, val.STRING, b.Kind())
		require.Equal(t, val.UTF8, b.Type())
		require.Equal(t, "value.uft8String", b.Class().String())
		require.Equal(t, hex, val.Hex(b))
		require.Equal(t, "\""+ str + "\"", val.Jsonify(b))
		require.Equal(t, str, b.String())

		testPackUnpack(t, b)

	}

}

func TestJsonString(t *testing.T) {

	src := "json\"val\"json"

	s := val.Utf8(src)

	require.Equal(t, src, s.String())
	require.Equal(t, "\"json\\\"val\\\"json\"", val.Jsonify(s))

	actual, _ := strconv.Unquote(val.Jsonify(s))
	require.Equal(t, src, actual)

}

func TestRawString(t *testing.T) {

	raw := []byte { 0, 1, 2, 3, 4, 5 }
	s := val.Raw(raw, false)

	require.Equal(t, val.STRING, s.Kind())
	require.Equal(t, val.RAW, s.Type())
	require.Equal(t, "value.rawString", s.Class().String())
	require.Equal(t, val.Base64Prefix + "AAECAwQF", s.String())
	require.Equal(t, "\"" + val.Base64Prefix + "AAECAwQF\"", val.Jsonify(s))
	require.Equal(t, "c406000102030405", val.Hex(s))
	require.Equal(t, 0, bytes.Compare(raw, s.Raw()))

	actual := val.ParseString(s.String())
	require.Equal(t, s.Raw(), actual.Raw())

	testPackUnpack(t, s)

}

type testStringStruct struct {
	S val.String
}

func TestStringMarshal(t *testing.T) {

	b := val.Utf8("a")

	j, _ := b.MarshalJSON()
	require.Equal(t, "\"a\"", string(j))

	bin, _ := b.MarshalBinary()
	require.Equal(t, []byte{0xa1, 0x61}, bin)

	b = val.Raw([]byte{0, 1}, false)

	j, _ = b.MarshalJSON()
	require.Equal(t, "\"base64,AAE\"", string(j))

	bin, _ = b.MarshalBinary()
	require.Equal(t, []byte{0xc4, 0x2, 0x0, 0x1}, bin)

	s := &testStringStruct{val.Utf8("b")}

	j, _ = json.Marshal(s)
	require.Equal(t, "{\"S\":\"b\"}", string(j))

}