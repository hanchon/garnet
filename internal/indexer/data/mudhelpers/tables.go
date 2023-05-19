package mudhelpers

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type TableSchema struct {
	TableId    string   `json:"id"`   // Table ID as it comes from chain.
	TableName  string   `json:"name"` // Table name is table ID but with naming adjustments to work with the database.
	FieldNames []string `json:"field_names"`
	KeyNames   []string `json:"key_names"` // Key names are separte from field names and are used for searching.

	SolidityTypes map[string]string `json:"solidity_types"` // Field name -> Solidity type
	PostgresTypes map[string]string `json:"postgres_types"` // Field name -> Postgres type

	IsKey map[string]bool `json:"is_key"` // Field name -> Is key?

	// Auxiliary data about the table.
	Namespace             string            `json:"namespace"`
	StoreCoreSchemaTypeKV *SchemaTypeKV     `json:"store_core_schema_type_kv"`
	PrimaryKey            string            `json:"primary_key"`
	OnChainReadableName   string            `json:"on_chain_readable_name"`
	OnChainColNames       map[string]string `json:"on_chain_col_names"`
}

const CONNECTOR string = "__"
const TABLE_PREFIX string = "mode"

func DefaultFieldName(index int) string {
	return "field_" + fmt.Sprint(index)
}

func DefaultKeyName(index int) string {
	return "key_" + fmt.Sprint(index)
}

func Namespace(chainId string, worldAddress string) string {
	var str strings.Builder
	str.WriteString(TABLE_PREFIX)
	if chainId != "" {
		str.WriteString(CONNECTOR + chainId)
	}
	if worldAddress != "" {
		str.WriteString(CONNECTOR + strings.ToLower(worldAddress))
	}
	return str.String()
}

func TableIdToTableName(tableId string) string {
	// Table ID comes in as a uint256 in string format comprised of two bytes16s
	// concatenated.
	// Get a byte array of the table ID.
	tableIdBytes := []byte(tableId[2:])

	b1 := tableIdBytes[:32]
	b2 := tableIdBytes[32:]
	namePart1 := string([]byte(b1))
	namePart2 := string([]byte(b2))

	p1 := string(common.FromHex("0x" + namePart1))
	p2 := string(common.FromHex("0x" + namePart2))

	return strings.Trim(p1, "\u0000") + CONNECTOR + strings.Trim(p2, "\u0000")
}

func TableNameToTableId(tableName string) string {
	// Assumes that table name was generated by TableIdToTableName.
	parts := strings.Split(tableName, CONNECTOR)
	if len(parts) != 2 {
		panic("Invalid table name: " + tableName)
	}
	return "0x" + common.Bytes2Hex(append(common.FromHex("0x"+parts[0]), common.FromHex("0x"+parts[1])...))
}
