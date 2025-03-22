package fetch

import (
	"fmt"
	"os/exec"
	"time"
)

func FetchPowerplay() {
	for {
		fmt.Println("Updating powerplay data...")
		time.Sleep(60 * time.Second)

		c := exec.Command("/home/dnee/Documents/Kod/LetsBetRNM/.venv/bin/python", "/home/dnee/Documents/Kod/LetsBetRNM/py_scripts/update_powerplay_html.py")
		if err := c.Run(); err != nil {
			fmt.Println("server: Error: ", err)
		}
	}
}
