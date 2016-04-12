package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	. "github.com/tendermint/go-common"
	pcm "github.com/tendermint/go-process"
	"github.com/tendermint/go-wire"
)

// Copy a file (or dir recursively) from srcPath (local machine) to
// dstPath in the tmcore container.
func copyToMachine(mach string, app string, srcPath string, dstPath string, copyContents bool) error {

	// First, copy the file to a temporary location
	// in the machine.
	tempFile := "temp_" + RandStr(12)
	args := []string{"scp", "-r", srcPath, mach + ":" + tempFile}
	if !runProcess("scp-file-"+mach, "docker-machine", args) {
		return errors.New("Failed to copy file to machine " + mach)
	}

	// Next, docker cp the file into the container
	if copyContents {
		tempFile = tempFile + "/."
	}
	args = []string{"ssh", mach, Fmt("docker cp %v %v_tmcommon:%v", tempFile, app, dstPath)}
	if !runProcess("docker-cp-file-"+mach, "docker-machine", args) {
		return errors.New("Failed to docker-cp file to container in machine " + mach)
	}

	// Next, change the ownership of the file to tmuser
	// TODO We don't really want to change all the permissions
	args = []string{"ssh", mach, Fmt(`docker run --rm --volumes-from %v_tmcommon -u root tendermint/tmbase chown -R tmuser:tmuser %v`, app, dstPath)}
	if !runProcess("docker-chmod-file-"+mach, "docker-machine", args) {
		return errors.New("Failed to docker-run(chmod) file in machine " + mach)
	}

	// TODO: remove tempFile
	return nil
}

// NOTE: returns false if any error
func checkFileExists(mach string, container string, path string) bool {
	args := []string{"ssh", mach, Fmt(`docker exec %v ls %v`, container, path)}
	_, ok := runProcessGetResult("check-file-exists-"+mach, "docker-machine", args)
	return ok
}

//--------------------------------------------------------------------------------

func runProcess(label string, command string, args []string) bool {
	outFile := NewBufferCloser(nil)
	proc, err := pcm.StartProcess(label, command, args, nil, outFile)
	if err != nil {
		fmt.Println(Red(err.Error()))
		return false
	}

	<-proc.WaitCh
	fmt.Println(Green(command), Green(args))
	if proc.ExitState.Success() {
		fmt.Println(Blue(string(outFile.Bytes())))
		return true
	} else {
		// Error!
		fmt.Println(Red(string(outFile.Bytes())))
		return false
	}
}

func runProcessGetResult(label string, command string, args []string) (string, bool) {
	outFile := NewBufferCloser(nil)
	proc, err := pcm.StartProcess(label, command, args, nil, outFile)
	if err != nil {
		return "", false
	}

	<-proc.WaitCh
	fmt.Println(Green(command), Green(args))
	if proc.ExitState.Success() {
		fmt.Println(Blue(string(outFile.Bytes())))
		return string(outFile.Bytes()), true
	} else {
		// Error!
		fmt.Println(Red(string(outFile.Bytes())))
		return string(outFile.Bytes()), false
	}
}

//--------------------------------------------------------------------------------

func eB(s string) string {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `$`, `\$`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	s = strings.Replace(s, `'`, `\'`, -1)
	s = strings.Replace(s, `!`, `\!`, -1)
	s = strings.Replace(s, `#`, `\#`, -1)
	s = strings.Replace(s, `%`, `\%`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	s = strings.Replace(s, "`", "\\`", -1)
	return s
}

func condenseBash(cmd string) string {
	cmd = strings.TrimSpace(cmd)
	lines := strings.Split(cmd, "\n")
	res := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		res = append(res, line)
	}
	return strings.Join(res, "; ")
}

//--------------------------------------------------------------------------------

func ReadJSONFile(o interface{}, filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	wire.ReadJSON(o, b, &err)
	if err != nil {
		return err
	}
	return nil
}

//--------------------------------------------------------------------------------
// for avoiding request limits in concurrent calls to amazonec2

func maybeSleep(n int, t int) {
	if n > 16 {
		sleepMS := int(rand.Int31n(int32(t)))
		time.Sleep(time.Millisecond * time.Duration(sleepMS))
	}
}
