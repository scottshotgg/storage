package test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pizzahutdigital/storage/object"
	"github.com/pizzahutdigital/storage/storage"
)

type Testerino struct {
	I  int
	F  float64
	B  bool
	S  string
	In interface{}
}

type Test struct {
	I  int
	F  float64
	B  bool
	S  string
	In interface{}
	St Testerino
}

const (
	AmountOfTests = 10000
	WorkerLimit   = 1000
)

var (
	WorkerChan = make(chan struct{}, WorkerLimit)
	DB         storage.Storage
	Objs       []*object.Object

	Testeroonis = []Test{
		Test{
			I:  6,
			F:  6.443,
			B:  true,
			S:  "yeah bro",
			In: "something here",
		},
	}
)

func init() {
	for i := 0; i < AmountOfTests; i++ {
		// Convert your custom object to binary first
		// Use JSON for ease
		Testeroonis[0].I = i
		Testeroonis[0].F += float64(i)
		bytes, err := json.Marshal(Testeroonis[0])
		if err != nil {
			log.Fatalf("err %+v", err)
		}

		// Create an object (item) to put in the database
		Objs = append(Objs, object.New(fmt.Sprintf("some_id_%d", i), bytes, map[string]interface{}{
			"another": i % 10,
		}))
	}
}
