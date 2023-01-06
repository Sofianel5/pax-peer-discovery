package main

import (
	"fmt"
	"os"
	"os/exec"
)

const MPC_DIR = "./workspace/mp-spdz"

func run2pc(app, input, myaddr, counterparty, player string) (err error) {
	// Write input to file in Player-Data/Input-P<party>-0
	f, err := os.OpenFile(fmt.Sprintf("%s/Player-Data/Input-P%s-0", MPC_DIR, player), os.O_WRONLY|os.O_CREATE, 0644)
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
	output_file := fmt.Sprintf("Player-Data/Output-P%s-0", player)
	logger.Info("Output file:", output_file)
	if player == "0" {
		cmd = exec.Command("./mascot-party.x", "-N", "2", "-h", myaddr, "-p", player, "-OF", output_file, app)
	} else {
		cmd = exec.Command("./mascot-party.x", "-N", "2", "-h", counterparty, "-p", player, "-OF", output_file, app)
	}
	cmd.Dir = MPC_DIR
	out, err := cmd.CombinedOutput()
	logger.Info("MPC Protocol Complete")
	logger.Info(string(out))
	if err != nil {
		logger.Error(err)
	}
	output, err := os.ReadFile(MPC_DIR + "/" + output_file)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info("Output:", string(output))
	return nil
}
