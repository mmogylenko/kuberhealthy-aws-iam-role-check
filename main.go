// Licensed to Mykola Mogylenko <mmogylenko@gmail.com> under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Mykola Mogylenko <mmogylenko@gmail.com> licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"

	"context"
	"syscall"

	"github.com/Comcast/kuberhealthy/v2/pkg/checks/external/checkclient"
	log "github.com/sirupsen/logrus"
)

var (
	buildVersion string = ""
	buildTime    string = ""

	sess  *session.Session
	debug bool

	// Environment Variables
	targetArnEnv = os.Getenv("TARGET_ARN")
	nodeName     = os.Getenv("NODE_NAME")
	debugEnv     = os.Getenv("DEBUG")

	ctx context.Context
	// Channel for interrupt signals
	signalChan chan os.Signal
)

func createAWSSession() *session.Session {
	// Build an AWS session
	log.Debugln("Building AWS session")
	awsConfig := aws.NewConfig().WithCredentialsChainVerboseErrors(debug)
	awsConfig.HTTPClient = &http.Client{
		Timeout: time.Duration(10) * time.Second,
	}
	awsConfig.Retryer = client.DefaultRetryer{
		NumMaxRetries: 1,
	}
	return session.Must(session.NewSession(awsConfig))
}

func init() {
	var err error

	ctx = context.Background()
	signalChan = make(chan os.Signal, 2)

	// Relay incoming OS interrupt signals to the signalChan
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	// Enable Debug just in case
	if len(debugEnv) != 0 {
		debug, err = strconv.ParseBool(debugEnv)
		if err != nil {
			log.Fatalln("Failed to parse DEBUG Environment variable:", err.Error())
		}
	}

	if debug {
		log.Infoln("Debug logging enabled")
		log.SetLevel(log.DebugLevel)
	}
	// APP Build information
	log.Debugln("Application Version:", buildVersion)
	log.Debugln("Application Build Time:", buildTime)
	// Good to know which worker is failing
	if len(nodeName) != 0 {
		log.Debugln("Running on:", nodeName)
	}

	// APP Environment
	log.Debugln(os.Args)
	// Check if ARN Environment variable is set
	if len(targetArnEnv) == 0 {
		log.Fatalln("ARN Environment variable is not set (TARGET_ARN)")
	}

}

func main() {

	var err error

	if arn.IsARN(targetArnEnv) != true {
		log.Errorln("ASSUMED_ROLE_ARN environment variable is not ARN:", targetArnEnv)
		err = checkclient.ReportFailure([]string{"ASSUMED_ROLE_ARN environment variable is not ARN"})
		if err != nil {
			log.Println(err.Error())
		}
		return

	}
	log.Debugln("TARGET_ARN environment variable matches ARN format")

	sess = createAWSSession()
	if sess == nil {
		err = fmt.Errorf("nil AWS session: %v", sess)
		err = checkclient.ReportFailure([]string{err.Error()})
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	go listenForInterrupts()

	// Catch panics
	var r interface{}
	defer func() {
		r = recover()
		if r != nil {
			log.Infoln("Recovered panic:", r)
			err = checkclient.ReportFailure([]string{r.(string)})
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()

	select {
	case err = <-runArnCheck():
		if err != nil {
			// Report a failure if there an error occurred during the check.
			err = fmt.Errorf("Error occurred during runArnCheck: %s", err)
			log.Infoln("IAM Role check failed. Set env DEBUG=1 for more verbosity")
			log.Debugln(err)
			err = checkclient.ReportFailure([]string{err.Error()})
			if err != nil {
				log.Println(err.Error())
			}
			return
		}
		log.Infoln("IAM Role check successful")
	case <-ctx.Done():
		return
	}

	err = checkclient.ReportSuccess()
	if err != nil {
		log.Println(err.Error())
	}

}

// listenForInterrupts watches the signal and done channels for termination.
func listenForInterrupts() {

	// Relay incoming OS interrupt signals to the signalChan.
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	sig := <-signalChan // This is a blocking operation -- the routine will stop here until there is something sent down the channel.
	log.Infoln("Received an interrupt signal from the signal channel")
	log.Debugln("Signal received was:", sig.String())

	// Clean up pods here.
	log.Infoln("Shutting down")

	os.Exit(0)
}
