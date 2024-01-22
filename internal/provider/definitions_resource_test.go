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

	defn := kmi.KMIDefinition{
		AutoGenerate: "True",
		Type:         "ssl_cert",
		Option: &kmi.KMIOption{
			Name: "is_ca",
			Text: "1",
		},
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

	defn := kmi.KMIDefinition{
		Type:          "symmetric_key",
		ExpirePeriod:  "3 months",
		RefreshPeriod: "1 months",
		AutoGenerate:  "True",
		Option: &kmi.KMIOption{
			Name: "key_size_bytes",
			Text: "16",
		},
	}
	out, _ := xml.MarshalIndent(defn, "", "")

	data := []byte(`<definition type="symmetric_key"><expire_period>3 months</expire_period><refresh_period>1 months</refresh_period><auto_generate>True</auto_generate><option name="key_size_bytes">16</option></definition>`)
	fmt.Println(string(out))
	fmt.Println(string(data))
	if !reflect.DeepEqual(string(out), string(data)) {
		t.Errorf("Marshalling() = %v, want %v", out, data)
	}

}
