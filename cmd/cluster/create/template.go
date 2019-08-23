package create

const kindConfigTemplate = `kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
nodes:
{{ range .Nodes }}
- role: {{ .Type }}
{{- end }}`
