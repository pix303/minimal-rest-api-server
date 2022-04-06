package domain

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            int32  `bun:"id,pk,autoincrement"`
	Username      string `bun:"username,notnull"`
}
