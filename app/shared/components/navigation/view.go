package navigation

import (
	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"
	sharedtui "solopg/app/shared/tui"
	"solopg/app/shared/types/primitive"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Navigation struct {
	primitive.Fallible
	cache

	static []string
	part   []Item
	size   size.Size
	view   viewport.Model
}

type cache struct {
	active int
}

type Template struct {
	Static []string
	Parts  []Item
}

type ItemTemplate struct {
	Focused bool
}

func New(template Template) *Navigation {
	return &Navigation{
		static: template.Static,
		part:   template.Parts,
		view:   viewport.New(),
	}
}

// Getters & Setters
func (n *Navigation) GetView() string {
	content, partStarts := n.makeContent()
	n.view.SetContent(content)

	if n.getActiveItem() != nil {
		n.view.EnsureVisible(partStarts[n.cache.active], 0, 0)
	}

	body := n.makeBody()

	return panelStyle().
		Width(n.size.GetWidth()).
		Height(n.size.GetHeight()).
		Render(body)
}

func (n *Navigation) getStaticView() string {
	return strings.Join(n.static, "\n")
}

func (n *Navigation) SetSize(panelSize size.Size) {
	n.size = panelSize
	style := panelStyle()

	w := max(0, n.size.GetWidth()-style.GetHorizontalFrameSize())
	h := max(0, n.size.GetHeight()-style.GetVerticalFrameSize())

	if static := n.getStaticView(); static != "" && len(n.part) > 0 {
		h -= lipgloss.Height(static) + 1
	}

	h = max(0, h)

	n.view.SetWidth(w)
	n.view.SetHeight(h)

	for _, item := range n.part {
		item.SetSize(size.New(size.Template{
			Width:  w,
			Height: h,
		}))
	}
}

func (n *Navigation) getActiveItem() Item {
	if len(n.part) == 0 {
		return nil
	}
	return n.part[n.cache.active]
}

// Methods
func (n *Navigation) Merge(items ...Item) *Navigation {
	n.part = append(n.part, items...)
	return n
}

func (n *Navigation) Add(item Item) *Navigation {
	n.part = append(n.part, item)
	return n
}

func (n *Navigation) HandleUpdate(msg tea.Msg) tea.Cmd {
	current := n.getActiveItem()
	if current == nil {
		return nil
	}
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.KEY_UP, sharedtui.KEY_DOWN:
			if openable, ok := current.(Openable); !ok || !openable.IsOpen() {
				n.handleDirectionInput(msg)
			}
		}
	}

	return current.HandleUpdate(msg)

}

func (n *Navigation) handleDirectionInput(key tea.KeyPressMsg) {
	currentMenu := n.getActiveItem()

	direction := HandleKeyArrow(Context{
		Menu:             currentMenu,
		CurrentMenuIndex: n.cache.active,
		SibblingCount:    len(n.part),
	}, key)

	n.updateMenuDirection(direction)

	// currentMenu.HandleDirectionInput(direction)
}

func (n *Navigation) updateMenuDirection(msg direction.Direction) {
	previousMenuIndex := n.cache.active
	menuLength := len(n.part)

	switch msg {
	case direction.PREVIOUS:
		if n.cache.active > 0 {
			n.cache.active--
		}
	case direction.NEXT:
		if n.cache.active < menuLength-1 {
			n.cache.active++
		}
	}

	if previousMenuIndex != n.cache.active {
		n.part[previousMenuIndex].SetFocus(false)
		n.part[n.cache.active].SetFocus(true)
	}
}

func (n *Navigation) makeContent() (string, []int) {
	blocks := make([]string, 0, len(n.part))
	partStarts := make([]int, len(n.part))
	offset := 0

	for i, item := range n.part {
		partStarts[i] = offset
		block := item.GetMenuView()
		if i == n.cache.active {
			if view := item.GetView(); view != "" {
				block = view
			}
		}
		blocks = append(blocks, block)
		offset += max(1, lipgloss.Height(block))
	}

	return strings.Join(blocks, "\n"), partStarts
}

func (n *Navigation) makeBody() string {
	body := n.view.View()
	if static := n.getStaticView(); static != "" {
		if body != "" {
			body = static + "\n" + body
		} else {
			body = static
		}
	}
	return body
}

// Helpers

func panelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("5"))
	// Background(lipgloss.Color("2")).
}
