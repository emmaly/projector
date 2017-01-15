package projector

// PowerOff turns off the projector
func (p *Projector) PowerOff() error {
	return p.Command("0500")
}

// PowerOn turns on the projector
func (p *Projector) PowerOn() error {
	return p.Command("0400")
}
