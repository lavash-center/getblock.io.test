package config

import (
	"os"

	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	Listen              string `env:"HTTP_LISTEN" envDefault:":8000"` // HTTP listen address
	ApiKey              string `env:"API_KEY" envDefault:""`          // API key of getblock.io account
	GetBlockRPCEndpoint string `env:"GET_BLOCK_RPC_ADDR" envDefault:"https://eth.getblock.io/mainnet/"`
}

func ParseConfig(v interface{}) error {
	defer os.Clearenv()

	return env.Parse(v)
}
