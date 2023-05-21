package data

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/logger"
)

type FieldData interface {
	String() string
}

type Field struct {
	Key  string
	Data FieldData
}

func (f Field) String() string {
	return fmt.Sprintf("\"%s\":%s", f.Key, f.Data.String())
}

type BytesField struct {
	Data []byte
}

func NewBytesField(data []byte) BytesField {
	return BytesField{Data: data}
}

func (f BytesField) String() string {
	return fmt.Sprintf("\"%s\"", hexutil.Encode(f.Data))
}

type StringField struct {
	Data string
}

func NewStringField(data []byte) StringField {
	return StringField{Data: string(data)}
}
func (f StringField) String() string {
	return fmt.Sprintf("\"%s\"", f.Data)
}

type ArrayField struct {
	Data []FieldData
}

func NewArrayField(size int) ArrayField {
	return ArrayField{
		Data: make([]FieldData, size),
	}
}
func (f ArrayField) String() string {
	// TODO: improve the string builder performance
	ret := "["
	for k, v := range f.Data {
		ret = ret + v.String()
		if k != len(f.Data)-1 {
			ret = ret + ","
		}
	}
	ret = ret + "]"
	return ret
}

type UintField struct {
	Data big.Int
}

func NewUintField(data []byte) UintField {
	return UintField{
		Data: *new(big.Int).SetBytes(data),
	}
}
func (f UintField) String() string {
	return f.Data.String()
}

type IntField struct {
	Data big.Int
}

func NewIntField(data []byte) IntField {
	return IntField{Data: *new(big.Int).SetBytes(data)}
}
func (f IntField) String() string {
	return f.Data.String()
}

type BoolField struct {
	Data bool
}

func NewBoolField(encoding byte) BoolField {
	return BoolField{Data: encoding == 1}
}
func (f BoolField) String() string {
	if f.Data == true {
		return "true"
	}
	return "false"
}

type AddressField struct {
	Data common.Address
}

func NewAddressField(encoding []byte) AddressField {
	return AddressField{Data: common.BytesToAddress(encoding)}
}
func (f AddressField) String() string {
	return fmt.Sprintf("\"%s\"", f.Data.Hex())
}

func FieldWithDefautValue(schemaType mudhelpers.SchemaType) FieldData {
	// UINT8 - UINT256 is the first range. We add one to the schema type to get the
	// number of bytes to read, since enums start from 0 and UINT8 is the first one.
	if schemaType >= mudhelpers.UINT8 && schemaType <= mudhelpers.UINT256 {
		return NewUintField([]byte{0})
	} else
	// INT8 - INT256 is the second range. We subtract UINT256 from the schema type
	// to account for the first range and re-set the bytes count to start from 1.
	if schemaType >= mudhelpers.INT8 && schemaType <= mudhelpers.INT256 {
		return NewIntField([]byte{0})
	} else
	// BYTES is the third range. We subtract INT256 from the schema type to account
	// for the previous ranges and re-set the bytes count to start from 1.
	if schemaType >= mudhelpers.BYTES1 && schemaType <= mudhelpers.BYTES32 {
		return NewBytesField([]byte{0})
	} else
	// BOOL is a standalone schema type.
	if schemaType == mudhelpers.BOOL {
		return NewBoolField(0)
	} else
	// ADDRESS is a standalone schema type.
	if schemaType == mudhelpers.ADDRESS {
		return NewAddressField([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	} else
	// STRING is a standalone schema type.
	if schemaType == mudhelpers.STRING {
		return NewStringField([]byte{})
	} else
	// BYTES
	if schemaType == mudhelpers.BYTES {
		return NewBytesField([]byte{})
	} else
	// ARRAYs
	if schemaType >= mudhelpers.UINT8_ARRAY && schemaType <= mudhelpers.ADDRESS_ARRAY {
		return NewArrayField(0)
	} else {
		logger.LogError(fmt.Sprintf("Unknown static field type %s", schemaType.String()))
		return nil
	}

}
