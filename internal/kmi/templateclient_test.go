package kmi

import (
	"encoding/xml"
	"fmt"

	"log"
	"reflect"
	"testing"
)

func Test_TemplateMarshalling(t *testing.T) {
	data := []byte(`<template name="cpc_client_template" source="restserv:user:hachandr_kmi_cert" add_date="2023-12-20 17:16:27" modified="354687073">
	<constraint type="common_name" warn="False" source="restserv:user:hachandr_kmi_cert" add_date="2023-12-20 17:16:27" modified="354687072">instance-validator</constraint>
	<collectionacl target="PIM_SECRETS" source="restserv:user:hachandr_kmi_cert" add_date="2023-12-20 17:16:28" modified="354687073"/>
  </template>`)
	var e1 TemplateResponse
	err := xml.Unmarshal(data, &e1)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		t.Errorf("Marshalling() error = %v", err)
		return
	}

	if !reflect.DeepEqual(e1.Name, "cpc_client_template") {
		t.Errorf("Marshalling() = %v, want %v", e1.Name, "cpc_client_template")
	}
	if !reflect.DeepEqual(e1.AddDate, "2023-12-20 17:16:27") {
		t.Errorf("Marshalling() = %v, want %v", e1.Name, "2023-12-20 17:16:27")
	}

	if !reflect.DeepEqual(e1.Collectionacl.Target, "PIM_SECRETS") {
		t.Errorf("Marshalling() = %v, want %v", e1.Collectionacl.Target, "PIM_SECRETS")
	}
	if !reflect.DeepEqual(e1.Constraints[0].Text, "instance-validator") {
		t.Errorf("Marshalling() = %v, want %v", e1.Collectionacl.Target, "instance-validator")
	}

}

func Test_SigningMarshalling(t *testing.T) {
	kmiSigner := Template{
		Collectionacl: &ApproveClientCollection{
			Target: "sometarget",
		},
	}

	out, _ := xml.MarshalIndent(kmiSigner, "", "")

	data := []byte(`<template><collectionacl target="sometarget"></collectionacl></template>`)
	fmt.Println(string(out))
	fmt.Println(string(data))
	if !reflect.DeepEqual(string(out), string(data)) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}
}
