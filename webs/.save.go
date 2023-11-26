package webs

type save struct {
	Desc string `json:"desc"`
	Mac  string `json:"mac"`
	Time string `json:"time"`
}

var SaveFile []save
