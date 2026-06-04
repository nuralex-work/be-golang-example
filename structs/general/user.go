package general

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" query:"id" form:"id" json:"id,omitempty" param:"id"`
	Name     string    `gorm:"column:name;size:100" query:"name" form:"name" json:"name" param:"name" validate:"required"`
	Email    string    `gorm:"column:email;uniqueIndex;size:100" query:"email" form:"email" json:"email" param:"email"`
	NoHp     string    `gorm:"column:nohp;size:45" query:"nohp" form:"nohp" json:"nohp" param:"nohp"`
	Address  string    `gorm:"column:address;size:100" query:"address" form:"address" json:"address" param:"address"`
	Username string    `gorm:"column:username;size:255" query:"username" form:"username" json:"username,omitempty" param:"username" validate:"required"`
	Password string    `gorm:"column:password;size:255" query:"password" form:"password" json:"password,omitempty" param:"password"`
	Status   string    `gorm:"column:status;size:45" query:"status" form:"status" json:"status" param:"status"`
	RoleId   uuid.UUID `gorm:"type:uuid;column:roleid" query:"roleid" form:"roleid" json:"roleid" param:"roleid" validate:"required"`
}

type UserPayload struct {
	User
	CreateUserTime
}
type UserPayloadBulk struct {
	UserPayload []UserPayload `gorm:"column:users;size:100" query:"users" form:"users" json:"users" param:"users"`
}
type GetUserParams struct {
	GetRequestDefault
	Name string `query:"name" form:"name" json:"name" param:"name"`
	Role string `query:"role" form:"role" json:"role" param:"role"`
}
