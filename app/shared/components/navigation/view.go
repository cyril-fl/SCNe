package navigation

import (
	"scne/app/cmdrun/types/direction"
	"scne/app/cmdrun/types/size"
	"scne/app/shared/services/style"
	"strings"

	plumprimitive "github.com/cyril-fl/plum/primitive"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

// sectionGap leaves one blank line between sections.
const sectionGap = "\n\n"

// Navigation is a panel of static sections above scrollable focusable ones.
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

// Template configures a Navigation.
type Template struct {
	Static     []SectionInterface
	Scrollable []FocusableSectionInterface
}

// New returns a Navigation.
func New(template Template) *Navigation {
	return &Navigation{
		static: template.Static,
		scroll: template.Scrollable,
		view:   viewport.New(),
	}
}

// Getters & Setters

// GetView renders the panel, scrolled to the active section.
func (n *Navigation) GetView() string {
	content, contentoffsets := makeNavigationContent(n.scroll)

	n.resizeView()
	n.view.SetContent(content)
	n.setOffset(contentoffsets)

	return style.Panel().
		Width(n.size.GetWidth()).
		Height(n.size.GetHeight()).
		Render(n.makeBody())
}

// GetFooter returns the shortcuts of the active section.
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
	return strings.Join(views, sectionGap)
}

// SetSize sets the panel size and resizes every section.
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
	w := max(0, n.size.GetWidth()-style.Panel().GetHorizontalFrameSize())
	h := max(0, n.size.GetHeight()-style.Panel().GetVerticalFrameSize())

	if static := n.getStaticView(); static != "" && len(n.scroll) > 0 {
		// +2 for the line break and the blank line of sectionGap.
		h -= lipgloss.Height(static) + 2
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

// Merge appends scrollable sections.
func (n *Navigation) Merge(items ...FocusableSectionInterface) *Navigation {
	n.scroll = append(n.scroll, items...)
	return n
}

// Add appends a scrollable section.
func (n *Navigation) Add(item FocusableSectionInterface) *Navigation {
	n.scroll = append(n.scroll, item)
	return n
}

// HandleUpdate moves to the adjacent section on up/down past the list edge,
// then forwards msg to the section that was active.
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

	return static + sectionGap + body
}

// Helpers
type context struct {
	section FocusableSectionInterface
}
