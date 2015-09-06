package structures

import "encoding/xml"

// Service set triple play services
type Service struct {
	ID   string `json:"Id"`
	Vlan string `json:"Vlan"`
	Vpi  string `json:"Vpi"`
	Vci  string `json:"Vci"`
}

// Port define DSLAM ports
type Port struct {
	ID               string    `json:"Index"`
	Name             string    `json:"Name"`
	AdmState         string    `json:"Adm-State"`
	OprStateTxRateDs string    `json:"Opr-State/Tx-Rate-Ds"`
	CurOpMode        string    `json:"Cur-Op-Mode"`
	Service          []Service `json:"Service"`
}

// Card define DSLAM card
type Card struct {
	Name         string `json:"Name"`
	Slot         string `json:"Slot"`
	OperStatus   string `json:"Opers-Status"`
	ErrorStatus  string `json:"Error-Status"`
	Availability string `json:"Availability"`
	Port         []Port `json:"Port"`
}

// DSLAM define DSLAM
type DSLAM struct {
	ID        string  `json:"Id"`
	Name      string  `json:"Name"`
	Status    string  `json:"Status"`
	Address   string  `json:"Address"`
	User      string  `json:"User"`
	Password  string  `json:"Password"`
	Card      []Card  `json:"Card"`
	Internet  Service `json:"Internet"`
	Telephony Service `json:"Telephony"`
	Video     Service `json:"Video"`
}

// Licence define Licence
type Licence struct {
	ID        string `json:"Id"`
	StartDate string `json:"StartDate"`
	EndDate   string `json:"EndDate"`
	User      string `json:"User"`
}

// Data define Data
type Data struct {
	DSLAM   []DSLAM `json:"DSLAM"`
	Licence Licence `json:"Licence"`
}

// SessionResponse define session response structure when we establish DSLAM connection
type SessionResponse struct {
	ID     string `json:"ID"`
	Status string `json:"Status"`
}

// CommandOut define DSLAM return
type CommandOut struct {
	CommandOut []string `json:"CommandOut"`
	ReturnCode string   `json:"ReturnCode"`
}

// InfoAttr define XML info of commandOut
type InfoAttr struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// Parameter define XML parameter of commandOut
type Parameter struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// XMLCard define XML DSLAM card struct
type XMLCard struct {
	XMLName xml.Name     `xml:"instance"`
	ResID   string       `xml:"res-id"`
	Info    [][]InfoAttr `xml:"info"`
}

// ShowEquipmentSlot Define DSLAM return of command show equipment slot
type ShowEquipmentSlot struct {
	XMLName xml.Name  `xml:"runtime-data"`
	Card    []XMLCard `xml:"hierarchy>hierarchy>hierarchy>instance"`
}

// XMLPort define XML DSLAM port struct
type XMLPort struct {
	XMLName xml.Name     `xml:"instance"`
	ResID   string       `xml:"res-id"`
	Info    [][]InfoAttr `xml:"info"`
}

// ShowXdslOperationalDataLine Define DSLAM return of command show xdsl operational-data line
type ShowXdslOperationalDataLine struct {
	XMLName xml.Name  `xml:"runtime-data"`
	Port    []XMLPort `xml:"hierarchy>hierarchy>hierarchy>hierarchy>instance"`
}

// XMLService define XML DSLAM Service struct
type XMLService struct {
	XMLName   xml.Name      `xml:"instance"`
	ResID     string        `xml:"res-id"`
	Parameter [][]Parameter `xml:"parameter"`
}

// ShowXdslOperDataPortIndexBridgePort Define DSLAM return of command show xdsl oper-data-port *Index* bridge-port xml
type ShowXdslOperDataPortIndexBridgePort struct {
	XMLName xml.Name     `xml:"runtime-data"`
	Service []XMLService `xml:"hierarchy>hierarchy>hierarchy>hierarchy>instance"`
}
