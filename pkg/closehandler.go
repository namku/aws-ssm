package pkg

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/briandowns/spinner"
)

func SetupCloseHandler(indicatorSpinner *spinner.Spinner) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			fmt.Println("\rThe interrupt got handled")
			indicatorSpinner.Stop()
			signal.Stop(c)
			os.Exit(0)
		}
	}()
}
