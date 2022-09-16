package main

import (
	"github.com/alecthomas/repr"
	"github.com/jonwinton/ddqp"
)

func main() {
	metricQueryParser := ddqp.NewMetricQueryParser()
	metricMonitorParser := ddqp.NewMetricMonitorParser()

	val, err := metricQueryParser.ParseString("", `sum:kubernetes.containers.state.terminated{reason:oomkilled-foo} by    {kube_cluster_name,kube_deployment}`)
	if err != nil {
		panic(err)
	}
	repr.Println(val)

	nextVal, err := metricMonitorParser.ParseString("", `avg(last_5m):max:system.disk.in_use{reason:oomkilled} by {host} < 1.2`)
	if err != nil {
		panic(err)
	}

	repr.Println(nextVal)
}
