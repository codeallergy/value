/*
 * Copyright (c) 2022-2023 Zander Schwid & Co. LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */


package value

/**
Base interface for the packing values
*/

type Ext byte

const (
	UnknownExt Ext = iota

	BigIntExt
	DecimalExt

	MaxExt
)

type Packer interface {
	PackNil()

	PackBool(bool)

	PackLong(int64)

	PackDouble(float64)

	PackStr(string)

	PackBin([]byte)

	PackExt(xtag Ext, data []byte)

	PackList(int)

	PackMap(int)

	PackRaw([]byte)

	Error() error
}

/**

Interface for write values

*/

type Writer interface {
	WriteNil() []byte

	WriteBool(val bool) []byte

	WriteLong(val int64) []byte

	WriteDouble(val float64) []byte

	WriteBinHeader(len int) []byte

	WriteStrHeader(len int) []byte

	WriteExtHeader(len int, xtag byte) []byte

	WriteArrayHeader(len int) []byte

	WriteMapHeader(len int) []byte
}

/**
Base interface for the unpacking values

*/

type Format int

const (
	EOF Format = iota
	UnexpectedEOF
	NilToken
	BoolToken
	LongToken
	DoubleToken
	FixExtToken
	BinHeader
	StrHeader
	ListHeader
	MapHeader
	ExtHeader
)

type Unpacker interface {
	Next() (Format, []byte)

	Read(int) ([]byte, error)
}

/**
	Parse value from Slice, the slice size must be enough to parse primitive value or header, slice always stars from code

    return - number of bytes read, value

	return 0 read bytes on error
*/

type Parser interface {
	ParseBool([]byte) bool

	ParseLong([]byte) int64

	ParseDouble([]byte) float64

	ParseBin([]byte) int

	ParseStr([]byte) int

	ParseList([]byte) int

	ParseMap([]byte) int

	ParseExt([]byte) (int, []byte)

	Error() error
}
