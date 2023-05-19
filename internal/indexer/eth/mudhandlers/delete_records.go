package mudhandlers

import (
	"fmt"

	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/logger"
	"go.uber.org/zap"
)

func HandleDeleteRecordEvent(event *mudhelpers.StorecoreStoreDeleteRecord, db *data.Database) {
	tableId := mudhelpers.PaddedTableId(event.TableId)
	logger.LogDebug(
		fmt.Sprintln(
			"handling delete record event",
			zap.String("table_id", tableId),
		),
	)

	table := db.GetTable(event.WorldAddress(), tableId)

	aggregateKey := data.AggregateKey(event.Key)

	db.DeleteRow(table, aggregateKey)
}
