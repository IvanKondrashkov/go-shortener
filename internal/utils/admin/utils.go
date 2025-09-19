package admin

import (
	"fmt"
	"net"
)

// IsIPInSubnet проверяет, принадлежит ли IP адрес указанной подсети
func IsIPInSubnet(ipStr, subnetCIDR string) (bool, error) {
	if subnetCIDR == "" {
		return false, nil
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, nil
	}

	_, subnet, err := net.ParseCIDR(subnetCIDR)
	if err != nil {
		return false, fmt.Errorf("invalid subnet CIDR: %w", err)
	}

	return subnet.Contains(ip), nil
}
