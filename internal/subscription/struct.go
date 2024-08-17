package subscription

import (
	"gorm.io/gorm"
	"nine-dubz/internal/user"
)

type Subscription struct {
	gorm.Model `json:"-"`
	ChannelID  uint
	Channel    user.User
	UserID     uint
	User       user.User
}
