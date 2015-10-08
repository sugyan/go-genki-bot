package main

import (
	"github.com/kurrik/twittergo"
)

// Users is users list
type Users []twittergo.User

// CursoredUsers is users list with cursor info
type CursoredUsers map[string]interface{}

// NextCursorStr returns next cursor
func (cl CursoredUsers) NextCursorStr() string {
	return cl["next_cursor_str"].(string)
}

// PreviousCursorStr returns previous cursor
func (cl CursoredUsers) PreviousCursorStr() string {
	return cl["previous_cursor_str"].(string)
}

// Users returns users list
func (cl CursoredUsers) Users() Users {
	var a []interface{} = cl["users"].([]interface{})
	b := make([]twittergo.User, len(a))
	for i, v := range a {
		b[i] = v.(map[string]interface{})
	}
	return b
}
