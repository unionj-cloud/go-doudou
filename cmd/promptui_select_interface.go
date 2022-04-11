package cmd

//go:generate mockgen -destination ./mock/mock_promptui_select_interface.go -package mock -source=./promptui_select_interface.go

type ISelect interface {
	Run() (int, string, error)
	RunCursorAt(cursorPos, scroll int) (int, string, error)
	ScrollPosition() int
}