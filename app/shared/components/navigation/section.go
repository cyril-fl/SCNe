package navigation

import (
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/style"
	"solopg/app/shared/types/primitive"

	pluminterface "github.com/cyril-fl/plum/interface"

	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/models"

	"charm.land/lipgloss/v2"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
	plumtransform "github.com/cyril-fl/plum/utils/transform"
)

// - Section - //
type Section struct {
	id      string
	footer  []string
	content func() string
}

type TemplateSection struct {
	ID      string
	Footer  []string
	Content func() string
}

func NewSection(params TemplateSection) *Section {
	return &Section{
		id:      params.ID,
		footer:  params.Footer,
		content: params.Content,
	}
}

// Getters and Setters
func (s *Section) ID() string {
	return s.id
}

func (s *Section) SetSize(size.Size) {}

// GetView renders the content on every call so it reflects the latest state.
func (s *Section) GetView() string {
	title := plumtransform.Uppercase(i19n.Localize(s.id))
	title = style.Title.Render(title)

	content := ""
	if s.content != nil {
		content = s.content()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		content,
	)
}

func (s *Section) GetFooter() []string {
	return s.footer
}

// Methods
func (s *Section) HandleUpdate(tea.Msg) tea.Cmd {
	return nil
}

// - Focusable Section - //
type FocusableSection struct {
	*Section
	primitive.List

	focused           bool
	handlerEnterInput func() tea.Cmd
	sendMsg           func() tea.Cmd
}

type TemplateFocusableSection[T pluminterface.Named] struct {
	ID               string
	Options          []T
	Focused          bool
	Footer           []string
	HandleEnterInput func(T) tea.Cmd
	Msg              tea.Msg
}

func NewFocusableSection[T pluminterface.Named](params TemplateFocusableSection[T]) *FocusableSection {
	model := makeSectionModel(params)
	s := &FocusableSection{
		List:              primitive.Listed(model),
		focused:           params.Focused,
		handlerEnterInput: makeEnterInputHandler(&model, params),
		sendMsg:           makeSendMsgHandler(params.Msg),
	}
	s.Section = NewSection(TemplateSection{
		ID:     params.ID,
		Footer: params.Footer,
		Content: func() string {
			return s.GetList().View()
		},
	})
	return s
}

// Getters and Setters
func (s *FocusableSection) SetFocus(focused bool) {
	s.focused = focused
	models.SetListFocus(s.GetList(), focused)

}

func (s *FocusableSection) SetSize(size size.Size) {
	list := s.GetList()
	list.SetSize(size.GetWidth(), models.ListContentHeight(list))
}

// Methods
func (s *FocusableSection) HandleUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.CMD_SHIFT_ENTER, plumtypes.CMD_CTRL_ENTER:
			cmds := []tea.Cmd{}
			if s.handlerEnterInput != nil {
				cmds = append(cmds, s.handlerEnterInput())
			}
			if s.sendMsg != nil {
				cmds = append(cmds, s.sendMsg())
			}
			return tea.Batch(cmds...)
		}
	}
	return nil
}

// 	Helper
