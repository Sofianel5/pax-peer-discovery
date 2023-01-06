package main

import (
	"fmt"
	"os"
	"os/exec"
)

const MPC_DIR = "./workspace/mp-spdz"

func run2pc(app, input, myaddr, counterparty string, player int) (err error) {
	// Write input to file in Player-Data/Input-P<party>-0
	f, err := os.OpenFile(fmt.Sprintf("%s/Player-Data/Input-P%d-0", MPC_DIR, player), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.WriteString(input)
	if err != nil {
		return
	}
	// Run MPC
	var cmd *exec.Cmd
	if player == 0 {
		cmd = exec.Command("ls")
		// cmd = exec.Command(fmt.Sprintf("./mascot-party.x -N 2 -h %s -p %d %s", myaddr, player, app))
	} else {
		cmd = exec.Command("ls")
		// cmd = exec.Command(fmt.Sprintf("./mascot-party.x -N 2 -h %s -p %d %s", counterparty, player, app))
	}
	cmd.Dir = MPC_DIR
	out, err := cmd.CombinedOutput()
	logger.Info("MPC Protocol Complete")
	logger.Info(string(out))
	if err != nil {
		logger.Error(err)
	}
	return nil
}
