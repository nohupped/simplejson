package simplejson

import (
	"fmt"
	"os"
	"testing"
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

	d, err := Loads(sampleJSON)
	if err != nil {
		panic(err)
	}
	t.Logf("%s\n\n", d)
	data := d.Get("Actors").Get("", 0).Get("name")
	t.Logf("%s\n", data)

}

// TestLoadsFail parses an invalid json to test the failure to parse.
func TestLoadsFail(t *testing.T) {
	_, err := Loads(fmt.Sprintf("%s%s",sampleJSON, "invalidJson"))
	if err != nil{
		t.Logf("Error parsing json: %s", err)
	}
}

// TestLoad does simplejson.Load() to load json from disk.
func TestLoad(t *testing.T) {
	fd, err := os.Open("samplejson.json")
	if err != nil {
		panic(err)
	}
	d1, err := Load(fd)
	if err != nil{
		panic(err)
	}
	t.Logf("Extracting tag from json file: %s", d1.Get("", 2).Get("tags").Get("", 3))

	fd.Close()
}


// TestLoadFail tests the ioutil.ReadAll failure after trying to read from a closed fd.
func TestLoadFail(t *testing.T) {
	fd, err := os.Open("samplejson.json")
	fd.Close()
	if err != nil {
		panic(err)
	}
	_, err = Load(fd)
	if err != nil{
		t.Logf("Error in ioutil.ReadAll :%s", err)
	}
	
	fd.Close()
}