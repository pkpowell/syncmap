package syncmap

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"
	"unique"

	"github.com/fxamacker/cbor/v2"
)

var jsonData = `
{
  "authTokens": [
    null
  ],
  "authorizationEndpoint": "",
  "capabilities": [
    {
      "default": false,
      "id": 1,
      "rules": [
        {
          "type": "ACTION_ACCEPT"
        }
      ]
    }
  ],
  "clientId": "",
  "creationTime": 1751222314511,
  "dns": [],
  "enableBroadcast": true,
  "id": "95987162f3023a29",
  "ipAssignmentPools": [],
  "mtu": 2800,
  "multicastLimit": 32,
  "name": "test network",
  "nwid": "95987162f3023a29",
  "objtype": "network",
  "private": true,
  "remoteTraceLevel": 0,
  "remoteTraceTarget": null,
  "revision": 2,
  "routes": [],
  "rules": [
    {
      "not": false,
      "or": false,
      "type": "ACTION_ACCEPT"
    }
  ],
  "rulesSource": "",
  "ssoEnabled": false,
  "tags": [],
  "v4AssignMode": {
    "zt": false
  },
  "v6AssignMode": {
    "6plane": false,
    "rfc4193": false,
    "zt": false
  }
}
`

type V4AssignMode struct {
	ZT bool `cbor:"ZT,omitempty" json:"zt,omitempty"`
}

type V6AssignMode struct {
	SixPlane bool `cbor:"6plane,omitempty" json:"6plane,omitempty"`
	Rfc4193  bool `cbor:"Rfc4193,omitempty" json:"rfc4193,omitempty"`
	ZT       bool `cbor:"ZT,omitempty" json:"zt,omitempty"`
}
type IPAssignmentPool struct {
	IPRangeStart string `cbor:"IPRangeStart,omitempty" json:"ipRangeStart,omitempty"`
	IPRangeEnd   string `cbor:"IPRangeEnd,omitempty" json:"ipRangeEnd,omitempty"`
}
type Rules []Rule
type Rule struct {
	Type       string `cbor:"Type,omitempty" json:"type,omitempty"`
	Or         bool   `cbor:"Or,omitempty" json:"or,omitempty"`
	Not        bool   `cbor:"Not,omitempty" json:"not,omitempty"`
	EtherType  int    `cbor:"EtherType,omitempty" json:"etherType,omitempty"`
	Start      int    `cbor:"Start,omitempty" json:"start,omitempty"`
	End        int    `cbor:"End,omitempty" json:"end,omitempty"`
	IPProtocol int    `cbor:"IPProtocol,omitempty" json:"ipProtocol,omitempty"`
	Value      int    `cbor:"Value,omitempty" json:"value,omitempty"`
}
type ZTRoutes struct {
	Target string `cbor:"Target" json:"target,omitempty"`
	Via    string `cbor:"Via" json:"via,omitempty"`
	Flags  int    `cbor:"Flags" json:"flags,omitempty"`
	Metric int    `cbor:"Metric" json:"metric,omitempty"`
}
type Capability struct {
	Default bool   `cbor:"Default,omitempty" json:"default,omitempty"`
	ID      int    `cbor:"ID,omitempty" json:"id,omitempty"`
	Rules   []Rule `cbor:"Rules,omitempty" json:"rules,omitempty"`
}
type Tag any
type Capabilities []Capability
type Tags []Tag
type ZTNetwork struct {
	ID                string             `cbor:"ID,omitempty" json:"id,omitempty"`
	NWID              string             `cbor:"NWID,omitempty" json:"nwid,omitempty"`
	Objtype           string             `cbor:"Objtype,omitempty" json:"objtype,omitempty"`
	Name              string             `cbor:"Name,omitempty" json:"name,omitempty"`
	CreationTime      int                `cbor:"CreationTime,omitempty" json:"creationTime,omitempty"`
	Private           bool               `cbor:"Private,omitempty" json:"private,omitempty"`
	EnableBroadcast   bool               `cbor:"EnableBroadcast,omitempty" json:"enableBroadcast,omitempty"`
	V4AssignMode      V4AssignMode       `cbor:"V4AssignMode,omitempty" json:"v4AssignMode,omitempty"`
	V6AssignMode      V6AssignMode       `cbor:"V6AssignMode,omitempty" json:"v6AssignMode,omitempty"`
	MTU               int                `cbor:"MTU,omitempty" json:"mtu,omitempty"`
	MulticastLimit    int                `cbor:"MulticastLimit,omitempty" json:"multicastLimit,omitempty"`
	Revision          int                `cbor:"Revision,omitempty" json:"revision,omitempty"`
	Routes            []ZTRoutes         `cbor:"Routes,omitempty" json:"routes,omitempty"`
	IPAssignmentPools []IPAssignmentPool `cbor:"IPAssignmentPools,omitempty" json:"ipAssignmentPools,omitempty"`
	Rules             Rules              `cbor:"Rules,omitempty" json:"rules,omitempty"`
	Capabilities      Capabilities       `cbor:"Capabilities,omitempty" json:"capabilities,omitempty"`
	Tags              Tags               `cbor:"Tags,omitempty" json:"tags,omitempty"`
	RemoteTraceTarget string             `cbor:"RemoteTraceTarget,omitempty" json:"remoteTraceTarget,omitempty"`
	RemoteTraceLevel  int                `cbor:"RemoteTraceLevel,omitempty" json:"remoteTraceLevel,omitempty"`

	// MemberCount int `cbor:"MemberCount,omitempty" json:"memberCount,omitempty"`
}

type TestType struct {
	Field string
	Array []int
}

type TestBool struct{}

func (t *TestType) GetID() string {
	return t.Field
}

func (t *TestType) IDX() string {
	return t.Field
}

func (t *TestType) Del(bool) {}

// func (t *TestType) GetMTX() *sync.RWMutex { return nil }
// func (t *Device) GetMTX() *sync.RWMutex   { return t.mtx }
func (t *ZTPeerID) GetID() string {
	return t.Address
}

func (t *ZTPeerID) IDX() string {
	return t.Address
}

func (t *ZTPeerID) Del(bool)  {}
func (t *ZTNetwork) Del(bool) {}
func (t *ZTNetwork) GetID() string {
	return t.NWID
}
func (t *Test4) Del(bool) {}
func (t *Test4) GetID() string {
	return t.one
}

// func (t *ZTPeerID) GetMTX() *sync.RWMutex { return t.mtx }
func (t *Device) GetID() string {
	return t.ID
}

func (t *Device) IDX() string {
	return t.ID
}

func (t *Device) Del(bool) {}

func (t *TestBool) IDX()     {}
func (t *TestBool) Del(bool) {}

// func (t *TestBool) GetMTX() *sync.RWMutex { return nil }
func (t *TestBool) GetID() string {
	return ""
}

var (
	p  = NewPointerMap[*TestType]()
	pc = NewCollection[*TestType, *TestBool]()
	c  = NewCollection[string, *TestType]()
	d  = NewCollection[string, *Device]()

	u = NewUniqueCollection[string, *TestType]()
	v = NewUniqueCollection[string, *ZTNetwork]()
	w = NewUniqueCollection[string, *Test4]()
)

func BenchmarkPointerMapAdd(b *testing.B) {
	for i := range b.N {
		p.Add(&TestType{
			Field: fmt.Sprintf("test-%d", i),
			Array: []int{i, 2, 3},
		})
	}
}

func BenchmarkCollBoolAdd(b *testing.B) {
	for i := range b.N {
		pc.Add(&TestType{
			Field: fmt.Sprintf("test-%d", i),
			Array: []int{i, 2, 3},
		}, &TestBool{})
	}
}

var s = &struct {
	TestType
	_field string
	_array []int
}{}

// var ss = []ZTNetwork{}

var t = &s.TestType

// var tt = &ss.ZTNetwork

func BenchmarkCollectionAdd(b *testing.B) {
	for i := range b.N {
		t.Field = fmt.Sprintf("test-%d", i)
		t.Array = []int{i, 2, 3}
		c.Add(t.Field, t)
	}
}
func BenchmarkCollectionAddComp(b *testing.B) {
	for i := range b.N {
		t.Field = fmt.Sprintf("test-%d", i)
		t.Array = []int{i, 2, 3}
		c.AddCompare(t.Field, t)
	}
}
func BenchmarkUniqueCollectionAdd(b *testing.B) {
	for i := range b.N {
		t.Field = fmt.Sprintf("test-%d", i)
		t.Array = []int{i, 2, 3}
		u.Add(t.Field, t)
	}

	var updated bool
	for i := range b.N {
		t.Field = fmt.Sprintf("test-%d", i)
		t.Array = []int{i, 2, 3}
		updated = u.Add(t.Field, t)
		if updated {
			b.Logf("%s updated!", t.Field)
		}
	}

	t.Field = fmt.Sprintf("test-%d", 500)
	t.Array = []int{500, 2, 3}
	if u.Add(t.Field, t) {
		b.Logf("%s updated!", t.Field)
	}
	t.Field = fmt.Sprintf("test-%d", 500)
	t.Array = []int{500, 50, 3}
	if u.Add(t.Field, t) {
		b.Logf("%s updated!", t.Field)
	}
}

func BenchmarkCollectionExists(b *testing.B) {
	for i := range b.N {
		c.Exists(fmt.Sprintf("test-%d", i))
	}
}

func BenchmarkCollectionGet(b *testing.B) {
	for i := range b.N {
		c.Get(fmt.Sprintf("test-%d", i))
	}
}

func BenchmarkCollectionGetP(b *testing.B) {
	for i := range b.N {
		var d *TestType
		c.GetP(fmt.Sprintf("test-%d", i), &d)

		//fmt.Println("d: ", d)
	}
}

func BenchmarkCollectionGetAll(b *testing.B) {
	for range b.N {
		for range c.All() {

		}
	}
}

func BenchmarkGet(b *testing.B) {
	b.Run("add", BenchmarkCollectionAdd)
	b.Run("get", BenchmarkCollectionGet)
	// b.Run("getp", BenchmarkCollectionGetP)
	// b.Run("getall", BenchmarkCollectionGetAll)
}
func BenchmarkUniqueGet(b *testing.B) {
	for i := range b.N {
		t.Field = fmt.Sprintf("test-%d", i)
		t.Array = []int{i, 2, 3}
		u.Add(t.Field, t)
	}
	d, ok := u.Get(fmt.Sprintf("test-%d", 20))
	if ok {
		b.Logf("got data %v", d)
	}
}
func BenchmarkUniqueGet2(b *testing.B) {
	var data ZTNetwork
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		b.Log("error", err)
		return
	}
	for i := range b.N {
		v.Add(strconv.Itoa(i), &data)
	}

	updated := v.Add("200", &data)
	if updated {
		b.Log("data changed")
	}
	d, ok := v.Get("123")
	if ok {
		b.Logf("got data %v", d)
	}
}
func BenchmarkUniqueGet3(b *testing.B) {
	var data ZTNetwork
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		b.Log("error", err)
		return
	}

	updated := v.Add("12345", &data)
	if updated {
		b.Log("1 data changed")
	}

	updated = v.Add("12345", &data)
	if updated {
		b.Log("2 data changed")
	}
	d, ok := v.Get("12345")
	if ok {
		b.Logf("got data %v", d)
	}
}
func BenchmarkUniqueGet5(b *testing.B) {
	var data = &ZTNetwork{
		ID:      "123456",
		NWID:    "123456",
		Private: true,
	}
	var data2 = &ZTNetwork{
		ID:      "123456",
		NWID:    "123456",
		Private: true,
	}

	h1 := unique.Make(data)
	h2 := unique.Make(data2)

	b.Log("h1==h2", h1 == h2)

	// updated := v.Add("12345", data)
	// if updated {
	// 	b.Log("1 data changed")
	// }

	// updated = v.Add("12345", data)
	// if updated {
	// 	b.Log("2 data changed")
	// }
	// d, ok := v.Get("12345")
	// if ok {
	// 	b.Logf("got data %v", d)
	// }
}
func BenchmarkEncToBytesJSON(b *testing.B) {
	var data ZTNetwork
	jsonBytes := []byte(jsonData)
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		b.Log("error", err)
		return
	}

	for range b.N {
		json.Marshal(data)
	}
}
func BenchmarkEncToBytesCBOR(b *testing.B) {
	var data ZTNetwork
	jsonBytes := []byte(jsonData)
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		b.Log("error", err)
		return
	}

	for range b.N {
		cbor.Marshal(data)
	}
}

// func BenchmarkEncToBytesGOB(b *testing.B) {
// 	var data ZTNetwork
// 	jsonBytes := []byte(jsonData)
// 	err := json.Unmarshal(jsonBytes, &data)
// 	if err != nil {
// 		b.Log("error", err)
// 		return
// 	}

// 	for range b.N {
// 		gob.GobEncoder(data)
// 	}
// }

type Seven struct {
	eight string
	nine  string
}
type Test4 struct {
	one   string `json:"one"`
	two   string
	three string
	four  int
	five  int
	six   int
	seven Seven
}

func BenchmarkUniqueGet4(b *testing.B) {
	var data = Test4{
		one:   "one",
		two:   "two",
		three: "three",
		four:  4,
		five:  5,
		six:   6,
		seven: Seven{
			eight: "eight",
		},
	}
	var data2 = Test4{
		one:   "one",
		two:   "two",
		three: "three",
		four:  4,
		five:  5,
		six:   6,
		seven: Seven{
			eight: "eight",
			nine:  "nine",
		},
	}

	h1 := unique.Make(data)
	b.Log("h1", h1)

	h2 := unique.Make(data2)
	b.Log("h2", h2)
	b.Log("h1==h2", h1 == h2)

}

func BenchmarkExists(b *testing.B) {
	b.Run("add", BenchmarkCollectionAdd)
	b.Run("exist", BenchmarkCollectionExists)
	// b.Run("getp", BenchmarkCollectionGetP)
	// b.Run("getall", BenchmarkCollectionGetAll)
}

// func TestPutGet(b *testing.T) {
// 	var tests = map[string]struct {
//         a, b *TestType
//         want *TestType
//     }{
//     "one": {
// 		a:   &TestType{
// 			Field: "test-1",
// 			Array:  []int{1, 2, 3},
// 		},
// 		b: &TestType{
// 			Field: "test-2",
// 			Array:  []int{4, 2, 3},
// 		},
// 		want: &TestType{
// 			Field: "test-1",
// 			Array:  []int{1, 2, 3},
// 		},
// 	},
// }

//     for name, tt := range tests {
//         testname := fmt.Sprintf("%d,%d", tt.a, tt.b)
//         t.Run(name, func(t *testing.T) {
//             ans := IntMin(tt.a, tt.b)
//             if ans != tt.want {
//                 t.Errorf("got %d, want %d", ans, tt.want)
//             }
//         })
//     }
// }

func BenchmarkCollectionPut(b *testing.B) {
	for i := range b.N {
		d.Add(fmt.Sprintf("id-%d", i), NewDevice())
	}
}
func TestCollectionGetP(t *testing.T) {
	d := NewDevice()
	for i := range 100 {
		d.ZTPeers.Add(fmt.Sprintf("ID-%d", i), &ZTPeerID{Address: fmt.Sprintf("ID-%d", i)})
	}
	for i := range 100 {
		var n *ZTPeerID
		d.ZTPeers.GetP(fmt.Sprintf("ID-%d", i), &n)
	}
}

func NewDevice() (dc *Device) {
	devContainer := struct {
		Device
		_ID                       string
		_IntuneUUID               string
		_DisplayName              string
		_Hostname                 string
		_Serial                   string
		_Added                    time.Time
		_LastSyncDateTime         time.Time
		_TotalStorageSpaceInBytes int
		_FreeStorageSpaceInBytes  int
		_WiFiMacAddress           string
		_MachineModel             string
		_ZTPeers                  *Collection[string, *ZTPeerID]
	}{}
	dc = &devContainer.Device
	(*dc).ID = devContainer._ID
	(*dc).IntuneUUID = devContainer._IntuneUUID
	(*dc).DisplayName = devContainer._DisplayName
	(*dc).Hostname = devContainer._Hostname
	(*dc).Serial = devContainer._Serial
	(*dc).Added = devContainer._Added
	(*dc).LastSyncDateTime = devContainer._LastSyncDateTime
	(*dc).TotalStorageSpaceInBytes = devContainer._TotalStorageSpaceInBytes
	(*dc).FreeStorageSpaceInBytes = devContainer._FreeStorageSpaceInBytes
	(*dc).WiFiMacAddress = devContainer._WiFiMacAddress
	(*dc).MachineModel = devContainer._MachineModel
	(*dc).ZTPeers = NewCollection[string, *ZTPeerID]()

	return
}

type InfoStat struct {
	Hostname             string `json:"hostname"`
	Uptime               uint64 `json:"uptime"`
	BootTime             uint64 `json:"bootTime"`
	Procs                uint64 `json:"procs"`           // number of processes
	OS                   string `json:"os"`              // ex: freebsd, linux
	Platform             string `json:"platform"`        // ex: ubuntu, linuxmint
	PlatformFamily       string `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion      string `json:"platformVersion"` // version of the complete OS
	KernelVersion        string `json:"kernelVersion"`   // version of the OS kernel (if available)
	KernelArch           string `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	VirtualizationSystem string `json:"virtualizationSystem"`
	VirtualizationRole   string `json:"virtualizationRole"` // guest or host
	HostID               string `json:"hostId"`             // ex: uuid
}

type SPSoftwareDataType struct {
	Name            string `cbor:"Name" json:"_name"`
	BootMode        string `cbor:"BootMode" json:"boot_mode"`
	BootVolume      string `cbor:"BootVolume" json:"boot_volume"`
	KernelVersion   string `cbor:"KernelVersion" json:"kernel_version"`
	LocalHostName   string `cbor:"LocalHostName" json:"local_host_name"`
	OSVersion       string `cbor:"OSVersion" json:"os_version"`
	SecureVM        string `cbor:"SecureVM" json:"secure_vm"`
	SystemIntegrity string `cbor:"SystemIntegrity" json:"system_integrity"`
	Uptime          string `cbor:"Uptime" json:"uptime"`
	UserName        string `cbor:"UserName" json:"user_name"`
}
type User struct {
	Name    string `cbor:"Name" json:"name,omitempty"`
	Admin   bool   `cbor:"Admin" json:"admin,omitempty"`
	Date    int64  `cbor:"Date" json:"date,omitempty"`
	Added   int64  `cbor:"Added" json:"added,omitempty"`
	Deleted bool   `cbor:"Deleted" json:"deleted,omitempty"`

	ImageSaved bool `cbor:"ImageSaved" json:"imageSaved,omitempty"`
}
type Zerotier struct {
	Address string `cbor:"Address" json:"address,omitempty"`
	Clock   int    `cbor:"Clock" json:"clock,omitempty"`
	// Settings             ZTSettings `cbor:"Settings" json:"settings,omitempty"`
	Online               *bool  `cbor:"Online" json:"online,omitempty"`
	PlanetWorldId        int    `cbor:"PlanetWorldId" json:"planetWorldId,omitempty"`
	PlanetWorldTimestamp int    `cbor:"PlanetWorldTimestamp" json:"planetWorldTimestamp,omitempty"`
	PublicIdentity       string `cbor:"PublicIdentity" json:"publicIdentity,omitempty"`
	TCPFallbackActive    *bool  `cbor:"TCPFallbackActive" json:"tcpFallbackActive,omitempty"`
	Version              string `cbor:"Version" json:"version,omitempty"`
	VersionBuild         int    `cbor:"VersionBuild" json:"versionBuild,omitempty"`
	VersionMajor         int    `cbor:"VersionMajor" json:"versionMajor,omitempty"`
	VersionMinor         int    `cbor:"VersionMinor" json:"versionMinor,omitempty"`
	VersionRev           int    `cbor:"VersionRev" json:"versionRev,omitempty"`
}

type ZTPeerID struct {
	Hostname string `cbor:"Hostname" json:"hostname"`
	Address  string `cbor:"Address" json:"address"`
	ID       string `cbor:"ID" json:"id"`

	// mtx *sync.RWMutex `cbor:"-" json:"-"`
	// TimeStamp int64 `cbor:"TimeStamp" json:"timeStamp"`
}
type Device struct {
	// mtx                      *sync.RWMutex                  `cbor:"-" json:"-"`
	ID                       string                         `cbor:"ID" json:"id"`
	Hostname                 string                         `cbor:"HostName" json:"hostname"`
	DisplayName              string                         `cbor:"DisplayName" json:"displayName,omitempty"`
	MachineModel             string                         `cbor:"MachineModel" json:"machineModel"`
	LastSyncDateTime         time.Time                      `cbor:"LastSyncDateTime" json:"lastSyncDateTime,omitempty"`
	Added                    time.Time                      `cbor:"Added" json:"added,omitempty"`
	Serial                   string                         `cbor:"Serial" json:"serial"`
	HASH                     string                         `cbor:"HASH" json:"hash" hash:"ignore"`
	HostInfo                 InfoStat                       `cbor:"Hostinfo" json:"hostinfo,omitempty"`
	Uptime                   uint64                         `cbor:"Uptime" json:"uptime,omitempty" hash:"ignore"`
	LastSeen                 int64                          `cbor:"LastSeen" json:"lastseen,omitempty" hash:"ignore"`
	OSSoftware               []SPSoftwareDataType           `cbor:"OSSoftware" json:"osSoftware,omitempty"`
	TotalMem                 uint64                         `cbor:"TotalMem" json:"totalMem,omitempty"`
	UserData                 map[string]*User               `cbor:"UserData" json:"userData"`
	AppVersion               string                         `cbor:"AppVersion" json:"appVersion,omitempty"`
	BuildDate                string                         `cbor:"BuildDate" json:"buildDate,omitempty"`
	GoVersion                string                         `cbor:"GoVersion" json:"goVersion,omitempty"`
	Zerotier                 *Zerotier                      `cbor:"Zerotier" json:"zerotier,omitempty"`
	ZTPeers                  *Collection[string, *ZTPeerID] `cbor:"-" json:"-" hash:"ignore"`
	ZTPeerData               map[string]*ZTPeerID           `cbor:"ZTPeers" json:"ztPeers,omitempty" hash:"ignore"`
	Retired                  bool                           `cbor:"Retired" json:"retired" hash:"ignore"`
	IntuneUUID               string                         `cbor:"IntuneUUID" json:"intuneUUID"`
	WiFiMacAddress           string                         `cbor:"WiFiMacAddress" json:"wiFiMacAddress"`
	EthernetMACAddress       string                         `cbor:"EthernetMACAddress" json:"ethernetMACAddress"`
	TotalStorageSpaceInBytes int                            `cbor:"TotalStorageSpaceInBytes" json:"totalStorageSpaceInBytes"`
	FreeStorageSpaceInBytes  int                            `cbor:"FreeStorageSpaceInBytes" json:"freeStorageSpaceInBytes"`
}

type ManagedDevice struct {
	ID                                          string              `cbor:"ID" json:"id"`
	UserId                                      string              `cbor:"UserId" json:"userId"`
	DeviceName                                  string              `cbor:"DeviceName" json:"deviceName"`
	OwnerType                                   string              `cbor:"OwnerType" json:"ownerType"`
	ManagedDeviceOwnerType                      string              `cbor:"ManagedDeviceOwnerType" json:"managedDeviceOwnerType"`
	ManagementState                             string              `cbor:"ManagementState" json:"managementState"`
	EnrolledDateTime                            string              `cbor:"EnrolledDateTime" json:"enrolledDateTime"`
	LastSyncDateTime                            string              `cbor:"LastSyncDateTime" json:"lastSyncDateTime"`
	ChassisType                                 string              `cbor:"ChassisType" json:"chassisType"`
	OperatingSystem                             string              `cbor:"OperatingSystem" json:"operatingSystem"`
	DeviceType                                  string              `cbor:"DeviceType" json:"deviceType"`
	ComplianceState                             string              `cbor:"ComplianceState" json:"complianceState"`
	JailBroken                                  string              `cbor:"JailBroken" json:"jailBroken"`
	ManagementAgent                             string              `cbor:"ManagementAgent" json:"managementAgent"`
	OSVersion                                   string              `cbor:"OSVersion" json:"osVersion"`
	EASActivated                                bool                `cbor:"EASActivated" json:"easActivated"`
	EASDeviceId                                 string              `cbor:"EASDeviceId" json:"easDeviceId"`
	EASActivationDateTime                       string              `cbor:"EASActivationDateTime" json:"easActivationDateTime"`
	AADRegistered                               any                 `cbor:"AADRegistered" json:"aadRegistered"`
	AzureADRegistered                           any                 `cbor:"AzureADRegistered" json:"azureADRegistered"`
	DeviceEnrollmentType                        string              `cbor:"DeviceEnrollmentType" json:"deviceEnrollmentType"`
	LostModeState                               string              `cbor:"LostModeState" json:"lostModeState"`
	ActivationLockBypassCode                    string              `cbor:"ActivationLockBypassCode" json:"activationLockBypassCode"`
	EmailAddress                                string              `cbor:"EmailAddress" json:"emailAddress"`
	AzureActiveDirectoryDeviceId                string              `cbor:"AzureActiveDirectoryDeviceId" json:"azureActiveDirectoryDeviceId"`
	AzureADDeviceId                             string              `cbor:"AzureADDeviceId" json:"azureADDeviceId"`
	DeviceRegistrationState                     string              `cbor:"DeviceRegistrationState" json:"deviceRegistrationState"`
	DeviceCategoryDisplayName                   string              `cbor:"DeviceCategoryDisplayName" json:"deviceCategoryDisplayName"`
	IsSupervised                                bool                `cbor:"IsSupervised" json:"isSupervised"`
	ExchangeLastSuccessfulSyncDateTime          string              `cbor:"ExchangeLastSuccessfulSyncDateTime" json:"exchangeLastSuccessfulSyncDateTime"`
	ExchangeAccessState                         string              `cbor:"ExchangeAccessState" json:"exchangeAccessState"`
	ExchangeAccessStateReason                   string              `cbor:"ExchangeAccessStateReason" json:"exchangeAccessStateReason"`
	RemoteAssistanceSessionUrl                  string              `cbor:"RemoteAssistanceSessionUrl" json:"remoteAssistanceSessionUrl"`
	RemoteAssistanceSessionErrorDetails         string              `cbor:"RemoteAssistanceSessionErrorDetails" json:"remoteAssistanceSessionErrorDetails"`
	IsEncrypted                                 bool                `cbor:"IsEncrypted" json:"isEncrypted"`
	UserPrincipalName                           string              `cbor:"UserPrincipalName" json:"userPrincipalName"`
	EnrolledByUserPrincipalName                 string              `cbor:"EnrolledByUserPrincipalName" json:"enrolledByUserPrincipalName"`
	Model                                       string              `cbor:"Model" json:"model"`
	Manufacturer                                string              `cbor:"Manufacturer" json:"manufacturer"`
	IMEI                                        string              `cbor:"IMEI" json:"imei"`
	ComplianceGracePeriodExpirationDateTime     string              `cbor:"ComplianceGracePeriodExpirationDateTime" json:"complianceGracePeriodExpirationDateTime"`
	SerialNumber                                string              `cbor:"SerialNumber" json:"serialNumber"`
	PhoneNumber                                 string              `cbor:"PhoneNumber" json:"phoneNumber"`
	AndroidSecurityPatchLevel                   string              `cbor:"AndroidSecurityPatchLevel" json:"androidSecurityPatchLevel"`
	UserDisplayName                             string              `cbor:"UserDisplayName" json:"userDisplayName"`
	ConfigurationManagerClientEnabledFeatures   []string            `cbor:"ConfigurationManagerClientEnabledFeatures" json:"configurationManagerClientEnabledFeatures"`
	WiFiMacAddress                              string              `cbor:"WiFiMacAddress" json:"wiFiMacAddress"`
	DeviceHealthAttestationState                any                 `cbor:"DeviceHealthAttestationState" json:"deviceHealthAttestationState"`
	SubscriberCarrier                           string              `cbor:"SubscriberCarrier" json:"subscriberCarrier"`
	MEID                                        string              `cbor:"MEID" json:"meid"`
	TotalStorageSpaceInBytes                    int                 `cbor:"TotalStorageSpaceInBytes" json:"totalStorageSpaceInBytes"`
	FreeStorageSpaceInBytes                     int                 `cbor:"FreeStorageSpaceInBytes" json:"freeStorageSpaceInBytes"`
	ManagedDeviceName                           string              `cbor:"ManagedDeviceName" json:"managedDeviceName"`
	PartnerReportedThreatState                  string              `cbor:"PartnerReportedThreatState" json:"partnerReportedThreatState"`
	RetireAfterDateTime                         string              `cbor:"RetireAfterDateTime" json:"retireAfterDateTime"`
	PreferMdmOverGroupPolicyAppliedDateTime     string              `cbor:"PreferMdmOverGroupPolicyAppliedDateTime" json:"preferMdmOverGroupPolicyAppliedDateTime"`
	AutopilotEnrolled                           bool                `cbor:"AutopilotEnrolled" json:"autopilotEnrolled"`
	RequireUserEnrollmentApproval               bool                `cbor:"RequireUserEnrollmentApproval" json:"requireUserEnrollmentApproval"`
	ManagementCertificateExpirationDate         string              `cbor:"ManagementCertificateExpirationDate" json:"managementCertificateExpirationDate"`
	ICCID                                       string              `cbor:"ICCID" json:"iccid"`
	UDID                                        string              `cbor:"UDID" json:"udid"`
	RoleScopeTagIds                             []string            `cbor:"RoleScopeTagIds" json:"roleScopeTagIds"`
	WindowsActiveMalwareCount                   int                 `cbor:"WindowsActiveMalwareCount" json:"windowsActiveMalwareCount"`
	WindowsRemediatedMalwareCount               int                 `cbor:"WindowsRemediatedMalwareCount" json:"windowsRemediatedMalwareCount"`
	Notes                                       string              `cbor:"Notes" json:"notes"`
	ConfigurationManagerClientHealthState       any                 `cbor:"ConfigurationManagerClientHealthState" json:"configurationManagerClientHealthState"`
	ConfigurationManagerClientInformation       any                 `cbor:"ConfigurationManagerClientInformation" json:"configurationManagerClientInformation"`
	EthernetMacAddress                          string              `cbor:"EthernetMacAddress" json:"ethernetMacAddress"`
	PhysicalMemoryInBytes                       int                 `cbor:"PhysicalMemoryInBytes" json:"physicalMemoryInBytes"`
	ProcessorArchitecture                       string              `cbor:"ProcessorArchitecture" json:"processorArchitecture"`
	SpecificationVersion                        string              `cbor:"SpecificationVersion" json:"specificationVersion"`
	JoinType                                    string              `cbor:"JoinType" json:"joinType"`
	SKUFamily                                   string              `cbor:"SKUFamily" json:"skuFamily"`
	SecurityPatchLevel                          string              `cbor:"SecurityPatchLevel" json:"securityPatchLevel"`
	SKUNumber                                   int                 `cbor:"SKUNumber" json:"skuNumber"`
	ManagementFeatures                          string              `cbor:"ManagementFeatures" json:"managementFeatures"`
	EnrollmentProfileName                       string              `cbor:"EnrollmentProfileName" json:"enrollmentProfileName"`
	BootstrapTokenEscrowed                      bool                `cbor:"BootstrapTokenEscrowed" json:"bootstrapTokenEscrowed"`
	DeviceFirmwareConfigurationInterfaceManaged bool                `cbor:"DeviceFirmwareConfigurationInterfaceManaged" json:"deviceFirmwareConfigurationInterfaceManaged"`
	DeviceIdentityAttestationDetail             bool                `cbor:"DeviceIdentityAttestationDetail" json:"deviceIdentityAttestationDetail"`
	HardwareInformation                         HardwareInformation `cbor:"HardwareInformation" json:"hardwareInformation"`
	DeviceActionResults                         []string            `cbor:"DeviceActionResults" json:"deviceActionResults"`
	UsersLoggedOn                               []string            `cbor:"UsersLoggedOn" json:"usersLoggedOn"`
	ChromeOSDeviceInfo                          []string            `cbor:"ChromeOSDeviceInfo" json:"chromeOSDeviceInfo"`
	SupplementalDeviceDetails                   []string            `cbor:"SupplementalDeviceDetails" json:"supplementalDeviceDetails"`
}

type HardwareInformation struct {
	SerialNumber      string `cbor:"SerialNumber" json:"serialNumber"`
	TotalStorageSpace int    `cbor:"TotalStorageSpace" json:"totalStorageSpace"`
	FreeStorageSpace  int    `cbor:"FreeStorageSpace" json:"freeStorageSpace"`
}
