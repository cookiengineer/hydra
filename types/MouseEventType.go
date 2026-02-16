package types

type MouseEventType int

const (
	MouseMove          MouseEventType = 0
	MouseButtonPress   MouseEventType = 1
	MouseButtonRelease MouseEventType = 2
	MouseScroll        MouseEventType = 3
)

