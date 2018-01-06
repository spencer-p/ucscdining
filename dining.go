/*
Package ucscdining implements a wrapper to retrieve UCSC dining hall menus on
the web.

This library is reverse engineered from observing UCSC dining's nutrition site,
available at http://nutrition.sa.ucsc.edu/.

The menus are exposed fluently. To get the current menu at Porter:

	menu, err := PorterKresge.GetMenu()

Or to get the menu on some other date:

	t := time.Parse("01/02/2006", "01/05/2018")
	menu, err := CollegesNineTen.On(t).GetMenu()

The library does not yet parse the result from UCSC.
*/
package ucscdining

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	// ID numbers for colleges
	CowellStevenson   = DiningHall{"5", "Cowell Stevenson Dining Hall"}
	CrownMerrill      = DiningHall{"20", "Crown Merril Dining Hall"}
	PorterKresge      = DiningHall{"25", "Porter Kresge Dining Hall"}
	RachelCarsonOakes = DiningHall{"30", "Rachel Carson Oakes Dining Hall"}
	CollegesNineTen   = DiningHall{"40", "Colleges Nine & Ten Dining Hall"}
)

const (
	// See the links on http://nutrition.sa.ucsc.edu/
	apiurl = "http://nutrition.sa.ucsc.edu/menuSamp.asp"

	// UCSC expects MM/DD/YYYY. Format string for time.Format.
	dateFormat = "01/02/2006"
)

// DiningHall represents a dining hall whose menu can be fetched.
type DiningHall struct {
	ID   string
	Name string
}

// RequestPayload is an object to deliver to UCSC dining's website. This struct
// should not necessarily be created manually. See DiningHall.On().
type RequestPayload struct {
	URLArgs url.Values
}

// On creates a request payload that encodes the given dining hall's menu
// *on* the given date. Typically you call GetMenu() on this object
// immediately after creation.
func (dh DiningHall) On(t time.Time) RequestPayload {
	return RequestPayload{
		// See the links on http://nutrition.sa.ucsc.edu/ for how these variable
		// names are found.
		URLArgs: url.Values{
			"locationNum":  []string{dh.ID},
			"locationName": []string{dh.Name},
			"dtdate":       []string{t.Format(dateFormat)},
			// Required to get breakfast formatted properly
			"myaction": []string{"read"},
			// sName is "UC Santa Cruz Dining" officially, but never checked
			"sName": nil,
			// naFlag is "1" officially, but also seems to not matter
			"naFlag": nil,
		},
	}
}

// GetMenu for a request payload returns the menu described by the payload.
// Use the On() method of DiningHall to get a payload.
// Returns the contents of the web page as a byte array.
func (r RequestPayload) GetMenu() ([]byte, error) {
	// Send request
	resp, err := http.Get(apiurl + "?" + r.URLArgs.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// GetMenu returns the dining hall's menu for the current calendar day.
func (dh DiningHall) GetMenu() ([]byte, error) {
	return dh.On(time.Now()).GetMenu()
}
