rem set GOOS=windows
rem set GOARCH=386


go build -ldflags -H=windowsgui main.go
del /Q "1C Starter.exe"
"C:\Program Files (x86)\Resource Hacker\ResourceHacker.exe" -open main.exe -save "1C Starter.exe" -action addskip -res C:\Users\Sanek\go\src\1C_Starter\1cv8_16x16.ico  -mask ICONGROUP,MAIN,
rem ren "main.exe" "1C Starter.exe"

