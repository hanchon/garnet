package mudhandlers

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/logger"
	"go.uber.org/zap"
)

func HandleSchemaTableEvent(event *mudhelpers.StorecoreStoreSetRecord, database *data.Database) {
	tableID := hexutil.Encode(event.Key[0][:])
	logger.LogDebug(
		fmt.Sprintln(
			"handling schema table event",
			zap.String("world_address", event.WorldAddress()),
			zap.String("table_id", tableID),
		),
	)
	world := database.GetWorld(event.WorldAddress())
	table := world.GetTable(tableID)

	// Parse out the schema types (both static and dynamic) for the table.
	keySchemaBytes32, valueSchemaBytes32 := event.Data[:32], event.Data[32:]
	valueStoreCoreSchemaTypePair := mudhelpers.DecodeSchemaTypePair(keySchemaBytes32)
	// The last 32 bytes are the table "key" schema.
	keyStoreCoreSchemaTypePair := mudhelpers.DecodeSchemaTypePair(valueSchemaBytes32)

	// Merge the two schemas into one, since the table schema is a combination of the key schema and the value schema.
	storeCoreSchemaTypeKV := mudhelpers.SchemaTypeKVFromPairs(keyStoreCoreSchemaTypePair, valueStoreCoreSchemaTypePair)
	table.Schema.Schema = mudhelpers.SchemaTypeKVFromPairs(keyStoreCoreSchemaTypePair, valueStoreCoreSchemaTypePair)

	fieldNames := *table.Schema.FieldNames
	keyNames := *table.Schema.KeyNames
	for idx := range storeCoreSchemaTypeKV.Value.Flatten() {
		columnName := mudhelpers.DefaultFieldName(idx)
		fieldNames = append(fieldNames, columnName)
	}

	for idx := range storeCoreSchemaTypeKV.Key.Flatten() {
		columnName := mudhelpers.DefaultKeyName(idx)
		keyNames = append(keyNames, columnName)
	}

	table.Schema.FieldNames = &fieldNames
	table.Schema.KeyNames = &keyNames

	// NOTE: we are overwritting every time this is called, but it is only called once so it is not a problem
	table.Metadata.TableName = "schema"
}
