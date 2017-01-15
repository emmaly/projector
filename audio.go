package projector

// Mute mutes the audio
func (p *Projector) Mute() error {
	return p.Command("fc13")
}

// Unmute unmutes the audio
func (p *Projector) Unmute() error {
	return p.Command("fd13")
}

// VolumeUp increases the audio volume level
func (p *Projector) VolumeUp() error {
	return p.Command("fa13")
}

// VolumeDown decreases the audio volume level
func (p *Projector) VolumeDown() error {
	return p.Command("fb13")
}
