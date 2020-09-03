package main

import (
	"reflect"
	"testing"
)

func getField(a *Arn, field string) string {
	r := reflect.ValueOf(a)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func TestParseArn(t *testing.T) {

	cases := []struct {
		component string
		want      string
	}{
		{
			component: "Service",
			want:      "sts",
		},
		{
			component: "AccountID",
			want:      "012345678910",
		},
		{
			component: "ResourceID",
			want:      "read-access-to-everything",
		},
	}

	for _, cs := range cases {
		t.Run(cs.component, func(t *testing.T) {
			if got := parseArn(targetArnEnv); getField(got, cs.component) != cs.want {
				t.Errorf("parseArn(exampleArn) %v = %v, want %v", cs.component, getField(got, cs.component), cs.want)
			}
		})
	}
}

func TestEqualWith(t *testing.T) {
	cases := []struct {
		caseName  string
		firstArn  Arn
		secondArn Arn
		want      bool
	}{
		{
			caseName:  "ARNsMatching",
			firstArn:  Arn{ARN: "arn", Partition: "aws", Service: "s3", Region: "", AccountID: "012345678910", ResourceID: "read-access-to-everything"},
			secondArn: Arn{ARN: "arn", Partition: "aws", Service: "s3", Region: "", AccountID: "012345678910", ResourceID: "read-access-to-everything"},
			want:      true,
		},
		{
			caseName:  "ARNsNotMatching",
			firstArn:  Arn{ARN: "arn", Partition: "aws", Service: "s3", Region: "", AccountID: "012345678910", ResourceID: "read-access-to-everything"},
			secondArn: Arn{ARN: "arn", Partition: "aws-us-gov", Service: "s3", Region: "", AccountID: "012345678910", ResourceID: "write-access-to-everything"},
			want:      false,
		},
	}

	for _, cs := range cases {
		t.Run(cs.caseName, func(t *testing.T) {
			if cs.firstArn.equalWith(&cs.secondArn, []string{"Service", "AccountID", "ResourceID"}) != cs.want {
				t.Errorf("Case %v, want == %v", cs.caseName, cs.want)
			}

		})
	}
}
