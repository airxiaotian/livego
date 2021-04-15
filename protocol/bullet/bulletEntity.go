package bullet

import (
	"time"
)

type Bullet struct {
	ID       int64 `gorm:"primary_key"`
	Content  string
	SentTime time.Time
}
