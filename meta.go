package projector

// Menu toggles the menu overlay
func (p *Projector) Menu() error {
	return p.Command("1d14")
}

// OK confirms the selection
func (p *Projector) OK() error {
	return p.Command("2314")
}

// CursorUp moves the cursor up
func (p *Projector) CursorUp() error {
	return p.Command("1e14")
}

// CursorDown moves the cursor down
func (p *Projector) CursorDown() error {
	return p.Command("1f14")
}

// CursorLeft moves the cursor right
func (p *Projector) CursorLeft() error {
	return p.Command("2014")
}

// CursorRight moves the cursor right
func (p *Projector) CursorRight() error {
	return p.Command("2114")
}
