package kmi

import "encoding/xml"

type Account struct {
	XMLName        xml.Name `xml:"account"`
	Text           string   `xml:",chardata"`
	Name           string   `xml:"name,attr"`
	Contact        string   `xml:"contact,attr"`
	AdminGroup     string   `xml:"admin_group,attr"`
	MaxCols        string   `xml:"max_cols,attr"`
	MaxGroups      string   `xml:"max_groups,attr"`
	MaxDefs        string   `xml:"max_defs,attr"`
	MaxSecrets     string   `xml:"max_secrets,attr"`
	MaxTemplates   string   `xml:"max_templates,attr"`
	MaxEngines     string   `xml:"max_engines,attr"`
	MaxProjections string   `xml:"max_projections,attr"`
	MaxOptions     string   `xml:"max_options,attr"`
	NumCollections string   `xml:"num_collections,attr"`
	NumGroups      string   `xml:"num_groups,attr"`
	NumUsers       string   `xml:"num_users,attr"`
	NumWorkloads   string   `xml:"num_workloads,attr"`
	NumMachines    string   `xml:"num_machines,attr"`
	NumSecrets     string   `xml:"num_secrets,attr"`
	NumEngines     string   `xml:"num_engines,attr"`
	Collection     []struct {
		Text            string `xml:",chardata"`
		Name            string `xml:"name,attr"`
		Source          string `xml:"source,attr"`
		Readers         string `xml:"readers,attr"`
		Adders          string `xml:"adders,attr"`
		Modifiers       string `xml:"modifiers,attr"`
		Modified        int64  `xml:"modified,attr"`
		Distributed     int64  `xml:"distributed,attr"`
		DistributedDate string `xml:"distributed_date,attr"`
		Keyspace        string `xml:"keyspace,attr"`
		Account         string `xml:"account,attr"`
	} `xml:"collection"`
	Group []struct {
		Text       string `xml:",chardata"`
		Name       string `xml:"name,attr"`
		Type       string `xml:"type,attr"`
		Source     string `xml:"source,attr"`
		Account    string `xml:"account,attr"`
		Engine     string `xml:"engine,attr"`
		Projection string `xml:"projection,attr"`
		Instance   string `xml:"instance,attr"`
	} `xml:"group"`
	Engine []struct {
		Text              string `xml:",chardata"`
		Name              string `xml:"name,attr"`
		Cloud             string `xml:"cloud,attr"`
		Type              string `xml:"type,attr"`
		Adders            string `xml:"adders,attr"`
		Modifiers         string `xml:"modifiers,attr"`
		Modified          int64  `xml:"modified,attr"`
		Source            string `xml:"source,attr"`
		Published         string `xml:"published,attr"`
		PublishedLocation string `xml:"published_location,attr"`
	} `xml:"engine"`
}

type IdentityEngine struct {
	XMLName           xml.Name `xml:"engine"`
	Text              string   `xml:",chardata"`
	Name              string   `xml:"name,attr"`
	Cloud             string   `xml:"cloud,attr"`
	Type              string   `xml:"type,attr"`
	Adders            string   `xml:"adders,attr"`
	Modifiers         string   `xml:"modifiers,attr"`
	Modified          string   `xml:"modified,attr"`
	Source            string   `xml:"source,attr"`
	Published         string   `xml:"published,attr"`
	PublishedLocation string   `xml:"published_location,attr"`
	Option            []struct {
		Text   string `xml:",chardata"`
		Name   string `xml:"name,attr"`
		Source string `xml:"source,attr"`
	} `xml:"option"`
	Workload []struct {
		// Text       string `xml:",chardata"`
		Projection string `xml:"projection,attr"`
	} `xml:"workload"`
}

type KMIEngine struct {
	XMLName   xml.Name      `xml:"engine"`
	Text      string        `xml:",chardata"`
	Cloud     string        `xml:"cloud,attr"`
	Type      string        `xml:"type,attr"`
	Option    []KMIOption   `xml:"option"`
	Workloads []KMIWorkload `xml:"workload"`
}
type KMIWorkload struct {
	Text                     string `xml:",chardata"`
	Projection               string `xml:"projection,attr"`
	KubernetesServiceAccount string `xml:"kubernetes_service_account"`
	Region                   string `xml:"region"`
}
type KMIOption struct {
	Text string `xml:",chardata"`
	Name string `xml:"name,attr"`
}

type Workload struct {
	XMLName    xml.Name `xml:"workload"`
	Text       string   `xml:",chardata"`
	Projection string   `xml:"projection,attr"`
	Source     string   `xml:"source,attr"`
	Region     struct {
		Text   string `xml:",chardata"`
		Source string `xml:"source,attr"`
	} `xml:"region"`
	KubernetesServiceAccount struct {
		Text   string `xml:",chardata"`
		Source string `xml:"source,attr"`
	} `xml:"kubernetes_service_account"`
}

type Collection struct {
	XMLName         xml.Name `xml:"collection"`
	Text            string   `xml:",chardata"`
	Name            string   `xml:"name,attr"`
	Source          string   `xml:"source,attr"`
	Readers         string   `xml:"readers,attr"`
	Adders          string   `xml:"adders,attr"`
	Modifiers       string   `xml:"modifiers,attr"`
	Modified        string   `xml:"modified,attr"`
	Distributed     string   `xml:"distributed,attr"`
	DistributedDate string   `xml:"distributed_date,attr"`
	Keyspace        string   `xml:"keyspace,attr"`
	Account         string   `xml:"account,attr"`
	Definition      []struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"definition"`
}

type CollectionRequest struct {
	XMLName   xml.Name `xml:"collection"`
	Text      string   `xml:",chardata"`
	Adders    string   `xml:"adders"`
	Modifiers string   `xml:"modifiers"`
	Readers   string   `xml:"readers"`
}

type KMIGroup struct {
	XMLName   xml.Name `xml:"group"`
	Text      string   `xml:",chardata"`
	Name      string   `xml:"name,attr"`
	Type      string   `xml:"type,attr"`
	Source    string   `xml:"source,attr"`
	Account   string   `xml:"account,attr"`
	Adders    string   `xml:"adders"`
	Modifiers string   `xml:"modifiers"`
}

type KMIDefinition struct {
	XMLName       xml.Name     `xml:"definition"`
	Text          string       `xml:",chardata"`
	Type          string       `xml:"type,attr"`
	ExpirePeriod  string       `xml:"expire_period,omitempty"`
	RefreshPeriod string       `xml:"refresh_period,omitempty"`
	AutoGenerate  string       `xml:"auto_generate,omitempty"`
	Options       []*KMIOption `xml:"option"`
}

type KMIDefinitionResponse struct {
	XMLName       xml.Name `xml:"definition"`
	Text          string   `xml:",chardata"`
	Name          string   `xml:"name,attr"`
	Source        string   `xml:"source,attr"`
	Type          string   `xml:"type,attr"`
	Modified      string   `xml:"modified,attr"`
	AutoGenerate  string   `xml:"auto_generate"`
	ExpirePeriod  string   `xml:"expire_period"`
	RefreshPeriod string   `xml:"refresh_period"`
	Option        []struct {
		Text   string `xml:",chardata"`
		Name   string `xml:"name,attr"`
		Source string `xml:"source,attr"`
	} `xml:"option"`
	Secret []struct {
		Text  string `xml:",chardata"`
		Index string `xml:"index,attr"`
	} `xml:"secret"`
}

type kmigenerator interface {
	RequestPayload() ([]byte, error)
}

type OpaqueSecret struct {
	XMLName xml.Name `xml:"secret"`
	Text    string   `xml:",chardata"`
	Block   struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"block"`
}
