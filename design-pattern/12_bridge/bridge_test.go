package bridge

import "testing"

func TestCommonSMS(t *testing.T) {
	cm := NewCommonMessage(ViaSMS())
	cm.SendMessage("have a drink?", "bob") // send have a drink? to bob via SMS
}

func TestCommonEmail(t *testing.T) {
	cm := NewCommonMessage(ViaEmail())
	cm.SendMessage("have a drink?", "bob") // send have a drink? to bob via Email
}

func TestUrgencySMS(t *testing.T) {
	um := NewUrgencyMessage(ViaSMS())
	um.SendMessage("have a drink?", "bob") // send [Urgency] have a drink? to bob via SMS
}

func TestUrgencyEmail(t *testing.T) {
	um := NewUrgencyMessage(ViaEmail())
	um.SendMessage("have a drink?", "bob") // send [Urgency] have a drink? to bob via Email
}
