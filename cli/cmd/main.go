package main

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	draw(app)
	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}

func draw(app *tview.Application) {
	servers := tview.NewList().ShowSecondaryText(false)
	servers.SetBorder(true).SetTitle("Servers")
	details := tview.NewTable().SetBorders(true)
	details.SetBorder(true).SetTitle("Queues")
	flex := tview.NewFlex().
		AddItem(servers, 0, 1, true).
		AddItem(details, 0, 1, false).Box

	servers.AddItem("10.0.0.1", "", 0, nil)
	servers.AddItem("10.0.0.2", "", 0, nil)
	servers.AddItem("10.0.0.3", "", 0, nil)
	servers.SetChangedFunc(func(i int, tableName string, t string, s rune) {
		details.Clear()
		adddetails(details, strconv.Itoa(i))
	})
	servers.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		switch key.Key() {
		case tcell.KeyEnter:
			return nil
		}
		return key
	})
	servers.SetDoneFunc(func() {
		app.Stop()
	})
	pages := tview.NewPages().
		AddPage("finderPage", flex, true, true)
	app.SetRoot(pages, true)
	servers.SetCurrentItem(3)
	servers.SetCurrentItem(0)
}

func addcomponents(components *tview.List, s string) {
	/*components.AddItem("Queues"+s, "", 0, nil)
	components.AddItem("Listener"+s, "", 0, nil)
	components.AddItem("Executors"+s, "", 0, nil)*/

}

func adddetails(details *tview.Table, d string) {
	/*color := tcell.ColorGreenYellow
	details.SetCell(0, 0, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 1, &tview.TableCell{Text: "Type", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 2, &tview.TableCell{Text: "Size", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 3, &tview.TableCell{Text: "Null", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
		SetCell(0, 4, &tview.TableCell{Text: "Constraint", Align: tview.AlignCenter, Color: tcell.ColorYellow})

	details.SetCell(1, 0, &tview.TableCell{Text: "columnName" + d, Color: color}).
		SetCell(1, 1, &tview.TableCell{Text: "dataType" + d, Color: color}).
		SetCell(1, 2, &tview.TableCell{Text: "sizeText" + d, Align: tview.AlignRight, Color: color}).
		SetCell(1, 3, &tview.TableCell{Text: "isNullable" + d, Align: tview.AlignRight, Color: color}).
		SetCell(1, 4, &tview.TableCell{Text: "constraintType.String" + d, Align: tview.AlignLeft, Color: color})
	*/
}
