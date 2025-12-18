package models

type Conversation struct {
	Message *Message     `json:"message"`
	User    *UserProfile `json:"user"`
}
