package entity

// @Model
type User struct {
	UpdateDeleteEntity
	Email    string `gorm:"unique type:varchar(255);not null;column:email" json:"email"`
	Password string `gorm:"type:varchar(255);not null;column:password_hash" json:"password"`
	Role     string `gorm:"type:varchar(20);default:user" json:"role"`
}

/*
	for filtering field use like this for [carts] table:
	- carts.quantity -> even for current table filtering, always call the table name like this
	- products.name -> filter using products table with field name ->
	remember to not using struct field -> always use real tables and field name
*/

func (User) TableName() string {
	return "users"
}
