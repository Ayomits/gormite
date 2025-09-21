package entities

import "github.com/KoNekoD/gormite/tests/docs_example/pkg/enums"

// Entities can have custom table names
//      \/

// User "app_user"
type User struct {
	ID           int                `db:"id" pk:"true"`
	Email        string             `db:"email" uniq:"email" length:"180" uniq_cond:"email:(identity_type = 'email')"`
	Phone        string             `db:"phone" uniq:"phone" length:"10" uniq_cond:"phone:(identity_type = 'phone')"`
	IdentityType enums.IdentityType `db:"identity_type" type:"varchar"`
	Code         *string            `db:"code" nullable:"true"`
	FullName     *string            `db:"full_name" default:"'Anonymous'" index:"index_full_name" index_cond:"index_full_name:(is_active = true)"`
	IsActive     bool               `db:"is_active"`
}

type UserProfile struct {
	ID   int   `db:"id" pk:"true"`
	User *User `db:"user_id"`
	Age  int   `db:"age"`
}
