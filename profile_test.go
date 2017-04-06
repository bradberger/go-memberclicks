package memberclicks

import (
	"encoding/json"
	"testing"

	"google.golang.org/appengine/datastore"

	"github.com/stretchr/testify/assert"
)

func TestProfileLoadSave(t *testing.T) {

	var p Profile
	var profileID int64

	assert.Equal(t, int64(0), p.ID())
	assert.Equal(t, "0", p.GetID())

	assert.NoError(t, p.Load([]datastore.Property{{Name: "[Profile ID]", Value: 1002625212}}))
	assert.Len(t, p.attributes, 1)
	assert.Equal(t, int64(1002625212), p.ID())
	assert.Equal(t, "1002625212", p.GetID())

	assert.NoError(t, p.Load([]datastore.Property{{Name: "[Profile ID]", Value: int64(1002625212)}}))
	assert.Len(t, p.attributes, 1)
	assert.Equal(t, int64(1002625212), p.ID())
	assert.Equal(t, "1002625212", p.GetID())

	assert.NoError(t, p.Load([]datastore.Property{{Name: "[Profile ID]", Value: 1002625212.0}}))
	assert.Len(t, p.attributes, 1)
	assert.Equal(t, int64(1002625212), p.ID())
	assert.Equal(t, "1002625212", p.GetID())

	assert.NoError(t, p.Load([]datastore.Property{{Name: "[Profile ID]", Value: int32(123)}}))
	assert.Len(t, p.attributes, 1)
	assert.Equal(t, int64(0), p.ID())
	assert.Equal(t, "0", p.GetID())

	// Now that we've set profile ID to a 64-bit int, try to get it again.
	p.Set("[Profile ID]", int64(1002625212))
	assert.NoError(t, p.Get("[Profile ID]", &profileID))
	assert.Equal(t, int64(1002625212), profileID)

	p.Set("[Profile ID]", "1002625212")
	assert.Equal(t, int64(1002625212), p.ID())

	list, _ := p.Save()
	assert.Len(t, list, len(p.attributes))
}

func TestProfileGetAttrErr(t *testing.T) {
	p := Profile{attributes: map[string]interface{}{}}
	assert.Error(t, p.Get("foobar", nil))
}

func TestProfileGetSet(t *testing.T) {
	var profileID int64
	p := Profile{attributes: map[string]interface{}{}}

	p.Set("[Profile ID]", int64(1002625212))
	assert.NoError(t, p.Get("[Profile ID]", &profileID))
	assert.Equal(t, int64(1002625212), profileID)
}

func TestProfileEmptyMap(t *testing.T) {
	var p Profile
	assert.Equal(t, ErrEmptyMap, p.Get("foo", nil))
	assert.NotPanics(t, func() {
		p.Set("foo", "bar")
	})
}

func TestProfileJSON(t *testing.T) {
	var p Profile
	b := []byte(`{"[Profile ID]":1002625212}`)
	assert.NoError(t, json.Unmarshal(b, &p))
	assert.Equal(t, "1002625212", p.GetID())
	res, err := json.Marshal(&p)
	assert.NoError(t, err)
	assert.EqualValues(t, b, res)
}
