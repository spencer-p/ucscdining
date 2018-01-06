# UCSC Dining Go Library

[![GoDoc](https://godoc.org/github.com/spencer-p/ucscdining?status.svg)](https://godoc.org/github.com/spencer-p/ucscdining)

Package ucscdining implements a wrapper to retrieve UCSC dining hall menus on
the web.

This library is reverse engineered from observing UCSC dining's nutrition site,
available at http://nutrition.sa.ucsc.edu/.

The menus are exposed fluently. To get the current menu at Porter:

```go
menu, err := PorterKresge.GetMenu()
```

Or to get the menu on some other date:

```go
t := time.Parse("01/02/2006", "01/05/2018")
menu, err := CollegesNineTen.On(t).GetMenu()
```

The library does not yet parse the result from UCSC.
