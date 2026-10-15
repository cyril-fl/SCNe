package primitive

import (
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

type Notifier struct {
	notify       bool
	notification string
	spinner      spinner.Model
}

func NewNotifier() Notifier {
	return Notifier{
		notify:       false,
		notification: "",
		spinner:      spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
}

// Getters & Setters
func (n *Notifier) SetNotification(notification string) {
	n.notification = notification
}

func (n *Notifier) GetNotification() string {
	return n.notification
}

func (n *Notifier) SetNotifying(notifying bool) {
	n.notify = notifying
}

func (n *Notifier) IsNotifying() bool {
	return n.notify
}

func (n *Notifier) GetSpinnerView() string {
	return n.spinner.View()
}
func (n *Notifier) GetSpinner() spinner.Model {
	return n.spinner
}

func (n *Notifier) SetSpinner(spinner spinner.Model) {
	n.spinner = spinner
}

// Methods
func (n *Notifier) HandleNotification(delay int, notification string) tea.Cmd {
	if n.notify {
		return nil
	}

	n.SetNotification(notification)

	n.notify = true
	return tea.Batch(
		n.spinner.Tick,
		tea.Tick(time.Duration(delay)*time.Millisecond, func(time.Time) tea.Msg {
			n.notify = false
			return nil
		}),
	)
}

func (n *Notifier) HandleSpinnerTick(msg tea.Msg) tea.Cmd {
	if !n.IsNotifying() {
		return nil
	}

	var cmd tea.Cmd
	n.spinner, cmd = n.spinner.Update(msg)
	return cmd
}
