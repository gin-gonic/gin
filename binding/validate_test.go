package binding

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

// Make string into io.ReadCloser
// this is just for convinience
func jsonFactory(s string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(s))
}

// Ensure all required fields are matching
func TestRequired(t *testing.T) {

	// // Test if STRING required is valid
	var testString struct {
		Test string `json:"something" validate:"required" `
	}

	testJSON := jsonFactory(`{"something": "hello"}`)

	req, _ := http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testString); err != nil {
		t.Error(err)
	}

	var testString2 struct {
		Test string `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testString2); err == nil {
		t.Error("Required string, empty JSON object should return error but did not.")
	}

	// Test if INT require is valid
	var testInt struct {
		Test int `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{"something": 2}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testInt); err != nil {
		t.Error(err)
	}

	// Test if BOOL required is valid
	var testBool struct {
		Test bool `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{"something": true}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testBool); err != nil {
		t.Error(err)
	}

	var testBool2 struct {
		Test string `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testBool2); err == nil {
		t.Error("Required bool, empty JSON object should return error but did not.")
	}

	// Test if ARRAY required is valid
	var testArray struct {
		Test []string `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{"something": ["test", "data"]}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testArray); err != nil {
		t.Error(err)
	}

	// Test is OBJECT required is valid
	type testObjectTP struct {
		Name string `json:"name" validate:"required" `
	}

	var testObject struct {
		Test testObjectTP `json:"something" validate:"required" `
	}

	testJSON = jsonFactory(`{"something": {"name": "test"}}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testObject); err != nil {
		t.Error(err)
	}

	type testObjectTP2 struct {
		Name string `json:"name" validate:"required" `
	}
}

func TestEmail(t *testing.T) {

	var testValEmail struct {
		Test string `json:"email" validate:"email" `
	}

	testJSON := jsonFactory(`{"email": "michaeljs@gmail.com"}`)

	req, _ := http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValEmail); err != nil {
		t.Error(err)
	}

	var testValEmail2 struct {
		Test string `json:"email" validate:"email" `
	}

	testJSON = jsonFactory(`{"email": "michaeljs@gail.edu"}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValEmail2); err != nil {
		t.Error(err)
	}

	var testValEmail3 struct {
		Test string `json:"email" validate:"email" `
	}

	testJSON = jsonFactory(`{"email": "michaeljs.edu"}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValEmail3); err == nil {
		t.Error("Email test failed, michaeljs.edu passed as a valid email.")
	}

	// This should not return an error since email is not required.
	var testValEmail4 struct {
		Test string `json:"email" validate:"email" `
	}

	testJSON = jsonFactory(`{"jeff": "really"}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValEmail4); err != nil {
		t.Error(err)
	}

}

// Ensure In is matching properly
// Supporting string and int currently
func TestIn(t *testing.T) {

	var testValIn struct {
		Test string `json:"special" validate:"in:admin,user,other" `
	}

	testJSON := jsonFactory(`{"special": "admin"}`)

	req, _ := http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValIn); err != nil {
		t.Error(err)
	}

	var testValIn2 struct {
		Test int `json:"special" validate:"in:1,3,2" `
	}

	testJSON = jsonFactory(`{"special": 3}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValIn2); err != nil {
		t.Error(err)
	}

	var testValIn3 struct {
		Test int `json:"special" validate:"in:1,3,2" `
	}

	testJSON = jsonFactory(`{"special": 6}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValIn3); err == nil {
		t.Error("6 is not in validate in call, err should not have been nil.")
	}

	var testValIn4 struct {
		Test2 string `json:"what" validate:in:this,that`
		Test  int    `json:"special" validate:"in:1,3,2" `
	}

	testJSON = jsonFactory(`{"special": 3,"what": "this"}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValIn4); err != nil {
		t.Error(err)
	}

	var testValIn5 struct {
		Test2 string `json:"what" validate:in:this,that`
		Test  int    `json:"special" validate:"in:1,3,2" `
	}

	testJSON = jsonFactory(`{"special": 3}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValIn5); err != nil {
		t.Error(err)
	}

	var testValIn6 struct {
		Test2 string `json:"what" validate:"in:this,that"`
		Test3 string `json:"what1" validate:"in:this,then"`
		Test4 string `json:"what2" validate:"in:this,that"`
		Test5 string `json:"what3" validate:"in:this,that"`
		Test  int    `json:"special" validate:"in:1,3,2"`
	}

	testJSON = jsonFactory(`{"sa": 34, "what":"this", "what1":"then", "what2":"this"}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValIn6); err != nil {
		t.Error(err)
	}
}

// Check if the entered JSON is a data matching the one in a string.
func TestDigit(t *testing.T) {

	var testValDigit struct {
		Test int `json:"digit" validate:"digit:5" `
	}

	testJSON := jsonFactory(`{"digit": 12345}`)

	req, _ := http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValDigit); err != nil {
		t.Error(err)
	}

	var testValDigit2 struct {
		Test int `json:"digit" validate:"digit:5" `
	}

	testJSON = jsonFactory(`{"digit": 123456}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValDigit2); err == nil {
		t.Error("Error should have been thrown, digits should be 5 but was 6.")
	}

	var testValDigit3 struct {
		Test int `json:"digit" validate:"digit:12" `
	}

	testJSON = jsonFactory(`{"digit": 111111111111}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValDigit3); err != nil {
		t.Error(err)
	}

}

// Check if the entered JSON is a data matching the one in a string.
func TestDigitBetween(t *testing.T) {

	var testValDigit struct {
		Test int `json:"digit" validate:"digits_between:5,7" `
	}

	testJSON := jsonFactory(`{"digit": 123456}`)

	req, _ := http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValDigit); err != nil {
		t.Error(err)
	}

	var testValDigit2 struct {
		Test int `json:"digit" validate:"digits_between:0,10" `
	}

	testJSON = jsonFactory(`{"digit": 1234564}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValDigit2); err != nil {
		t.Error("Error should have been thrown, digit was not between 0 and 10.")
	}

	var testValDigit3 struct {
		Test int `json:"digit" validate:"digits_between:0,1" `
	}

	testJSON = jsonFactory(`{"digit": 1234564}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValDigit3); err == nil {
		t.Error(err)
	}

}

// Check if the entered JSON is a data matching the one in a string.
func TestMin(t *testing.T) {

	var testValMin struct {
		Test int `json:"digit" validate:"min:23" `
	}

	testJSON := jsonFactory(`{"digit": 24}`)

	req, _ := http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValMin); err != nil {
		t.Error(err)
	}

	var testValMin2 struct {
		Test int `json:"digit" validate:"min:20" `
	}

	testJSON = jsonFactory(`{"digit": 19}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValMin2); err == nil {
		t.Error("Min was 20 digit of 19 should not have validated properly.")
	}

	var testValMin3 struct {
		Test int `json:"digit" validate:"min:20" `
	}

	testJSON = jsonFactory(`{"jeff":"greg"}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValMin3); err != nil {
		t.Error("Nothing was entered but min was not required. No error should be thrown.")
	}
}

func TestMax(t *testing.T) {

	var testValMin struct {
		Test int `json:"digit" validate:"max:23" `
	}

	testJSON := jsonFactory(`{"digit": 23}`)

	req, _ := http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValMin); err != nil {
		t.Error(err)
	}

	var testValMin2 struct {
		Test int `json:"digit" validate:"max:20" `
	}

	testJSON = jsonFactory(`{"digit": 21}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValMin2); err == nil {
		t.Error("Max was 20 digit of 21 should not have validated properly.")
	}

	var testValMin3 struct {
		Test int `json:"digit" validate:"max:20" `
	}

	testJSON = jsonFactory(`{"jeff":"greg"}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValMin3); err != nil {
		t.Error("Nothing was entered but max was not required. No error should be thrown.")
	}
}

func TestRegex(t *testing.T) {

	var testValDigit struct {
		Test int `json:"digit" validate:"regex:\d+" `
	}

	testJSON := jsonFactory(`{"digit": 23}`)

	req, _ := http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValDigit); err != nil {
		t.Error(err)
	}

	var testValDigit2 struct {
		Test int `json:"digit" validate:"regex:\d+" `
	}

	testJSON = jsonFactory(`{"digit": 2dsa3}`)

	req, _ = http.NewRequest("POST", "/", testJSON)

	if err := JSON.Bind(req, &testValDigit2); err == nil {
		t.Error("\\d+ regex should not match the string 2dsa3.")
	}
}
