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
		Test *string `json:"something" validate:"required" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"something": "hello"}`))

	if err := JSON.Bind(req, &testString); err != nil {
		t.Error(err)
	}

	var testString2 struct {
		Test *string `json:"something" validate:"required" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{}`))

	if err := JSON.Bind(req, &testString2); err == nil {
		t.Error("Required string, empty JSON object should return error but did not.")
	}

	// Test if INT require is valid
	var testInt struct {
		Test *int `json:"something" validate:"required" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"something": 2}`))

	if err := JSON.Bind(req, &testInt); err != nil {
		t.Error(err)
	}

	// Test if BOOL required is valid
	var testBool struct {
		Test *bool `json:"something" validate:"required" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"something": true}`))

	if err := JSON.Bind(req, &testBool); err != nil {
		t.Error(err)
	}

	var testBool2 struct {
		Test *string `json:"something" validate:"required" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{}`))

	if err := JSON.Bind(req, &testBool2); err == nil {
		t.Error("Required bool, empty JSON object should return error but did not.")
	}

	// Test if ARRAY required is valid
	var testArray struct {
		Test *[]string `json:"something" validate:"required" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"something": ["test", "data"]}`))

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

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"something": {"name": "test"}}`))

	if err := JSON.Bind(req, &testObject); err != nil {
		t.Error(err)
	}
}

func TestEmail(t *testing.T) {

	var testValEmail struct {
		Test *string `json:"email" validate:"email" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"email": "michaeljs@gmail.com"}`))

	if err := JSON.Bind(req, &testValEmail); err != nil {
		t.Error(err)
	}

	var testValEmail2 struct {
		Test *string `json:"email" validate:"email" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"email": "michaeljs@gail.edu"}`))

	if err := JSON.Bind(req, &testValEmail2); err != nil {
		t.Error(err)
	}

	var testValEmail3 struct {
		Test *string `json:"email" validate:"email" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"email": "michaeljs.edu"}`))

	if err := JSON.Bind(req, &testValEmail3); err == nil {
		t.Error("Email test failed, michaeljs.edu passed as a valid email.")
	}

	// This should not return an error since email is not required.
	var testValEmail4 struct {
		Test *string `json:"email" validate:"email" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"jeff": "really"}`))

	if err := JSON.Bind(req, &testValEmail4); err != nil {
		t.Error(err)
	}

}

// Ensure In is matching properly
// Supporting string and int currently
func TestIn(t *testing.T) {

	var testValIn struct {
		Test *string `json:"special" validate:"in:admin,user,other" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"special": "admin"}`))

	if err := JSON.Bind(req, &testValIn); err != nil {
		t.Error(err)
	}

	var testValIn2 struct {
		Test *string `json:"special" validate:"in:1,3,2" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"special": "3"}`))

	if err := JSON.Bind(req, &testValIn2); err != nil {
		t.Error(err)
	}

	var testValIn3 struct {
		Test *int `json:"special" validate:"in:1,3,2" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"special": 6}`))

	if err := JSON.Bind(req, &testValIn3); err == nil {
		t.Error("6 is not in validate in call, err should not have been nil.")
	}

	var testValIn4 struct {
		Test2 *string `json:"what" validate:in:this,that`
		Test  *string `json:"special" validate:"in:1,3,2" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"special": "3","what": "this"}`))

	if err := JSON.Bind(req, &testValIn4); err != nil {
		t.Error(err)
	}

	var testValIn5 struct {
		Test2 *string `json:"what" validate:in:this,that`
		Test  *string `json:"special" validate:"in:1,3,2" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"special": "3"}`))

	if err := JSON.Bind(req, &testValIn5); err != nil {
		t.Error(err)
	}

	var testValIn6 struct {
		Test2 *string `json:"what" validate:"in:this,that"`
		Test3 *string `json:"what1" validate:"in:this,then"`
		Test4 *string `json:"what2" validate:"in:this,that"`
		Test5 *string `json:"what3" validate:"in:this,that"`
		Test  *string `json:"special" validate:"in:1,3,2"`
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"sa": 34, "what":"this", "what1":"then", "what2":"this"}`))

	if err := JSON.Bind(req, &testValIn6); err != nil {
		t.Error(err)
	}
}

// Check if the entered JSON is a data matching the one in a string.
func TestMin(t *testing.T) {

	var testValMin struct {
		Test *int `json:"digit" validate:"min:23" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"digit": 24}`))

	if err := JSON.Bind(req, &testValMin); err != nil {
		t.Error(err)
	}

	var testValMin2 struct {
		Test *int `json:"digit" validate:"min:20" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"digit": 19}`))

	if err := JSON.Bind(req, &testValMin2); err == nil {
		t.Error("Min was 20 digit of 19 should not have validated properly.")
	}

	var testValMin3 struct {
		Test *int `json:"digit" validate:"min:20" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"jeff":"greg"}`))

	if err := JSON.Bind(req, &testValMin3); err != nil {
		t.Error("Nothing was entered but min was not required. No error should be thrown.")
	}
}

func TestMax(t *testing.T) {

	var testValMin struct {
		Test *int `json:"digit" validate:"max:23" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"digit": 23}`))

	if err := JSON.Bind(req, &testValMin); err != nil {
		t.Error(err)
	}

	var testValMin2 struct {
		Test *int `json:"digit" validate:"max:20" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"digit": 21}`))

	if err := JSON.Bind(req, &testValMin2); err == nil {
		t.Error("Max was 20 digit of 21 should not have validated properly.")
	}

	var testValMin3 struct {
		Test *int `json:"digit" validate:"max:20" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"jeff":"greg"}`))

	if err := JSON.Bind(req, &testValMin3); err != nil {
		t.Error("Nothing was entered but max was not required. No error should be thrown.")
	}
}

func TestRegex(t *testing.T) {

	var testValDigit struct {
		Test *int `json:"digit" validate:"regex:\\d+" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"digit": 23}`))

	if err := JSON.Bind(req, &testValDigit); err != nil {
		t.Error(err)
	}

	var testValDigit2 struct {
		Test *int `json:"digit" validate:"regex:\\d+" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"digit": 2dsa3}`))

	if err := JSON.Bind(req, &testValDigit2); err == nil {
		t.Error("\\d+ regex should not match the string 2dsa3.")
	}
}

func TestMultiple(t *testing.T) {

	var testValMulti struct {
		Test *int `json:"digit" validate:"regex:\\d+|required|max:23" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"digit": 23}`))

	if err := JSON.Bind(req, &testValMulti); err != nil {
		t.Error(err)
	}

	var testValMulti2 struct {
		Test *string `json:"digit" validate:"email|required|regex:\\d+" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"digit": "m@g.com"}`))

	if err := JSON.Bind(req, &testValMulti2); err == nil {
		t.Error("Should have returned error but did not.")
	}

}

func TestPointers(t *testing.T) {

	var testValMulti struct {
		Test *string `json:"digit" validate:"in:3,4,5" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"invalid": "23"}`))

	if err := JSON.Bind(req, &testValMulti); err != nil {
		t.Error(err)
	}

	var testValMulti2 struct {
		Test *string `json:"digit" validate:"in:3,4,5" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"digit": "23"}`))

	if err := JSON.Bind(req, &testValMulti2); err == nil {
		t.Error("Value was passed in but did not match in:3,4,5 error should have been returned.")
	}

	var testValMulti3 struct {
		Test *string `json:"digit" validate:"in:3,4,5" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"digit": "4"}`))

	if err := JSON.Bind(req, &testValMulti3); err != nil {
		t.Error(err)
	}

}

func TestLength(t *testing.T) {

	var testValLength struct {
		Test *string `json:"username" validate:"length:5" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"username": "aaaaa"}`))

	if err := JSON.Bind(req, &testValLength); err != nil {
		t.Error(err)
	}

	var testValLength2 struct {
		Test *string `json:"username" validate:"length:4" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"username": "aaa"}`))

	if err := JSON.Bind(req, &testValLength2); err == nil {
		t.Error("Value was passed in but did not match length of 4 error should have been returned.")
	}

	var testValLength3 struct {
		Test *string `json:"username" validate:"length:23" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"digit": "4"}`))

	if err := JSON.Bind(req, &testValLength3); err != nil {
		t.Error(err)
	}

}

func TestLengthBetween(t *testing.T) {

	var testValLength struct {
		Test *string `json:"username" validate:"length_between:5,10" `
	}

	req, _ := http.NewRequest("POST", "/", jsonFactory(`{"username": "aaaaaa"}`))

	if err := JSON.Bind(req, &testValLength); err != nil {
		t.Error(err)
	}

	var testValLength2 struct {
		Test *string `json:"username" validate:"length_between:4,5" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"username": "aaa"}`))

	if err := JSON.Bind(req, &testValLength2); err == nil {
		t.Error("Value was passed in but was not inbetween 4,5 should have returned error.")
	}

	var testValLength3 struct {
		Test *string `json:"username" validate:"length_between:2,3" `
	}

	req, _ = http.NewRequest("POST", "/", jsonFactory(`{"digit": "4"}`))

	if err := JSON.Bind(req, &testValLength3); err != nil {
		t.Error(err)
	}

}
