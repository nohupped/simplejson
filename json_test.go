package simplejson

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleJSON string = `{
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

func TestLoadsEmpty(t *testing.T) {

	d, err := Loads([]byte(sampleJSON))
	assert.Nil(t, err)
	data := d.Get("Actors").Get("", 0).Get("names").Get("Foo").Get("Bar")
	assert.Equal(t, data.String(), "", "Should equals \"\"")
}

func TestEmpty(t *testing.T) {

	d, err := Loads([]byte(sampleJSON))
	assert.Nil(t, err)
	data := d.Get("Actors").Get("", 0).Get("names")
	e := new(empty)
	if data != nil {
		assert.IsType(t, e, data)
	} else {
		assert.IsType(t, e, data)
	}
}

// TestLoadsFail parses an invalid json to test the failure to parse.
func TestLoadsFail(t *testing.T) {
	_, err := Loads([]byte(fmt.Sprintf("%s%s", sampleJSON, "invalidJson")))
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
	assert.EqualError(t, err, "read samplejson.json: file already closed")
}

// TestLoadKeyFail tests the key error.
func TestLoadsKeyFail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected recovery from Key Error. Error is: %s", r)
			assert.EqualError(t, r.(error), "key error: names not found")
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
			t.Logf("Expected recovery from Json Marshal in String() method. Error is: %s", r)
			assert.EqualError(t, r.(error), "json: unsupported type: chan int")
		}
	}()
	d := new(data)
	d.jsonData = make(chan int)
	t.Logf("%s", d.String())
}

// TestDumpAsBytesFail tests the json Marshal error in the String method.
func TestDumpAsBytesFail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected recovery from Json Marshal in DumpAsBytes() method. Error is: %s", r)
			assert.EqualError(t, r.(error), "json: unsupported type: chan int")
		}
	}()
	d := new(data)
	d.jsonData = make(chan int)
	t.Logf("%s", d.Bytes())
}

// TestGetFail tests "Not Implemented" part of the Get method.
func TestGetFail(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected recovery from String() method. Error is: %s", r)
			assert.EqualError(t, r.(error), "not implemented for chan int")
		}
	}()
	d := new(data)
	d.jsonData = make(chan int)
	t.Logf("Output is %s", d.Get("Invalid"))
}

func TestJsonTypes(t *testing.T) {
	assert := assert.New(t)
	val, err := Dumps(10000)
	assert.Nil(err)
	data, err := Loads(val)
	assert.Nil(err)
	assert.EqualValuesf(data.Bytes(), []byte{0x31, 0x30, 0x30, 0x30, 0x30}, "10000 in byte slice")
	val, err = Dumps(true)
	assert.Nil(err)
	data, err = Loads(val)
	assert.Nil(err)
	assert.EqualValuesf(data.Bytes(), []byte{0x74, 0x72, 0x75, 0x65}, "true in byte slice")
	val, err = Dumps(nil)
	assert.Nil(err)
	data, err = Loads(val)
	assert.Nil(err)
	assert.EqualValuesf(data.Bytes(), []byte{0x6e, 0x75, 0x6c, 0x6c}, "null in byte slice")
}
