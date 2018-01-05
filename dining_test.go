package ucscdining

import (
	"fmt"
	"time"
)

func ExampleGetMenu() {
	t, _ := time.Parse(dateFormat, "01/05/2018")
	menu, _ := PorterKresge.On(t).GetMenu()
	fmt.Printf("%s\n", menu)
	// Output:
}
