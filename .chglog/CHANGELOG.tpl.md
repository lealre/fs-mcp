{{ range .Versions }}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }} ({{ datetime "2006-01-02" .Tag.Date }})

{{ range .Commits -}}
* {{ .Header }}
{{ end }}

{{- if .RevertCommits }}
### Reverts
{{ range .RevertCommits -}}
* {{ .Revert.Header }}
{{ end }}
{{ end -}}

{{- if .MergeCommits }}
### Pull Requests
{{ range .MergeCommits -}}
* {{ .Header }}
{{ end }}
{{ end -}}

{{ end -}}
