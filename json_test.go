package simplejson

import (
	"fmt"
	"os"
	"testing"
	"reflect"
	"github.com/stretchr/testify/assert"
)
var sampleJSON string =`{
	"Actors": [
	  {
		"name": "Tom Cruise",
		"age": 56,
		"Born At": "Syracuse, NY",
		"Birthdate": "July 3, 1962",
		"hasChildren":true,
		  "hasTwitterAccount":true,
		  "hasGreyHair":false,
		"photo": "https://jsonformatter.org/img/tom-cruise.jpg"
	  },
	  {
		"name": "Robert Downey Jr.",
		"age": 53,
		"Born At": "New York City, NY",
		"Birthdate": "April 4, 1965",
		"photo": "https://jsonformatter.org/img/Robert-Downey-Jr.jpg"
	  }
	]
	}
`

// TestLoads parses a valid json and gets keys.
func TestLoads(t *testing.T) {

	d, err := Loads([]byte(sampleJSON))
	assert.Nil(t, err)
	data := d.Get("Actors").Get("", 0).Get("name")
	assert.Equal(t, data.String(), "\"Tom Cruise\"", "Should equals Tom Cruise")
}

// TestLoadsFail parses an invalid json to test the failure to parse.
func TestLoadsFail(t *testing.T) {
	_, err := Loads([]byte(fmt.Sprintf("%s%s",sampleJSON, "invalidJson")))
	assert.EqualError(t, err, "invalid character 'i' after top-level value")
}

// TestLoad does simplejson.Load() to load json from disk.
func TestLoad(t *testing.T) {	
	fd, err := os.Open("samplejson.json")
	assert.Nil(t, err)
	d1, err := Load(fd)
	assert.Nil(t, err)
	assert.Equal(t, d1.Get("", 2).Get("tags").Get("", 3).String(), "\"incididunt\"")
	fd.Close()
}


// TestLoadFail tests the ioutil.ReadAll failure after trying to read from a closed fd.
func TestLoadFail(t *testing.T) {
	fd, err := os.Open("samplejson.json")
	fd.Close()
	assert.Nil(t, err)
	_, err = Load(fd)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "read samplejson.json: file already closed")
}

// TestLoadKeyFail tests the key error.
func TestLoadsKeyFail(t *testing.T) {
	defer func() {
        if r := recover(); r != nil {
			t.Logf("Expected recovery from Key Error. Error is: %s", r)
			assert.Equal(t, r, "Key error: names not found\n")
        }
	}()
	
	d, err := Loads([]byte(sampleJSON))
	if err != nil {
		panic(err)
	}
	d.Get("Actors").Get("", 0).Get("names")

}

// TestStringFail tests the json Marshal error in the String method.
func TestStringFail(t *testing.T) {
	defer func() {
        if r := recover(); r != nil {
			t.Logf("Expected recovery from Json Marshal. Error is: %s", r)
        }
	}()
	d := new(data)
	d.jsonData = make(chan int)
}


// TestGetFail tests "Not Implemented" part of the Get method.
func TestGetFail(t *testing.T) {
	defer func() {
        if r := recover(); r != nil {
			t.Logf("Expected recovery from implementation error. Error is: %s", r)
			assert.Equal(t, r, "Not implemented for chan int\n")
        }
	}()
	d := new(data)
	d.jsonData = make(chan int)
	t.Logf("Output is %s", d.Get("Invalid"))
}


func TestJsonTypes(t *testing.T) {
	val1, err := Dumps(10000)
	if err != nil {
		panic(err)
	}
	t.Logf("%s: %s: %x", val1, reflect.TypeOf(val1), val1)
	// data1, err := Loads(val1)
}
