package simplejson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// Getter is the interface that is returned by the Loads() function. Any new methods that can be used to get/validate values from the json can be used inside Getter.
type Getter interface {
	// Get method is to get the value of a json key and also accepts an optional index number.
	// If the *data object is a map, it returns the value of the key as x.Get("key1"). The value of index is irrelevant in this case.
	// If the *data object is a slice, we need to select the index in one call to Get() as `x.Get("", 0)`
	// where the first parameter is irrelevant and 0 indicates the index of the element in the json slice.
	// Eg: data := d.Get("Actors").Get("", 0).Get("name")
	Get(string, ...int) Getter
	// Return the json string representation of the selected key's value (Equicalent to json.Dumps).
	String() string
}

// jsonData Holds the actual unmarshalled data
type jsonData interface{}

// data holds jsonData and a mutex.
type data struct {
	sync.Mutex
	jsonData
}

// Loads load a json string and returns a Getter
func Loads(s string) (Getter, error) {
	j := new(data)
	err := j.unmarshalJSON([]byte(s))
	if err != nil {
		return nil, err
	}
	return j, nil
}

// Get method is to get the value of a json key. Get() also accepts an optional index number.
// If the *data object is a map, it returns the value of the key as x.Get("key1"). The value of index is irrelevant in this case.
// If the *data object is a slice, we need to select the index in one call to Get() as `x.Get("", 0)`
// where the first parameter is irrelevant and 0 indicates the index of the element in the json slice.
// Eg: data := d.Get("Actors").Get("", 0).Get("name")
func (d *data) Get(key string, index ...int) Getter {
	d.Lock()
	defer d.Unlock()
	sliceData := new(data)
	sliceData.Lock()
	defer sliceData.Unlock()

	switch d.jsonData.(type) {

	case map[string]interface{}:
		if d.jsonData.(map[string]interface{})[key] == nil {
			panic(fmt.Sprintf("Key error: %s not found\n", key))
		}
		sliceData.jsonData = d.jsonData.(map[string]interface{})[key]
	case []interface{}:
		sliceData.jsonData = d.jsonData.([]interface{})[index[0]]

	default:
		panic(fmt.Sprintf("Not implemented for %s\n", reflect.TypeOf(d.jsonData)))
	}

	return sliceData
}

func (d *data) String() string {
	d.Lock()
	defer d.Unlock()
	data, err := json.Marshal(d.jsonData)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", data)
}

func (d *data) unmarshalJSON(b []byte) error {
	d.Lock()
	defer d.Unlock()
	err := json.Unmarshal(b, &d.jsonData)
	if err != nil {
		return err
	}
	return nil
}
