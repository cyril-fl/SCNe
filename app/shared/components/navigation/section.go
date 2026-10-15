package navigation

import (
	"solopg/app/shared/services/i19n"
	interfass "solopg/app/shared/types/interface"
	"solopg/app/shared/types/primitive"
	"solopg/app/shared/utils/transform"

	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/models"

	"charm.land/lipgloss/v2"

	sharedtui "solopg/app/shared/tui"

	tea "charm.land/bubbletea/v2"
)

type Section struct {
	primitive.List

	id                string
	focused           bool
	footer            []string
	handlerEnterInput func() tea.Cmd
	sendMsg           func() tea.Cmd
}

type TemplateSection[T interfass.Named] struct {
	ID               string
	Options          []T
	Focused          bool
	Footer           []string
	HandleEnterInput func(T) tea.Cmd
	Msg              tea.Msg
}

func NewSection[T interfass.Named](params TemplateSection[T]) *Section {
	model := makeSectionModel(params)
	return &Section{
		List:              primitive.Listed(model),
		id:                params.ID,
		focused:           params.Focused,
		footer:            params.Footer,
		handlerEnterInput: makeEnterInputHandler(&model, params),
		sendMsg:           makeSendMsgHandler(params.Msg),
	}
}

// Getters and Setters
func (s *Section) ID() string {
	return s.id
}

func (s *Section) SetFocus(focused bool) {
	s.focused = focused
	models.SetListFocus(s.GetList(), focused)

}

func (s *Section) SetSize(size size.Size) {
	list := s.GetList()
	list.SetSize(size.GetWidth(), models.ListContentHeight(list))
}

func (s *Section) GetView() string {
	list := s.GetList()
	title := transform.Uppercase(i19n.Localize(s.id))
	title = lipgloss.NewStyle().Bold(true).Render(title)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		list.View(),
	)
}

func (s *Section) GetFooter() []string {
	return s.footer
}

// Methods
func (s *Section) HandleUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch sharedtui.EventKey(msg.String()) {
		case sharedtui.CMD_SHIFT_ENTER, sharedtui.CMD_ALT_ENTER:
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
