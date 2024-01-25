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
	err = xml.Unmarshal(data, &e1)
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
	data := []byte(`<engine name="pi-qa-automation-webapp-spa" cloud="linode" type="kubernetes" adders="PIM_TEST_admins" modifiers="PIM_TEST_admins" modified="355711849" source="restserv:user:hachandr_kmi_cert" published="2024-01-16 03:43:40" published_location="/secret/Col=kmi_identity_engines/Def=PIM_TEST.pi-qa-automation-webapp-spa/Idx=355711849">
	<option name="cas_base64" source="restserv:user:hachandr_kmi_cert"></option>
	<option name="endpoint_url" source="restserv:user:hachandr_kmi_cert"></option>
	<workload projection="instance_validator"/>
	<workload projection="instance_validator-1"/>
  </engine>`)
	var e1 IdentityEngine
	err := xml.Unmarshal(data, &e1)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(e1.Name, "pi-qa-automation-webapp-spa") {
		t.Errorf("Marshalling() = %v, want %v", e1.Name, "pi-dev-usiad-l1-2023-1127-121839")
	}
	if !reflect.DeepEqual(e1.Cloud, "linode") {
		t.Errorf("Marshalling() = %v, want %v", e1.Cloud, "linode")
	}

	if !reflect.DeepEqual(e1.Workload[0].Projection, "instance_validator") {
		t.Errorf("Marshalling() = %v, want %v", e1.Workload[0].Projection, "instance_validator")
	}

	if !reflect.DeepEqual(e1.Workload[1].Projection, "instance_validator-1") {
		t.Errorf("Marshalling() = %v, want %v", e1.Workload[0].Projection, "instance_validator")
	}

}

func Test_IdentityEngineMarshallingSecond(t *testing.T) {
	data := []byte(`<engine name="pi-dev-usiad-l1-2023-1127-121839" cloud="linode" type="kubernetes" adders="PIM_TEST_admins" modifiers="PIM_TEST_admins" modified="355322586" source="restserv:user:hachandr_kmi_cert" published="2024-01-13 17:30:12" published_location="/secret/Col=kmi_identity_engines/Def=PIM_TEST.pi-dev-usiad-l1-2023-1127-121839/Idx=355322586"><option name="cas_base64" source="restserv:user:hachandr_kmi_cert"></option><option name="endpoint_url" source="restserv:user:hachandr_kmi_cert"></option><workload projection="instance_validator"/></engine>`)
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

func Test_CollectionMarshalling(t *testing.T) {
	data := []byte(`<collection name="PIM_SECRETS" source="restserv:user:hachandr_kmi_cert" readers="PIM_READERS" adders="PIM_TEST_admins" modifiers="PIM_TEST_admins" modified="354310919" distributed="354318302" distributed_date="2024-01-08 20:46:23" keyspace="k3_PIM_SECRETS" account="PIM_TEST">
	<definition name="ETP_S3"/>
	<definition name="instance_validator_definition"/>
	<definition name="master_service_api_definition"/>
	<definition name="pi_api_definition"/>
	<definition name="pim_ssl_definition"/>
	<definition name="pim_ssl_definition2"/>
  </collection>`)
	var e1 Collection
	err := xml.Unmarshal(data, &e1)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(e1.Name, "PIM_SECRETS") {
		t.Errorf("Marshalling() = %v, want %v", e1.Name, "PIM_SECRETS")
	}
}

func Test_CollectionRequestMarshalling(t *testing.T) {
	data := []byte(`<collection><adders>ADMINS</adders><modifiers>ADMINS</modifiers><readers>GROUP</readers></collection>`)
	var e1 CollectionRequest
	err := xml.Unmarshal(data, &e1)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(e1.Adders, "ADMINS") {
		t.Errorf("Marshalling() = %v, want %v", e1.Adders, "ADMINS")
	}
	if !reflect.DeepEqual(e1.Modifiers, "ADMINS") {
		t.Errorf("Marshalling() = %v, want %v", e1.Modifiers, "ADMINS")
	}
	if !reflect.DeepEqual(e1.Readers, "GROUP") {
		t.Errorf("Marshalling() = %v, want %v", e1.Readers, "GROUP")
	}
}

func Test_GroupsMarshalling(t *testing.T) {
	data := []byte(`<group name="PIM_ADMIN" type="union" source="reqserv:user:hachandr_kmi_cert" account="PIM_TEST"><adders>superusers</adders><modifiers>superusers</modifiers></group>`)
	var e1 KMIGroup
	err := xml.Unmarshal(data, &e1)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(e1.Adders, "superusers") {
		t.Errorf("Marshalling() = %v, want %v", e1.Adders, "superusers")
	}
	if !reflect.DeepEqual(e1.Modifiers, "superusers") {
		t.Errorf("Marshalling() = %v, want %v", e1.Modifiers, "superusers")
	}
	if !reflect.DeepEqual(e1.Account, "PIM_TEST") {
		t.Errorf("Marshalling() = %v, want %v", e1.Account, "PIM_TEST")
	}

	if !reflect.DeepEqual(e1.Name, "PIM_ADMIN") {
		t.Errorf("Marshalling() = %v, want %v", e1.Name, "PIM_ADMIN")
	}
}
func Test_GroupsUnMarshalling(t *testing.T) {
	group := GroupRequest{
		Account: "PIM_TEST",
		Type:    "union",
	}
	out, _ := xml.MarshalIndent(group, " ", "  ")

	data := []byte(` <group type="union" account="PIM_TEST"></group>`)
	if !reflect.DeepEqual(out, data) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}

}

func Test_BlockSecretUnmarshalling(t *testing.T) {
	data := []byte(`<secret><block name="opaque" b64encoded="true">dGVzdC1lbmNvZGluZwo=</block></secret>`)
	var unmarshalled BlockSecret
	err := xml.Unmarshal(data, &unmarshalled)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(unmarshalled.Block.Text, "dGVzdC1lbmNvZGluZwo=") {
		t.Errorf("Marshalling() = %v, want %v", unmarshalled.Block.Text, "dGVzdC1lbmNvZGluZwo=")
	}
	if !reflect.DeepEqual(unmarshalled.Text, "") {
		t.Errorf("Marshalling() = %v, want %v", unmarshalled.Text, "\"\"")
	}
	if !reflect.DeepEqual(unmarshalled.Block.B64Encoded, "true") {
		t.Errorf("Marshalling() = %v, want %v", unmarshalled.Block.B64Encoded, "true")
	}
}

func Test_BlockSecretMarshalling(t *testing.T) {
	blockSecret := BlockSecret{
		Text: "",
		Block: struct {
			Text       string `xml:",chardata"`
			Name       string `xml:"name,attr"`
			B64Encoded string `xml:"b64encoded,attr"`
		}{
			Text:       "dGVzdC1lbmNvZGluZwo=",
			Name:       "opaque",
			B64Encoded: "true",
		},
	}

	out, err := xml.Marshal(blockSecret)
	if err != nil {
		log.Fatal(err)
	}

	data := []byte(`<secret><block name="opaque" b64encoded="true">dGVzdC1lbmNvZGluZwo=</block></secret>`)
	if !reflect.DeepEqual(out, data) {
		t.Errorf("Marshalling() = %v, want %v", string(out), string(data))
	}
}
