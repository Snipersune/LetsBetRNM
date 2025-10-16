package fetch

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func FetchPowerplay() {
	var pwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	var py_path = fmt.Sprintf("%s/.venv/bin/python", pwd)
	var script_path = fmt.Sprintf("%s/py_scripts/update_powerplay_html.py", pwd)

	for {
		fmt.Println("Updating powerplay data...")
		c := exec.Command(py_path, script_path)
		if err := c.Run(); err != nil {
			fmt.Println("server: Error: ", err)
		}

		time.Sleep(60 * time.Second)
	}
}
