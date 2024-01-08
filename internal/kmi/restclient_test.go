package kmi

import (
	"encoding/xml"
	"os"
	"strings"

	"log"
	"reflect"
	"testing"
)

func Test_AccountMarshalling(t *testing.T) {
	data, err := os.ReadFile("account_test.xml")

	if err != nil {
		log.Fatal(err)
	}
	var e1 Account
	err = xml.Unmarshal([]byte(data), &e1)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		t.Errorf("Marshalling() error = %v", err)
		return
	}

	if !reflect.DeepEqual(e1.AdminGroup, "PIM_TEST_admins") {
		t.Errorf("Marshalling() = %v, want %v", e1.AdminGroup, "PIM_TEST_admins")
	}
	if !reflect.DeepEqual(e1.Contact, "sme:2756") {
		t.Errorf("Marshalling() = %v, want %v", e1.AdminGroup, "sme:2756")
	}
	if !reflect.DeepEqual(e1.Contact, "sme:2756") {
		t.Errorf("Marshalling() = %v, want %v", e1.AdminGroup, "sme:2756")
	}

	if !reflect.DeepEqual(e1.Engine[len(e1.Engine)-1].Cloud, "linode") {
		t.Errorf("Marshalling() = %v, want %v", e1.Engine[len(e1.Engine)-1].Cloud, "linode")
	}

}
func Test_EngineMarshalling(t *testing.T) {
	data := []byte(`<workload projection="instance_validator" source="restserv:user:hachandr_kmi_cert">
	<region source="restserv:user:hachandr_kmi_cert">us-iad</region>
	<kubernetes_service_account source="restserv:user:hachandr_kmi_cert">system:serviceaccount:app:SA1</kubernetes_service_account>
  </workload>`)
	var e1 Workload
	err := xml.Unmarshal(data, &e1)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		t.Errorf("Marshalling() error = %v", err)
		return
	}

	if !reflect.DeepEqual(e1.Projection, "instance_validator") {
		t.Errorf("Marshalling() = %v, want %v", e1.Projection, "instance_validator")
	}

	if !reflect.DeepEqual(e1.KubernetesServiceAccount.Text, "system:serviceaccount:app:SA1") {
		t.Errorf("Marshalling() = %v, want %v", e1.KubernetesServiceAccount, "system:serviceaccount:app:SA1")
	}

	kmiserviceAcc := e1.KubernetesServiceAccount.Text

	k8ServiceAccount := strings.Split(kmiserviceAcc, ":")[3]
	k8Namepace := strings.Split(kmiserviceAcc, ":")[2]
	if !reflect.DeepEqual(k8ServiceAccount, "SA1") {
		t.Errorf("Marshalling() = %v, want %v", k8ServiceAccount, "SA1")
	}

	if !reflect.DeepEqual(k8Namepace, "app") {
		t.Errorf("Marshalling() = %v, want %v", k8Namepace, "app")
	}
}

func Test_IdentityEngineMarshalling(t *testing.T) {
	data := []byte(`<engine name="pi-dev-usiad-l1-2023-1127-121839" cloud="linode" type="kubernetes" adders="PIM_TEST_admins" modifiers="PIM_TEST_admins" modified="353954223" source="restserv:user:hachandr_kmi_cert" published="2024-01-07 07:08:56" published_location="/secret/Col=kmi_identity_engines/Def=PIM_TEST.pi-dev-usiad-l1-2023-1127-121839/Idx=353954027">
	<option name="cas_base64" source="restserv:user:hachandr_kmi_cert">LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMvakNDQWVhZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJek1URXlOekV5TWpBMU1sb1hEVE16TVRFeU5ERXlNakExTWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS2ZLCjhmcjhqYjlsdFlRVittdkc4K1RSVFZiQ1VWb1NZZlJ4OVg4VGY2NVE1QzdYTFhZWTNTUjhacm5QZG9mdSttdEoKajZ2YmplQ1prOVh0K2ZZR1hhOFdKZHYxbjQ2YzUybTIzejJnQTNLQTUxaGRVWEUxQ0NzUlF6eVVldzA5K0grQwpRQjMyeVZZOVo5SjZRQ3RGVG9makN1OHUrRHllM0UrMWlyYy9UcmRlbGNCRm9xS255R2VsMi9XZWJ6UnIyVmhuCmpJUklRSEpKWCtOVE1GcEo5SFVKMUQyQTltMHZBazRmNkFIRHZKWkcxVUZkaHArbUFGajBsYVRsRTZTTytEMWMKTmFZVGZnSmtpcVIraWhKMEtCdzNWb3FSL1hqb20wMXZPcEkzZWRObFdqUjNIWXJmbnRFN1VjQTFaR2g5Q1pVOQpKZThoWmFDUThiUDJHNzRBY2prQ0F3RUFBYU5aTUZjd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZPelRKN3FOdTBMY2RnZXFpSTlWZlhQaExyK0dNQlVHQTFVZEVRUU8KTUF5Q0NtdDFZbVZ5Ym1WMFpYTXdEUVlKS29aSWh2Y05BUUVMQlFBRGdnRUJBQWpRRTdHZCtLYWhhME9lY0hlawo2dk5ZUmpUeGVaSjFETGRDc2dxZFErVVRudkF3RWRzWDBra2x5L05wWHVraWpBL011QlNYaTRMS242ZFcyalRUClIyVUpQRkswTzFoVjk3OG9USXF2QVJOV1czbStoclh3SGJVRnIrbE1MZ2s4MUgyWTBGTk5vNnlFUlAvSWtRWC8KNFpzWmNocFhVei9sU05mMm0rdVhKdUYxUVFYSDBNS3RFU1JxVnIwR3hhZmVadDNzVDE2S1ZoTTZIRThTZEk5WApadFk4VVhyYnR3ZEp0bVhjVjdmMjlZQVptcVRyVXpPa2ZjdFBrVzNVT25Ga1cwNHhkZ1Y2d2krZjlWallRbEFRCmk5dTEwMWJITjB1aEo5emVWb3JlSDdmblIwRHVPSHBqUndlcXdIK0xrZ2FBU2tjQ3ExMTJMdlVoUkRYa1h1dTEKV1E0PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==</option>
	<option name="endpoint_url" source="restserv:user:hachandr_kmi_cert">https://fbfcd97b-cf7b-4ce6-8736-1abbd7e072b1.us-iad-1.linodelke.net:443</option>
	<workload projection="instance_validator"/>
  </engine>`)
	var e1 IdentityEngine
	err := xml.Unmarshal(data, &e1)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(e1.Name, "pi-dev-usiad-l1-2023-1127-121839") {
		t.Errorf("Marshalling() = %v, want %v", e1.Name, "pi-dev-usiad-l1-2023-1127-121839")
	}
	if !reflect.DeepEqual(e1.Cloud, "linode") {
		t.Errorf("Marshalling() = %v, want %v", e1.Cloud, "linode")
	}

	if !reflect.DeepEqual(e1.Workload[0].Projection, "instance_validator") {
		t.Errorf("Marshalling() = %v, want %v", e1.Workload[0].Projection, "instance_validator")
	}

}
