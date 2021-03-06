package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/dns/dnsmessage"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	//"strconv"
)

type DNSServer interface {
	Listen()
	Query(Packet)
}

// DNSService is an implementation of DNSServer interface.
type DNSService struct {
	conn *net.UDPConn
	//book       store
	//memo       addrBag
	forwarders []net.UDPAddr
}

// Packet carries DNS packet payload and sender address.
type Packet struct {
	addr    net.UDPAddr
	message dnsmessage.Message
}

const (
	// DNS server default port
	udpPort int = 53
	// DNS packet max length
	packetLen int = 512
)

func buildAns(q dnsmessage.Question) []dnsmessage.Resource {
	//	var rType dnsmessage.Type
	var rBody dnsmessage.ResourceBody
	var r []dnsmessage.Resource = nil
	var ip net.IP = nil
	var it []net.IP = nil
	fmt.Println(q.Name)
	sName := strings.TrimRight(string(q.Name.Data[:q.Name.Length]), ".")
	rType := q.Type
	hoster := nameKey{"A", sName}
	fmt.Println(sName)
	fmt.Println(q.Name)
	fmt.Println(nameKey{"A", sName})
	fmt.Println(mapNameId)
	id := mapNameId[hoster]
	//	if (id =! nil) {
	if (id == 0 && sName == "localhost") || id > 0 {
		r = make([]dnsmessage.Resource, len(ghosts.HostRecords[id].IP4))
		fmt.Println(id)
		it = ghosts.HostRecords[id].getIP4()
		for n := range it {
			ip = it[n]
			fmt.Println(ip)
			rBody = &dnsmessage.AResource{A: [4]byte{ip[12], ip[13], ip[14], ip[15]}}
			r[n] = dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  q.Name,
					Type:  rType,
					Class: q.Class,
					TTL:   300,
				},
				Body: rBody,
			}
		}
		return r
		//	} else {
	} else {
		it, _ := net.LookupHost(sName)
		fmt.Println(it)

		r = make([]dnsmessage.Resource, len(it))
		//ip := ghosts.HostRecords[id].IP4

		for n := range it {
			fmt.Println(it[n])
			ip = net.ParseIP(it[n]).To16()
			fmt.Println(ip)
			if ip.To4() == nil {
				rBody = &dnsmessage.AAAAResource{AAAA: [16]byte{ip[0], ip[1], ip[2], ip[3], ip[4], ip[5], ip[6], ip[7], ip[8], ip[9], ip[10], ip[11], ip[12], ip[13], ip[14], ip[15]}}
			} else {
				rBody = &dnsmessage.AResource{A: [4]byte{ip[12], ip[13], ip[14], ip[15]}}
			}
			r[n] = dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  q.Name,
					Type:  rType,
					Class: q.Class,
					TTL:   300,
				},
				Body: rBody,
			}
		}
		//}	}}
		return r
	}
}

func toHeader(name string, sType string) (h dnsmessage.ResourceHeader, err error) {
	h.Name, err = dnsmessage.NewName("localhost.")
	if err != nil {
		return
	}
	h.Type = dnsmessage.TypeA
	return h, err
}

func (s *DNSService) Listen() {
	var err error
	fmt.Println("DNSService")
	//s.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: udpPort})
	s2, err := net.ResolveUDPAddr("udp", "10.1.10.27:53")
	s.conn, err = net.ListenUDP("udp", s2)
	fmt.Println(s)
	if err != nil {
		log.Fatal(err)
	}
	defer s.conn.Close()

	for {
		buffer := make([]byte, 1024)
		//		_, addr, err := s.conn.ReadFromUDP(buf)
		var m dnsmessage.Message
		for {
			n, addr, err := s.conn.ReadFromUDP(buffer)
			if err != nil {
				log.Println(err)
				continue
			}
			err = m.Unpack(buffer)
			if err != nil {
				log.Println(err)
				continue
			}
			if len(m.Questions) == 0 {
				continue
			}
			fmt.Println(m.Header.GoString())
			fmt.Println(m.Questions[0].GoString())
			fmt.Println(m.Questions[0].Type.GoString())
			fmt.Println(m.Questions[0].Type)
			for i := range m.Questions {
				q := m.Questions[i]
				var newM dnsmessage.Message
				switch q.Type {
				case dnsmessage.TypeA:
					resource := buildAns(q)
					//ans, _ := toHeader("localhost.", "TypeA") //rType := dnsmessage.TypeA
					//data := []byte(buffer[0:n])
					newM.Header = m.Header
					//newM.Answers[0].Header = dnsmessage.ResourceHeader{Name: q.Name, Type: dnsmessage.TypeA, Class: q.Class, TTL: 1, Length: 1024}
					for x := range resource {
						newM.Answers = append(newM.Answers, resource[x])
					}
					packed, _ := newM.Pack()
					_, err = s.conn.WriteToUDP(packed, addr)

				//					ip := nil
				//					if ip == nil {
				//						return none, errIPInvalid
				//					}
				//					rBody = &dnsmessage.AResource{A: [4]byte{ip[12], ip[13], ip[14], ip[15]}}
				default:
					fmt.Println("1")
				}
			}

			data := []byte(buffer[0:n])
			_, err = s.conn.WriteToUDP(data, addr)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(100)
		}

		for i := range m.Questions {
			m.Questions[i].GoString()
			m.Questions[i].GoString()
			fmt.Println(m.Questions[i].GoString())
		}
	}
}

// Query lookup answers for DNS message.

////
///
///
type HostRecords struct {
	//Slice [Array] of monsters
	HostRecords []HostRecord `json:"HostRecords"`
}

type HostRecord struct {
	HostName string   `json:"HostName"`
	IP4      []net.IP `json:"Address4"`
	IP6      []net.IP `json:"Address6"`
	Id       int      `json:"Id"`
	DNSType  string   `json:"DNSType"` //We only need the name and XP of the monsters, however other critical stats are included for future expansion
}
type nameKey struct {
	DNSType  string `json:"DNSType"` //We only need the name and XP of the monsters, however other critical stats are included for future expansion
	HostName string `json:"HostName"`
}

func (h HostRecord) getIP4() []net.IP {
	return h.IP4
}

func (h HostRecord) getIP6() []net.IP {
	return h.IP6
}

func (h HostRecord) getDNSType() string {
	return h.DNSType
}

func (h HostRecord) getName() string {
	return h.HostName
}

func (h HostRecord) setId(i int) {
	h.Id = i
}

func New(forwarders []net.UDPAddr) DNSService {
	return DNSService{
		//		book:       store{data: make(map[string]entry), rwDirPath: rwDirPath},
		//		memo:       addrBag{data: make(map[string][]net.UDPAddr)},
		forwarders: forwarders,
	}
}

func jsonDns() HostRecords {
	jsonFile, err := os.Open("DNS.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var hosts HostRecords

	json.Unmarshal(byteValue, &hosts)
	return hosts
}

var ghosts = jsonDns()

var DNSTypes = [10]string{"CNAME", "A", "AAAA", "ALIAS", "PTR", "SOA", "TXT", "SRV", "MX", "NS"}

var m = make(map[string][]HostRecord)
var mapNameId = make(map[nameKey]int)

func main() {

	jsonFile, err := os.Open("DNS.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var hosts HostRecords

	json.Unmarshal(byteValue, &hosts)
	fmt.Println(hosts.HostRecords)
	fmt.Println(hosts.HostRecords[1].getIP4())
	fmt.Println(len(hosts.HostRecords))
	fmt.Println(DNSTypes[9])
	for n := range DNSTypes {
		for i := range hosts.HostRecords {
			if hosts.HostRecords[i].getDNSType() == DNSTypes[n] {
				//numberI, _ := strconv.Atoi(i)
				hosts.HostRecords[i].setId(i)
				m[DNSTypes[n]] = append(m[DNSTypes[n]], hosts.HostRecords[i])
				mapNameId[nameKey{DNSTypes[n], hosts.HostRecords[i].getName()}] = i

			}
		}
	}
	fmt.Println(m["A"])
	fmt.Println(m["MX"])
	fmt.Println(m["A"][1])
	var s = DNSService{}
	fmt.Println("Before")
	s.Listen()
	fmt.Println("After")
}
