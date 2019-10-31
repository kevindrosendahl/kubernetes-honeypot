package kubelet

type HoneypotConfig struct {
	ConnectionString string                 `toml:"connection_string"`
	PodStorePath     string                 `toml:"pod_store_path"`
	Capacity         HoneypotCapacityConfig `toml:"capacity"`
}

type HoneypotCapacityConfig struct {
	Cpu    string `toml:"cpu"`
	Memory string `toml:"memory"`
	Pods   string `toml:"pods"`
}
