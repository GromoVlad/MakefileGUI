package service

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func Bootstrap() {
	/** Подгружаем данные из .env */
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки переменных из .env: %s", err.Error())
	}
	_, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}

	createFile("group.txt")
	createFile("password.txt")
}

func createFile(filename string) {
	/** Проверка существования файла */
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		/** Файл не существует, создаем его */
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Ошибка при создании файла:", err)
			return
		}
		/** Закрываем файл после завершения работы с ним */
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func CreateWindowAndApp() (fyne.Window, fyne.App) {
	myApp := app.New()
	window := myApp.NewWindow("Command panel")
	window.Resize(fyne.Size{Width: 1000, Height: 1000})

	return window, myApp
}

func CreateLabel() *widget.Label {
	label := widget.NewLabel("Статус выполнения команд")
	label.TextStyle = fyne.TextStyle{Bold: true}
	return UpdateLabel(label)
}

func UpdateLabel(clock *widget.Label) *widget.Label {
	readFile, err := os.Open("./log.txt")
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		isInfo := (strings.Contains(fileScanner.Text(), "[START]") &&
			!strings.Contains(fileScanner.Text(), "[END]")) ||
			(!strings.Contains(fileScanner.Text(), "[START]") &&
				strings.Contains(fileScanner.Text(), "[END]"))

		if isInfo {
			clock.Wrapping = fyne.TextWrapOff
			clock.TextStyle = fyne.TextStyle{Bold: true}
			clock.SetText(fileScanner.Text())
		}

	}
	err = readFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	return clock
}

func StartingListeners(label *widget.Label) {
	go func() {
		for range time.Tick(time.Second) {
			label = UpdateLabel(label)
		}
	}()
}

func AddContent(label *widget.Label, window fyne.Window, app fyne.App) fyne.Window {
	gui := container.NewHBox(getButtonBox(label, app), getTextBox())
	gui.Resize(fyne.Size{Width: 1000, Height: 1000})
	window.SetContent(gui)

	return window
}

func getButtonBox(label *widget.Label, app fyne.App) *fyne.Container {
	myLog, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := getLines()
	groups := getGroups()
	widgets := make(map[string]*widget.Button, 10)

	for _, line := range lines {
		if strings.Contains(line, "_gui") {
			originalLine := strings.ReplaceAll(line, ":", "")
			header := strings.ReplaceAll(strings.SplitAfterN(line, "_", 2)[0], "_", "")

			widgets[header] = widget.NewButton(header, func() {
				cmd := exec.Command("make", originalLine)
				os.Stdout = myLog
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				_ = cmd.Run()
			})
		}
	}

	keys := make([]string, 0, len(widgets))
	for k := range widgets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for groupKey, groupValue := range groups {
		for _, index := range keys {
			if strings.Contains(index, groupKey) {
				groupValue.Add(widgets[index])
			}
		}
	}

	buttonBox := container.NewVBox(label)
	passwordButton := createPasswordButton(app)
	buttonBox.Add(passwordButton)

	for key, value := range groups {
		buttonBox.Add(container.NewAppTabs(container.NewTabItem(key, value)))
	}

	buttonBox.Resize(fyne.Size{Width: 600, Height: 1000})

	return buttonBox
}

func createPasswordButton(app fyne.App) *container.AppTabs {
	dialogInfo := widget.NewButton("Пароль/Логин", func() {
		readFile2, _ := os.Open("./password.txt")
		readFile2Data, _ := io.ReadAll(readFile2)
		secondWindow := app.NewWindow("Login and Password")
		loginAndPassword := widget.NewMultiLineEntry()
		loginAndPassword.Text = string(readFile2Data)
		secondWindow.SetContent(loginAndPassword)
		secondWindow.Resize(fyne.Size{Width: 600, Height: 600})
		secondWindow.Show()
	})

	containerDialogBox := container.NewVBox()
	containerDialogBox.Add(dialogInfo)

	return container.NewAppTabs(
		container.NewTabItem("Passwords", containerDialogBox),
	)
}

func getTextBox() *fyne.Container {
	indent := "                                                                                  "
	text := indent + "Справочная информация" + indent
	tabs := container.NewAppTabs(container.NewTabItem(text, getTextInput()))
	card := widget.NewCard("", "", tabs)
	textBox := container.NewHBox(card)

	return textBox
}

func getTextInput() *widget.Entry {
	readFile, _ := os.Open(os.Getenv("TEXT_INPUT"))
	data, _ := io.ReadAll(readFile)
	textInput := widget.NewMultiLineEntry()
	textInput.Text = string(data)
	err := readFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	return textInput
}

func getLines() []string {
	readFile, err := os.Open(os.Getenv("MAKEFILE"))
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	err = readFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	return lines
}

func getGroups() map[string]*fyne.Container {
	groupFile, err := os.Open(os.Getenv("GROUP_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	groupScanner := bufio.NewScanner(groupFile)
	groupScanner.Split(bufio.ScanLines)
	groups := make(map[string]*fyne.Container)
	for groupScanner.Scan() {
		groups[groupScanner.Text()] = container.NewVBox()
	}
	err = groupFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	return groups
}
