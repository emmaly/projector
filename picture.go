package projector

// Freeze toggles the frame pauses state of the video output
func (p *Projector) Freeze() error {
	return p.Command("f013")
}

// Blank toggles the blank state of the screen (or just strongly (un)dims it)
func (p *Projector) Blank() error {
	return p.Command("ee13")
}

// ContrastUp raises the contrast level
func (p *Projector) ContrastUp() error {
	return p.Command("f613")
}

// ContrastDown lowers the contrast level
func (p *Projector) ContrastDown() error {
	return p.Command("f713")
}

// BrightnessUp raises the brightness level
func (p *Projector) BrightnessUp() error {
	return p.Command("f413")
}

// BrightnessDown lowers the brightness level
func (p *Projector) BrightnessDown() error {
	return p.Command("f513")
}

// SaturationUp raises the saturation level
func (p *Projector) SaturationUp() error {
	return p.Command("f213")
}

// SaturationDown lowers the saturation level
func (p *Projector) SaturationDown() error {
	return p.Command("f313")
}

// SharpnessUp raises the sharpness level
func (p *Projector) SharpnessUp() error {
	return p.Command("f813")
}

// SharpnessDown lowers the sharpness level
func (p *Projector) SharpnessDown() error {
	return p.Command("f913")
}

// ZoomIn raises the zoom level
func (p *Projector) ZoomIn() error {
	return p.Command("3914")
}

// ZoomOut lowers the zoom level
func (p *Projector) ZoomOut() error {
	return p.Command("3a14")
}

// PanUp moves the viewport to the up (when zoomed)
func (p *Projector) PanUp() error {
	return p.Command("3b14")
}

// PanDown moves the viewport to the down (when zoomed)
func (p *Projector) PanDown() error {
	return p.Command("3c14")
}

// PanLeft moves the viewport to the left (when zoomed)
func (p *Projector) PanLeft() error {
	return p.Command("3d14")
}

// PanRight moves the viewport to the right (when zoomed)
func (p *Projector) PanRight() error {
	return p.Command("3e14")
}
