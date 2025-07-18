package cluster

func Cluster(results map[string]uint32) map[uint32][]string {
	clusterMap := make(map[uint32][]string)
	for domain, hash := range results {
		clusterMap[hash] = append(clusterMap[hash], domain)
	}
	return clusterMap
}
