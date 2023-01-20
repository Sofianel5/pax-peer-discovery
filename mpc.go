package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const MPC_DIR = "./workspace/mp-spdz"

func run2pc(app, input, myaddr, counterparty, player string) (output string, err error) {
	logger.Info("Running protocol as player ", player, " at ", myaddr, ", with input: ", input, ", and counterparty: ", counterparty)
	// Write input to file in Player-Data/Input-P<party>-0
	f, err := os.OpenFile(fmt.Sprintf("%s/Player-Data/Input-P%s-0", MPC_DIR, player), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
	if player == "0" {
		cmd = exec.Command("./mascot-party.x", "-N", "2", "-h", myaddr, "-p", player, "-OF", "Output", app)
	} else {
		cmd = exec.Command("./mascot-party.x", "-N", "2", "-h", counterparty, "-p", player, "-OF", "Output", app)
	}
	cmd.Dir = MPC_DIR
	out, err := cmd.CombinedOutput()
	logger.Info("MPC Protocol Complete")
	logger.Info(string(out))
	if err != nil {
		logger.Error(err)
		return
	}
	bytes_output, err := os.ReadFile(MPC_DIR + "/Output-P" + player + "-0")
	if err != nil {
		logger.Error(err)
		return
	}
	output = strings.TrimSpace(string(bytes_output))
	if output == "0" {
		logger.Info("No match found. Find another counterparty.")
	} else if output == "1" {
		logger.Info("Match found! Execute trade with counterparty.")
	} else {
		logger.Error("Unexpected output: ", string(output))
	}
	return
}
