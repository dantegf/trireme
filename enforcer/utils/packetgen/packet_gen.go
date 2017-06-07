package packetgen

import (
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	err error
)

func NewPacket() PacketManipulator {

	return &Packet{}
}

// func NewIPPacket(sip, dip string) PacketManipulator {
// 	p := NewPacket()
// 	p.AddIPLayer(sip, dip)
// 	return p
// }
//
// func NewIPTCPPacket(sip, dip, sport, dport string) PacketManipulator {
// 	p := NewIPPacket(sip, dip)
// 	// p.AddTCPLayer(sport, dport)
// 	return p
// }

//Use this function to create an IP packet (IPv4)
func (p *Packet) AddIPLayer(srcIPstr string, dstIPstr string) error {

	if p.IPLayer != nil {
		return fmt.Errorf("IP Layer already exists")
	}

	var srcIP, dstIP net.IP

	//IP address of the source
	srcIP = net.ParseIP(srcIPstr)
	if srcIP == nil {
		log.Printf("non-ip target: %q\n", srcIPstr)
	}

	//IP address of the destination
	dstIP = net.ParseIP(dstIPstr)
	if dstIP == nil {
		log.Printf("non-ip target: %q\n", dstIPstr)
	}

	//IP packet header
	p.IPLayer = &layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}
	//fmt.Println(p.IPLayer)
	return nil
}

//Get IP checksum
func (p *Packet) GetIPChecksum() uint16 {
	return p.IPLayer.Checksum

}

//Use this function to generate a single IP or TCP or a complete packet in both layers
func (p *Packet) AddTCPLayer(srcPort layers.TCPPort, dstPort layers.TCPPort) error {

	if p.TCPLayer != nil {
		return fmt.Errorf("TCP Layer already exists")
	}
	//Port number of the source
	srcport := srcPort
	//Port number of the destination
	dstport := dstPort
	//TCP packet header
	p.TCPLayer = &layers.TCP{
		SrcPort: srcport,
		DstPort: dstport,
		Window:  1505,
		Urgent:  0,
		Seq:     11050,
		Ack:     0,
		ACK:     false,
		SYN:     false,
		FIN:     false,
		RST:     false,
		URG:     false,
		ECE:     false,
		CWR:     false,
		NS:      false,
		PSH:     false,
	}
	//Checksum cannot be computed without network layer. Set IP protocol to TCP
	p.TCPLayer.SetNetworkLayerForChecksum(p.IPLayer)
	//fmt.Println(p.TCPLayer)
	return nil

}

//Change TCP sequence number
func (p *Packet) ChangeTCPSequenceNumber(seqNum uint32) error {
	p.TCPLayer.Seq = seqNum

	return nil
}

//Change TCP acknowledgement number
func (p *Packet) ChangeTCPAcknowledgementNumber(ackNum uint32) error {
	p.TCPLayer.Ack = ackNum
	return nil
}

//Change TCP window
func (p *Packet) ChangeTCPWindow(window uint16) error {
	p.TCPLayer.Window = window
	return nil
}

//Set TCP SYN flag to true
func (p *Packet) SetTCPSyn() error {
	p.TCPLayer.SYN = true
	p.TCPLayer.ACK = false
	p.TCPLayer.FIN = false
	return nil
}

//Set TCP SYN and ACK flag to true
func (p *Packet) SetTCPSynAck() error {
	p.TCPLayer.SYN = true
	p.TCPLayer.ACK = true
	p.TCPLayer.FIN = false
	return nil
}

//Set TCP ACK flag to true
func (p *Packet) SetTCPAck() error {
	p.TCPLayer.SYN = false
	p.TCPLayer.ACK = true
	p.TCPLayer.FIN = false
	return nil
}

func (p *Packet) GetTCPSyn() bool {

	return p.TCPLayer.SYN
}

func (p *Packet) GetTCPAck() bool {

	return p.TCPLayer.ACK
}

func (p *Packet) GetTCPFin() bool {

	return p.TCPLayer.FIN
}

//Add new payload to TCP layer
func (p *Packet) NewTCPPayload(newPayload string) error {
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	tcpPayloadBuf := gopacket.NewSerializeBuffer()
	payload := gopacket.Payload([]byte(newPayload))
	err = gopacket.SerializeLayers(tcpPayloadBuf, opts, p.IPLayer, p.TCPLayer, payload)
	if err != nil {
		panic(err)
	}

	return nil
}

//Get TCP Checksum
func (p *Packet) GetTCPChecksum() uint16 {
	return p.TCPLayer.Checksum

}

//Get IP Packet created
func (p *Packet) GetIPPacket() layers.IPv4 {
	return *p.IPLayer
}

//Get TCP Packet created
func (p *Packet) GetTCPPacket() layers.TCP {
	return *p.TCPLayer
}

//Display method
func (p *Packet) DisplayTCPPacket() {
	fmt.Println(p.TCPLayer)
}

//Create a new buffer and
//Convert it into bytes
func (p *Packet) ToBytes() [][]byte {

	ethernetLayer := layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0xFF, 0xAA, 0xFA, 0xAA, 0xFF, 0xAA},
		DstMAC:       net.HardwareAddr{0xBD, 0xBD, 0xBD, 0xBD, 0xBD, 0xBD},
		EthernetType: layers.EthernetTypeIPv4,
	}

	opts := gopacket.SerializeOptions{
		FixLengths:       true, //fix lengths based on the payload (data)
		ComputeChecksums: true, //compute checksum based on the payload during serialization
	}
	p.TCPLayer.SetNetworkLayerForChecksum(p.IPLayer)
	//serializing the layers to create a packet
	packetBuf := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(packetBuf, opts, &ethernetLayer, p.IPLayer, p.TCPLayer)
	bytes1 := packetBuf.Bytes()
	bytes1WithoutEthernet := bytes1[14:]
	var finalBytes [][]byte
	finalBytes = append(finalBytes, bytes1WithoutEthernet)

	return finalBytes
}

func NewTCPPacketFlow(sip string, dip string, sport layers.TCPPort, dport layers.TCPPort) PacketFlowManipulator {
	pf := &PacketFlow{
		SIP:   sip,
		DIP:   dip,
		SPort: sport,
		DPort: dport,
		Flow:  make([]PacketManipulator, 0),
	}
	return pf
}

func (p *PacketFlow) GenerateTCPFlow(bytePacket [][]byte) []PacketManipulator {

	// Create TCP Syn Packet
	np := NewPacket()
	np.AddIPLayer(p.SIP, p.DIP)
	np.AddTCPLayer(p.SPort, p.DPort)
	np.SetTCPSyn()
	np.ChangeTCPSequenceNumber(0)
	np.ChangeTCPAcknowledgementNumber(0)
	packet, ok := np.(*Packet)
	if !ok {
		return nil
	}
	//fmt.Println(packet.IPLayer)
	p.Flow = append(p.Flow, packet)

	// Create TCP SynAck Packet
	np2 := NewPacket()
	np2.AddIPLayer(p.DIP, p.SIP)
	np2.AddTCPLayer(p.DPort, p.SPort)
	np2.SetTCPSynAck()
	np2.ChangeTCPSequenceNumber(0)
	np2.ChangeTCPAcknowledgementNumber(packet.TCPLayer.Seq + 1)
	packet2, ok := np2.(*Packet)
	if !ok {
		return nil
	}
	//fmt.Println(packet2.IPLayer)
	p.Flow = append(p.Flow, packet2)

	// Create TCP Ack Packet
	np3 := NewPacket()
	np3.AddIPLayer(p.SIP, p.DIP)
	np3.AddTCPLayer(p.SPort, p.DPort)
	np3.SetTCPAck()
	np3.ChangeTCPSequenceNumber(packet2.TCPLayer.Ack)
	np3.ChangeTCPAcknowledgementNumber(packet2.TCPLayer.Seq + 1)
	packet3, ok := np3.(*Packet)
	if !ok {
		return nil
	}
	//fmt.Println(packet3.IPLayer)
	p.Flow = append(p.Flow, packet3)
	//fmt.Println(p.Flow[2].ToBytes())

	return p.Flow

}
func (p *PacketFlow) GenerateTCPFlowPayload(newPayload string) [][]byte {
	return nil
}

func (p *PacketFlow) GetSynPackets() []PacketManipulator {
	var synPackets []PacketManipulator
	for j := 0; j < len(p.Flow); j++ {
		if p.Flow[j].GetTCPSyn() == true && p.Flow[j].GetTCPAck() == false && p.Flow[j].GetTCPFin() == false {
			synPackets = append(synPackets, p.Flow[j])
		}
	}
	return synPackets
}
func (p *PacketFlow) GetSynAckPackets() [][]byte {
	return nil
}
func (p *PacketFlow) GetAckPackets() [][]byte {
	return nil
}

func (p *PacketFlow) GetNthPacket(index int) PacketManipulator {
	for i := 0; i < len(p.Flow); i++ {
		if index == i {
			return p.Flow[i]
		} else {
			return nil
		}
	}
	panic("Index out of range")
}
