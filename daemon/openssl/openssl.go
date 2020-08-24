package openssl

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strings"
)

// This function is used
func HmacSha256Signature( params string, key string ) (string, error) {
	// formalize the param then output to stdin of openssl command
	formalizeParamsCmd := exec.Command("echo", "-n", strings.TrimSpace(params))
	// openssl command: sha256 message digest algorithm, create hashed MAC with key
	signatureMsgCmd := exec.Command("openssl", "dgst", "-sha256", "-hmac", key)

	// create new pipe.
	pipeReader, pipeWriter := io.Pipe()

	// wire up the pipe b/w 2 commands as below:

	// assign stdout of first cmd to pipe Writer
	formalizeParamsCmd.Stdout = pipeWriter

	// assign stdin of second cmd to pipe reader.
	signatureMsgCmd.Stdin = pipeReader

	// assign the os stdout to the second cmd.
	signatureMsgCmd.Stdout = os.Stdout

	// run the first cmd.
	err := formalizeParamsCmd.Start()
	if err != nil {
		logrus.Errorf("Failed to execute the echo command: param = %s, err = %s", params, err)
		return "", err
	}

	// run the second cmd.
	var b bytes.Buffer
	signatureMsgCmd.Stdout = &b
	err = signatureMsgCmd.Start()
	if err != nil {
		logrus.Errorf("Failed to execute the openssl dgst command: param = %s, err = %s", params, err)
		return "", err
	}

	// make a new go routine to wait for the first command finished.
	go func() {
		// defer util the go routine done.
		defer pipeWriter.Close()
		// wait util finished.
		err = formalizeParamsCmd.Wait()
		if err != nil {
			logrus.Errorf("Failed to run the echo finished with err %s %s", formalizeParamsCmd.Stdout, err)
		}
		// done now can close the pipeWriter.
	}()
	// wait util the second done.
	err = signatureMsgCmd.Wait()
	if err != nil {
		logrus.Errorf("Failed to run the openssl finished with %s err %s", signatureMsgCmd.Stdout, err)
		return "", err
	}
	// return the result from stdout.
	return  strings.Trim(b.String(), "\n"), nil
}

