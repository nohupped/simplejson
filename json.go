package simplejson

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
	// String return the json string representation of the selected key's value (Equivalent to json.Dumps()).
	String() string
	// Bytes return the marshalled byte representation of the value
	Bytes() []byte
}

// jsonData Holds the actual unmarshalled data
type jsonData interface{}

// data holds jsonData and a mutex.
type data struct {
	sync.Mutex
	jsonData
}

// Loads load a json string and returns a Getter
func Loads(b []byte) (Getter, error) {
	j := new(data)
	err := j.unmarshalJSON(b)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// Dumps returns the string representation of a type.
func Dumps(o interface{}) ([]byte, error) {
	return json.Marshal(o)
}

// Load accepts an io.Reader to read the content and unmarshall the json. The called must close the handler passed to this function.
func Load(i io.Reader) (Getter, error) {
	d, err := ioutil.ReadAll(i)
	if err != nil {
		return nil, err
	}
	return Loads(d)
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
			panic(fmt.Errorf("key error: %s not found", key))
		}
		sliceData.jsonData = d.jsonData.(map[string]interface{})[key]

	case []interface{}:
		sliceData.jsonData = d.jsonData.([]interface{})[index[0]]

	default:
		panic(fmt.Errorf("not implemented for %s", reflect.TypeOf(d.jsonData)))
	}

	return sliceData
}

// String return the json string representation of the selected key's value (Equivalent to json.Dumps()).
func (d *data) String() string {
	d.Lock()
	defer d.Unlock()
	data, err := json.Marshal(d.jsonData)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", data)
}

// Bytes return the marshalled byte representation of the value
func (d *data) Bytes() []byte {
	d.Lock()
	defer d.Unlock()
	data, err := json.Marshal(d.jsonData)
	if err != nil {
		panic(err)
	}
	return data
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
