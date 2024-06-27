package data_xml

import "encoding/xml"

type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr
	Children []Node `xml:",any"`
	Text     string `xml:",chardata"`
}

type RequestXML struct {
	DBTable string `json:"db_table"`
	Root    string `json:"root"`
	Tags    []Tag  `json:"tags"`
}

type Tag struct {
	Parent string `json:"parent"`
	DB     string `json:"db"`
	Tag    string `json:"tag"`
}

type TagValue struct {
	Parent string
	DB     string
	Tag    string
	Value  interface{}
}
