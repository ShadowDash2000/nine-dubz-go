package payload

import "nine-dubz/app/model"

type UserPayload struct {
	ID        omit `json:"ID,omitempty"`
	CreatedAt omit `json:"CreatedAt,omitempty"`
	UpdatedAt omit `json:"UpdatedAt,omitempty"`
	DeletedAt omit `json:"DeletedAt,omitempty"`
	Roles     omit `json:"roles,omitempty"`
	Tokens    omit `json:"tokens,omitempty"`
	Picture   omit `json:"picture,omitempty"`
	*model.User
}
