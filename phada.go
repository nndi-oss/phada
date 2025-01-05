package phada

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// / UssdRequestSession
// /
// / go representation of the structure of an AfricasTalking USSD call
type UssdRequestSession struct {
	PhoneNumber string `json:"phoneNumber"`
	SessionID   string `json:"sessionID"`
	Text        string `json:"text"`
	Channel     string `json:"channel"`
	// The State of the request
	State     int       `json:"-"`
	startedAt time.Time `json:"-"`
	// Number of hops for this ussd session
	hops int `json:"-"`
	// The offset in the text from the last hop
	textOffset int `json:"-"`
}

// / SessionStore
// /
// / Interface for storing session data
type SessionStore interface {
	Get(sessionID string) (*UssdRequestSession, error)
	PutHop(*UssdRequestSession) error
	Delete(sessionID string)
}

func (u *UssdRequestSession) RecordHop(text string) {
	if len(text) <= 1 {
		u.Text = text
		u.textOffset = 0
		u.hops++
		return
	}

	// add one to account for the asterisk
	u.textOffset = len(u.Text) + 1
	u.Text = text
	u.hops++
}

// RecordHopAndReadIn
//
// Records input string for the ussd session and immediately returns
// the new input string
func (u *UssdRequestSession) RecordHopAndReadIn(text string) string {
	u.RecordHop(text)
	return u.ReadIn()
}

// ReadIn
//
// Reads the last input string recorded for this session
func (u *UssdRequestSession) ReadIn() string {
	if u.textOffset > len(u.Text) {
		u.textOffset = 0
		return ""
	}
	currentHopInput := u.Text[u.textOffset:]

	return currentHopInput
}

// SetState
//
// Set the state for this Ussd session
func (u *UssdRequestSession) SetState(state int) {
	u.State = state
}

// GetHopN
//
// Get the data provided at the nth hop
func (u *UssdRequestSession) GetHopN(n int) string {
	if n == 0 {
		return ""
	}
	a := strings.Split(u.Text, "*")
	if n > len(a) {
		// TODO(zikani): errors.New(fmt.Sprintf("Cannot read hop %d, session only had %s hops", n, hops))
		return ""
	}

	hopText := a[n-1]
	return hopText
}

// Count Hops
//
// Count the number of hops (interactions) for the Ussd session
// the number of hops is based on the asterisk count so it's
// approximate
func (u *UssdRequestSession) CountHops() int {
	a := strings.Split(u.Text, "*")

	return len(a) + 1
}

// ToJSON
//
// Convert the UssdRequestSession to JSON string or empty string on error
func (u *UssdRequestSession) ToJSON() string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	}

	return string(b)
}

// ParseUssdRequestSession
//
// Parse the Request data to a UssdRequestSession if the parameters
// are present in the body
func ParseUssdRequest(req *http.Request) (*UssdRequestSession, error) {
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	return parseUrlValuesToUssdRequestSession(req.Form)
}

func parseUrlValuesToUssdRequestSession(form url.Values) (*UssdRequestSession, error) {
	UssdRequestSession := &UssdRequestSession{
		PhoneNumber: form.Get("phoneNumber"),
		SessionID:   form.Get("sessionId"),
		Text:        form.Get("text"),
		Channel:     form.Get("channel"),
		textOffset:  0,
		startedAt:   time.Now(),
	}

	err := validateUssdRequestSession(UssdRequestSession)

	return UssdRequestSession, err
}

func validateUssdRequestSession(req *UssdRequestSession) error {
	if req.PhoneNumber == "" {
		return errors.New("UssdRequestSession PhoneNumber cannot be empty")
	}

	if req.SessionID == "" {
		return errors.New("UssdRequestSession SessionID cannot be empty")
	}

	// req.Text can be empty
	return nil
}
