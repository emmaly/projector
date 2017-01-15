package projector

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// Properties are the projector properties
type Properties struct {
	Name        string
	Hostname    string
	Location    string
	AssignedTo  string
	ColorMode   string
	EcoMode     bool
	Orientation string
	Firmware    string
	Language    string
	MAC         string
	Input       string
	Resolution  string
	Refresh     string
	BulbHours   int
}

// RefreshProperties requests all properties from the projector.
// The results will be found at projector.Properties.
// This will not block.
func (p *Projector) RefreshProperties() error {
	return p.RawCommand("050005000002031e")
}

// 00060000 seems to be a declaration of capabilities
// 000c0000091513be returns "User1" which is my current color mode
// 0009000006151388 returns "On" (but absent altogether when power is off)
// 0009000006151389 returns "On" (but absent altogether when power is off)
// 000a000007151388 returns "Off" (not sure what it is; the projector is on)
// 000a000007151389 returns "Off" (not sure what it is; the projector is on)
// 000a00000715138a returns "OFF" is the EcoMode
// 00170000141513cd returns "Computer1/YPbPr1" which is the name of input "cd13"
// 00170000141513ce returns "Computer2/YPbPr2" which is the name of input "ce13"
// 000c0000091513cf returns "HDMI1" which is the name of "cf13"
// 000c0000091513d0 returns "HDMI2" which is the name of "d013"
// 000c0000091513d1 returns "Video"
// 000e00000b1513d2 returns "S-Video"
// 000e00000b151391 returns "S-Video" (again?  is it associated with the 0x0 res below?)
// 00140000111513bc returns "Front Ceiling" which is my display orientation (it's hanging on the ceiling with front projection)
// 00090000061513e2 returns "en" which I presume is the language
// 00180000151513b3 returns "00:26:47:01:f1:12" which matches the NIC MAC address
// 000e00000b1513bf returns "VS-1.21" which is the firmware level
// 001000000d1513bd returns "0 x 0 0Hz" which is ?
// 00170000141513bd returns "1920 x 1080 60Hz" which is the resolution and refresh rate
// 00180000151513b9 returns "ViewSonic-Pro8400" which is projector name field
// 001100000e1513bb returns "LivingRoom" which is the location field
// 000d00000a1513ba returns "Wilson" which is the AssignedTo field
// 001000000d1513ae returns "projector" which is the hostname
// 000e00000b1513af returns "0.0.0.0"
// 000e00000b1513b0 returns "0.0.0.0"
// 000e00000b1513b1 returns "0.0.0.0"
// 000e00000b1513b2 returns "0.0.0.0"
// 000c0000091513b6 returns "41794" which is the port number used for this communication
// 001200000f1513b4 returns "192.168.0.2" which is the Crestron target IP (but it's disabled in my config)
// 000a00000715138b returns "182" which is my bulb hours
// 000a000007150004 returns "182" which is my bulb hours

// 1513b5 returns "5" (unknown)
// 000c000009151391 returned "HDMI2" when I switched to that input.  And same for "HDMI1" which I switched to that one.
// 00170000141513bd changes to current res on input changes as well
// 00060000 re-appears with a bunch of values on each input change also, which I think is a capabilities declaration of some kind

// 001000000d1513be is set to the color mode when the color mode changes
// 000f00000c1513be is the same when changing to "Standard" as the above
// 000e00000b1513be same for "Theater" color mode
// 000a00000715138a (ecomode?) changed to "OFF" when changing color modes
// 000900000615138a kicked ON when changing to "Dark Room" color mode
// 000a00000715138a and 000900000615138a are set to ON and OFF in association with "blank" video dimming mode; see also same codes above... because it does set ECO on or off

// 001200000f1513bb updated with Location when changed
// 000d00000a1513bb updated with Location when changed
// 000f00000c1513bb updated with Location when changed
// 001000000d1513bb updated with Location when changed
// 00130000101513bb ... location
// 00180000151513b9 updated with Name when changed

// 000002 seems to indicate that the properties request has completed, though it also appears near the top (but not *at* the top)

// ParseProperty parses the (potential) property supplied from the projector
func (p *Projector) ParseProperty(r []byte) []string {
	bs := bytes.Split(r, []byte("\x03"))
	q := make([]string, len(bs))
	m := make([]string, len(bs))
	for i, v := range bs {
		vs := bytes.Split(v, []byte{byte(5)})
		m[i] = hex.EncodeToString(vs[0])
		q[i] = string(vs[0])
	}
	key := m[0]
	if len(m[0]) >= 6 {
		key = key[len(key)-6 : len(key)]
	}
	if len(q) > 1 {
		if key == "1513be" {
			if p.Properties.ColorMode != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "ColorMode",
					ChangeFrom: p.Properties.ColorMode,
					ChangeTo:   q[1],
				}
				p.Properties.ColorMode = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "15138a" {
			v := strings.ToUpper(q[1]) != "OFF"
			if p.Properties.EcoMode != v {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "EcoMode",
					ChangeFrom: p.Properties.EcoMode,
					ChangeTo:   v,
				}
				p.Properties.EcoMode = v
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513bc" {
			if p.Properties.Orientation != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Orientation",
					ChangeFrom: p.Properties.Orientation,
					ChangeTo:   q[1],
				}
				p.Properties.Orientation = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513e2" {
			if p.Properties.Language != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Language",
					ChangeFrom: p.Properties.Language,
					ChangeTo:   q[1],
				}
				p.Properties.Language = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513b3" {
			if p.Properties.MAC != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "MAC",
					ChangeFrom: p.Properties.MAC,
					ChangeTo:   q[1],
				}
				p.Properties.MAC = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513bf" {
			if p.Properties.Firmware != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Firmware",
					ChangeFrom: p.Properties.Firmware,
					ChangeTo:   q[1],
				}
				p.Properties.Firmware = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513bd" {
			s := strings.Split(strings.Replace(q[1], " x ", "x", 1), " ")
			if p.Properties.Resolution != s[0] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Resolution",
					ChangeFrom: p.Properties.Resolution,
					ChangeTo:   q[1],
				}
				p.Properties.Resolution = s[0]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
			if len(s) > 1 {
				if p.Properties.Refresh != s[1] {
					e := Event{
						EventType:  EventPropertyChange,
						Field:      "Refresh",
						ChangeFrom: p.Properties.Refresh,
						ChangeTo:   q[1],
					}
					p.Properties.Refresh = s[1]
					e.Properties = *p.Properties
					p.emitEvent(e)
				}
			} else {
				if len(p.Properties.Refresh) > 0 {
					e := Event{
						EventType:  EventPropertyChange,
						Field:      "Refresh",
						ChangeFrom: p.Properties.Refresh,
						ChangeTo:   q[1],
					}
					p.Properties.Refresh = ""
					e.Properties = *p.Properties
					p.emitEvent(e)
				}
			}
		}
		if key == "1513b9" {
			if p.Properties.Name != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Name",
					ChangeFrom: p.Properties.Name,
					ChangeTo:   q[1],
				}
				p.Properties.Name = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513bb" {
			if p.Properties.Location != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Location",
					ChangeFrom: p.Properties.Location,
					ChangeTo:   q[1],
				}
				p.Properties.Location = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513ba" {
			if p.Properties.AssignedTo != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "AssignedTo",
					ChangeFrom: p.Properties.AssignedTo,
					ChangeTo:   q[1],
				}
				p.Properties.AssignedTo = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513e2" {
			if p.Properties.Language != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Language",
					ChangeFrom: p.Properties.Language,
					ChangeTo:   q[1],
				}
				p.Properties.Language = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "1513ae" {
			if p.Properties.Hostname != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Hostname",
					ChangeFrom: p.Properties.Hostname,
					ChangeTo:   q[1],
				}
				p.Properties.Hostname = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "151391" {
			if p.Properties.Input != q[1] {
				e := Event{
					EventType:  EventPropertyChange,
					Field:      "Input",
					ChangeFrom: p.Properties.Input,
					ChangeTo:   q[1],
				}
				p.Properties.Input = q[1]
				e.Properties = *p.Properties
				p.emitEvent(e)
			}
		}
		if key == "15138b" || key == "150004" {
			i, err := strconv.Atoi(q[1])
			if err == nil {
				if p.Properties.BulbHours != i {
					e := Event{
						EventType:  EventPropertyChange,
						Field:      "BulbHours",
						ChangeFrom: p.Properties.BulbHours,
						ChangeTo:   i,
					}
					p.Properties.BulbHours = i
					e.Properties = *p.Properties
					p.emitEvent(e)
				}
			}
		}
		if p.DebugOutput && key != "1513be" && key != "1513bc" {
			fmt.Printf("[%s] [%s]\n\n", m[0], q[1])
		}
	}
	return nil
}
