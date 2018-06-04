package phada

import (
	"testing"
)

// TestInMemorySessionStore_Get
func TestInMemorySessionStore_Get(t *testing.T) {
	sessionStore := NewInMemorySessionStore()

	_, err := sessionStore.Get("1234567890")
	if err == nil {
		t.Errorf("Expected error in sessionStore.Get() call")
	}

	session := &UssdRequestSession{
		PhoneNumber: "265888981234",
		Text:        "",
		Channel:     "384",
		SessionID:   "1234567890",
	}
	sessionStore.PutHop(session)

	_, err = sessionStore.Get("1234567890")
	if err != nil {
		t.Errorf("Expected session '1234567890' in sessionStore.Get() call")
	}
}

func write_input(session *UssdRequestSession, in string) *UssdRequestSession {
	if session.Text == "" {
		session.Text = in
		return session
	}
	session.Text = session.Text + "*" + in
	return session
}

// TestInMemorySessionStore_Put
func TestInMemorySessionStore_PutHop(t *testing.T) {
	sessionStore := NewInMemorySessionStore()
	session := &UssdRequestSession{
		PhoneNumber: "265888981234",
		Text:        "",
		Channel:     "384",
		SessionID:   "1234567890",
	}

	sessionStore.PutHop(write_input(session, "1"))
	sessionStore.PutHop(write_input(session, "1"))
	sessionStore.PutHop(write_input(session, "2"))
	sessionStore.PutHop(write_input(session, "Foo"))
	sessionStore.PutHop(write_input(session, "Foo Bar"))
	sessionStore.PutHop(write_input(session, "Baz"))
	sessionStore.PutHop(write_input(session, "Foo:Bar"))

	s, err := sessionStore.Get("1234567890")
	if err != nil {
		t.Errorf("Test Failed:%s", err)
	}
	assertEquals(t, "Foo:Bar", s.ReadIn())
	assertEquals(t, "1*1*2*Foo*Foo Bar*Baz*Foo:Bar", s.Text)
}
