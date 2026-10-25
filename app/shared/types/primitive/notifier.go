package primitive

import (
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

// Notifier shows a temporary message with a spinner.
type Notifier struct {
	notify       bool
	notification string
	spinner      spinner.Model
}

// NewNotifier returns an idle Notifier with a dot spinner.
func NewNotifier() Notifier {
	return Notifier{
		notify:       false,
		notification: "",
		spinner:      spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
}

// Getters & Setters

// SetNotification sets the message.
func (n *Notifier) SetNotification(notification string) {
	n.notification = notification
}

// GetNotification returns the message.
func (n *Notifier) GetNotification() string {
	return n.notification
}

// SetNotifying sets whether the message is shown.
func (n *Notifier) SetNotifying(notifying bool) {
	n.notify = notifying
}

// IsNotifying reports whether the message is shown.
func (n *Notifier) IsNotifying() bool {
	return n.notify
}

// GetSpinnerView renders the spinner.
func (n *Notifier) GetSpinnerView() string {
	return n.spinner.View()
}

// GetSpinner returns the spinner.
func (n *Notifier) GetSpinner() spinner.Model {
	return n.spinner
}

// SetSpinner replaces the spinner.
func (n *Notifier) SetSpinner(spinner spinner.Model) {
	n.spinner = spinner
}

// Methods

// HandleNotification shows notification for delay milliseconds.
// It does nothing while another one is shown.
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

// HandleSpinnerTick advances the spinner while notifying.
func (n *Notifier) HandleSpinnerTick(msg tea.Msg) tea.Cmd {
	if !n.IsNotifying() {
		return nil
	}

	var cmd tea.Cmd
	n.spinner, cmd = n.spinner.Update(msg)
	return cmd
}
