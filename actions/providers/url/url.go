package url

import (
	"bytes"
	"encoding/json"
	"github.com/bigblind/marvin/actions/domain"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

func init() {
	p := domain.NewProvider("url", "URL", "Actions triggering and triggered-by HTTP requests.")
	a := CallURL{}
	(&a).SetMeta("call_url", "Send a HTTP request", "Send an HtTP (web) request to a URL", false, true)
	p.Add(a)
}

// CallURL Is an Action that sends a HTTP request
type CallURL struct {
	domain.ActionMeta
}

// urlInput is a struct that'll receive input data for the action.
type urlInput struct {
	URL     string
	Body    string
	Method  string `enum:"GET|PUT|POST|PATCH|DELETE"`
	Headers []struct {
		Name  string
		Value string
	}
}

// Setup takes setup data from the frontend.
// In this case, we don't need any other data than ourn input, so do nothing.
func (a CallURL) Setup(data string, c domain.ActionContext) error {
	return nil
}

// InputType returns the type that input json will be deserialized into.
func (a CallURL) InputType(c domain.ActionContext) interface{} {
	return urlInput{}
}

// GlobaConfigType returns the type used to store global configuration.
// CallURL doesn't need any global config.
func (a CallURL) GlobalConfigType() interface{} {
	return nil
}

// Execute actually executes the action
func (a CallURL) Execute(input interface{}, c domain.ActionContext) error {
	var bodyr io.Reader
	inp := input.(urlInput)
	resp, err := a.makeRequest(inp)
	if err != nil {
		return err
	}
	if c.IsTestCall() {
		buf := bytes.NewBuffer([]byte{})
		_, err = buf.ReadFrom(resp)
		if err != nil {
			return err
		}

		bodys := buf.String()
		err = c.InstanceStore().Put("outputSchema", bodys)
		if err != nil {
			return err
		}

		bodyr = strings.NewReader(bodys)
	} else {
		bodyr = resp
	}

	r, err := arbitraryJSONToTypet(bodyr)
	if err != nil {
		return err
	}

	//todo: handle slices as the root element of r, and turn them into multiple outputs.

	c.Output(r)
	return nil
}

// OutputType returns an instance of the type that output from this action will have
func (a CallURL) OutputType(c domain.ActionContext) interface{} {
	json, _ := c.InstanceStore().Get("outputSchema")
	resp, _ := arbitraryJSONToTypet(strings.NewReader(json.(string)))
	return resp
}

// makeRequest actually makes the HTTP request
func (a CallURL) makeRequest(u urlInput) (io.Reader, error) {
	var body *strings.Reader
	if u.Method != "GET" && u.Method != "DELETE" {
		body = strings.NewReader(u.Body)
	}

	req, err := http.NewRequest(u.Method, u.URL, body)
	if err != nil {
		return nil, err
	}

	for _, h := range u.Headers {
		req.Header.Set(h.Name, h.Value)
	}

	c := http.Client{}
	resp, err := c.Do(req)
	return resp.Body, err
}

func arbitraryJSONToTypet(r io.Reader) (interface{}, error) {
	var i interface{}
	d := json.NewDecoder(r)
	err := d.Decode(i)
	if err != nil {
		return nil, err
	}
	return interfaceToType(i), nil
}

func interfaceToType(i interface{}) interface{} {
	switch t := i.(type) {
	case *interface{}:
		return interfaceToType(*t)
	case []interface{}:
		return handleSlice(t)
	case map[string]interface{}:
		return handleMap(t)
	default:
		return handleValue(t)
	}
}

func handleSlice(s []interface{}) interface{} {
	et := reflect.TypeOf(s).Elem()
	a := reflect.MakeSlice(et, 0, len(s))
	for _, o := range s {
		v := reflect.ValueOf(o)
		a = reflect.Append(a, v)
	}

	return a.Interface()
}

func handleMap(m map[string]interface{}) interface{} {
	// Initialize a list of struct fields
	fs := []reflect.StructField{}
	// Initialize a list, mapping from map keys to struct field names
	fm := map[string]string{}

	for k, v := range m {
		// convert the value into the type we want:
		v = interfaceToType(v)
		m[k] = v
		// Create a valid field name from the key
		vf := ensureValidFieldName(k)
		// Create a struct field
		f := reflect.StructField{
			Name:      vf,
			PkgPath:   "",
			Type:      reflect.TypeOf(v),
			Tag:       reflect.StructTag("json:\"" + k + "\""),
			Anonymous: false,
		}
		fs = append(fs, f)
		fm[k] = vf
	}

	// Creates A new struct type with the given fields, and instantiate it.
	s := reflect.New(reflect.StructOf(fs))

	// For each of the fields, set the value from the map
	for k, fn := range fm {
		s.FieldByName(fn).Set(reflect.ValueOf(m[k]))
	}

	return s
}

func ensureValidFieldName(k string) string {
	// This regexp matches anything that's not a unicode letter or digit
	r := regexp.MustCompile(`\P{L}|\D`)
	// Anything in the key that's not a letter or digit, becomes a space, splitting the string into words
	k = r.ReplaceAllString(k, "")
	// We uppercase the first letter of each word, so it looks kind of like a Go identifier (but with spaces)
	k = strings.Title(k)
	// And then remove the spaces
	k = strings.Replace(k, " ", "", -1) // remove ALL the spaces
	return k
}

func handleValue(i interface{}) interface{} {
	switch t := i.(type) {
	case float64:
		if t == float64(int64(t)) {
			return int64(t)
		}
		return t
	case bool, string, nil:
		return t
	default:
		log.Printf("Wasn't expecting %v (%v) when converting json to a struct. We should handle this.")
		return t
	}
}
