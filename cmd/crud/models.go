package main

type Person struct {
	Uid       string    `json:"uid,omitempty"`
	Name      string    `json:"name,omitempty"`
	Type      int       `json:"type,omitempty"`
	CreatedAt int64     `json:"created_at,omitempty"`
	UpdatedAt int64     `json:"updated_at,omitempty"`
	FriendOf  []*Person `json:"friend_of,omitempty"`
}

type QueryPerson struct {
	Everyone []*Person `json:"everyone,omitempty"`
}
