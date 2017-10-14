package classic

import (
	"strings"
	"time"
)

type UserList struct {
	Users []Profile `json:"user" xml:"user"`
}

type AttributeType int

func (a AttributeType) String() string {
	switch int(a) {
	case 1:
		return "E-mail Address (Contact Center)"
	case 2:
		return "First Name"
	case 3:
		return "Last Name"
	case 4:
		return "Picture"
	case 5:
		return " Number"
	case 6:
		return "Fax Number (Contact Center)"
	case 7:
		return "Plain Text"
	case 8:
		return "User Group"
	case 9:
		return "Selection Set"
	case 10:
		return "Login Name"
	case 11:
		return "Password"
	case 12:
		return "Notes"
	case 13:
		return "E-mail Address"
	case 14:
		return "Fax Number"
	case 15:
		return "Web Page"
	case 16:
		return "Contact Center Greeting"
	case 19:
		return "Date"
	case 20:
		return " Date & Time"
	case 21:
		return " Numeric"
	case 24:
		return "Expiration Date"
	case 25:
		return "Hidden Email Address"
	case 26:
		return "Hidden Email Address (Contact Center)"
	case 27:
		return "Attachment"
	case 28:
		return "Address Line 1"
	case 29:
		return "Address Line 2"
	case 30:
		return "City"
	case 31:
		return "State"
	case 32:
		return "Zipcode"
	case 33:
		return "Country"
	default:
		return "Unknown"
	}
}

type Profile struct {
	UserID      string      `json:"userId" xml:"userId"`
	GroupID     string      `json:"groupId" xml:"groupId"`
	ContactName string      `json:"contactName" xml:"contactName"`
	Active      bool        `json:"active,string" xml:"active"`
	Validated   bool        `json:"validated,string" xml:"validated"`
	Deleted     bool        `json:"deleted,string" xml:"deleted"`
	Attributes  []Attribute `json:"attribute" xml:"attribute"`
}

type Attribute struct {
	UserID     string        `json:"userId" xml:"userId"`
	AttID      string        `json:"attId" xml:"attId"`
	AttTypeID  AttributeType `json:"attTypeId" xml:"attTypeId"`
	AttName    string        `json:"attName" xml:"attName"`
	AttData    string        `json:"attData" xml:"attData"`
	LastModify time.Time     `json:"lastModifiy" xml:"lastModify"`
}

func (a Attribute) String() string {
	return a.AttData
}

func (p *Profile) Get(attName string) (attData string) {
	attName = strings.ToLower(attName)
	for i := range p.Attributes {
		if strings.ToLower(p.Attributes[i].AttName) == attName {
			return p.Attributes[i].AttData
		}
	}
	return
}
