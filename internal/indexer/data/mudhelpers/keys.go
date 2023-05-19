package mudhelpers

import "github.com/ethereum/go-ethereum/common"

func SchemaTableId() string {
	return "0x" + common.Bytes2Hex(append(RightPadId("mudstore"), RightPadId("schema")...))
}

func MetadataTableId() string {
	return "0x" + common.Bytes2Hex(append(RightPadId("mudstore"), RightPadId("StoreMetadata")...))
}

func SchemaTableName() string {
	return "mudstore__schema"
}

func MetadataTableName() string {
	return "mudstore__storemetadata"
}

func RightPadId(id string) []byte {
	return common.RightPadBytes([]byte(id), 16)
}

func PaddedTableId(id [32]byte) string {
	return "0x" + common.Bytes2Hex(id[:])
}

func RowFromDecodedData(decodedKeyData *DecodedData, decodedFieldData *DecodedData, keynames *[]string, fieldnames *[]string) map[string]interface{} {
	// Create a row for the table.
	row := map[string]interface{}{}
	// Add the keys.
	for idx, key_name := range *keynames {
		// Skip the key if it's not in the row.
		row[key_name] = decodedKeyData.DataAt(idx)
	}
	// Add the fields.
	for idx, field_name := range *fieldnames {
		row[field_name] = decodedFieldData.DataAt(idx)
	}

	return row
}
