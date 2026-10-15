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
	scroll []SectionInterface
	size   size.Size
	view   viewport.Model
}

type cache struct {
	active int
}

type Template struct {
	Static     []string
	Scrollable []SectionInterface
}

func New(template Template) *Navigation {
	return &Navigation{
		static: template.Static,
		scroll: template.Scrollable,
		view:   viewport.New(),
	}
}

// Getters & Setters
func (n *Navigation) GetView() string {
	content, contentoffsets := makeNavigationContent(n.scroll)

	n.view.SetContent(content)
	n.setOffset(contentoffsets)

	return panelStyle().
		Width(n.size.GetWidth()).
		Height(n.size.GetHeight()).
		Render(n.makeBody())
}

func (n *Navigation) GetFooter() []string {
	footer := []string{}
	if current := n.getActiveItem(); current != nil {
		footer = append(footer, current.GetFooter()...)
	}
	return footer
}

func (n *Navigation) getStaticView() string {
	return strings.Join(n.static, "\n")
}

func (n *Navigation) SetSize(newsize size.Size) {
	n.size = newsize
	style := panelStyle()

	w := max(0, n.size.GetWidth()-style.GetHorizontalFrameSize())
	h := max(0, n.size.GetHeight()-style.GetVerticalFrameSize())

	if static := n.getStaticView(); static != "" && len(n.scroll) > 0 {
		h -= lipgloss.Height(static) + 1
	}

	h = max(0, h)

	n.view.SetWidth(w)
	n.view.SetHeight(h)

	for _, item := range n.scroll {
		item.SetSize(size.New(size.Template{
			Width:  w,
			Height: h,
		}))
	}
}

func (n *Navigation) setOffset(offsets []int) {
	if n.getActiveItem() != nil {
		n.view.EnsureVisible(offsets[n.cache.active], 0, 0)
	}
}

func (n *Navigation) getActiveItem() SectionInterface {
	if len(n.scroll) == 0 {
		return nil
	}
	return n.scroll[n.cache.active]
}

// Methods
func (n *Navigation) Merge(items ...SectionInterface) *Navigation {
	n.scroll = append(n.scroll, items...)
	return n
}

func (n *Navigation) Add(item SectionInterface) *Navigation {
	n.scroll = append(n.scroll, item)
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
			n.handleDirectionInput(msg)
		}
	}

	return current.HandleUpdate(msg)

}

func (n *Navigation) handleDirectionInput(key tea.KeyPressMsg) {
	currentMenu := n.getActiveItem()

	direction := handleKeyArrow(context{
		section: currentMenu,
	}, key)

	n.updateMenuDirection(direction)
}

func (n *Navigation) updateMenuDirection(msg direction.Direction) {
	previousMenuIndex := n.cache.active
	menuLength := len(n.scroll)

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
		n.scroll[previousMenuIndex].SetFocus(false)
		n.scroll[n.cache.active].SetFocus(true)
	}
}

func (n *Navigation) makeBody() string {
	static := n.getStaticView()
	body := n.view.View()

	if static == "" {
		return body
	}

	if body == "" {
		return static
	}

	return static + "\n" + body
}

// Helpers
type context struct {
	section SectionInterface
}

// LOW deplacer ça dans un factory qui gerer les style
func panelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("5"))
	// Background(lipgloss.Color("2")).
}
