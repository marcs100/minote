package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Define a function that animates the button size
func animateButton(button *widget.Button) {
	originalSize := button.Size()
	fmt.Println(originalSize)
	targetSize := fyne.NewSize(originalSize.Width*1.05, originalSize.Height*1.05)

	// Create a new animation
	animation := fyne.NewAnimation(time.Millisecond*200, func(val float32) {
		// Calculate the new size based on the animation progress (val)
		newWidth := originalSize.Width + (targetSize.Width-originalSize.Width)*val
		newHeight := originalSize.Height + (targetSize.Height-originalSize.Height)*val

		// Update the button's size
		button.Resize(fyne.NewSize(newWidth, newHeight))
		button.Refresh() // Redraw the button
		fmt.Println(fyne.NewSize(newWidth, newHeight))
	})
	animation.AutoReverse = true
	// Start the animation
	animation.Start()
}
