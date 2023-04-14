package launch

import (
	"flag"
	"os"

	"github.com/wlf92/torch/pkg/log"
	"gopkg.in/yaml.v3"
)

var Config = &Info{}

func init() {
	def := flag.String("launch", "../launch.yaml", "Specify the launch file path")

	bts, err := os.ReadFile(*def)
	if err != nil {
		log.Fatalw("load launch file fail")
	}
	err = yaml.Unmarshal(bts, Config)
	if err != nil {
		log.Fatalw("load launch file fail")
	}
}

// Yaml2Go
type Info struct {
	Node Node `yaml:"node"`
	Gate Gate `yaml:"gate"`
}

// Node
type Node struct {
	RpcPort int `yaml:"rpc_port"`
}

// Gate
type Gate struct {
	RpcPort int `yaml:"rpc_port"`
}
