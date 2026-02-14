package lssnmp

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
)

// Snmp is wrapped some convinince SNMP interface for SNMP operation
type Snmp struct {
	Host    string
	Oid     string
	Version string
	CommStr string
	Port    int
	Timeout int
	Retries int
}

// init
func (s *Snmp) init() {
	s.Port = 161
	s.Timeout = 2
	s.Version = "v2c"
	s.CommStr = "public"
	s.Retries = 1
}

// GetReturnInterface  Get to get SNMP value by oid string
func (s *Snmp) GetReturnInterface(oid string) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
		}
	}()

	gs := s.newGoSNMP()
	if err := gs.Connect(); err != nil {
		return "", fmt.Errorf("error connecting snmp: %v", err)
	}
	defer gs.Conn.Close()

	// Perform GET
	pkt, err := gs.Get([]string{oid})
	if err != nil {
		return "", err
	}
	if len(pkt.Variables) == 0 {
		return "", fmt.Errorf("no variable returned for OID: %s", oid)
	}
	v := pkt.Variables[0]

	// Check for SNMP errors encoded in the variable type
	switch v.Type {
	case gosnmp.NoSuchObject, gosnmp.NoSuchInstance:
		return "", fmt.Errorf("no Such Object available on this agent at this OID: %s", oid)
	case gosnmp.Null:
		return "", fmt.Errorf("value is null for OID: %s", oid)
	}
	return v.Value, nil
}

// Get to get SNMP value by oid string
func (s *Snmp) Get(oid string) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered", r)
		}
	}()

	sSnmpValue, err := s.GetReturnInterface(oid)
	if sSnmpValue == nil {
		return "", err
	}
	return interfaceToString(sSnmpValue)

	/*
		value := strings.Trim(fmt.Sprintf("%v", sSnmpValue), ".")
		if value == "[128 0]" { // [128 0]  (kept for compatibility with original WapSNMP behavior)
			return "", fmt.Errorf("no Such Object available on this agent at this OID: %s", oid)
		}
		return value, nil*/
}

func interfaceToString(v interface{}) (string, error) {
	// 使用 type switch 判断 b 的实际类型
	switch v := v.(type) {
	case nil:
		return "", errors.New("返回值为 nil")
	case string:
		return fmt.Sprintf("%s", v), nil
	case []byte:
		return string(v), nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case float64:
		fmt.Printf("类型是 float64: %.2f\n", v)
	case bool:
		return fmt.Sprintf("%t", v), nil
	case []interface{}:
		//fmt.Printf("类型是 []interface{}: %v\n", v)
		return fmt.Sprintf("%v", v), nil
	}
	return fmt.Sprintf("%v", v), errors.New(fmt.Sprintf("未知数据类型 %T", v))
}

// GetMultiple to get SNMP values by oid slice
func (s *Snmp) GetMultiple(oids []string) (map[string]string, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in GetMultiple", r)
		}
	}()

	gs := s.newGoSNMP()
	if err := gs.Connect(); err != nil {
		return nil, fmt.Errorf("error connecting snmp: %v", err)
	}
	defer gs.Conn.Close()

	pkt, err := gs.Get(oids)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for i, v := range pkt.Variables {
		// Use the original OID from the request to preserve the key format
		oid := oids[i]
		if v.Type == gosnmp.NoSuchObject || v.Type == gosnmp.NoSuchInstance {
			res[oid] = ""
			continue
		}
		res[oid] = fmt.Sprintf("%v", v.Value)
	}
	return res, nil
}

// WalkResult to stroe walk result with key-value pair in struct
type WalkResult struct {
	Oid   string
	Value string
}

// SortedWalkSlice sorted slice by WalkResult.Oid
type SortedWalkSlice []WalkResult

func (s SortedWalkSlice) Len() int {
	return len(s)
}
func (s SortedWalkSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortedWalkSlice) Less(i, j int) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered Less", r)
		}
	}()
	v1 := strings.Split(s[i].Oid, ".")
	v2 := strings.Split(s[j].Oid, ".")
	if len(v1) != len(v2) {
		return false
	}
	n1, err := strconv.Atoi(v1[len(v1)-1])
	if err != nil {
		return false
	}
	n2, err := strconv.Atoi(v2[len(v2)-1])
	if err != nil {
		return false
	}
	return n1 < n2
}

// SortWalk snmp walk function return with a sorted result slice
func (s *Snmp) SortWalk(oid string) (SortedWalkSlice, error) {
	result, err := s.WalkReturnInterface(oid)
	if err != nil {
		return nil, err
	}
	res := make(SortedWalkSlice, 0, len(result))
	for k, v := range result {
		res = append(res, WalkResult{Oid: strings.Trim(k, "."), Value: fmt.Sprintf("%+v", v)})
	}
	sort.Sort(res)
	return res, nil
}

// WalkReturnInterface Walk to get do snmp walk
func (s *Snmp) WalkReturnInterface(oid string) (map[string]interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Walk(oid)", oid)
		}
	}()

	gs := s.newGoSNMP()
	if err := gs.Connect(); err != nil {
		return nil, fmt.Errorf("error connecting snmp: %v", err)
	}
	defer gs.Conn.Close()

	results := make(map[string]interface{})
	err := gs.Walk(oid, func(pdu gosnmp.SnmpPDU) error {
		// pdu.Name is the full OID (including the root part)
		results[pdu.Name] = pdu.Value
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Walk to get do snmp walk
func (s *Snmp) Walk(oid string) (map[string]string, error) {
	result, err := s.WalkReturnInterface(oid)
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for k, v := range result {
		tv, _ := interfaceToString(v)
		res[strings.Trim(k, ".")] = fmt.Sprintf("%v", tv)
	}
	return res, nil
}

// Set use gosnmp to do snmp set
func (s *Snmp) Set(oid string, value interface{}) (interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered set(oid)", oid)
		}
	}()

	gs := s.newGoSNMP()
	if err := gs.Connect(); err != nil {
		return nil, fmt.Errorf("error connecting snmp: %v", err)
	}
	defer gs.Conn.Close()

	pdu, err := buildPdu(oid, value)
	if err != nil {
		return nil, err
	}

	pkt, err := gs.Set([]gosnmp.SnmpPDU{pdu})
	if err != nil {
		return nil, err
	}
	if len(pkt.Variables) == 0 {
		return nil, fmt.Errorf("no response variable for OID %s", oid)
	}
	return pkt.Variables[0].Value, nil
}

// MultiSet use gosnmp to do multi set
/*
	toset := make(map[string]interface{})
	toset[oid1]="abc"
	toset[oid2]=1
	ret, err := MultiSet1(toset)
*/
func (s *Snmp) MultiSet(toset map[string]interface{}) (map[string]interface{}, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered MultiSet", r)
		}
	}()

	gs := s.newGoSNMP()
	if err := gs.Connect(); err != nil {
		return nil, fmt.Errorf("error connecting snmp: %v", err)
	}
	defer gs.Conn.Close()

	pdus := make([]gosnmp.SnmpPDU, 0, len(toset))
	for oid, val := range toset {
		pdu, err := buildPdu(oid, val)
		if err != nil {
			return nil, fmt.Errorf("failed to build PDU for OID %s: %v", oid, err)
		}
		pdus = append(pdus, pdu)
	}

	pkt, err := gs.Set(pdus)
	if err != nil {
		return nil, err
	}

	// Build result map from response variables, keyed by original OIDs.
	// Note: The response order should match the request order.
	result := make(map[string]interface{})
	for i, oid := range pdus {
		if i < len(pkt.Variables) {
			result[oid.Name] = pkt.Variables[i].Value
		} else {
			result[oid.Name] = nil
		}
	}
	return result, nil
}

// newGoSNMP creates a gosnmp.GoSNMP instance configured from the Snmp struct.
func (s *Snmp) newGoSNMP() *gosnmp.GoSNMP {
	version := gosnmp.Version2c
	switch strings.ToLower(s.Version) {
	case "1", "v1", "snmpv1":
		version = gosnmp.Version1
	case "2", "2c", "v2", "v2c", "snmpv2", "snmpv2c":
		version = gosnmp.Version2c
	default:
		// fallback to v2c (the original code returned an error in such cases,
		// but we handle version mapping earlier; here we just keep the default)
	}

	port := s.Port
	if port == 0 {
		port = 161
	}
	timeout := s.Timeout
	if timeout == 0 {
		timeout = 2
	}
	retries := s.Retries
	if retries == 0 {
		retries = 1
	}
	community := s.CommStr
	if community == "" {
		community = "public"
	}

	return &gosnmp.GoSNMP{
		Target:    s.Host,
		Port:      uint16(port),
		Community: community,
		Version:   version,
		Timeout:   time.Duration(timeout) * time.Second,
		Retries:   retries,
		// Set MaxRepetitions to a reasonable value for Walk/BulkWalk
		MaxRepetitions: 50,
	}
}

// buildPdu creates a gosnmp.SnmpPDU from an OID and a Go value.
// It attempts to infer the correct ASN.1 type based on the value's kind.
func buildPdu(oid string, value interface{}) (gosnmp.SnmpPDU, error) {
	pdu := gosnmp.SnmpPDU{
		Name: oid,
	}

	switch v := value.(type) {
	case string:
		pdu.Type = gosnmp.OctetString
		pdu.Value = v
	case int:
		pdu.Type = gosnmp.Integer
		pdu.Value = v
	case int32:
		pdu.Type = gosnmp.Integer
		pdu.Value = v
	case int64:
		pdu.Type = gosnmp.Counter64
		pdu.Value = v
	case uint:
		pdu.Type = gosnmp.Gauge32
		pdu.Value = v
	case uint32:
		pdu.Type = gosnmp.Gauge32
		pdu.Value = v
	case uint64:
		pdu.Type = gosnmp.Counter64
		pdu.Value = v
	case []byte:
		pdu.Type = gosnmp.OctetString
		pdu.Value = v
	case nil:
		pdu.Type = gosnmp.Null
		pdu.Value = nil
	default:
		return pdu, fmt.Errorf("unsupported value type for SNMP set: %T", value)
	}
	return pdu, nil
}
