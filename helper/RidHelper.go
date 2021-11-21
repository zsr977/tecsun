package helper

import (
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	AnswerT int64 = iota
	AnswerDetailT
	AreaT
	BillboardT
	BillboardItemT
	BookingT
	BookingEmployeeT
	BookingMemberT
	CompanyT
	DeviceT
	FilesT
	HelperT
	ManagerT
	ManagerTenantT
	MecT
	MessageT
	PassageT
	QuestionListT
	QuestionNaireT
	TenantT
	UnitT
	UnitMemberT
	UrlT
	UserT
	UserTenantT
	UserUnionT
	VisitorT
	WorkflowT
	WorkflowFormT
	WorkflowNodeT
	WorkflowRecordT
	WorkflowRecordDetailT
	WorkflowRecordTaskT
	WorkflowTaskT
)

var (
	vpc     int64
	ipc     int64
	epoch   int64
	nodeMap map[int64]Business
	once    sync.Once
)

func Rid(table int64) string {
	n := nodeMap[table]
	n.node = table
	res, err := snowflake(n)
	for err != nil {
		res, err = snowflake(n)
	}
	return strconv.Itoa(int(res))
}

type Business struct {
	node     int64
	lastTime int64
	lastId   int64
	step     int64
}

func snowflake(n Business) (int64, error) {
	once.Do(initId)
	now := time.Now().Unix()
	if now == n.lastTime {
		n.step = n.step + 1
		if n.step > 4095 {
			return 0, errors.New("序列号溢出")
		}
	} else {
		n.step = 0
	}
	n.lastTime = now
	res := (now-epoch)<<29 | (vpc << 26) | (ipc << 18) | (n.node << 12) | n.step
	if n.lastId >= res {
		res += 1
	}
	n.lastId = res
	nodeMap[n.node] = n
	return res, nil
}

func initId() {
	vpcStr := os.Getenv("VPC")
	if vpcStr != "" {
		vpcInt, err := strconv.Atoi(vpcStr)
		if err != nil {
			panic(err)
		}
		vpc = int64(vpcInt)
	}
	ips := getLocalIP()
	if len(ips) == 0 {
		panic(errors.New("获取本地IP失败"))
	}
	ipArr := strings.Split(ips[0], ".")
	if len(ipArr) == 4 {
		ipInt, _ := strconv.Atoi(ipArr[3])
		ipc = int64(ipInt)
	}
	epoch = 1633017600
	nodeMap = make(map[int64]Business)
}

func getLocalIP() []string {
	var ipStr []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return ipStr
	}
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addr, _ := netInterfaces[i].Addrs()
			for _, address := range addr {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						ipStr = append(ipStr, ipNet.IP.String())
					}
				}
			}
		}
	}
	return ipStr
}
