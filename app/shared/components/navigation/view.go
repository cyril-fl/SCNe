package navigation

import (
	"solopg/app/cmdrun/types/direction"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/services/style"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

type Navigation struct {
	plumprimitive.Fallible
	cache

	static []SectionInterface
	scroll []FocusableSectionInterface
	size   size.Size
	view   viewport.Model
}

type cache struct {
	active int
}

type Template struct {
	Static     []SectionInterface
	Scrollable []FocusableSectionInterface
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

	n.resizeView()
	n.view.SetContent(content)
	n.setOffset(contentoffsets)

	return style.Panel.
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
	views := make([]string, 0, len(n.static))
	for _, item := range n.static {
		views = append(views, item.GetView())
	}
	return strings.Join(views, "\n")
}

func (n *Navigation) SetSize(newsize size.Size) {
	n.size = newsize

	w, h := n.resizeView()
	itemsize := size.New(size.Template{
		Width:  w,
		Height: h,
	})

	for _, item := range n.static {
		item.SetSize(itemsize)
	}
	for _, item := range n.scroll {
		item.SetSize(itemsize)
	}
}

// resizeView fits the viewport under the static sections, whose height can change with their content.
func (n *Navigation) resizeView() (int, int) {
	w := max(0, n.size.GetWidth()-style.Panel.GetHorizontalFrameSize())
	h := max(0, n.size.GetHeight()-style.Panel.GetVerticalFrameSize())

	if static := n.getStaticView(); static != "" && len(n.scroll) > 0 {
		h -= lipgloss.Height(static) + 1
	}

	h = max(0, h)

	n.view.SetWidth(w)
	n.view.SetHeight(h)

	return w, h
}

func (n *Navigation) setOffset(offsets []int) {
	if n.getActiveItem() != nil {
		n.view.EnsureVisible(offsets[n.cache.active], 0, 0)
	}
}

func (n *Navigation) getActiveItem() FocusableSectionInterface {
	if len(n.scroll) == 0 {
		return nil
	}
	return n.scroll[n.cache.active]
}

// Methods
func (n *Navigation) Merge(items ...FocusableSectionInterface) *Navigation {
	n.scroll = append(n.scroll, items...)
	return n
}

func (n *Navigation) Add(item FocusableSectionInterface) *Navigation {
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
		switch plumtypes.Key(msg.String()) {
		case plumtypes.KEY_UP, plumtypes.KEY_DOWN:
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
	section FocusableSectionInterface
}
