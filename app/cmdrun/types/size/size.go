package size

import tea "charm.land/bubbletea/v2"

type Size struct {
	width  int
	height int
}

type Template struct {
	Width  int
	Height int
}

func New(params Template) Size {
	return Size{
		width:  params.Width,
		height: params.Height,
	}
}
func NewFromWindow(msg *tea.WindowSizeMsg) Size {
	if msg == nil {
		return Size{}
	}

	return Size{
		width:  msg.Width,
		height: msg.Height,
	}
}

func Nil() Size {
	return Size{
		width:  0,
		height: 0,
	}
}

// Getters & Setters
func (s Size) GetWidth() int {
	return max(0, s.width)
}

func (s Size) GetHeight() int {
	return max(0, s.height)
}

// Methods

// Helpers
