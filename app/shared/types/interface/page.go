package interfass

import (
	"solopg/app/cmdrun/types/size"

	tea "charm.land/bubbletea/v2"
)

type Page interface {
	GetView() tea.View
	SetBody(body string)
	SetFooter(footer string)
	GetSize() size.Size
	SetSize(size *tea.WindowSizeMsg)
	GetAvailableSize() size.Size
}

type Presentable interface {
	SetPageSize(size *tea.WindowSizeMsg)
}
