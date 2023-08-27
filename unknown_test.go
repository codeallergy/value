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
)

func TestUnknown(t *testing.T) {

	tagAndData := []byte { byte(val.MaxExt),  1 }

	v := val.Unknown(tagAndData)

	require.Equal(t, val.UNKNOWN, v.Kind())
	require.Equal(t, val.UnknownPrefix+ val.Base64Prefix + "AwE", v.String())
	require.Equal(t, "\"" + v.String() + "\"", val.Jsonify(v))
	require.Equal(t, "d40301", val.Hex(v))
	require.Equal(t, 0, bytes.Compare(tagAndData, v.Native()))

	mp, err := val.Pack(v)
	require.Nil(t, err)

	a, err := val.Unpack(mp, false)
	require.Nil(t, err)

	require.Equal(t, val.UNKNOWN, v.Kind())
	require.Equal(t, 0, bytes.Compare(tagAndData, a.(val.Extension).Native()))

}