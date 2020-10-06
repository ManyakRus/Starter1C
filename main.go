////go:generate goversioninfo -icon=C:\Users\Sanek\go\src\1C_Starter\1cv8_16x16.png

package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/gen2brain/go-unarr"
	"gopkg.in/ini.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"sort"
	"strconv"
	"strings"
)

var Bases map[string]Base
var LabelBasename *widget.Label
var WidgetRadio *widget.Radio
var Application1 fyne.App
var Window1 fyne.Window
var FileNameStart1C1 string
var FileNameStart1C2 string
var CatalogNameUpgrade1C string
var URLUpgrade1C string
var MassBases []string
var CatalogStarter1C string
var StringDownload = "Download"
var ParametersInstall1C []string
var ShowStandart1CBases bool
var MassIBsases []string

//type enterEntry struct {
//	widget.Entry
//}

type Base struct {
	Name string
	//File string
	//Srvr string
	//Ref string
	Path             string
	ConnectionString string
}

func init() {
	Bases = make(map[string]Base)
	ParametersInstall1C = make([]string, 0)
	MassIBsases = make([]string, 0)
}

func main() {
	FillCatalogStarter1C()

	ReadINIFile()

	CreateGUI()
}

func FillCatalogStarter1C() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	CatalogStarter1C = usr.HomeDir
	CatalogStarter1C = CatalogStarter1C + `\AppData\Roaming\Starter1C`

	if _, err := os.Stat(CatalogStarter1C); err != nil {
		os.MkdirAll(CatalogStarter1C, os.ModePerm)
	}

	CatalogDownload := CatalogStarter1C + `\` + StringDownload
	if _, err := os.Stat(CatalogDownload); err != nil {
		os.MkdirAll(CatalogDownload, os.ModePerm)
	}

	CatalogStarter1C = CatalogStarter1C + `\`

}

func FileExists(fileName string) bool {

	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false

	//if _, err := os.Stat(fileName); err != nil {
	//	return true
	//}
	//
	//return false

}

func ReadINIFile() {
	DirectoryApp := FindDirectoryApp()
	fileName := DirectoryApp + `1C starter.ini`
	cfg, err := ini.Load(fileName)
	if err == nil {
		FileNameStart1C1 = cfg.Section("Main").Key("FileNameStart1C1").String()
		FileNameStart1C2 = cfg.Section("Main").Key("FileNameStart1C2").String()
		CatalogNameUpgrade1C = cfg.Section("Main").Key("CatalogNameUpgrade1C").String()
		URLUpgrade1C = cfg.Section("Main").Key("URLUpgrade1C").String()

		//ParametersInstall1C
		for i := 1; i <= 20; i++ {
			Param := cfg.Section("Main").Key("ParametrInstall1C" + strconv.Itoa(i)).String()
			if Param != "" {
				ParametersInstall1C = append(ParametersInstall1C, Param)
			}
		}

		//add ibases.v8i
		ShowStandart1CBases, _ = cfg.Section("Main").Key("ShowStandart1CBases").Bool()
		if ShowStandart1CBases == true {
			HomeDir, _ := os.UserHomeDir()
			Filename1C := HomeDir + `\AppData\Roaming\1C\1CEStart\ibases.v8i`
			MassIBsases = append(MassIBsases, Filename1C)
		}

		//ParametersInstall1C
		for i := 1; i <= 20; i++ {
			Param := cfg.Section("Main").Key("IbasesV8i" + strconv.Itoa(i)).String()
			if Param != "" {
				MassIBsases = append(MassIBsases, Param)
			}
		}

		//os.Exit(1)
	} else {
		println("Fail to read file: " + fileName + " Error: " + err.Error())
	}

	if FileNameStart1C1 == "" {
		FileNameStart1C1 = `C:\Program Files (x86)\1cv8\common\1cestart.exe`
	}

	if FileNameStart1C2 == "" {
		FileNameStart1C2 = `C:\Program Files\1cv8\common\1cestart.exe`
	}

}

func FindDirectoryApp() string {
	DirectoryApp, _ := os.Getwd()
	DirectoryApp = DirectoryApp + `\`

	return DirectoryApp
}

func CreateGUI() {
	Application1 := app.New()

	DirectoryApp := FindDirectoryApp()

	FilenameDownload := DirectoryApp + "download.png"
	Filename1C := DirectoryApp + "1cv8_16x16.png"

	Resource1C, _ := fyne.LoadResourceFromPath(Filename1C)
	//canvas.
	Application1.SetIcon(Resource1C)

	Window1 := Application1.NewWindow("1С Предприятие стартер")
	//w.SetContent(widget.NewLabel("Hello Fyne!"))
	Window1.CenterOnScreen()
	Window1.Canvas().SetOnTypedKey(KeyEventWindow)

	//entry := newEnterEntry()
	//Window1.SetContent(entry)

	LabelBasename = widget.NewLabel("путь")
	VBoxBases := CreateVBoxBases()

	ButtonDownload := widget.NewButton("Загрузить новую версию платформы 1С", DownloadClick)
	ResourceDownload, _ := fyne.LoadResourceFromPath(FilenameDownload)
	ButtonDownload.SetIcon(ResourceDownload)
	if CatalogNameUpgrade1C == "" && URLUpgrade1C == "" {
		ButtonDownload.Hidden = true
	}

	ButtonOpen := widget.NewButton("Открыть", OpenClick)
	ButtonOpen.SetIcon(Resource1C)

	VBoxButtons := widget.NewVBox(
		ButtonOpen,
		ButtonDownload)

	LayoutAll := layout.NewBorderLayout(nil, LabelBasename, nil, VBoxButtons)

	//AllElements2 := fyne.NewContainerWithLayout(LayoutAll, VBoxBases, LabelBasename, VBoxButtons)
	//Window1.SetContent(AllElements2)

	ScrollContainer := widget.NewVScrollContainer(VBoxBases)
	ScrollContainer.SetMinSize(fyne.Size{100, 300})
	AllElements2 := fyne.NewContainerWithLayout(LayoutAll, ScrollContainer, LabelBasename, VBoxButtons)
	Window1.SetContent(AllElements2)

	//(widget.NewButton("Загрузить новую версию", DownloadClick))

	Window1.ShowAndRun()

}

func OpenClick() {
	TextRadio := WidgetRadio.Selected
	Base1 := Bases[TextRadio]

	arg := make([]string, 0)
	arg = append(arg, "ENTERPRISE")
	////parametr := ``
	//if Base1.File != "" {
	//	//parametr = parametr + ` /F "` + Base1.File + `"`
	//	arg = append(arg, "/F")
	//	arg = append(arg, Base1.File)
	//}
	//if Base1.Srvr != "" {
	//	//parametr = parametr + ` /S "` + Base1.Srvr + `"\"` + Base1.Ref + `"`
	//	arg = append(arg, "/S")
	//	var s string
	//	s = Base1.Srvr + `\` + Base1.Ref
	//	arg = append(arg, s)
	//}

	arg = append(arg, "/IBConnectionString")

	ConnectionString := Base1.ConnectionString
	ConnectionString = strings.ReplaceAll(ConnectionString, `"`, "'")
	arg = append(arg, ConnectionString)

	//parametrCMD := ""
	Filename1C := FileNameStart1C1
	if FileExists(Filename1C) == false {
		Filename1C = FileNameStart1C2
	}

	if FileExists(Filename1C) == false {
		return
	}

	defer os.Exit(0)

	cmnd := exec.Command(Filename1C, arg...)

	cmnd.Run()
	//if err != nil {
	//	println("Error: " + err.Error())
	//}
	cmnd = nil

	Window1.Close()
	Application1.Quit()
}

func DownloadClick() {

	if CatalogStarter1C == "" {
		return
	}

	ClearDir(CatalogStarter1C + StringDownload)

	if CatalogNameUpgrade1C != "" {
		if FileExists(CatalogNameUpgrade1C) {
			Install1CFromCatalog()
			return
		}
	}

	if URLUpgrade1C != "" {
		Install1CFromURL()
	}

}

func Install1CFromURL() {
	//CopyDir(CatalogNameUpgrade1C, CatalogStarter1C + StringDownload + `\`)

	//download 1c installer
	var CatalogDownload string
	CatalogDownload = CatalogStarter1C + StringDownload + `\`
	FilenameRar := CatalogDownload + `install.rar`
	err := DownloadFile(FilenameRar, URLUpgrade1C)
	if err != nil {
		return
	}

	//unpack .rar
	a, err := unarr.NewArchive(FilenameRar)
	if err != nil {
		return
	}
	defer a.Close()

	_, err = a.Extract(CatalogDownload)
	if err != nil {
		return
	}

	//
	Install1C()
}

func Install1CFromCatalog() {
	CopyDir(CatalogNameUpgrade1C, CatalogStarter1C+StringDownload+`\`)

	Install1C()
}

func Install1C() {
	FileName := CatalogStarter1C + StringDownload + `\1CEnterprise 8.msi`
	if FileExists(FileName) {
		//var MassParams []string
		MassParams := make([]string, 0)
		MassParams = append(MassParams, "/C")
		MassParams = append(MassParams, `msiexec`)
		MassParams = append(MassParams, `/a`)
		MassParams = append(MassParams, FileName)
		MassParams = append(MassParams, ParametersInstall1C...)

		cmnd := exec.Command("cmd", MassParams...)
		cmnd.Run()
		//if err != nil {
		//	println(err.Error())
		//}
		cmnd = nil
		return
	}

	FileName = CatalogStarter1C + StringDownload + `\setup.exe`
	if FileExists(FileName) {
		cmnd := exec.Command(FileName)
		cmnd.Run()
		//if err != nil {
		//	println(err.Error())
		//}
		cmnd = nil
	}

}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// File copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

//Copy a directory recursively
// Dir copies a whole directory recursively
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func ClearDir(dir string) error {
	names, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, entery := range names {
		os.RemoveAll(path.Join([]string{dir, entery.Name()}...))
	}
	return nil
}

func AddSlash(s string) string {
	Symbol := s[len(s)-1:]
	if Symbol == `\` {
		return s
	} else {
		return s + `\`
	}

}

func RadioClick(TextRadion string) {
	Base1 := Bases[TextRadion]
	s := Base1.Path
	//if Base1.File != "" {
	//	s = s + `File="` + Base1.File + `";`
	//}
	//if Base1.Ref != "" {
	//	s = s + `File="` + Base1.Ref + `";`
	//}
	//if Base1.Srvr != "" {
	//	s = s + `File="` + Base1.Srvr + `";`
	//}
	LabelBasename.SetText(s)

}

func FillBases1(Filename1C string) {
	//HomeDir, _ := os.UserHomeDir()
	//Filename1C := HomeDir + `\AppData\Roaming\1C\1CEStart\ibases.v8i`
	content, _ := ioutil.ReadFile(Filename1C)
	text := string(content)
	arr := strings.Split(text, "[")
	for _, arr1 := range arr {
		//Name := StrBetween(arr1, "[", "]")
		Name := ""
		pos1 := strings.Index(arr1, "]")
		if pos1 > 0 {
			Name = arr1[0:pos1]
		}

		if Name == "" {
			continue
		}

		File := StrBetween(arr1, `File="`, `"`)
		Srvr := StrBetween(arr1, `Srvr="`, `"`)
		Ref := StrBetween(arr1, `Ref="`, `"`)
		Ws := StrBetween(arr1, `ws="`, `"`)
		Path := ""
		if File != "" {
			Path = File
		} else if Srvr != "" {
			Path = Srvr
		} else if Ref != "" {
			Path = Ref
		} else if Ws != "" {
			Path = Ws
		}
		ConnectionString := StrBetween(arr1, `Connect=`, "\r")

		Base1 := Base{Name, Path, ConnectionString}
		Bases[Name] = Base1

		MassBases = append(MassBases, Name)
	}

}

func CreateVBoxBases() fyne.CanvasObject {

	MassBases = make([]string, 0)
	//заполним все файлы .v8i
	for _, Filename := range MassIBsases {
		if FileExists(Filename) == true {
			FillBases1(Filename)
			//} else if FileExists(CatalogStarter1C + Filename) {
			//	FillBases1(CatalogStarter1C + Filename)
		}
	}

	//
	sort.Strings(MassBases)

	//MassBases = append(MassBases, "База1", "База2")
	WidgetRadio = widget.NewRadio(MassBases, RadioClick)
	if len(MassBases) > 0 {
		WidgetRadio.Selected = MassBases[0]
		RadioClick(MassBases[0])
	}
	VBoxBases := widget.NewVBox(WidgetRadio)

	return VBoxBases
}

// GetStringInBetween Returns empty string if no start string found
func StrBetween(str string, start string, end string) (result string) {
	pos1 := strings.Index(str, start)
	if pos1 == -1 {
		return
	}
	pos1 += len(start)

	//var str2 string
	str2 := str[pos1:]

	pos2 := strings.Index(str2, end)
	if pos2 == -1 {
		return
	}
	return str2[0:pos2]
}

func KeyEventWindow(k *fyne.KeyEvent) {
	//println(k)
	if len(MassBases) == 0 {
		return
	}

	BaseName := WidgetRadio.Selected
	KeyName := k.Name
	switch KeyName {
	case "Down":
		{
			BaseName = NextBaseName(BaseName)
			WidgetRadio.Selected = BaseName
			WidgetRadio.Refresh()
			RadioClick(BaseName)

		}
	case "Up":
		{
			BaseName = PreviousBaseName(BaseName)
			WidgetRadio.Selected = BaseName
			WidgetRadio.Refresh()
			RadioClick(BaseName)
		}
	case "Home":
		{
			BaseName = MassBases[0]
			WidgetRadio.Selected = BaseName
			WidgetRadio.Refresh()
			RadioClick(BaseName)
		}
	case "End":
		{
			Count := len(MassBases)
			BaseName = MassBases[Count-1]
			WidgetRadio.Selected = BaseName
			WidgetRadio.Refresh()
			RadioClick(BaseName)
		}
	case "Return":
		OpenClick() //Enter
	}
}

func NextBaseName(s string) string {
	i := 0
	for _, Mass1 := range MassBases {
		if Mass1 == s {
			break
		}
		i++
	}

	Count := len(MassBases) - 1
	if i == Count {
		Base1 := MassBases[0]
		return Base1
	} else {
		Base1 := MassBases[i+1]
		return Base1
	}
}

func PreviousBaseName(s string) string {
	i := 0
	for _, Mass1 := range MassBases {
		if Mass1 == s {
			break
		}
		i++
	}

	Count := len(MassBases) - 1
	if i == 0 {
		Base1 := MassBases[Count]
		return Base1
	} else {
		Base1 := MassBases[i-1]
		return Base1
	}
}
