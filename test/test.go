package test

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/scottshotgg/storage/object"
	"github.com/scottshotgg/storage/storage"
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
	AmountOfTests = 100000
	WorkerLimit   = 1000
)

var (
	DB storage.Storage

	WorkerChan = make(chan struct{}, WorkerLimit)
	Objs       = make([]*object.Object, AmountOfTests)

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
	var (
		bytes []byte
		err   error
	)

	for i := 0; i < AmountOfTests; i++ {
		// Convert your custom object to binary first
		// Use JSON for ease
		Testeroonis[0].I = i
		Testeroonis[0].F += float64(i)
		bytes, err = json.Marshal(Testeroonis[0])
		if err != nil {
			log.Fatalf("err %+v", err)
		}

		// Create an object (Item) to put in the database
		Objs[i] = object.New(fmt.Sprintf("some_id_%d", i), bytes, []string{"another"})
		// map[string]interface{}{
		// 	"another": i % 10,
		// }))
	}
}
