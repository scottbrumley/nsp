package nsp

import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"strings"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"bytes"
	"golang.org/x/crypto/ssh"
)

// Parameters from command line
type ParamStruct struct{
	UserName string
	UserPass string
	Test bool
	Cmd string
	HostName string
	HostFile string
	HostPort string
}

func getPasswd()(string){
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err == nil {
		fmt.Println("\nPassword typed: " + string(bytePassword))
	}
	password := string(bytePassword)

	return strings.TrimSpace(password)
}

func getUser()(string){
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	return strings.TrimSpace(username)
}

// Collect parameters from the command line
func GetParams()(retParams ParamStruct){
	var userFlag = flag.String("user","","IPS Username")
	var passFlag = flag.String("password","","IPS Password")
	var cmdFlag = flag.String("cmd","ls","Command String")
	var testFlag = flag.Bool("test",false,"Testing Mode")
	var hostFlag = flag.String("hostname","localhost","Hostname or IP Address")
	var fileFlag = flag.String("hostfile","","HostsFile Name")
	var portFlag = flag.String("port","22","Host Port")
	flag.Parse()

	retParams.UserName = *userFlag
	retParams.UserPass = *passFlag
	retParams.Cmd = *cmdFlag
	retParams.Test = *testFlag
	retParams.HostName = *hostFlag
	retParams.HostFile = *fileFlag
	retParams.HostPort = *portFlag

	// Test Params
	if retParams.UserName == "" {
		retParams.UserName = getUser()
	}
	if retParams.UserPass == "" {
		retParams.UserPass = getPasswd()
	}

	return retParams
}

func SshCommand(lclParms ParamStruct)(string){
	config := &ssh.ClientConfig{
		User: lclParms.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(lclParms.UserPass),
		},
	}
	client, err := ssh.Dial("tcp", lclParms.HostName + ":" + lclParms.HostPort , config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(lclParms.Cmd); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	return b.String()
}
