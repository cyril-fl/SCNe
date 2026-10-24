package navigation

import (
	pluminterface "github.com/cyril-fl/plum/interface"
	"solopg/app/cmdrun/types/direction"
	"solopg/app/shared/components/models"
	"solopg/app/shared/messages"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Helpers
func makeNavigationContent(part []FocusableSectionInterface) (string, []int) {
	blocks := make([]string, 0, len(part))
	blockoffsets := make([]int, len(part))
	offset := 0

	for i, item := range part {
		block := item.GetView()
		blockoffsets[i] = offset
		blocks = append(blocks, block)
		offset += max(1, lipgloss.Height(block))
	}

	return strings.Join(blocks, "\n"), blockoffsets
}

func makeSectionModel[T pluminterface.Named](params TemplateFocusableSection[T]) list.Model {
	items := make([]list.Item, 0, len(params.Options))

	for _, option := range params.Options {
		names := i19n.Localize(option.GetName())
		items = append(items, models.NewItem(names, "", option))
	}

	model := list.New(items, list.NewDefaultDelegate(), 0, 0)
	models.ConfigureList(&model)
	model.SetSize(0, models.ListContentHeight(&model))

	return model
}

func makeEnterInputHandler[T pluminterface.Named](model *list.Model, params TemplateFocusableSection[T]) func() tea.Cmd {
	if params.HandleEnterInput == nil {
		return nil
	}

	return func() tea.Cmd {
		selected := model.SelectedItem()
		item, ok := selected.(models.Item[T])
		if !ok {
			return messages.Fail(logs.Error("error.invalid", map[string]any{
				"Subject": params.ID,
				"Value":   selected,
			}))
		}

		return params.HandleEnterInput(item.Value())
	}
}

func handleKeyArrow(metadata context, key tea.KeyPressMsg) direction.Direction {
	list, newdirection := direction.GetListDirection(metadata.section.GetList(), key)

	metadata.section.SetList(&list)

	return newdirection
}

func makeSendMsgHandler(msg tea.Msg) func() tea.Cmd {
	if msg == nil {
		return nil
	}

	return func() tea.Cmd {
		return func() tea.Msg {
			return msg
		}
	}
}
