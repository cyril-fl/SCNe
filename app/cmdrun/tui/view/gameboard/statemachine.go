package gameboard

import (
	"slices"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/layout/footer"
	"solopg/app/shared/services/logs"
	"solopg/app/shared/types/primitive"

	tea "charm.land/bubbletea/v2"
)

type state = string
var (
	CHAT_SECTION state = "chat"
	FORM state = "form"
)

type statemachine struct {
	primitive.Fallible
	curent state
	transition map[state][]state
	view map[state]stateview
}

type stateview interface {
	GetView () string
	Update (msg tea.Msg) tea.Cmd
	SetSize(size size.Size)
}


func newstate(params UiParams) statemachine {
	return statemachine{
		curent: CHAT_SECTION,
		transition: map[state][]state{
			CHAT_SECTION : {FORM},
			FORM : {CHAT_SECTION},
		},
		view: map[state]stateview{
			CHAT_SECTION: makeChat(params.Engine),
		},
	}
}

func (s *statemachine) HandleTransition(from, to state) {
	if slices.Contains(s.transition[from], to) {
		s.curent = to
		return
	} 
	// i18n
	err := logs.MildError("Transition pas autorise")
	s.SetErr(err)
}
func (s statemachine) GetView() string {
	return s.view[s.curent].GetView()
}
func (s statemachine) Update (msg tea.Msg) tea.Cmd {
	return s.view[s.curent].Update(msg)
}
func (s statemachine) SetSize(size size.Size) {
	s.view[s.curent].SetSize(size)
}

func (s statemachine) getFooter() []string {
	var f[]string

	if view, ok := s.view[s.curent].(footer.HasFooter); ok {
		f = append(f, view.GetFooter()...)
	}
	return f
}