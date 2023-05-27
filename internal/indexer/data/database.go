package data

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/hanchon/garnet/internal/indexer/data/mudhelpers"
	"github.com/hanchon/garnet/internal/logger"
)

type Event struct {
	// TODO: add world here
	// World string `json:"world"`
	Table string `json:"table"`
	Row   string `json:"row"`
	Value string `json:"value"`
}

type TableMetadata struct {
	TableID          string
	TableName        string
	OnChainTableName string
	WorldAddress     string
}

type TableSchema struct {
	FieldNames  *[]string
	KeyNames    *[]string
	Schema      *mudhelpers.SchemaTypeKV
	NamedFields *map[string]mudhelpers.SchemaType
}

type Table struct {
	Metadata *TableMetadata
	Schema   *TableSchema
	Rows     *map[string][]Field
}

type World struct {
	Address string
	Tables  map[string]*Table
}

// TODO: add a cache layer here to avoid the loop
func (w *World) GetTableByName(tableName string) *Table {
	for tableID := range w.Tables {
		if w.Tables[tableID].Metadata != nil {
			if w.Tables[tableID].Metadata.TableName == tableName {
				return w.Tables[tableID]
			}
		}
	}
	return nil
}

func (w *World) GetTable(tableID string) *Table {
	if table, ok := w.Tables[tableID]; ok {
		return table
	}
	w.Tables[tableID] = &Table{
		Metadata: &TableMetadata{TableID: tableID, TableName: "", OnChainTableName: "", WorldAddress: w.Address},
		Schema:   &TableSchema{FieldNames: &[]string{}, KeyNames: &[]string{}, Schema: &mudhelpers.SchemaTypeKV{}, NamedFields: &map[string]mudhelpers.SchemaType{}},
		Rows:     &map[string][]Field{},
	}
	table := w.Tables[tableID]
	return table
}

type Database struct {
	Worlds     map[string]*World
	Events     []Event
	LastUpdate time.Time
	LastHeight uint64
	ChainID    string
}

func NewDatabase() *Database {
	return &Database{
		Worlds:     map[string]*World{},
		Events:     make([]Event, 0),
		LastUpdate: time.Now(),
		LastHeight: 0,
		ChainID:    "",
	}
}

func (db *Database) AddEvent(tableName string, key string, fields *[]Field) {
	value := ""
	if fields != nil {
		value = "{"
		for i, v := range *fields {
			value += v.String()
			if i != len(*fields)-1 {
				value += ","
			}
		}
		value += "}"
	}
	db.Events = append(db.Events, Event{Table: tableName, Row: key, Value: value})
	db.LastUpdate = time.Now()
}

func (db *Database) GetWorld(worldID string) *World {
	if world, ok := db.Worlds[worldID]; ok {
		return world
	}
	db.Worlds[worldID] = &World{Address: worldID, Tables: map[string]*Table{}}
	logger.LogInfo(fmt.Sprintf("new world registered %s", worldID))
	world := db.Worlds[worldID]
	return world
}

func (db *Database) GetTable(worldID string, tableID string) *Table {
	world := db.GetWorld(worldID)
	return world.GetTable(tableID)
}

func (db *Database) AddRow(table *Table, key []byte, fields *[]Field) {
	// Use the database to add and remove info so we can broadcast events to subs
	// keyAsString := string(key)
	keyAsString := hexutil.Encode(key)
	// TODO: add locks here
	(*table.Rows)[keyAsString] = *fields
	db.AddEvent(table.Metadata.TableName, keyAsString, fields)
}

func (db *Database) SetField(table *Table, key []byte, event *mudhelpers.StorecoreStoreSetField) {
	// TODO: add locks here

	// keyAsString := string(key)
	keyAsString := hexutil.Encode(key)
	fields, modified := BytesToFieldWithDefaults(event.Data, *table.Schema.Schema.Value, event.SchemaIndex, table.Schema.FieldNames)

	_, ok := (*table.Rows)[keyAsString]
	if ok {
		// Edit the row because it already exists
		for i := range (*table.Rows)[keyAsString] {
			if (*table.Rows)[keyAsString][i].Key == modified.Key {
				(*table.Rows)[keyAsString][i].Data = modified.Data
				break
			}
		}
	} else {
		// Create an empty row with defaults but the event index that uses event.Data
		(*table.Rows)[keyAsString] = *fields
	}

	db.AddEvent(table.Metadata.TableName, keyAsString, fields)
}

func (db *Database) DeleteRow(table *Table, key []byte) {
	// keyAsString := string(key)
	keyAsString := hexutil.Encode(key)
	// TODO: add locks here
	delete((*table.Rows), keyAsString)
	db.AddEvent(table.Metadata.TableName, keyAsString, nil)
}
