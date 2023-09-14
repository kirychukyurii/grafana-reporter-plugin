package templates

var Index = `
{{ $tmpDir := .TmpDir }}
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{ .Dashboard.Model.Title }}</title>
	</head>
	<body>
		{{ range .Dashboard.Model.Panels }}
		<img src="{{ $tmpDir }}/panel-{{ .Type }}-{{ .Id }}.png" 
			style="position:fixed; right:{{ .GridPos.X }}px; bottom:{{ .GridPos.Y }}px; width:{{ .GridPos.W }}px; height:{{ .GridPos.H }}px; border:none;"
			alt="{{ .Title }}" />
		{{ end }}
	</body>
</html>
`
