package spaceballoon

import "testing"

const phoneNum = "xxxxxxxxxxx"

var msg = []string{
	"Your 3G dongle",
	"is talking to you",
}

func TestSMS(t *testing.T) {
	dongle, err := NewDongle(phoneNum)
	if err != nil {
		t.FailNow()
	}
	defer dongle.Stop()

	if err = dongle.Test(); err != nil {
		t.FailNow()
	}

	if err = dongle.SendSMS(msg); err != nil {
		t.FailNow()
	}
}
