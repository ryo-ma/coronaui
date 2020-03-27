package main

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
	"github.com/ryo-ma/coronaui/lib"
	"github.com/ryo-ma/coronaui/ui"
)

var client *lib.Client
var countryPanel *ui.CountryPanel
var textPanel *ui.TextPanel
var statusPanel *ui.StatusPanel
var searchPanel *ui.SearchPanel
var loadingPanel *ui.LoadingPanel
var cursor *ui.Cursor

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)

	client, _ = lib.NewClient()
	countryPanel, _ = ui.NewCountryPanel()
	textPanel, _ = ui.NewTextPanel()
	statusPanel, _ = ui.NewStatusPanel()
	searchPanel, _ = ui.NewSearchPanel()
	loadingPanel, _ = ui.NewLoadingPanel()
	cursor = &ui.Cursor{}

	countries, _ := client.GetCountries()
	countryPanel.Countries = countries

	countryPanel.DrawView(g)
	textPanel.DrawView(g)
	statusPanel.DrawView(g)
	textPanel.DrawText(g, &countryPanel.Countries[0])
	g.SetCurrentView(countryPanel.ViewName)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(countryPanel.ViewName, 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(textPanel.ViewName, 'q', gocui.ModNone, exit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(textPanel.ViewName, 'x', gocui.ModNone, exit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(countryPanel.ViewName, 'k', gocui.ModNone, cursorMovement(-1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(countryPanel.ViewName, 'j', gocui.ModNone, cursorMovement(1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(textPanel.ViewName, 'k', gocui.ModNone, cursorMovement(-1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(textPanel.ViewName, 'j', gocui.ModNone, cursorMovement(1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(countryPanel.ViewName, 'c', gocui.ModNone, copyLineData); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding(textPanel.ViewName, 'c', gocui.ModNone, copyLineData); err != nil {
		log.Panicln(err)
	}
	// if err := g.SetKeybinding(countryPanel.ViewName, 'o', gocui.ModNone, openBrowser); err != nil {
	// 	log.Panicln(err)
	// }
	//if err := g.SetKeybinding(textPanel.ViewName, 'o', gocui.ModNone, openBrowser); err != nil {
	//	log.Panicln(err)
	//}
	if err := g.SetKeybinding("", gocui.KeyCtrlU, gocui.ModNone, cursorMovement(-5)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlD, gocui.ModNone, cursorMovement(5)); err != nil {
		log.Panicln(err)
	}
	// if err := g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, drawSearchEditor); err != nil {
	// 	log.Panicln(err)
	// }
	// if err := g.SetKeybinding(searchPanel.ViewName, gocui.KeyEnter, gocui.ModNone, searchCountry); err != nil {
	// 	log.Panicln(err)
	// }
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorMovement(-1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorMovement(1)); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

func layout(g *gocui.Gui) error {
	render(g)
	return nil
}
func render(g *gocui.Gui) {
	countryPanel.DrawView(g)
	statusPanel.DrawView(g)
	textPanel.DrawView(g)
}

func cursorMovement(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		cursor.Move(g, v, d, func(yOffset int, yCurrent int) error {
			if g.CurrentView().Name() == countryPanel.ViewName {
				if yOffset+yCurrent >= len(countryPanel.Countries) {
					return nil
				}
				textPanel.DrawText(g, &countryPanel.Countries[yOffset+yCurrent])
			}
			return nil
		})
		return nil
	}
}

func copyLineData(g *gocui.Gui, _ *gocui.View) error {
	//yOffset, yCurrent, _ := cursor.FindPosition(g, countryPanel.ViewName)
	//currentItem := countryPanel.Result.Items[yCurrent+yOffset]

	err := clipboard.WriteAll("")
	if err != nil {
		statusPanel.DrawText(g, "Failed to copy. Please copy \033[32m \033[0m")
		return nil
	}

	statusPanel.DrawText(g, "Copied successfully! \033[32m \033[0m")
	return nil
}

func openBrowser(g *gocui.Gui, _ *gocui.View) error {
	//yOffset, yCurrent, _ := cursor.FindPosition(g, countryPanel.ViewName)
	//currentItem := countryPanel.Result.Items[yCurrent+yOffset]
	//url := currentItem.GetcountryURL()
	url := ""

	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		statusPanel.DrawText(g, "Failed to open URL. Unsupported platform.")
	}
	if err != nil {
		statusPanel.DrawText(g, "Failed to open URL.")
	}
	statusPanel.DrawText(g, "Success to open URL. "+url)
	return nil

}
func drawSearchEditor(g *gocui.Gui, _ *gocui.View) error {
	err := searchPanel.DrawView(g)
	if err != nil {
		return err
	}
	return nil
}

func searchCountry(g *gocui.Gui, v *gocui.View) error {
	topic, _ := v.Line(0)
	if topic == "" {
		g.DeleteView(searchPanel.ViewName)
		g.SetCurrentView(countryPanel.ViewName)
		return nil
	}
	vr, err := g.View(countryPanel.ViewName)
	if err != nil {
		return err
	}
	//g.DeleteView(searchPanel.ViewName)
	//loadingPanel.ShowLoading(g, func() {
	//	countryPanel.Result, err = client.Searchcountry(topic)
	//	if err != nil {
	//		statusPanel.DrawText(g, "Failed to search repositories.")
	//	} else {
	//		vr.Clear()
	//		vr.Title = " Search [" + topic + "]"
	//		countryPanel.Result.Draw(vr)
	//		g.SetCurrentView(countryPanel.ViewName)
	//		if len(countryPanel.Result.Items) != 0 {
	//			textPanel.DrawText(g, &countryPanel.Result.Items[0])
	//		}
	//	}
	//})
	cursor.MoveToFirst(g, vr)
	return nil
}

func exit(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == textPanel.ViewName {
		cursor.MoveToFirst(g, v)
		v.Highlight = false
		g.SetCurrentView(countryPanel.ViewName)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
