package classic

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testProfile = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<user>
    <userId>12345678</userId>
    <groupId>1234</groupId>
    <contactName>Test User</contactName>
    <active>true</active>
    <validated>true</validated>
    <deleted>false</deleted>
    <attribute>
        <userId>12345678</userId>
        <attId>4567</attId>
        <attName>Email</attName>
        <attTypeId>1</attTypeId>
        <attData>test@memberclicks.com</attData>
        <lastModify>2010-10-21T16:46:31-04:00</lastModify>
    </attribute>
    ...
</user>`)

	testUserList = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<userList>
    <user>
        <userId>12345678</userId>
        <groupId>1234</groupId>
        <contactName>Test User</contactName>    
    </user>
</userList>`)
)

func TestProfileXMLDecode(t *testing.T) {
	var p Profile
	r := bytes.NewBuffer(testProfile)
	if !assert.NoError(t, xml.NewDecoder(r).Decode(&p)) {
		return
	}
	assert.Equal(t, "4567", p.Attributes[0].AttID)
	assert.Equal(t, 2010, p.Attributes[0].LastModify.Year())
	assert.Equal(t, "test@memberclicks.com", p.Get("email"))
	assert.Equal(t, "Test User", p.ContactName)
	assert.Equal(t, "E-mail Address (Contact Center)", p.Attributes[0].AttTypeID.String())
}

func TestUserListXMLDecode(t *testing.T) {
	var l UserList
	r := bytes.NewBuffer(testUserList)
	if !assert.NoError(t, xml.NewDecoder(r).Decode(&l)) {
		return
	}
	assert.Len(t, l.Users, 1)
	assert.Equal(t, "Test User", l.Users[0].ContactName)
}

func TestUserList(t *testing.T) {

}
