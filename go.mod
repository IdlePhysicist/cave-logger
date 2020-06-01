module github.com/idlephysicist/cave-logger

go 1.13

replace github.com/rivo/tview => ../tview

require (
	github.com/bvinc/go-sqlite-lite v0.6.1
	github.com/gdamore/tcell v1.3.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/rivo/tview v0.0.0-20190829161255-f8bc69b90341
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.2
)
