package main

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"
	"github.com/jonwinton/dotodag-ql"
)

func main() {
	parser := participle.MustBuild[dotodag.MetricQuery](
		participle.Unquote("String"),
	)

	query, err := parser.ParseString("", `sum:kubernetes.containers.state.terminated{reason:oomkilled} by    {kube_cluster_name,kube_deployment}`)
	if err != nil {
		panic(err)
	}
	repr.Println(query)
}
