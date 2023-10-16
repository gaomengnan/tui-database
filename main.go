package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Name       string
	Host       string
	Port       string
	Username   string
	Password   string
	connection *gorm.DB
}

var (
	mainLayout *tview.Flex
	// modal      *tview.Modal
	connections = map[string]DBConfig{}
)

func main() {
	app := tview.NewApplication()
	main := createMain(app)
	if err := app.SetRoot(main, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func createMain(app *tview.Application) *tview.Flex {
	if mainLayout != nil {
		return mainLayout
	}

	flex := tview.NewFlex()
	form := tview.NewForm().
		AddButton("New +", func() {
			form := createForm(app)
			app.SetRoot(form, true).SetFocus(form)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Edit").SetTitleAlign(tview.AlignLeft)
	form.SetButtonBackgroundColor(tcell.ColorGrey)
	flex = flex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Connected"), 0, 7, false), 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Top"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Middle (3 x height of Top)"), 0, 3, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Bottom (5 rows)"), 5, 1, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)
	mainLayout = flex
	return flex
}

// func createNewConnectionBtn(app *tview.Application) *tview.Button {
// 	btn := tview.NewButton("New+")
// 	btnStyle := tcell.Style{}
// 	btnStyle = btnStyle.Background(tcell.ColorGreen)
// 	btn.SetStyle(btnStyle).SetRect(0, 0, 22, 3)
// 	return btn
// }

// func createModal(app *tview.Application) *tview.Modal {
// 	modal := tview.NewModal().
// 		SetText("Do you want to quit the application?").
// 		AddButtons([]string{"Quit", "Cancel"}).
// 		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
// 			if buttonLabel == "Quit" {
// 				main := createMain(app)
// 				if err := app.SetRoot(main, true).EnableMouse(true).Run(); err != nil {
// 					panic(err)
// 				}
// 			}
// 		})

// 	return modal
// }

func createForm(app *tview.Application) *tview.Form {
	var (
		name     string
		host     string
		port     string
		userName string
		password string
	)
	form := tview.NewForm().
		// AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
		AddInputField("Connection Name:", "", 40, nil, func(text string) {
			name = text
		}).
		AddInputField("Host:", "", 40, nil, func(text string) {
			host = text
		}).
		AddInputField("Port", "", 40, nil, func(text string) {
			port = text
		}).
		AddInputField("Username", "", 40, nil, func(text string) {
			userName = text
		}).
		AddInputField("Password", "", 40, nil, func(text string) {
			password = text
		}).
		// AddCheckbox("Age 18+", false, nil).
		// AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", func() {
			fileName := "/tmp/host.txt"
			f, err := os.Create(fileName)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = f.Write([]byte(fmt.Sprintf("host:%s,port:%s,user:%s,pass:%s", host, port, userName, password)))
			if err != nil {
				panic(err)
			}
			connectDB(DBConfig{
				Name:     name,
				Host:     host,
				Port:     port,
				Username: userName,
				Password: password,
			})

		}).
		AddButton("Quit", func() {
			root := createMain(app)
			app.SetRoot(root, true).SetFocus(root)
		})
	form.SetBorder(true).SetTitle("New Connection").SetTitleAlign(tview.AlignLeft)
	form.SetButtonBackgroundColor(tcell.ColorGrey)
	form.SetFieldTextColor(tcell.ColorBlack)
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	return form
}

func connectDB(conf DBConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Host, conf.Port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_db, err := db.DB()
	if err != nil {
		panic(err)
	}
	if err := _db.Ping(); err != nil {
		panic(err)
	}
	conf.connection = db
	connections[conf.Name] = conf
}
