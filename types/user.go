package types

// Login is the request body binding for login
type Login struct {
	Email    string `form:"email" json:"email" bson:"email" binding:"required"`
	Password string `form:"password" json:"password,omitempty" bson:"password" binding:"required"`
}

// GetEmail returns the user's email
func (auth *Login) GetEmail() string {
	return auth.Email
}

// GetPassword returns the user's password
// The password will be hashed if retrieving from database
func (auth *Login) GetPassword() string {
	return auth.Password
}

// Email stores email id of a user
type Email struct {
	Email string `form:"email" json:"email" bson:"email" binding:"required"`
}

// User stores user related information
type User struct {
	Email    string `form:"email" json:"email" binding:"required" bson:"email" valid:"required~Field 'email' is required but was not provided,email"`
	Password string `form:"password" json:"password,omitempty" bson:"password" binding:"required" valid:"required~Field 'password' is required but was not provided"`
	Username string `form:"username" json:"username" bson:"username" binding:"required" valid:"required~Field 'username' is required but was not provided,alphanum~Field 'username' should only have alphanumeric characters,stringlength(3|40)~Field 'username' should have length between 3 to 40 characters,lowercase~Field 'username' should have only lowercase characters"`
	Admin    bool   `form:"admin" json:"admin" bson:"admin"`
	Success  bool   `json:"success,omitempty" bson:"-"`
	GctlUUID string `json:"gctl_uuid" bson:"gctl_uuid"`
}

// GetName returns the user's username
func (user *User) GetName() string {
	return user.Username
}

// GetEmail returns the user's email
func (user *User) GetEmail() string {
	return user.Email
}

// SetEmail sets the user's email in its context
func (user *User) SetEmail(email string) {
	user.Email = email
}

// GetPassword returns the user's password
// The password will be hashed if retrieving from database
func (user *User) GetPassword() string {
	return user.Password
}

// SetPassword sets a password in the user's context
func (user *User) SetPassword(password string) {
	user.Password = password
}

// SetAdmin grants/revokes superuser privileges in the user's context
func (user *User) SetAdmin(admin bool) {
	user.Admin = admin
}

// IsAdmin checks whether a user has superuser privileges or not
func (user *User) IsAdmin() bool {
	return user.Admin
}

// SetSuccess defines the success of user creation
func (user *User) SetSuccess(success bool) {
	user.Success = success
}

// GetUUID returns the user's uuid
func (user *User) GetUUID() string {
	return user.GctlUUID
}

// SetUUID sets a unique id for authentication through gctl
func (user *User) SetUUID(uuid string) {
	user.GctlUUID = uuid
}

// PasswordUpdate is the request body for updating a user's password
type PasswordUpdate struct {
	OldPassword string `form:"old_password" json:"old_password,omitempty" binding:"required"`
	NewPassword string `form:"new_password" json:"new_password,omitempty" binding:"required"`
}

// GetOldPassword returns the user's old password
func (pw *PasswordUpdate) GetOldPassword() string {
	return pw.OldPassword
}

// GetNewPassword returns the user's new password
// which will replace the old password
func (pw *PasswordUpdate) GetNewPassword() string {
	return pw.NewPassword
}

// GCTLToken is the request body for generating GCTLToken
type GCTLToken struct {
	Email string `form:"email" json:"email" bson:"email" binding:"required"`
	Token string `form:"token" json:"token" bson:"token" binding:"required"`
}

// GetEmail returns the email of the user
func (auth *GCTLToken) GetEmail() string {
	return auth.Email
}
