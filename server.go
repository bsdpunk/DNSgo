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
	"sync"
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

	sName := strings.TrimRight(string(q.Name.Data[:q.Name.Length]), ".")
	rType := dnsmessage.TypeA
	hoster := nameKey{"A", sName}
	id := mapNameId[hoster]
	//	if (id =! nil) {
	if (id == 0 && sName == "localhost") || id > 0 {
		it = ghosts.HostRecords[id].getIP4()
		r = make([]dnsmessage.Resource, len(it))
		for n := range it {
			ip = it[n]
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
		//	} else
		return r
	} else {
		it, _ := net.LookupHost(sName)

		r = make([]dnsmessage.Resource, len(it))
		//ip := ghosts.HostRecords[id].IP4

		for n := range it {
			ip = net.ParseIP(it[n]).To16()
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
		return r
	}
}

func buildMX(q dnsmessage.Question) []dnsmessage.Resource {

	var rBody dnsmessage.ResourceBody
	var r []dnsmessage.Resource = nil
	var it []net.IP = nil //	var ip net.IP = nil
	sName := strings.TrimRight(string(q.Name.Data[:q.Name.Length]), ".")
	hoster := nameKey{"MX", sName}
	//rType := dnsmessage.TypeMX
	id := mapNameId[hoster]
	if (id == 0 && sName == "localhost") || id > 0 {
		it = ghosts.HostRecords[id].getIP4()
		fmt.Println(it)
		r = make([]dnsmessage.Resource, len(it))
		fmt.Println(it)
		fmt.Println(len(it))
		for g := range it {
			fmt.Println(g)
			mxName, _ := dnsmessage.NewName(ghosts.HostRecords[id].MailNames[g])
			var numeral uint16 = ghosts.HostRecords[id].Priority[g]
			fmt.Println(g)

			fmt.Println(ghosts.HostRecords[id].IP4[g])
			fmt.Println(ghosts.HostRecords[id].MailNames[g])
			fmt.Println(ghosts.HostRecords[id].Priority[g])
			fmt.Println(g)

			rBody = &dnsmessage.MXResource{Pref: numeral, MX: mxName}
			r[g] = dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  q.Name,
					Type:  dnsmessage.TypeMX,
					Class: q.Class,
					TTL:   300,
				},
				Body: rBody,
			}
			fmt.Println(g)

		}
		fmt.Println("----------------------------------------------")
		return r
	} else {
		var it, _ = net.LookupMX(string(q.Name.Data[:q.Name.Length]))

		r = make([]dnsmessage.Resource, len(it))
		for n := range it {
			fmt.Println(it[n])
			mxName, _ := dnsmessage.NewName(it[n].Host)
			numeral := it[n].Pref
			rBody = &dnsmessage.MXResource{Pref: numeral, MX: mxName}
			r[n] = dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  q.Name,
					Type:  dnsmessage.TypeMX,
					Class: q.Class,
					TTL:   300,
				},
				Body: rBody,
			}
			fmt.Println(r)
		}
		return r
	}

}

func buildNS(q dnsmessage.Question) []dnsmessage.Resource {

	var rBody dnsmessage.ResourceBody
	var r []dnsmessage.Resource = nil
	//	var ip net.IP = nil
	var it []net.NS = nil
	fmt.Println(q.Name)
	sName := strings.TrimRight(string(q.Name.Data[:q.Name.Length]), ".")
	hoster := nameKey{"NS", sName}
	//rType := dnsmessage.TypeMX
	id := mapNameId[hoster]

	fmt.Println(it)
	if (id == 0 && sName == "localhost") || id > 0 {
		fmt.Println("Do nothing")
	} else {
		it, _ := net.LookupNS(sName)
		fmt.Println("here")
		fmt.Println(it)
		fmt.Println("here")

		r = make([]dnsmessage.Resource, len(it))
		for n := range it {
			fmt.Println(it[n])
			nsName, _ := dnsmessage.NewName(it[n].Host)
			rBody = &dnsmessage.NSResource{NS: nsName}
			r[n] = dnsmessage.Resource{
				Header: dnsmessage.ResourceHeader{
					Name:  q.Name,
					Type:  dnsmessage.TypeMX,
					Class: q.Class,
					TTL:   300,
				},
				Body: rBody,
			}
			fmt.Println(r)
		}
		return r
	}
	return r
}

func buildPTR(q dnsmessage.Question) []dnsmessage.Resource {
	sAddr := strings.TrimRight(string(q.Name.Data[:q.Name.Length]), ".")
	sName := strings.TrimRight(string(q.Name.Data[:q.Name.Length]), ".")
	hoster := nameKey{"NS", sName}
	//var rBody dnsmessage.ResourceBody
	//var r []dnsmessage.Resource = nil //	var ip net.IP = nil //var it []net.NS = nil fmt.Println(q.Name) sAddr := strings.TrimRight(string(q.Name.Data[:q.Name.Length]), ".") hoster := nameKey{"PTR", sAddr} //rType := dnsmessage.TypeMX
	id := mapNameId[hoster]

	fmt.Println(q)
	//fmt.Println(it)
	var r []dnsmessage.Resource = nil
	if (id == 0 && sAddr == "127.0.0.1") || id > 0 {
		fmt.Println("Do nothing")
	} else {
		//it, _ := net.LookupNS(sName)
		ptr, _ := net.LookupAddr(sAddr)
		for _, ptrvalue := range ptr {
			fmt.Println(ptrvalue)
		}
		//fmt.Println(it[n])
		//var dot dnsmessage.Name
		//dot, _ = dnsmessage.NewName(".")
		ptrName := q.Name
		//rBody = &dnsmessage.PTRResource{PTR: dot}

		//r = make([]dnsmessage.Resource, 1)

		msg := dnsmessage.Message{
			Header: dnsmessage.Header{Response: true, Authoritative: true},
			Questions: []dnsmessage.Question{
				{
					Name:  ptrName,
					Type:  dnsmessage.TypeA,
					Class: dnsmessage.ClassINET,
				},
				{
					Name:  ptrName,
					Type:  dnsmessage.TypeA,
					Class: dnsmessage.ClassINET,
				},
			},
			Answers: []dnsmessage.Resource{
				{
					Header: dnsmessage.ResourceHeader{
						Name:  ptrName,
						Type:  dnsmessage.TypePTR,
						Class: dnsmessage.ClassINET,
					},
					Body: &dnsmessage.PTRResource{PTR: ptrName},
				},
				{
					Header: dnsmessage.ResourceHeader{
						Name:  ptrName,
						Type:  dnsmessage.TypePTR,
						Class: dnsmessage.ClassINET,
					},
					Body: &dnsmessage.PTRResource{PTR: ptrName},
				},
			},
		}

		return msg.Answers
	}
	return r
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
	s2, err := net.ResolveUDPAddr("udp", "127.0.0.1:53")
	s.conn, err = net.ListenUDP("udp", s2)
	fmt.Println(s)
	if err != nil {
		log.Fatal(err)
	}

	var count int64
	count = 0
	defer s.conn.Close()
	for {
		//		_, addr, err := s.conn.ReadFromUDP(buf)
		var m dnsmessage.Message
		buffer := make([]byte, 1024)

		_, addr, err := s.conn.ReadFromUDP(buffer)
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
		//fmt.Println(addr)
		for i := range m.Questions {
			q := m.Questions[i]
			var wg sync.WaitGroup
			wg.Add(len(m.Questions))
			go func(q dnsmessage.Question) {
				//data := []byte(buffer[0:n])
				count++
				defer wg.Done()
				fmt.Println(count)
				//	fmt.Println(len(m.Questions))
				var newMX dnsmessage.Message
				var newM dnsmessage.Message
				switch q.Type {
				case dnsmessage.TypeA:
					resource := buildAns(q)
					theParse(buffer)
					//ans, _ := toHeader("localhost.", "TypeA") //rType := dnsmessage.TypeA
					newM.Header = m.Header
					//newM.Answers[0].Header = dnsmessage.ResourceHeader{Name: q.Name, Type: dnsmessage.TypeA, Class: q.Class, TTL: 1, Length: 1024}
					for x := range resource {
						newM.Answers = append(newM.Answers, resource[x])
					}
					packed, _ := newM.Pack()
					_, err = s.conn.WriteToUDP(packed, addr)

				case dnsmessage.TypeMX:

					resource := buildMX(q)
					theParse(buffer)
					newMX.Header = m.Header
					//p = dnsmessage.Parser
					for x := range resource {
						newMX.Answers = append(newMX.Answers, resource[x])
					}
					fmt.Println(newMX)
					fmt.Println(newMX.Answers)
					fmt.Println(newMX.Answers)

					fmt.Println(q.Name)
					fmt.Println(newMX.Answers)

					fmt.Println(q.Name.GoString())
					fmt.Println(q.GoString())
					fmt.Println(q.Type.GoString())
					fmt.Println(q.Class.GoString())
					fmt.Println(q.Type.GoString())
					fmt.Println(q.Class.GoString())
					fmt.Println("--------------------------------------")
					packed, _ := newMX.Pack()
					_, err = s.conn.WriteToUDP(packed, addr)

				case dnsmessage.TypeNS:
					resource := buildNS(q)
					theParse(buffer)
					newM.Header = m.Header
					for x := range resource {
						newM.Answers = append(newM.Answers, resource[x])
					}
					packed, _ := newM.Pack()
					_, err = s.conn.WriteToUDP(packed, addr)

				case dnsmessage.TypePTR:
					msg := buildPTR(q)
					theParse(buffer)
					newM.Header = m.Header
					for x := range msg {
						newM.Answers = append(newM.Answers, msg[x])
					}

					fmt.Println(m.GoString())
					fmt.Println(newM.GoString())
					fmt.Println(newM.Answers)
					fmt.Println(msg)
					packed, _ := newM.Pack()
					_, err = s.conn.WriteToUDP(packed, addr)
					if err != nil {
						fmt.Println(err)
						os.Exit(100)
					}

				default:
					break
				}

			}(q)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(100)
		}

	}
}

// Query lookup answers for DNS Amessage.

func theParse(buf []byte) {

	//	var p dnsmessage.Parser
	//	if _, err := p.Start(m.Unpack()); err != nil {
	//		panic(err)
	//	}

	wantName := "localhost."

	var m dnsmessage.Message //buf := make([]byte, 2, 514)
	var err error = m.Unpack(buf)
	if err != nil {
		panic(err)
	}
	var p dnsmessage.Parser
	if _, err := p.Start(buf); err != nil {
		panic(err)
	}
	for {
		q, err := p.Question()
		if err == dnsmessage.ErrSectionDone {
			break
		}
		if err != nil {
			panic(err)
		}

		//		if q.Name.String() != "localhost" {
		//			continue
		//		}

		////		fmt.Println("Found question for name", "localhost")
		//		if err := p.SkipAllQuestions(); err != nil {
		//			panic(err)
		//}
		fmt.Println(q.GoString())
		fmt.Println(q.GoString())
		fmt.Println(q.GoString())
		fmt.Println(q.GoString())
		fmt.Println(q.GoString())
		fmt.Println(q.GoString())

		break

	}

	q, err := p.Question()

	var gotIPs []net.IP
	var mail []dnsmessage.MXResource
	var pref []dnsmessage.PTRResource
	//	h, err := p.AnswerHeader()
	//	if err == dnsmessage.ErrSectionDone {
	//		break
	//	}
	//	if err != nil {
	//		panic(err)
	//	}

	//	if (h.Type != dnsmessage.TypeA && h.Type != dnsmessage.TypeAAAA) || h.Class != dnsmessage.ClassINET {
	//		continue
	//	}
	//
	//		if !strings.EqualFold(h.Name.String(), wantName) {
	//			if err := p.SkipAnswer(); err != nil {
	//				panic(err)
	//			}
	//			continue
	//		}

	switch q.Type {
	case dnsmessage.TypeA:
		r, err := p.AResource()
		if err != nil {
			panic(err)
		}
		gotIPs = append(gotIPs, r.A[:])
	case dnsmessage.TypeAAAA:
		r, err := p.AAAAResource()
		if err != nil {
			panic(err)
		}
		gotIPs = append(gotIPs, r.AAAA[:])

	case dnsmessage.TypeMX:
		r, err := p.MXResource()
		if err != nil {
			panic(err)
		}
		mail = append(mail, r)
	case dnsmessage.TypePTR:
		r, err := p.PTRResource()
		if err != nil {
			panic(err)
		}
		pref = append(pref, r)

	default:
		break
	}

	fmt.Printf("Found A/AAAA records for name %s: %v\n", wantName, gotIPs)
	fmt.Println(mail)
	fmt.Println(pref)
}

////
///
///
type HostRecords struct {
	//Slice [Array] of monsters
	HostRecords []HostRecord `json:"HostRecords"`
}

type HostRecord struct {
	HostName  string   `json:"HostName"`
	IP4       []net.IP `json:"Address4"`
	IP6       []net.IP `json:"Address6"`
	Id        int      `json:"Id"`
	MailNames []string `json:"MailNames"`
	Priority  []uint16 `json:"Priority"`
	DNSType   string   `json:"DNSType"` //We only need the name and XP of the monsters, however other critical stats are included for future expansion
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
