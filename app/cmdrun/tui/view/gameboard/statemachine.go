package gameboard

import (
	"slices"

	"solopg/app/cmdrun/tui/view/sidemenu/codexmenu"
	"solopg/app/cmdrun/types/size"
	"solopg/app/shared/components/layout/footer"
	"solopg/app/shared/services/logs"
	interfass "solopg/app/shared/types/interface"
	"solopg/app/shared/types/primitive"

	tea "charm.land/bubbletea/v2"
)

type state = string

var CHAT_SECTION state = "chat"

type statemachine struct {
	primitive.Fallible
	curent     state
	transition map[state][]state
	view       map[state]stateview
	size       size.Size
}

type stateview interface {
	GetView() string
	Update(msg tea.Msg) tea.Cmd
	SetSize(size size.Size)
}

func newstate(params UiParams) statemachine {
	params.Engine.State.Codex.EnsureInitialized()
	pages := codexmenu.NewFormPages(params.Engine.State.Codex)
	views := map[state]stateview{CHAT_SECTION: makeChat(params.Engine)}
	transitions := map[state][]state{CHAT_SECTION: {}}

	for id, page := range pages {
		step := state(id)
		views[step] = page
		transitions[CHAT_SECTION] = append(transitions[CHAT_SECTION], step)
		transitions[step] = []state{CHAT_SECTION}
	}

	return statemachine{curent: CHAT_SECTION, transition: transitions, view: views}
}

func (s *statemachine) HandleTransition(from, to state) {
	if !slices.Contains(s.transition[from], to) {
		s.SetErr(logs.MildError("Transition pas autorise"))
		return
	}
	s.curent = to
	if view := s.view[to]; view != nil {
		view.SetSize(s.size)
	}
}

func (s statemachine) Current() state { return s.curent }

func (s statemachine) GetView() string { return s.view[s.curent].GetView() }

func (s *statemachine) Update(msg tea.Msg) tea.Cmd {
	return s.view[s.curent].Update(msg)
}

func (s *statemachine) SetSize(value size.Size) {
	s.size = value
	if view := s.view[s.curent]; view != nil {
		view.SetSize(value)
	}
}

func (s statemachine) Amend(content string) {
	if view, ok := s.view[s.curent].(interfass.Amandable); ok {
		view.Amend(content)
	}
}

func (s statemachine) getFooter() []string {
	var f []string
	if view, ok := s.view[s.curent].(footer.HasFooter); ok {
		f = append(f, view.GetFooter()...)
	}
	return f
}
