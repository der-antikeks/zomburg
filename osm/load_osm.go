package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"os"
	"time"
)

type Tag struct {
	Key   string `xml:"k,attr"`
	Value string `xml:"v,attr"`
}

type Node struct {
	Id        int64     `xml:"id,attr"`
	Lat       float64   `xml:"lat,attr"`
	Lon       float64   `xml:"lon,attr"`
	User      string    `xml:"user,attr"`
	Uid       int64     `xml:"uid,attr"`
	Visible   bool      `xml:"visible,attr"`
	Version   int       `xml:"version,attr"`
	Timestamp time.Time `xml:"timestamp,attr"`
	Tags      []Tag     `xml:"tag"`
}

type Nd struct {
	Ref int64 `xml:"ref,attr"`
}

type Way struct {
	Id        int64     `xml:"id,attr"`
	User      string    `xml:"user,attr"`
	Uid       int64     `xml:"uid,attr"`
	Visible   bool      `xml:"visible,attr"`
	Version   int       `xml:"version,attr"`
	Changeset int64     `xml:"changeset,attr"`
	Timestamp time.Time `xml:"timestamp,attr"`
	Nds       []Nd      `xml:"nd"`
	Tags      []Tag     `xml:"tag"`
}

type Member struct {
	Type string `xml:"type,attr"`
	Ref  int64  `xml:"ref,attr"`
	Role string `xml:"role,attr"`
}

type Relation struct {
	Id        int64     `xml:"id,attr"`
	User      string    `xml:"user,attr"`
	Uid       int64     `xml:"uid,attr"`
	Visible   bool      `xml:"visible,attr"`
	Version   int       `xml:"version,attr"`
	Changeset int64     `xml:"changeset,attr"`
	Timestamp time.Time `xml:"timestamp,attr"`
	Members   []Member  `xml:"member"`
	Tags      []Tag     `xml:"tag"`
}

func main() {
	var filename string
	flag.StringVar(&filename, "i", "", ".osm input file")
	flag.Parse()

	if filename == "" {
		log.Println("load_osm.go usage:")
		flag.PrintDefaults()
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	defer file.Close()

	nodes := make(map[int64]Node)
	ways := make(map[int64]Way)
	relations := make(map[int64]Relation)
	tags := make(map[string]bool)

	decoder := xml.NewDecoder(file)
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("err decoding: %v", err)
		}

		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "node":
				// a block of nodes containing location and tags
				// http://wiki.openstreetmap.org/wiki/Node
				var n Node
				decoder.DecodeElement(&n, &se)
				if _, ok := nodes[n.Id]; ok != false {
					log.Fatalf("duplicate node: %v", n)
				}
				nodes[n.Id] = n

				for _, t := range n.Tags {
					tags[t.Key] = true
				}
			case "way":
				// a block of ways containing references to its nodes for each way and tags
				// http://wiki.openstreetmap.org/wiki/Way
				var w Way
				decoder.DecodeElement(&w, &se)
				if _, ok := ways[w.Id]; ok != false {
					log.Fatalf("duplicate way: %v", w)
				}
				ways[w.Id] = w

				for _, t := range w.Tags {
					tags[t.Key] = true
				}
			case "relation":
				// a block of relations containing references to its members for each relation and tags
				// http://wiki.openstreetmap.org/wiki/Relation
				var r Relation
				decoder.DecodeElement(&r, &se)
				if _, ok := relations[r.Id]; ok != false {
					log.Fatalf("duplicate relation: %v", r)
				}
				relations[r.Id] = r

				for _, t := range r.Tags {
					tags[t.Key] = true
				}
			}
		}
	}

	log.Println("fin")
	log.Println("nodes:", len(nodes))
	log.Println("ways:", len(ways))
	log.Println("relations:", len(relations))

	out, err := os.Create("out.csv")
	if err != nil {
		log.Fatalf("could not create file: %v", err)
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	for tag := range tags {
		writer.Write([]string{tag})
	}
}
