package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

//cron expression to execute job - see: [logfx cron --help]

var expr string
var exprFlagStr = `(required) A cron expression represents a set of times to execute a job. Use 5 space-separated fields, around quotes. 
Like: logfx cron --expr="* * * * *"

field name   | mandatory? | allowed values  | allowed special characters
----------   | ---------- | --------------  | --------------------------
minutes      | yes        | 0-59            | * / , -
hours        | yes        | 0-23            | * / , -
day of month | yes        | 1-31            | * / , - ?
month        | yes        | 1-12 or jan-dec | * / , -
day of week  | yes        | 0-6 or sun-sat  | * / , - ?
	`

func CronCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cron",
		Short: "cron creates a execution job for logfx.",
		Run:   cronjob,
	}

	cmd.Flags().StringVar(&expr, "expr", "", exprFlagStr)
	cmd.MarkFlagRequired("expr")

	return cmd
}

func cronjob(cmd *cobra.Command, args []string) {
	logger.Info("logfx - starting cron job", "from dir", fromDir, "to dir", toDir, "cron expression", expr)

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		fmt.Println("Error loading the location!")
		return
	}

	c := cron.New(
		cron.WithLocation(location),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	id, err := c.AddFunc(strings.Trim(expr, "\""), func() {
		err := targz(&fromDir, &toDir)
		if err != nil {
			logger.Error("logfx error - task fails", "err", err)
		}
	})
	if err != nil {
		fmt.Println("Error scheduling job:", err)
		return
	}

	job := c.Entry(id).Job
	job.Run()

	go c.Start()

	listenSignal()
}

func listenSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-ch
	fmt.Println("Received signal, shutting down...")
}

func task() {
	fmt.Println(time.Now().String() + " - Start Task")
	time.Sleep(4 * time.Second)
	fmt.Println(fromDir, expr)
	fmt.Println(time.Now().String() + " - End Task")
}
