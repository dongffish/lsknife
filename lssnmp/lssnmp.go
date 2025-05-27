package lssnmp

import (
	"errors"
	"fmt"
	"github.com/cdevr/WapSNMP"
	"sort"
	"strconv"
	"strings"
	"time"
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

// Get to get a SNMP value by oid string
func (s *Snmp) Get(oid string) (string, error) {
	sOid, err := wapsnmp.ParseOid(oid)
	if err != nil {
		//lslog.Logger.Debugf("%s snmpget -c %s -v %s -t %d -r %d %s %s", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid)
		return "", err
	}
	version := wapsnmp.SNMPv2c
	switch strings.ToLower(s.Version) {
	case "1", "v1", "snmpv1":
		version = wapsnmp.SNMPv1
	case "2", "2c", "v2", "v2c", "snmpv2", "snmpv2c":
		version = wapsnmp.SNMPv2c
	default:
		return "", errors.New("Not supported version: " + s.Version)
	}
	wsnmp, err := wapsnmp.NewWapSNMP(s.Host, s.CommStr, version, time.Duration(s.Timeout)*time.Second, s.Retries)
	if err != nil {
		//lslog.Logger.Debugf("%s snmpget -c %s -v %s -t %d -r %d %s %s", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid)
		return "", fmt.Errorf("error creating wsnmp => %v", err)
	}
	defer wsnmp.Close()
	sSnmpValue, err := wsnmp.Get(sOid)
	if sSnmpValue == nil {
		//lslog.Logger.Debugf("%s snmpget -c %s -v %s -t %d -r %d %s %s", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid)
		return "", err
	}
	/*
		tp := reflect.TypeOf(sSnmpValue)
		Logger.Infoln(oid, "返回类型：", tp.Kind(), "数据:", fmt.Sprintf("%v", sSnmpValue))
			switch tp.Kind() {
			case reflect.String:
				//
			case reflect.Slice:
				//
				//case reflect.
			}
	*/
	value := strings.Trim(fmt.Sprintf("%v", sSnmpValue), ".")
	if value == "[128 0]" { // [128 0]
		//lslog.Logger.Debugf("%s snmpget -c %s -v %s -t %d -r %d %s %s", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid)
		return "", fmt.Errorf("no Such Object available on this agent at this OID: %s", oid)
	}
	return value, nil
}

// GetMultiple to get SNMP values by oid slice
func (s *Snmp) GetMultiple(oids []string) (map[string]string, error) {
	version := wapsnmp.SNMPv2c
	switch strings.ToLower(s.Version) {
	case "1", "v1", "snmpv1":
		version = wapsnmp.SNMPv1
	case "2", "2c", "v2", "v2c", "snmpv2", "snmpv2c":
		version = wapsnmp.SNMPv2c
	default:
		return nil, errors.New("Not supported version: " + s.Version)
	}
	var oidSlice []wapsnmp.Oid
	var err error
	for _, v := range oids {
		sOid, err := wapsnmp.ParseOid(v)
		if err != nil {
			//lslog.Logger.Debugf("%s GetMultiple -c %s -v %s -t %d -r %d %s %v", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oids)
			return nil, err
		}
		oidSlice = append(oidSlice, sOid)
	}
	wsnmp, err := wapsnmp.NewWapSNMP(s.Host, s.CommStr, version, time.Duration(s.Timeout)*time.Second, s.Retries)
	if err != nil {
		//lslog.Logger.Debugf("%s GetMultiple -c %s -v %s -t %d -r %d %s %v", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oids)
		return nil, fmt.Errorf("error creating wsnmp => %v", err)
	}
	defer wsnmp.Close()
	result, err := wsnmp.GetMultiple(oidSlice)
	if err != nil {
		//lslog.Logger.Debugf("%s GetMultiple -c %s -v %s -t %d -r %d %s %v", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oids)
		return nil, err
	}
	res := make(map[string]string)
	for k, v := range result {
		res[strings.Trim(k, ".")] = fmt.Sprintf("%v", v)
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
	version := wapsnmp.SNMPv2c
	switch strings.ToLower(s.Version) {
	case "1", "v1", "snmpv1":
		version = wapsnmp.SNMPv1
	case "2", "2c", "v2", "v2c", "snmpv2", "snmpv2c":
		version = wapsnmp.SNMPv2c
	default:
		return nil, errors.New("Not supported version: " + s.Version)
	}
	sOid, err := wapsnmp.ParseOid(oid)
	if err != nil {
		//lslog.Logger.Debugf("%s snmpwalk -c %s -v %s -t %d -r %d %s %s", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid)
		return nil, err
	}

	wsnmp, err := wapsnmp.NewWapSNMP(s.Host, s.CommStr, version, time.Duration(s.Timeout)*time.Second, s.Retries)
	if err != nil {
		//lslog.Logger.Debugf("%s snmpwalk -c %s -v %s -t %d -r %d %s %s", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid)
		return nil, fmt.Errorf("error creating wsnmp => %v", err)
	}
	defer wsnmp.Close()
	result, err := wsnmp.GetTable(sOid)
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

// Walk to get do snmp walk
func (s *Snmp) Walk(oid string) (map[string]string, error) {
	version := wapsnmp.SNMPv2c
	switch strings.ToLower(s.Version) {
	case "1", "v1", "snmpv1":
		version = wapsnmp.SNMPv1
	case "2", "2c", "v2", "v2c", "snmpv2", "snmpv2c":
		version = wapsnmp.SNMPv2c
	default:
		return nil, errors.New("Not supported version: " + s.Version)
	}
	sOid, err := wapsnmp.ParseOid(oid)
	if err != nil {
		//lslog.Logger.Debugf("%s snmpwalk -c %s -v %s -t %d -r %d %s %s", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid)
		return nil, err
	}

	wsnmp, err := wapsnmp.NewWapSNMP(s.Host, s.CommStr, version, time.Duration(s.Timeout)*time.Second, s.Retries)
	if err != nil {
		//lslog.Logger.Debugf("%s snmpwalk -c %s -v %s -t %d -r %d %s %s", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid)
		return nil, fmt.Errorf("error creating wsnmp => %v", err)
	}
	defer wsnmp.Close()
	result, err := wsnmp.GetTable(sOid)
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	for k, v := range result {
		res[strings.Trim(k, ".")] = fmt.Sprintf("%+v", v)
	}
	return res, nil
}

// Set use wapsnmp to do snmp set
func (s *Snmp) Set(oid string, value interface{}) (interface{}, error) {
	version := wapsnmp.SNMPv2c
	switch strings.ToLower(s.Version) {
	case "1", "v1", "snmpv1":
		version = wapsnmp.SNMPv1
	case "2", "2c", "v2", "v2c", "snmpv2", "snmpv2c":
		version = wapsnmp.SNMPv2c
	default:
		return "", errors.New("Not supported version: " + s.Version)
	}
	sOid, err := wapsnmp.ParseOid(oid)
	if err != nil {
		//lslog.Logger.Debugf("%s snmpset -c %s -v %s -t %d -r %d %s %s %+v", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid, value)
		return nil, err
	}
	wsnmp, err := wapsnmp.NewWapSNMP(s.Host, s.CommStr, version, time.Duration(s.Timeout)*time.Second, s.Retries)
	if err != nil {
		//lslog.Logger.Debugf("%s snmpset -c %s -v %s -t %d -r %d %s %s %+v", lslog.GetLogPrefix(), s.CommStr, s.Version, s.Timeout, s.Retries, s.Host, oid, value)
		return nil, fmt.Errorf("error creating wsnmp => %v", err)
	}
	defer wsnmp.Close()
	return wsnmp.Set(sOid, value)
}

// MultiSet use wsnmp to do multi set
/*
	toset := make(map[string]interface{})
	toset[oid1]="abc"
	toset[oid2]=1
	ret, err := MultiSet1(toset)
*/
func (s *Snmp) MultiSet(toset map[string]interface{}) (map[string]interface{}, error) {
	version := wapsnmp.SNMPv2c
	switch strings.ToLower(s.Version) {
	case "1", "v1", "snmpv1":
		version = wapsnmp.SNMPv1
	case "2", "2c", "v2", "v2c", "snmpv2", "snmpv2c":
		version = wapsnmp.SNMPv2c
	default:
		return nil, errors.New("Not supported version: " + s.Version)
	}
	wsnmp, err := wapsnmp.NewWapSNMP(s.Host, s.CommStr, version, time.Duration(s.Timeout)*time.Second, s.Retries)
	if err != nil {
		return nil, fmt.Errorf("error creating wsnmp => %v", err)
	}
	defer wsnmp.Close()
	return wsnmp.SetMultiple(toset)
}
