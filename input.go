package projector

// InputComputer1 switches to Computer1 input
func (p *Projector) InputComputer1() error {
	return p.Command("cd13")
}

// InputComputer2 switches to Computer2 input
func (p *Projector) InputComputer2() error {
	return p.Command("ce13")
}

// InputHDMI1 switches to HDMI1 input
func (p *Projector) InputHDMI1() error {
	return p.Command("cf13")
}

// InputHDMI2 switches to HDMI2 input
func (p *Projector) InputHDMI2() error {
	return p.Command("d013")
}

// InputComposite switches to Composite/Video input
func (p *Projector) InputComposite() error {
	return p.Command("d113")
}

// InputSVideo switches to SVideo input
func (p *Projector) InputSVideo() error {
	return p.Command("d213")
}

// InputComponent switches to Component input
func (p *Projector) InputComponent() error {
	return p.Command("d313")
}

// InputUSBA switches to USB-A input
func (p *Projector) InputUSBA() error {
	return p.Command("d413")
}

// InputUSBB switches to USB-B input
func (p *Projector) InputUSBB() error {
	return p.Command("d513")
}

// InputLAN switches to LAN input
func (p *Projector) InputLAN() error {
	return p.Command("d613")
}
