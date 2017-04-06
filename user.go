package memberclicks

// type User struct {
// 	UserID      int64         `xml:"userID"`
// 	GroupID     int64         `xml:"groupID"`
// 	ContactName string        `xml:"contactName"`
// 	Active      bool          `xml:"active"`
// 	Validated   bool          `xml:"validated"`
// 	Deleted     bool          `xml:"deleted"`
// 	Attribute   UserAttribute `xml:"attribute"`
// }
//
// type UserAttribute struct {
// 	UserID     int64     `xml:"userID"`
// 	AttID      int64     `xml:"attID"`
// 	AttName    string    `xml:"attName"`
// 	AttTypeID  int64     `xml:"attTypeId"`
// 	AttData    string    `xml:"attData"`
// 	LastModify time.Time `xml:"lastModify"`
// }
//
// // GetUser returns the given user's data
// func (a *API) GetProfile(ctx context.Context, userID int64) (*User, error) {
// 	var u User
// 	if _, err := a.Get(ctx, fmt.Sprintf("/api/v1/profile/%d", userID), &u); err != nil {
// 		return nil, err
// 	}
// 	return &u, nil
// }
