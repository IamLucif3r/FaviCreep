package cluster

func Cluster(data map[string]uint32) map[uint32][]string {
	clusters := make(map[uint32][]string)

	for domain, hash := range data {
		clusters[hash] = append(clusters[hash], domain)
	}

	return clusters
}
