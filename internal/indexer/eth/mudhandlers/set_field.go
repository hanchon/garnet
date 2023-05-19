package mudhandlers

import (
	"fmt"

	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/logger"
	"go.uber.org/zap"
)

func HandleSetFieldEvent(event *mudhelpers.StorecoreStoreSetField, db *data.Database) {
	tableId := mudhelpers.PaddedTableId(event.TableId)
	logger.LogDebug(
		fmt.Sprintln(
			"handling set field (StoreSetFieldEvent) event",
			zap.String("table_id", tableId),
		),
	)

	table := db.GetTable(event.WorldAddress(), tableId)

	// Handle the following scenarios:
	// 1. The setField event is modifying a row that doesn't yet exist (i.e. key doesn't match anything),
	//    in which case we insert a new row with default values for each column.
	//
	// 2. The setField event is modifying a row that already exists, in which case we update the
	//    row by constructing a partial row with the new value for the field that was modified.

	key := data.AggregateKey(event.Key)

	db.SetField(table, key, event)

}
