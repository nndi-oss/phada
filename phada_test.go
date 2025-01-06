package phada

import (
	"net/url"
	"testing"
)

// TestValidateUssdRequest_PhoneNumber
func TestValidateUssdRequest_PhoneNumber(t *testing.T) {
	formData := url.Values{
		"phoneNumber": []string{""},
		"text":        []string{"ABCDE"},
		"sessionId":   []string{"jjjg"},
		"channel":     []string{"hhgh"},
	}

	_, err := parseUrlValuesToUssdRequestSession(formData)

	if err != nil && err.Error() != "UssdRequestSession PhoneNumber cannot be empty" {
		t.Errorf("Expected error message")
	}
}

// TestValidateUssdRequest_SessionID
func TestValidateUssdRequest_SessionID(t *testing.T) {
	formData := url.Values{
		"phoneNumber": []string{"265888123456"},
		"text":        []string{"ABCDE"},
		"sessionId":   []string{""},
		"channel":     []string{"xyz"},
	}

	_, err := parseUrlValuesToUssdRequestSession(formData)

	if err != nil && err.Error() != "UssdRequestSession SessionID cannot be empty" {
		t.Errorf("Expected error message")
	}
}

// TestUssdRequest_CountHops
func TestUssdRequest_CountHops(t *testing.T) {
	u := &UssdRequestSession{
		PhoneNumber: "265888981234",
		Text:        "1*1*1*1*1",
		Channel:     "384",
		SessionID:   "1234567890",
	}
	want := 6
	got := u.CountHops()
	if got != want {
		t.Errorf("Failed to count the number of hops. Got: %d, Want: %d", got, want)
	}
}

func TestUssdTextProcessor_RecordHop(t *testing.T) {
	u := &UssdRequestSession{
		PhoneNumber: "265888981234",
		Text:        "",
		Channel:     "384",
		SessionID:   "1234567890",
	}

	assertEquals(t, "", u.RecordHopAndReadIn(""))
	assertEquals(t, "1", u.RecordHopAndReadIn("1"))
	assertEquals(t, "1", u.RecordHopAndReadIn("1*1"))
	assertEquals(t, "2", u.RecordHopAndReadIn("1*1*2"))
	assertEquals(t, "Foo", u.RecordHopAndReadIn("1*1*2*Foo"))
	assertEquals(t, "Foo Bar", u.RecordHopAndReadIn("1*1*2*Foo*Foo Bar"))
	assertEquals(t, "Baz", u.RecordHopAndReadIn("1*1*2*Foo*Foo Bar*Baz"))
	assertEquals(t, "Foo:Bar", u.RecordHopAndReadIn("1*1*2*Foo*Foo Bar*Baz*Foo:Bar"))
	// repeating the request results in empty input response
	assertEquals(t, "", u.RecordHopAndReadIn("1*1*2*Foo*Foo Bar*Baz*Foo:Bar"))
	// but we can add more data to the request and read it in
	assertEquals(t, "Phada", u.RecordHopAndReadIn("1*1*2*Foo*Foo Bar*Baz*Foo:Bar*Phada"))
}

// TestUssdTextProcessor_GetHopN
func TestUssdTextProcessor_GetHopN(t *testing.T) {
	u := &UssdRequestSession{
		PhoneNumber: "265888981234",
		Text:        "",
		Channel:     "384",
		SessionID:   "1234567890",
	}

	assertEquals(t, "", u.GetHopN(0))
	// trying to read data for a hop that hasn't happened yet, get empty string
	assertEquals(t, "", u.GetHopN(1))

	u.RecordHopAndReadIn("1")
	assertEquals(t, "1", u.GetHopN(1))

	u.RecordHopAndReadIn("1*1")
	u.RecordHopAndReadIn("1*1*2")
	u.RecordHopAndReadIn("1*1*2*Foo")
	u.RecordHopAndReadIn("1*1*2*Foo*Foo Bar")

	assertEquals(t, "Foo Bar", u.GetHopN(5))
	// this duplication doesn't register as a hop
	u.RecordHopAndReadIn("1*1*2*Foo*Foo Bar")
	u.RecordHopAndReadIn("1*1*2*Foo*Foo Bar*Baz*Foo:Bar*Phada")
	assertEquals(t, "Baz", u.GetHopN(6))
	assertEquals(t, "Foo:Bar", u.GetHopN(7))
	assertEquals(t, "Phada", u.GetHopN(8))
}

// assertEquals utility function for asserting two strings are equal
func assertEquals(t *testing.T, want string, got string) {
	if want != got {
		t.Errorf("Test failed. Want: %s\tGot: %s", want, got)
	}
}
