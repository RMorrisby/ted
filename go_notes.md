# Notes on Go

'go mod init <project name>' to initialise your project & create a go.mod file
'go mod tidy' will try to fetch any remote modules/packages in your imports, and add them to the go.mod file

'go build' will build your application & create an .exe
'go run <go file, not .exe>' will run your application

## Localhost

Launching a local web-application will have it running on HTTP port 8080, i.e.
http://localhost:8080/

## Websockets

https://blog.markvincze.com/programmatically-refreshing-a-browser-tab-from-a-golang-application/

## Datetime formats

Golang is -eccentric- borderline-insane. To define a datetime format, instead of the usual yyyyMMdd format-string, the format-string is this datetime, modified accordingly : `2006-01-02T15:04:05.000Z`

You will notice that in the full datetime, no digit-section shares its value with another section. So you can convert to the sane format by treating `2006 == yyyy`, `01 == MM`, and so on.

Therefore these are all valid format-strings :

```
2006-01-02 15:04:05
2006-01-02T15:04
15:04:05 2006-01-02
etc.
```
