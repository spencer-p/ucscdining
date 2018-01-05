/*
Package ucsc-dining implements a wrapper to retrieve UCSC dining hall menus on
the web.

This library is reverse engineered from observing how UCSC's dining hall
websites interact with their backend. There are JavaScript examples at
http://eat.ucsc.edu/scripts/.

The menus are exposed fluently. To get the current menu at Porter:

	menu, err := PorterKresge.GetMenu()

Or to get the menu on some other date:

	t := time.Parse("01/02/2006", "01/05/2018")
	menu, err := CollegesNineTen.On(t).GetMenu()

*/
package ucscdining

import (
	"bytes"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// ID numbers for colleges
	CowellStevenson   DiningHall = 5
	CrownMerrill      DiningHall = 20
	PorterKresge      DiningHall = 25
	CollegeEightOakes DiningHall = 30
	CollegesNineTen   DiningHall = 40

	// See the javascript samples
	apiurl = "http://eat.ucsc.edu/menu.php"

	// Format for UCSC's api
	dateFormat = "01/02/2006"
)

// DiningHall represents a dining hall whose menu can be fetched.
type DiningHall int

// Request object to deliver to UCSC dining's website. This struct should not
// necessarily be created manually. See DiningHall.On().
type RequestPayload struct {
	// Date must be in the format "MM/DD/YYYY"
	Date string `json:"serve_date"`

	// LocationID is from the idByLocation lookup table.
	LocationID DiningHall `json:"location_num"`

	// All UCSC dining halls should have true for this field.
	UseDiningDB bool `json:"foodproDB"`
}

// On creates a request payload that encodes the given dining hall's menu
// *on* the given date. Typically you call GetMenu() on this object
// immediately after creation.
func (dh DiningHall) On(t time.Time) RequestPayload {
	return RequestPayload{
		Date:        t.Format(dateFormat),
		LocationID:  dh,
		UseDiningDB: true,
	}
}

// GetMenu for a request payload returns the menu described by the payload.
// Use the On() method of DiningHall to get a payload.
func (r RequestPayload) GetMenu() ([]byte, error) {
	// Create payload
	// I would prefer JSON, but the server doesn't seem to like that.
	payload := fmt.Sprintf("serve_date=\"%s\"&location_num=%d&foodproDB=%t",
		r.Date,
		r.LocationID,
		r.UseDiningDB,
	)

	// Send request
	resp, err := http.Post(
		apiurl,
		"application/x-www-form-urlencoded",
		bytes.NewBufferString(payload),
	)
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
