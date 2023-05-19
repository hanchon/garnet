package data

import (
	"fmt"
	"math/big"

	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/logger"
	"go.uber.org/zap"
)

func BytesToDynamicField(schemaType mudhelpers.SchemaType, encodingSlice []byte) FieldData {
	switch schemaType {
	case mudhelpers.BYTES:
		return NewBytesField(encodingSlice)
	case mudhelpers.STRING:
		return NewStringField(encodingSlice)
	default:
		// Try to decode as an array.
		staticSchemaType := (schemaType - 98)
		if staticSchemaType > 97 {
			logger.LogError(
				fmt.Sprintln(
					"Unknown dynamic field type",
					zap.String("type", schemaType.String()),
				),
			)
			return nil
		}

		// Allocate an array of the correct size.
		fieldLength := mudhelpers.GetStaticByteLength(staticSchemaType)
		arrayLength := len(encodingSlice) / int(fieldLength)
		array := NewArrayField(arrayLength)
		// Iterate and decode each element as a static field.
		for i := 0; i < arrayLength; i++ {
			array.Data[i] = BytesToStaticField(staticSchemaType, encodingSlice, uint64(i)*fieldLength)
		}

		return array
	}
}

func BytesToStaticField(schemaType mudhelpers.SchemaType, encoding []byte, bytesOffset uint64) FieldData {
	// UINT8 - UINT256 is the first range. We add one to the schema type to get the
	// number of bytes to read, since enums start from 0 and UINT8 is the first one.
	if schemaType >= mudhelpers.UINT8 && schemaType <= mudhelpers.UINT256 {
		return NewUintField(encoding[bytesOffset : bytesOffset+uint64(schemaType)+1])
	} else
	// INT8 - INT256 is the second range. We subtract UINT256 from the schema type
	// to account for the first range and re-set the bytes count to start from 1.
	if schemaType >= mudhelpers.INT8 && schemaType <= mudhelpers.INT256 {
		return NewIntField(encoding[bytesOffset : bytesOffset+uint64(schemaType-mudhelpers.UINT256)])
	} else
	// BYTES is the third range. We subtract INT256 from the schema type to account
	// for the previous ranges and re-set the bytes count to start from 1.
	if schemaType >= mudhelpers.BYTES1 && schemaType <= mudhelpers.BYTES32 {
		return NewBytesField(encoding[bytesOffset : bytesOffset+uint64(schemaType-mudhelpers.INT256)])
	} else
	// BOOL is a standalone schema type.
	if schemaType == mudhelpers.BOOL {
		return NewBoolField(encoding[bytesOffset])
	} else
	// ADDRESS is a standalone schema type.
	if schemaType == mudhelpers.ADDRESS {
		return NewAddressField(encoding[bytesOffset : bytesOffset+20])
	} else
	// STRING is a standalone schema type.
	if schemaType == mudhelpers.STRING {
		return NewStringField(encoding[bytesOffset:])
	} else {
		logger.LogError(
			fmt.Sprintln(
				"Unknown static field type",
				zap.String("type", schemaType.String()),
			),
		)
		return nil
	}
}

func AggregateKey(key [][32]byte) []byte {
	aggregateKey := []byte{}
	for _, keyElement := range key {
		aggregateKey = append(aggregateKey, keyElement[:]...)
	}
	return aggregateKey
}

func BytesToFields(encoding []byte, schemaTypePair mudhelpers.SchemaTypePair, fieldnames *[]string) *[]Field {
	var bytesOffset uint64 = 0
	ret := []Field{}

	// Decode static fields.
	for _, fieldType := range schemaTypePair.Static {
		value := BytesToStaticField(fieldType, encoding, bytesOffset)
		bytesOffset += mudhelpers.GetStaticByteLength(fieldType)
		ret = append(ret, Field{Key: "", Data: value})
	}

	// Decode dynamic fields.
	if len(schemaTypePair.Dynamic) > 0 {
		dynamicDataSlice := encoding[schemaTypePair.StaticDataLength : schemaTypePair.StaticDataLength+32]
		bytesOffset += 32

		for i, fieldType := range schemaTypePair.Dynamic {
			offset := 4 + i*2
			dataLength := new(big.Int).SetBytes(dynamicDataSlice[offset : offset+2]).Uint64()
			value := BytesToDynamicField(fieldType, encoding[bytesOffset:bytesOffset+dataLength])
			bytesOffset += dataLength
			ret = append(ret, Field{Key: "", Data: value})
		}
	}

	// Add the fields.
	for idx, field_name := range *fieldnames {
		ret[idx].Key = field_name
	}

	return &ret
}

func BytesToFieldWithDefaults(encoding []byte, schemaTypePair mudhelpers.SchemaTypePair, index uint8, fieldnames *[]string) (*[]Field, Field) {
	ret := []Field{}
	modified := Field{}

	// Try to decode either as a static or dynamic field.
	for idx, fieldType := range schemaTypePair.Static {
		if uint8(idx) == index {
			value := BytesToStaticField(fieldType, encoding, 0)
			modified = Field{Key: "", Data: value}
			ret = append(ret, modified)
		} else {
			// Create empty fieldType
			ret = append(ret, Field{Key: "", Data: FieldWithDefautValue(fieldType)})
		}
	}
	for idx, fieldType := range schemaTypePair.Dynamic {
		// Offset by the static data length.
		if uint8(idx+len(schemaTypePair.Static)) == index {
			value := BytesToDynamicField(fieldType, encoding)
			modified = Field{Key: "", Data: value}
			ret = append(ret, modified)
		} else {
			// Create empty with fieldType
			ret = append(ret, Field{Key: "", Data: FieldWithDefautValue(fieldType)})
		}
	}

	// Add the fields.
	for idx, field_name := range *fieldnames {
		ret[idx].Key = field_name
	}
	modified.Key = (*fieldnames)[index]

	return &ret, modified

}
