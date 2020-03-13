package web

import (
	"github.com/soyum2222/slog"
	"gonat/server/config"
	"gonat/server/conn"
	"html/template"
	"net/http"
)

func Run() {

	go func() {
		if config.CFG.UIP != "" {
			http.HandleFunc("/", Show)
			err := http.ListenAndServe(":"+config.CFG.UIP, nil)
			if err != nil {
				panic(err)
			}
		}
	}()

}

type Table struct {
	Addr string
	Port string
}

type CT struct {
	ClientTables []Table
}

var temp = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>gonat</title>
</head>
<body>
<table border="1">
        <tr>
            <th>addr</th>
            <th>port</th>
        </tr>
        {{range .ClientTables}}
        <tr>
            <td>{{.Addr}}</td>
            <td>{{.Port}}</td>
        </tr>
        {{end}}
    </table>
</body>
</html>
`

func Show(writer http.ResponseWriter, request *http.Request) {

	var ct CT

	conn.ClientTabel.Range(func(key, value interface{}) bool {
		ct.ClientTables = append(ct.ClientTables, Table{key.(string), value.(string)})
		return true
	})

	index, err := template.New("index").Parse(temp)
	if err != nil {
		slog.Logger.Error(err)
		writer.WriteHeader(500)
	}

	err = index.Execute(writer, ct)
	if err != nil {
		slog.Logger.Error(err)
	}
}
