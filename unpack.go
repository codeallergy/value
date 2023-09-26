/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */


package value

import (
	"github.com/pkg/errors"
	"io"
	"strconv"
)


func doParse(unpacker Unpacker, parser Parser) (Value, error) {

	format, header := unpacker.Next()

	switch format {
	case EOF:
		return nil, io.EOF
	case UnexpectedEOF:
		return nil, io.ErrUnexpectedEOF
	case NilToken:
		return Null, nil
	case BoolToken:
		return Boolean(parser.ParseBool(header)), parser.Error()
	case LongToken:
		return Long(parser.ParseLong(header)), parser.Error()
	case DoubleToken:
		return Double(parser.ParseDouble(header)), parser.Error()
	case FixExtToken:
		_, tagAndData := parser.ParseExt(header)
		return doParseExt(tagAndData)
	case BinHeader:
		size := parser.ParseBin(header)
		if parser.Error() != nil {
			return nil, parser.Error()
		}
		raw, err := unpacker.Read(size)
		if err != nil {
			return nil, err
		}
		return Raw(raw, false), nil
	case StrHeader:
		len := parser.ParseStr(header)
		if parser.Error() != nil {
			return nil, parser.Error()
		}
		str, err := unpacker.Read(len)
		if err != nil {
			return nil, err
		}
		return Utf8(string(str)), nil
	case ListHeader:
		return doParseList(header, unpacker, parser)
	case MapHeader:
		return doParseMap(header, unpacker, parser)
	case ExtHeader:
		n, _ := parser.ParseExt(header)
		if parser.Error() != nil {
			return nil, parser.Error()
		}
		tagAndData, err := unpacker.Read(n+1)
		if err != nil {
			return nil, err
		}
		return doParseExt(tagAndData)
	default:
		return nil, errors.Errorf("parse: invalid format %v", format)
	}

}

func doParseList(header []byte, unpacker Unpacker, parser Parser) (List, error) {
	cnt := parser.ParseList(header)
	if parser.Error() != nil {
		return nil, parser.Error()
	}
	if cnt == 0 {
		return EmptyImmutableList(), nil
	}
	list := make([]Value, cnt)
	for i := 0; i < cnt; i++ {
		el, err := doParse(unpacker, parser)
		if err != nil {
			return nil, err
		}
		list[i] = el
	}
	return ImmutableList(list), nil
}

func doParseMap(header []byte, unpacker Unpacker, parser Parser) (Value, error) {
	cnt := parser.ParseMap(header)
	if parser.Error() != nil {
		return nil, parser.Error()
	}
	if cnt == 0 {
		return EmptyImmutableMap(), nil
	}
	var sparseListItems []ListItem
	mayBeList := false
	var sortedMapEntries []MapEntry
	sorted := true
	var prevListKey int64
	var prevMapKey string

	for i := 0; i < cnt; i++ {
		key, err := doParse(unpacker, parser)
		if err != nil {
			return nil, err
		}
		value, err := doParse(unpacker, parser)
		if err != nil {
			return nil, err
		}

		if key == nil {
			// nothing to do with this
			continue
		}

		// first element
		if i == 0 {
			if key.Kind() == NUMBER {
				// try to build sparse list
				mayBeList = true
				sparseListItems = make([]ListItem, cnt)
			} else {
				mayBeList = false
				sortedMapEntries = make([]MapEntry, cnt)
			}
		}

		if mayBeList {

			if key.Kind() == NUMBER {
				k := key.(Number).Long()
				if i > 0 && prevListKey > k {
					sorted = false
				}
				sparseListItems[i] = ImmutableItem(int(k), value)
				prevListKey = k
			} else {
				// not a list
				mayBeList = false
				sortedMapEntries = make([]MapEntry, cnt)
				for j := 0; j < i; j++ {
					item := sparseListItems[i]
					sortedMapEntries[i] = ImmutableEntry(strconv.Itoa(item.Key()), item.Value())
				}
				k := key.String()
				if i > 0 && prevMapKey > k {
					sorted = false
				}
				sortedMapEntries[i] = ImmutableEntry(k, value)
				prevMapKey = k
			}

		} else {
			k := key.String()
			if i > 0 && prevMapKey > k {
				sorted = false
			}
			sortedMapEntries[i] = ImmutableEntry(k, value)
			prevMapKey = k
		}

	}

	if mayBeList {
		return SparseList(sparseListItems, sorted), nil
	} else {
		return ImmutableMap(sortedMapEntries, sorted), nil
	}

}

func doParseExt(tagAndData []byte) (Value, error) {
	xtag := Ext(tagAndData[0])
	ext := tagAndData[1:]
	switch xtag {

	case BigIntExt:
		v, err := UnpackBigInt(ext)
		return BigInt(v), err
	case DecimalExt:
		v, err := UnpackDecimal(ext)
		return Decimal(v), err

	}
	return Unknown(tagAndData), nil
}
