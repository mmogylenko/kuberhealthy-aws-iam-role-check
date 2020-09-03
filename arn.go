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
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	log "github.com/sirupsen/logrus"
)

// Arn represents Amazon
// Resource Name
type Arn struct {
	ARN          string
	Partition    string
	Service      string
	Region       string
	AccountID    string
	ResourceType string
	ResourceID   string
	ResourcePath string
}

// String is pretty-print for Arn struct
func (a Arn) String() string {
	return fmt.Sprintf("{Service:%s, AccountID:%s, ResourceID:%s}", a.Service, a.AccountID, a.ResourceID)
}

// parseArn is decomposess ARN
// into components accordingly to ARN Format
func parseArn(arn string) *Arn {
	sections := strings.Split(arn, ":")

	result := &Arn{
		ARN:          sections[0],
		Partition:    sections[1],
		Service:      sections[2],
		Region:       sections[3],
		AccountID:    sections[4],
		ResourcePath: "",
	}

	if n := strings.Count(sections[5], ":"); n > 0 {
		parts := strings.Split(sections[5], ":")
		result.ResourceType = parts[0]
		result.ResourceID = parts[1]
	} else {
		slashes := strings.Count(sections[5], "/")
		if slashes == 0 {
			result.ResourceID = sections[5]
		} else {
			parts := strings.Split(sections[5], "/")
			result.ResourceType = parts[0]
			result.ResourceID = parts[1]
			result.ResourcePath = strings.Join(parts[2:], "/")
		}
	}
	return result
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

// equalWith compares two ARNs components
func (a *Arn) equalWith(second *Arn, CompareComponents []string) bool {
	val := reflect.ValueOf(a).Elem()
	secondFields := reflect.Indirect(reflect.ValueOf(second))

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		if !isValueInList(typeField.Name, CompareComponents) {
			continue
		}
		value := val.Field(i)
		secondValue := secondFields.FieldByName(typeField.Name)

		if value.Interface() != secondValue.Interface() {
			return false
		}
	}
	return true
}

// createStsClient creates and returns an AWS STS client.
func createStsClient(s *session.Session) *sts.STS {
	// Build a STS client.
	log.Debugln("Building STS client")
	return sts.New(s)
}

func runArnCheck() chan error {
	// Create a channel for the check.
	checkChan := make(chan error, 0)

	log.Debugln("Starting ARN Check")

	go func() {
		defer close(checkChan)
		c := createStsClient(sess)

		identity, err := c.GetCallerIdentity(nil)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					checkChan <- fmt.Errorf(aerr.Error())
				}
			} else {
				checkChan <- fmt.Errorf(err.Error())
			}
			return
		}

		// Decompose ARNs
		assumeArn := parseArn(*identity.Arn)
		targetArn := parseArn(targetArnEnv)

		log.Debugln("Comparing Service, AccountID and ResourceID ARN Components")
		// Output pretty ARNs for DEBUG purposes
		log.Debugln("ARN from environment", targetArn.String())
		log.Debugln("ARN from assumed-role", assumeArn.String())

		if targetArn.equalWith(assumeArn, []string{"Service", "AccountID", "ResourceID"}) {
			log.Debugln("assumed-role ARN is matching target ARN")
			checkChan <- nil
		} else {
			log.Println("assumed-role ARN is not matching Target ARN. Set DEBUG=1 to see the ARNs difference")
			checkChan <- fmt.Errorf("assumed-role ARN is not matching Target ARN")
		}
		return
	}()

	return checkChan
}
