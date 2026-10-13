package navigation

import (
	"solopg/app/cmdrun/types/direction"
	"solopg/app/shared/components/models"
	"solopg/app/shared/services/i19n"
	"solopg/app/shared/services/logs"
	sharedtui "solopg/app/shared/tui"
	interfass "solopg/app/shared/types/interface"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Helper
func makeNavigationContent(part []SectionInterface) (string, []int) {
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

func makeSectionModel[T interfass.Named](params TemplateSection[T]) list.Model {
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

func makeEnterInputHandler[T interfass.Named](model *list.Model, params TemplateSection[T]) func() tea.Cmd {
	if params.HandleEnterInput == nil {
		return nil
	}

	return func() tea.Cmd {
		selected := model.SelectedItem()
		item, ok := selected.(models.Item[T])
		if !ok {
			return sharedtui.SendErrorMsg(logs.Error("error.invalid", map[string]any{
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
