package memberclicks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/appengine/datastore"
)

// Error messages
var (
	ErrNoSuchField = errors.New("no such field")
	ErrEmptyMap    = errors.New("map is empty")
)

var (
	_ json.Marshaler              = (*Profile)(nil)
	_ json.Unmarshaler            = (*Profile)(nil)
	_ datastore.PropertyLoadSaver = (*Profile)(nil)
)

// Profile is a memberclicks user profile. It stores the attributes in a private
// attributes map, and implmenets the interface of json.Marshaler, json.Unmarshaler,
// and datastore.PropertyLoadSaver to make the values accessible
type Profile struct {
	attributes map[string]interface{}
}

// ID returns the ID of the profile
func (p *Profile) ID() int64 {
	val, ok := p.attributes["[Profile ID]"]
	if !ok {
		return 0
	}
	switch val.(type) {
	case string:
		i, _ := strconv.ParseInt(val.(string), 10, 64)
		return i
	case int:
		return int64(val.(int))
	case float64:
		return int64(val.(float64))
	case int64:
		return val.(int64)
	}
	return 0
}

// GetID implements the goaedstorm.EntityID interface
func (p *Profile) GetID() string {
	return fmt.Sprintf("%v", p.ID())
}

// Get retrieves the profile attribute with the given name into dstVal
func (p *Profile) Get(name string, dstVal interface{}) error {
	if p.attributes == nil {
		return ErrEmptyMap
	}
	val, ok := p.attributes[name]
	if !ok {
		return ErrNoSuchField
	}
	return copy(val, dstVal)
}

// Set sets the profile attribute with the given name to val
func (p *Profile) Set(name string, val interface{}) {
	if p.attributes == nil {
		p.attributes = map[string]interface{}{}
	}
	p.attributes[name] = val
}

// MarshalJSON implements the json.Marshaler interface
func (p *Profile) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := json.NewEncoder(buf).Encode(p.attributes)
	return buf.Bytes(), err
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (p *Profile) UnmarshalJSON(data []byte) error {
	return json.NewDecoder(bytes.NewBuffer(data)).Decode(&p.attributes)
}

// Load implements the datastore.PropertyLoadSaver interface
func (p *Profile) Load(ps []datastore.Property) error {
	if p.attributes == nil {
		p.attributes = make(map[string]interface{}, 0)
	}
	for i := range ps {
		p.attributes[ps[i].Name] = ps[i].Value
	}
	return nil
}

// Save implements the PropertyLoadSaver interface
func (p *Profile) Save() (list []datastore.Property, err error) {
	for i := range p.attributes {
		list = append(list, datastore.Property{Name: i, Value: p.attributes[i]})
	}
	return
}
