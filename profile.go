package memberclicks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

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

func (p *Profile) Groups() []string {
	list := []string{}
	grps := []interface{}{}
	p.Get("[Group]", &grps)
	for i := range grps {
		list = append(list, grps[i].(string))
	}
	return list
}

func (p *Profile) Attributes() map[string]interface{} {
	return p.attributes
}

// DeleteAttr deletes a given attribute
func (p *Profile) DeleteAttr(names ...string) {
	if p.attributes == nil {
		return
	}
	for i := range names {
		delete(p.attributes, names[i])
	}
}

// MemberType returns the profile member type
func (p *Profile) MemberType() string {
	if p.attributes == nil || p.attributes["[Member Type]"] == nil {
		return ""
	}
	return p.attributes["[Member Type]"].(string)
}

// GetID implements the aedstorm.EntityID interface
func (p *Profile) GetID() string {
	return fmt.Sprintf("%v", p.ID())
}

// Entity implements the aedstorm.EntityName interface
func (p *Profile) Entity() string {
	return "profile"
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

	if id, ok := p.attributes["[Profile ID]"]; ok {
		switch id.(type) {
		case int64:
		case float64:
			p.attributes["[Profile ID]"] = int64(id.(float64))
		}
	}

	// Normalize the key names here for easier use in javascript, etc.
	// m := map[string]interface{}{}
	// for k, v := range p.attributes {
	// 	k = strings.Trim(k, " []")
	// 	k = strings.Replace(k, " | ", "_", -1)
	// 	k = strings.Replace(k, " ", "_", -1)
	// 	m[k] = v
	// }

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

// Me returns the profile associated with the accessToken
func (a *API) Me(ctx context.Context, accessToken string) (*Profile, error) {

	var p Profile
	req, err := http.NewRequest("GET", a.makeURL("/api/v1/profile/me"), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	if err := a.Do(ctx, req, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// Profile retrieves a profile with the given id
func (a *API) Profile(ctx context.Context, id string) (*Profile, error) {
	var p Profile
	if err := a.Get(ctx, "/api/v1/profile/"+id, &p); err != nil {
		return nil, err
	}
	return &p, nil
}
