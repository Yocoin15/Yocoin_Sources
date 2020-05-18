// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package http

import (
	"html/template"
	"path"

	"github.com/Yocoin15/Yocoin_Sources/swarm/api"
)

type htmlListData struct {
	URI  *api.URI
	List *api.ManifestList
}

var htmlListTemplate = template.Must(template.New("html-list").Funcs(template.FuncMap{"basename": path.Base}).Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Swarm index of {{ .URI }}</title>
</head>

<body>
  <h1>Swarm index of {{ .URI }}</h1>
  <hr>
  <table>
    <thead>
      <tr>
	<th>Path</th>
	<th>Type</th>
	<th>Size</th>
      </tr>
    </thead>

    <tbody>
      {{ range .List.CommonPrefixes }}
	<tr>
	  <td><a href="{{ basename . }}/">{{ basename . }}/</a></td>
	  <td>DIR</td>
	  <td>-</td>
	</tr>
      {{ end }}

      {{ range .List.Entries }}
	<tr>
	  <td><a href="{{ basename .Path }}">{{ basename .Path }}</a></td>
	  <td>{{ .ContentType }}</td>
	  <td>{{ .Size }}</td>
	</tr>
      {{ end }}
  </table>
  <hr>
</body>
`[1:]))
