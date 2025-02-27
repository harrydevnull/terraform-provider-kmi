package provider

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"terraform-provider-kmi/internal/kmi"
	"testing"
)

func Test_Definition_AzureDirectly(t *testing.T) {
	defn := kmi.KMIDefinition{
		AutoGenerate: "True",
		Type:         "azure_sp",
	}
	out, _ := xml.MarshalIndent(defn, "", "")
	data := []byte(`<definition type="azure_sp"><auto_generate>True</auto_generate></definition>`)
	fmt.Println(string(out))
	fmt.Println(string(data))
	if !reflect.DeepEqual(string(out), string(data)) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}
}

func Test_Definition_SSL(t *testing.T) {
	defn := kmi.KMIDefinition{
		AutoGenerate: "True",
		Type:         "ssl_cert",
	}
	out, _ := xml.MarshalIndent(defn, "", "")
	data := []byte(`<definition type="ssl_cert"><auto_generate>True</auto_generate></definition>`)
	fmt.Println(string(out))
	fmt.Println(string(data))
	if !reflect.DeepEqual(string(out), string(data)) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}
}

func Test_Definition_SSLCA(t *testing.T) {
	var options []*kmi.KMIOption
	options = append(options, &kmi.KMIOption{
		Name: "is_ca",
		Text: "1",
	})

	defn := kmi.KMIDefinition{
		AutoGenerate: "True",
		Type:         "ssl_cert",
		Options:      options,
	}
	out, _ := xml.MarshalIndent(defn, "", "")
	data := []byte(`<definition type="ssl_cert"><auto_generate>True</auto_generate><option name="is_ca">1</option></definition>`)
	fmt.Println(string(out))
	fmt.Println(string(data))
	if !reflect.DeepEqual(string(out), string(data)) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}
}

func Test_Definition_Opaque(t *testing.T) {
	defn := kmi.KMIDefinition{
		Type: "opaque",
	}
	out, _ := xml.MarshalIndent(defn, "", "")

	data := []byte(`<definition type="opaque"></definition>`)
	fmt.Println(string(out))
	fmt.Println(string(data))
	if !reflect.DeepEqual(string(out), string(data)) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}
}

func Test_Definition_SymetricKey(t *testing.T) {
	var options []*kmi.KMIOption
	options = append(options, &kmi.KMIOption{
		Name: "key_size_bytes",
		Text: "16",
	})

	defn := kmi.KMIDefinition{
		Adders:        "test_adder",
		Type:          "symmetric_key",
		ExpirePeriod:  "3 months",
		RefreshPeriod: "1 months",
		AutoGenerate:  "True",
		Options:       options,
	}
	out, _ := xml.MarshalIndent(defn, "", "")

	data := []byte(`<definition type="symmetric_key"><adders>test_adder</adders><expire_period>3 months</expire_period><refresh_period>1 months</refresh_period><auto_generate>True</auto_generate><option name="key_size_bytes">16</option></definition>`)
	fmt.Println(string(out))
	fmt.Println(string(data))
	if !reflect.DeepEqual(string(out), string(data)) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}
}

func Test_Definition_SSL_subject(t *testing.T) {
	var options []*kmi.KMIOption
	options = append(options, &kmi.KMIOption{
		Name: "subject",
		Text: "/CN=test-user",
	})

	defn := kmi.KMIDefinition{
		AutoGenerate: "True",
		Type:         "ssl_cert",
		Options:      options,
	}
	out, _ := xml.MarshalIndent(defn, "", "")
	data := []byte(`<definition type="ssl_cert"><auto_generate>True</auto_generate><option name="subject">/CN=test-user</option></definition>`)
	fmt.Println(string(out))
	fmt.Println(string(data))
	if !reflect.DeepEqual(string(out), string(data)) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}
}
