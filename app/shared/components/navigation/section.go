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

type Section struct {
	primitive.List

	id                string
	focused           bool
	footer            []string
	handlerEnterInput func() tea.Cmd
	sendMsg           func() tea.Cmd
}

type TemplateSection[T pluminterface.Named] struct {
	ID               string
	Options          []T
	Focused          bool
	Footer           []string
	HandleEnterInput func(T) tea.Cmd
	Msg              tea.Msg
}

func NewSection[T pluminterface.Named](params TemplateSection[T]) *Section {
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
	title := plumtransform.Uppercase(i19n.Localize(s.id))
	title = style.Title.Render(title)

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
