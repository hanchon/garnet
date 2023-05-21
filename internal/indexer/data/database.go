package data

import (
	"fmt"
	"time"

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
	TableId          string
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

func (w *World) GetTable(tableId string) *Table {
	if table, ok := w.Tables[tableId]; ok {
		return table
	}
	w.Tables[tableId] = &Table{
		Metadata: &TableMetadata{TableId: tableId, TableName: "", OnChainTableName: "", WorldAddress: w.Address},
		Schema:   &TableSchema{FieldNames: &[]string{}, KeyNames: &[]string{}, Schema: &mudhelpers.SchemaTypeKV{}, NamedFields: &map[string]mudhelpers.SchemaType{}},
		Rows:     &map[string][]Field{},
	}
	table, _ := w.Tables[tableId]
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
			value = value + v.String()
			if i != len(*fields)-1 {
				value = value + ","
			}
		}
		value = value + "}"
	}
	db.Events = append(db.Events, Event{Table: tableName, Row: key, Value: value})
	db.LastUpdate = time.Now()
}

func (db *Database) GetWorld(worldId string) *World {
	if world, ok := db.Worlds[worldId]; ok {
		return world
	}
	db.Worlds[worldId] = &World{Address: worldId, Tables: map[string]*Table{}}
	logger.LogInfo(fmt.Sprintf("new world registered %s", worldId))
	world, _ := db.Worlds[worldId]
	return world
}

func (db *Database) GetTable(worldId string, tableId string) *Table {
	world := db.GetWorld(worldId)
	return world.GetTable(tableId)
}

func (db *Database) AddRow(table *Table, key []byte, fields *[]Field) {
	// Use the database to add and remove info so we can broadcast events to subs
	keyAsString := string(key)
	// TODO: add locks here
	(*table.Rows)[keyAsString] = *fields
	db.AddEvent(table.Metadata.TableName, keyAsString, fields)
	db.LastUpdate = time.Now()
}

func (db *Database) SetField(table *Table, key []byte, event *mudhelpers.StorecoreStoreSetField) {
	// TODO: add locks here

	keyAsString := string(key)
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
	db.LastUpdate = time.Now()
}

func (db *Database) DeleteRow(table *Table, key []byte) {
	keyAsString := string(key)
	// TODO: add locks here
	if _, ok := (*table.Rows)[keyAsString]; ok {
		delete((*table.Rows), keyAsString)
	}
	db.AddEvent(table.Metadata.TableName, keyAsString, nil)
	db.LastUpdate = time.Now()
}
