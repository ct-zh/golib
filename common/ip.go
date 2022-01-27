package common

import (
	"errors"
	"fmt"
	"math/big"
	"net"
)

// IpToInt take IP address conversion to Int64
func IpToInt(ip string) (int64, error) {
	n := net.ParseIP(ip)
	if n == nil {
		return 0, errors.New("invalid ip")
	}
	trans := big.NewInt(0)
	return trans.SetBytes(n.To4()).Int64(), nil
}

// IpToStr converts an IP address of type in64 to a string
func IpToStr(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
