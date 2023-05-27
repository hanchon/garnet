package mudhandlers

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hanchon/garnet/internal/indexer/data"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/logger"
	"github.com/umbracle/ethgo/abi"
	"go.uber.org/zap"
)

func HandleMetadataTableEvent(event *mudhelpers.StorecoreStoreSetRecord, db *data.Database) {
	tableID := hexutil.Encode(event.Key[0][:])
	logger.LogDebug(
		fmt.Sprintln(
			"handling metadata table event",
			zap.String("world_address", event.WorldAddress()),
			zap.String("table_id", tableID),
		),
	)
	world := db.GetWorld(event.WorldAddress())
	table := world.GetTable(tableID)
	metadata := world.GetTable(mudhelpers.MetadataTableId())

	decodedMetadata := mudhelpers.DecodeData(event.Data, *metadata.Schema.Schema.Value)

	// Since we know the structure of the metadata, we decode it directly into types and handle.
	tableReadableName := decodedMetadata.DataAt(0).(string)
	table.Metadata.TableName = tableReadableName

	tableColumnNamesHexString := decodedMetadata.DataAt(1).(string)
	tableColumnNamesBytes, err := hexutil.Decode(tableColumnNamesHexString)
	if err != nil {
		logger.LogError(fmt.Sprintf("error decoding hex value: %s", err))
		return
	}

	// For some reason just string[] doesn't work with abi decoding here, so we use a tuple.
	_type := abi.MustNewType("tuple(string[] cols)")
	outStruct := struct {
		Cols []string
	}{
		Cols: []string{},
	}
	err = _type.DecodeStruct(tableColumnNamesBytes, &outStruct)
	if err != nil {
		logger.LogDebug(fmt.Sprintln("failed to decode table column names", zap.Error(err)))
		return
	}

	newTableFieldNames := []string{}
	for idx, schemaType := range table.Schema.Schema.Value.Flatten() {
		columnName := strings.ToLower(outStruct.Cols[idx])
		newTableFieldNames = append(newTableFieldNames, columnName)
		(*table.Schema.NamedFields)[columnName] = schemaType

		// TODO: just itereate once after getting all the new names
		for i := range *table.Rows {
			for j := range (*table.Rows)[i] {
				if (*table.Rows)[i][j].Key == (*table.Schema.FieldNames)[idx] {
					(*table.Rows)[i][j].Key = columnName
					break
				}
			}
		}
	}
	table.Schema.FieldNames = &newTableFieldNames

	// Save it as a row in the metadata table
	fields := data.BytesToFields(event.Data, *metadata.Schema.Schema.Value, metadata.Schema.FieldNames)
	// key := tableId
	// (*metadata.Rows)[key] = *fields
	db.AddRow(metadata, []byte(tableID), fields)
}
