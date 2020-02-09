package json


import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// Getter is the interface that is returned by the Loads() function. Any new methods that can be used to get/validate values from the json can be used inside Getter. 
type Getter interface{
	Get() (Getter, error)
}

type jsonData interface{}
type data struct{
	sync.Mutex
	jsonData
}

// Loads load string and returns a Getter
func Loads(s string) {
	fmt.Println(reflect.TypeOf(s))
	j, err := getAsJSON([]byte(s))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v", j)
}


func getAsJSON(b []byte) (*data, error){

	d := new(data)

	err := json.Unmarshal(b, &d.jsonData)
	if err != nil {
		return nil, err
	}
	return d, nil
}