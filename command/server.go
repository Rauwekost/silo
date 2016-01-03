package command

import (
	"net/http"

	log "github.com/rauwekost/silo/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/rauwekost/silo/Godeps/_workspace/src/github.com/jacobstr/confer"
	"github.com/rauwekost/silo/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/rauwekost/silo/Godeps/_workspace/src/github.com/spf13/pflag"
	web "github.com/rauwekost/silo/http"
)

func NewServerCommand() *cobra.Command {
	c := cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			srv, err := web.NewServer(GetServerConfiguration(cmd.Flags()))
			if err != nil {
				return err
			}

			//http-server
			addr := srv.Config.GetString("silo.address")
			httpsrv := &http.Server{
				Addr:    addr,
				Handler: srv.HTTPHandler(),
			}

			//serve
			log.Infof("http-server listening on: %s", addr)
			log.Fatal(httpsrv.ListenAndServe())
			return nil
		},
	}
	c.Flags().StringP("config", "c", "", "path to config.yml")
	c.Flags().StringP("addr", "", "127.0.0.1:3001", "address to serve on")
	return &c
}

//GetServerConfiguration formats the configuration
func GetServerConfiguration(flags *pflag.FlagSet) *confer.Config {
	//get flags
	path, _ := flags.GetString("config")
	addr, _ := flags.GetString("addr")
	store_type, _ := flags.GetString("storage-type")
	store_location, _ := flags.GetString("storage-location")

	//define a new configuration with flags as defaults
	conf := confer.NewConfig()
	conf.SetDefault("silo.address", addr)
	conf.SetDefault("silo.storage.type", store_type)
	conf.SetDefault("silo.storage.location", store_location)

	//if config path is set overwrite defaults from flags
	if path != "" {
		conf.ReadPaths(path)
	}

	return conf
}
