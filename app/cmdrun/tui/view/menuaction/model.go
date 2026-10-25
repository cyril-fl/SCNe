package menuaction

import (
	"fmt"
	"maps"
	"slices"

	"scne/app/cmdrun/domain/card/characters"
	"scne/app/cmdrun/domain/gameplay/portal"
	"scne/app/cmdrun/domain/gameplay/session"
	"scne/app/cmdrun/tui/shortcut"
	"scne/app/cmdrun/types/id"
	"scne/app/shared/components/form"
	"scne/app/shared/components/form/field"
	"scne/app/shared/components/models"
	"scne/app/shared/components/navigation"
	"scne/app/shared/services/i19n"

	tea "charm.land/bubbletea/v2"
	plumtypes "github.com/cyril-fl/plum/types"
)

const (
	FIGHT_PAGE  = "action.fight"
	TRAVEL_PAGE = "action.travel"
	// COMBAT_PAGE is the combat chat opened by the fight page.
	COMBAT_PAGE = "action.combat"
)

func NewSection(focused bool) *actionMenu {
	return &actionMenu{
		FocusableSection: navigation.NewFocusableSection(navigation.TemplateFocusableSection[*actionMenuItem]{
			ID:      "actions",
			Options: makeActionList(),
			Focused: focused,
		}),
	}
}

type actionMenu struct {
	*navigation.FocusableSection
	isSubviewOpen bool
}

// Getters & Setters

func (m *actionMenu) GetFooter() []string {
	if !m.isSubviewOpen {
		return []string{shortcut.OPEN.String()}
	}

	return nil
}

func (m *actionMenu) getSelectedPage() *actionMenuItem {
	if item, ok := m.GetList().SelectedItem().(models.Item[*actionMenuItem]); ok {
		return item.Value()
	}
	return nil
}

// Methods

// The action side menu only selects an action; its page is opened in the
// gameboard state machine.
func (m *actionMenu) HandleUpdate(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case form.Close:
		m.isSubviewOpen = false
	case form.Open:
		m.isSubviewOpen = true
	case tea.KeyPressMsg:
		switch plumtypes.Key(msg.String()) {
		case plumtypes.CMD_CTRL_ENTER:
			if selected := m.getSelectedPage(); selected != nil {
				return form.OpenPage(selected.id)
			}
		}
	}
	return nil
}

// Helpers

// NewFormPages creates the action pages owned by the gameboard state machine.
func NewFormPages(gameSession *session.GameSession) map[string]*actionPage {
	return map[string]*actionPage{
		FIGHT_PAGE:  newFightPage(gameSession),
		TRAVEL_PAGE: newTravelPage(gameSession),
	}
}

// newFightPage picks one or more codex NPCs (ctrl+"+" adds one) and opens the combat chat against them.
func newFightPage(gameSession *session.GameSession) *actionPage {
	return newPage(PageTemplate{
		Label: "fight",
		NewForm: func() *form.Form {
			options := getOpponentOptions(gameSession)
			return form.New(
				form.Repeat(form.RepeatTemplate{
					ID: "opponents",
					Make: func(index int) []form.Field {
						return []form.Field{
							field.SelectField(field.SelectTemplate[string]{Template: field.Template[string]{ID: fmt.Sprintf("opponent.%d", index), Label: "field.opponent", Required: true}, Options: options}),
						}
					},
				}),
			).Focus()
		},
		Submit: func(values *form.Form) tea.Cmd {
			opponents := findOpponents(gameSession, values.GetValuesAsMappedString())
			gameSession.StartCombat(opponents)
			return form.OpenPage(COMBAT_PAGE)
		},
	})
}

// newTravelPage sets the chosen codex location as the session's current location.
func newTravelPage(gameSession *session.GameSession) *actionPage {
	return newPage(PageTemplate{
		Label: "travel",
		NewForm: func() *form.Form {
			return form.New(
				field.SelectField(field.SelectTemplate[id.ID]{Template: field.Template[id.ID]{ID: "location", Label: "field.location", Required: true}, Options: getLocationOptions(gameSession)}),
			).Focus()
		},
		Submit: func(values *form.Form) tea.Cmd {
			entryID := values.GetFieldValue("location").(id.ID)
			portal.Travel(gameSession, entryID)
			return nil
		},
	})
}

// getOpponentOptions lists the codex NPCs by entry ID, as Repeat reads its values as strings.
func getOpponentOptions(gameSession *session.GameSession) []models.Item[string] {
	items := make([]models.Item[string], 0)
	for _, entry := range gameSession.State.Codex.Npcs.Entries {
		if entry.Character == nil {
			continue
		}
		items = append(items, models.NewItem(i19n.Localize(entry.Character.GetName()), "", entry.ID.String()))
	}
	return items
}

// findOpponents returns the codex NPCs whose entry ID is among values, once each, in codex order.
func findOpponents(gameSession *session.GameSession, values map[string]string) []*characters.Character {
	selected := slices.Collect(maps.Values(values))
	opponents := make([]*characters.Character, 0, len(selected))

	for _, entry := range gameSession.State.Codex.Npcs.Entries {
		if entry.Character != nil && slices.Contains(selected, entry.ID.String()) {
			opponents = append(opponents, entry.Character)
		}
	}
	return opponents
}

func getLocationOptions(gameSession *session.GameSession) []models.Item[id.ID] {
	items := make([]models.Item[id.ID], 0)
	for _, entry := range gameSession.State.Codex.Locations.Entries {
		if entry.Location == nil {
			continue
		}
		items = append(items, models.NewItem(i19n.Localize(entry.Location.GetName()), "", entry.ID))
	}
	return items
}
