package leases

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	dateFormat = "2006/01/02 15:04:05"
)

//Lease is the definition of an IPv4 lease
//file described on dhcpd.leases(5) man page
type Lease struct {
	IP         string
	Start      time.Time
	End        time.Time
	Tstp       time.Time
	Cltt       time.Time
	Binding    string
	Next       string
	Hardware   string
	UID        string
	Identifier string
}

var (
	errInvalidAttr = errors.New("Invalid attribute")
	errMalformed   = errors.New("malformed attribute")
)

//ParseLeases parse dhcpd.leases file and return the
//slice of Leases
func ParseLeases(reader io.ReadCloser) ([]Lease, error) {
	defer reader.Close()
	scan := bufio.NewScanner(reader)
	leases := make(map[string]Lease)
	lease := Lease{}
	for scan.Scan() {
		line := scan.Text()
		if isComment(line) {
			continue
		}
		if isEndOfLease(line) {
			if !strings.Contains(lease.Binding, "active") {
				lease = Lease{}
				continue
			}
			leases[lease.IP] = lease

			lease = Lease{}
			continue
		}
		var err error
		switch {
		case strings.Contains(line, "lease"):
			err = parseLease(strings.TrimSpace(line), &lease)
		case strings.Contains(line, "starts"):
			err = parseStart(strings.TrimSpace(line), &lease)
		case strings.Contains(line, "ends"):
			err = parseEnd(strings.TrimSpace(line), &lease)
		case strings.Contains(line, "tstp"):
			err = parseTstp(strings.TrimSpace(line), &lease)
		case strings.Contains(line, "cltt"):
			err = parseCltt(strings.TrimSpace(line), &lease)
		case strings.Contains(line, "binding"):
			err = parseBinding(strings.TrimSpace(line), &lease)
		case strings.Contains(line, "uid"):
			err = parseUID(strings.TrimSpace(line), &lease)
		case strings.Contains(line, "set vendor-class"):
			err = parseIdentifier(strings.TrimSpace(line), &lease)
		case strings.Contains(line, "hardware"):
			err = parseEthernet(strings.TrimSpace(line), &lease)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	var nl []Lease
	for _, l := range leases {
		nl = append(nl, l)
	}
	return nl, nil
}

func sanitize(attr string) string {
	attr = strings.Replace(attr, ";", "", 1)
	attr = strings.Replace(attr, `"`, "", 2)
	return attr
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "#")
}

func isEndOfLease(line string) bool {
	return strings.Contains(line, "}")
}

func parseStart(line string, lease *Lease) error {
	parts := strings.Split(line, " ")
	date := sanitize(parts[2] + " " + parts[3])
	fmt.Println(date)
	dt, err := time.Parse(dateFormat, date)
	if err != nil {
		return err
	}
	lease.Start = dt
	return err
}

func parseEnd(line string, lease *Lease) error {
	parts := strings.Split(line, " ")
	date := sanitize(parts[2] + " " + parts[3])
	dt, err := time.Parse(dateFormat, date)
	if err != nil {
		return err
	}
	lease.End = dt
	return err
}

func parseTstp(line string, lease *Lease) error {
	parts := strings.Split(line, " ")
	date := sanitize(parts[2] + " " + parts[3])
	dt, err := time.Parse(dateFormat, date)
	if err != nil {
		return err
	}
	lease.Tstp = dt
	return err
}

func parseCltt(line string, lease *Lease) error {
	parts := strings.Split(line, " ")
	date := sanitize(parts[2] + " " + parts[3])
	dt, err := time.Parse(dateFormat, date)
	if err != nil {
		return err
	}
	lease.Cltt = dt
	return err
}

func parseUID(line string, lease *Lease) error {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return errMalformed
	}
	lease.UID = sanitize(parts[1])
	return nil
}

func parseIdentifier(line string, lease *Lease) error {
	parts := strings.Split(line, " ")
	if len(parts) < 4 {
		return errMalformed
	}
	lease.Identifier = sanitize(parts[3])
	return nil
}

func parseEthernet(line string, lease *Lease) error {
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return errMalformed
	}
	lease.Hardware = sanitize(parts[2])
	return nil
}

func parseLease(line string, lease *Lease) error {
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return errMalformed
	}
	lease.IP = parts[1]
	return nil
}

func parseBinding(line string, lease *Lease) error {
	if lease.Binding != "" {
		return nil
	}
	parts := strings.Split(line, " ")
	lease.Binding = sanitize(parts[len(parts)-1])
	return nil
}
