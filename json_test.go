package simplejson

import (
	"reflect"
	"testing"
)

func TestLoads(t *testing.T) {
	sampleJSON := `{
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

	d, err := Loads(sampleJSON)
	if err != nil {
		panic(err)
	}
	t.Logf("%s\n\n", d)
	data := d.Get("Actors").Get("", 0).Get("name")
	t.Logf("%s, %s\n", data, reflect.TypeOf(data))

}
