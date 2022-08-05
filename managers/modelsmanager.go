package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
)

var Chunks []*chunkModel

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func isArray(s string) bool {
	r, _ := regexp.MatchString(`(\d|\[(\d|,\s*)*])`, s)
	return r
}

func defineModelsType(value any) string {
	valueType := reflect.TypeOf(value).Name()
	if valueType == "string" {
		valToString := fmt.Sprintf("%v", value)
		if isJSON(valToString) {
			valueType = "json"
		} else if isArray(valToString) {
			valueType = "array"
		}
	}
	return valueType
}

// Returns empty dataModel if couldn't find data with provided name
func FindData(name string) (*dataModel, bool) {
	for _, chunk := range Chunks {
		for _, d := range chunk.data {
			if d.Name == name {
				return d, true
			}
		}
	}
	return &dataModel{}, false
}

// Creates data struct and appends it to chunk. If no chunk exists then creates one before
func CreateDataModel(name string, value any) (dataModel, bool) {
	d, exists := FindData(name)
	if exists {
		return *d, false
	}

	valueType := defineModelsType(value)
	dm := dataModel{Type: valueType, Name: name, Value: value}
	var chunk chunkModel

	if len(Chunks) == 0 {
		chunk = CreateNewChunk()
		log.Println("Creating first chunk")
		chunk.AddDataToChunk(dm)
		Chunks = append(Chunks, &chunk)
	} else if len(chunk.data)+1 >= CHUNK_SIZE {
		chunk = CreateNewChunk()
		chunk.AddDataToChunk(dm)
		Chunks = append(Chunks, &chunk)
		log.Println("Chunk's data size has achieved max length\nCreating new one")
	} else {
		chunk = *Chunks[len(Chunks)-1]
		chunk.AddDataToChunk(dm)
		Chunks[len(Chunks)-1] = &chunk
		log.Println("Last chunk is going to be used")
	}

	if len(chunk.data)%10 == 0 {
		chunk.SaveDataChunk()
		log.Println("Saving chunk...")
	}

	return dm, true
}

// Removes data by name
func RemoveData(name string) (bool, dataModel) {
	data, exists := FindData(name)
	if !exists {
		return false, *data
	}
	log.Printf("Clearing data from \"%s\".\n", data.Name)
	z := data.clear()
	return z, dataModel{}
}

// Updates data by name
func UpdateData(name string, value any) (bool, dataModel) {
	data, exists := FindData(name)
	if !exists {
		return false, dataModel{}
	}
	log.Printf("Updating data in \"%s\" from \"%v\" to \"%v\".\n", data.Name, data.Value, value)
	z, dm := data.update(value)
	return z, dm
}

// Creates new data chunk
func CreateNewChunk() chunkModel {
	var chunkId int
	if len(Chunks) == 0 {
		chunkId = 0
	} else {
		chunkId = len(Chunks)
	}
	return chunkModel{chunkId: chunkId}
}
