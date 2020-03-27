package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/ryo-ma/coronaui/lib"
	"strconv"
)

type CountryPanel struct {
	ViewName     string
	viewPosition ViewPosition
	Countries    []lib.Country
}

func NewCountryPanel() (*CountryPanel, error) {
	countryPanel := CountryPanel{
		ViewName: "country",
		viewPosition: ViewPosition{
			x0: Position{0.0, 0},
			y0: Position{0.0, 0},
			x1: Position{0.3, 2},
			y1: Position{0.9, 2},
		},
	}
	return &countryPanel, nil
}

func (countryPanel *CountryPanel) DrawView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	x0, y0, x1, y1 := countryPanel.viewPosition.GetCoordinates(maxX, maxY)
	if v, err := g.SetView(countryPanel.ViewName, x0, y0, x1, y1); err != nil {
		v.SelFgColor = gocui.ColorBlack
		v.SelBgColor = gocui.ColorGreen
		v.Highlight = true
		for _, country := range countryPanel.Countries {
			casesText := " 🦠 " + strconv.Itoa(country.Cases)
			fmt.Fprintf(v, "%-10.10s\033[32m%s\033[0m\n", casesText, country.Name)
		}
		v.Title = " Countries "
	}
	return nil
}
