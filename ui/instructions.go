package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Instructions struct {
	widget.BaseWidget
	InstructionsText string
	OnTapped         func()
}

func NewInstructions(instructionsText string, onTapped func()) *Instructions {

	instructions := &Instructions{
		InstructionsText: instructionsText,
		OnTapped:         onTapped,
	}
	instructions.ExtendBaseWidget(instructions)
	return instructions
}

func (i *Instructions) CreateRenderer() fyne.WidgetRenderer {
	con := container.NewCenter(widget.NewLabelWithStyle(i.InstructionsText, fyne.TextAlignCenter, widget.RichTextStyleHeading.TextStyle))

	return widget.NewSimpleRenderer(con)
}
func (i *Instructions) Tapped(ev *fyne.PointEvent) {

	i.OnTapped()
}