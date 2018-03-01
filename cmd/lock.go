package cmd

import (
	"os"
	"time"

	"github.com/jessedearing/go-distributed/lock"
	_ "github.com/jessedearing/go-distributed/lock/mongo"
	_ "github.com/jessedearing/go-distributed/lock/mysql"
	_ "github.com/jessedearing/go-distributed/lock/postgres"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Aquire a distributed lock",
	Long: `This is the CLI interface to the go-distributed library for acquring distributed locks from the command line.

This can be used to make sure tasks only run one at a time or that a single instance runs a task without having to fully boot an application.

For example:
    go-distributed lock --type mongo --db-connection "mongodb://localhost:27017/mydb?replSet=local" --non-blocking && echo "only run this once"
`,
	Run: func(cmd *cobra.Command, args []string) {
		{
			err := viper.BindPFlags(cmd.Flags())
			if err != nil {
				log.Fatal(err)
			}
		}

		nonblocking := viper.GetBool("non-blocking")
		connectionString := viper.GetString("db-connection")
		lockType := viper.GetString("type")

		if lockType == "" {
			log.Fatal("No lock type was specified")
		}

		if connectionString == "" {
			log.Fatal("No valid connection string was specified")
		}

		l, err := lock.New("postgres", connectionString)
		if err != nil {
			log.Error(err)
			os.Exit(2)
		}
		defer l.Close()

		tim := time.Now()
		log.Info("preparing to run at the nearest 15 second mark")
		log.Infof("Current time: %s", tim.String())

		// run at the nearest 15 second mark
		<-time.After(time.Duration(14-(tim.Second()%15)) * time.Second)

		if nonblocking {
			if !l.NonBlockLock("jesse") {
				log.Warn("Failed to acquire lock")
				os.Exit(1)
			}
		} else {
			l.Lock("jesse")
		}
		log.Info("Got the lock")
		time.Sleep(2 * time.Second)
		log.Info("Releasing lock")
		l.Unlock("jesse")
	},
}

func init() {
	RootCmd.AddCommand(lockCmd)
	lockCmd.Flags().String("type", "", "REQUIRED The type of database to connect to (must be one of `mysql`, `postgres`, or `mongo`)")
	lockCmd.Flags().Bool("non-blocking", false, "Non blocking lock acquisition (i.e. for a leader election)")
	lockCmd.Flags().String("db-connection", "", "REQUIRED The database connection string to be passed to the locker")
}
