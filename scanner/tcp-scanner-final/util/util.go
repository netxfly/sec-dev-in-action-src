package util

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/malfunkt/iprange"
	"github.com/urfave/cli"

	"sec-dev-in-action-src/scanner/tcp-scanner-final/scanner"
	"sec-dev-in-action-src/scanner/tcp-scanner-final/vars"
)

func GetPorts(selection string) ([]int, error) {
	ports := []int{}
	if selection == "" {
		return ports, nil
	}

	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid port selection segment: '%s'", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[0])
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[1])
			}

			if p1 > p2 {
				return nil, fmt.Errorf("Invalid port range: %d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}

		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", r)
			} else {
				ports = append(ports, port)
			}
		}
	}
	return ports, nil
}

func GetIpList(ips string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}

	list := addressList.Expand()
	return list, err
}

func IsRoot() bool {
	return os.Geteuid() == 0
}

func CheckRoot() {
	if !IsRoot() {
		fmt.Println("must run with root")
		os.Exit(0)
	}
}

func Scan(ctx *cli.Context) error {
	if ctx.IsSet("iplist") {
		vars.Host = ctx.String("iplist")
	}

	if ctx.IsSet("port") {
		vars.Port = ctx.String("port")
	}

	if ctx.IsSet("mode") {
		vars.Mode = ctx.String("mode")
	}

	if ctx.IsSet("timeout") {
		vars.Timeout = ctx.Int("timeout")
	}

	if ctx.IsSet("concurrency") {
		vars.ThreadNum = ctx.Int("concurrency")
	}

	if strings.ToLower(vars.Mode) == "syn" {
		CheckRoot()
	}

	ips, err := GetIpList(vars.Host)
	ports, err := GetPorts(vars.Port)
	tasks, n := scanner.GenerateTask(ips, ports)
	_ = n
	scanner.RunTask(tasks)
	scanner.PrintResult()
	return err
}
