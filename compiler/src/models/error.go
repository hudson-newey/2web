package models

type Error struct {
	Message  string
	FilePath string
	Position Position
}
