package command

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
)

func NewStatsCommand() *cobra.Command {
	c := cobra.Command{
		Use: "stats",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf := GetServerConfiguration(cmd.Flags())
			db, err := bolt.Open(conf.GetString("imgd.database.location"), 0600, &bolt.Options{Timeout: 1 * time.Second})
			if err != nil {
				return err
			}
			defer db.Close()

			prev := db.Stats()
			for {
				tm.Clear()

				stats := db.Stats()
				diff := stats.Sub(&prev)
				totals := tm.NewTable(0, 10, 5, ' ', 0)
				fmt.Fprintf(totals, "FreeAlloc\tFreePageN\tFreeListInUse\tOpenTxN\n")
				fmt.Fprintf(totals, "%d\t%d\t%d\t%d\n", diff.FreeAlloc, diff.FreePageN, diff.FreelistInuse, diff.OpenTxN)

				prev = stats
				time.Sleep(10 * time.Second)

				tm.Print(totals)
				tm.Flush()
			}
		},
	}
	c.Flags().StringP("config", "c", "", "path to config.yml")
	return &c
}
