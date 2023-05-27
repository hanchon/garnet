package mudhandlers

import (
	"fmt"

	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/logger"
	"go.uber.org/zap"
)

func HandleGenericTableEvent(event *mudhelpers.StorecoreStoreSetRecord, db *data.Database) {
	tableID := mudhelpers.PaddedTableId(event.TableId)
	logger.LogDebug(
		fmt.Sprintln(
			"handling generic table event",
			zap.String("world_address", event.WorldAddress()),
			zap.String("table_id", tableID),
		),
	)

	table := db.GetTable(event.WorldAddress(), tableID)

	// Decode the row record data
	fields := data.BytesToFields(event.Data, *table.Schema.Schema.Value, table.Schema.FieldNames)

	// Decode the row key
	aggregateKey := data.AggregateKey(event.Key)

	// Save it
	db.AddRow(table, aggregateKey, fields)

	// TODO: do we need the info of each key or is it always going to match the complete expresion
	// decodedKeyData := mudhelpers.DecodeData(aggregateKey, *table.Schema.Schema.Key)
	// decodedKeyDataNew := data.BytesToFields(aggregateKey, *table.Schema.Schema.Key, table.Schema.KeyNames)
}
