package nmcli

import (
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

// TODO: Update all code to use bash -c for exec.Command

type ConnDetails interface {
	// empty interface for package specific structs
	// TODO: Is this the best way to have common methods?
	// TODO: Remove this and just pass an empty interface to the generate commands method. It works on any struct.
	// Connection
	// Address
}

// TODO: Expand field set captured. Included State, etc in here.
type Connection struct {
	Name      string `cmd:"con-name"`
	Uuid      string
	Conn_type string `cmd:"type"`
	Device    string `cmd:"ifname"`
	Addr      *AddressDetail
}

type AddressDetail struct {
	Ipv4_method  string   `cmd:"ipv4.method"`
	Ipv4_address string   `cmd:"ipv4.address"`
	Ipv4_gateway string   `cmd:"ipv4.gateway"`
	Ipv4_dns     []string `cmd:"ipv4.dns"`
}

// Deletes the connection.
// Returns nmcli success message and error.
func (c Connection) Delete() (msg string, err error) {
	res, err := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("nmcli connection del %v", c.Name),
	).Output()
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// Clones an existing connnection and gives it a new name
// Equivalent to: nmcli con clone {name|uuid} {new_name}
func (c Connection) Clone(new_name string) (msg string, err error) {
	res, err := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("nmcli connection clone %v %v", c.Uuid, new_name),
	).Output()
	if err != nil {
		return string(res), err
	}
	return string(res), nil
}

// Enables the current connection
// Equivalent to: nmcli con up {name|uuid}
func (c Connection) Up() (msg string, err error) {
	res, err := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("nmcli connection up %v", c.Uuid),
	).Output()
	if err != nil {
		return string(res), err
	}
	return string(res), nil
}

// Disables the current connection
// Equivalen to: nmcli con down {name|uuid}
func (c Connection) Down() (msg string, err error) {
	res, err := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("nmcli connection down %v", c.Uuid),
	).Output()
	if err != nil {
		return string(res), err
	}
	return string(res), nil
}


// Modifies the connection with given parameters.
func (c *Connection) Modify(new_c Connection) (msg string, err error) {
	// conn_name := &c.Name
	cmds := new_c.construct_commands()
	// if address details provided then include
	if !reflect.DeepEqual(new_c.Addr, AddressDetail{}) {
		cmds = append(cmds, new_c.Addr.construct_commands()...)
		fmt.Println("Address present.")
	}
	cmds_str := strings.Join(cmds, " ")
	command:=fmt.Sprintf(`nmcli connection mod "%v" %v`, c.Name, cmds_str)
	fmt.Println(command)
	return
	res, err := exec.Command(
		"bash",
		"-c",
		command,
	).Output()
	if err != nil {
		return string(res), err
	}
	// update original connection with new details
	new_conn, err := GetConnectionByName(new_c.Name)
	if err != nil {
		// multiple connections by that name / not exists
		return "", err
	}
	*c = new_conn[0]

	return string(res), nil
}


// Returns all connections defined in nmcli
// Equivalent to: nmcli connection
func Connections() []Connection {
	res, err := exec.Command("bash", "-c", "nmcli connection").Output()
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	// process result
	results := make([]Connection, 0)
	input := strings.Split(strings.TrimSpace(string(res[:])), "\n")
	// fmt.Printf("%+v\n", input)
	// pop first row (headers)
	for _, line := range input[1:] {
		// fmt.Println(line)
		results = append(results, parseConnection(line))
	}
	for index, conn := range results {
		addr,err:=GetAddrDetail(conn.Name)
		if err==nil{
			results[index].Addr = &addr
		}
	}
	return results
}

// Finds connection by con-name, if it exists.
// Returns a list of Connections and an error.
// Equivalent to: nmcli connection show {name}
func GetConnectionByName(conn string) ([]Connection, error) {
	// get connections
	conns := Connections()
	// check if connection with name exists
	existingConns := []Connection{}
	for _, c := range conns {
		if c.Name == conn {
			addr,err:=GetAddrDetail(c.Name)
			if err==nil{
				c.Addr = &addr
			}
			existingConns = append(existingConns, c)
		}
	}
	// single conn = OK, multi conn = ERROR, no conn = ERROR
	switch len(existingConns) {
	case 0:
		return existingConns, errors.New("no connection found")
	case 1:
		return existingConns, nil
	default:
		return existingConns, errors.New("multiple connections found")
	}
}

// Creates a new connection
// Equivalent to: nmcli con add con-name {name} type {type} ifname {ifname}
// Returns nmcli message and error
func AddConnection(conn *Connection) (msg string, err error) {
	// Create new connection
	// TODO: Is it worth doing this in two parts? Or should execute as one command?
	res, err := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("nmcli connection add con-name %v type %v ifname %v", conn.Name, conn.Conn_type, conn.Device),
	).Output()
	if err != nil {
		return string(res), err
	}

	// Update connection with address details
	cmds := conn.Addr.construct_commands()
	// fmt.Println(append([]string{"connection", "mod", conn.Name}, cmds...))
	res, err = exec.Command("nmcli", append([]string{"connection", "mod", conn.Name}, cmds...)...).CombinedOutput()
	if err != nil {
		return string(res), err
	}
	return string(res), nil
}

func GetAddrDetail(connName string) (AddressDetail, error) {
	res, err := exec.Command(
		"bash",
		"-c",
		`nmcli connection show "`+connName+`"`,
	).Output()
	if err != nil {
		return AddressDetail{}, err
	}
	lines := strings.Split(strings.TrimSpace(string(res)), "\n")
	addr := AddressDetail{}
	for _, line := range lines {
		if strings.Contains(line, "ipv4.method:") {
			addr.Ipv4_method = strings.TrimSpace(strings.TrimPrefix(line, "ipv4.method:"))
		}
		if strings.Contains(line, "IP4.ADDRESS[1]:") {
			addr.Ipv4_address = strings.TrimSpace(strings.TrimPrefix(line, "IP4.ADDRESS[1]:"))
		}
		if strings.Contains(line, "IP4.GATEWAY:") {
			addr.Ipv4_gateway = strings.TrimSpace(strings.TrimPrefix(line, "IP4.GATEWAY:"))
		}
		if strings.Contains(line, "IP4.DNS[1]:") {
			dns := strings.TrimSpace(strings.TrimPrefix(line, "IP4.DNS[1]:"))
			addr.Ipv4_dns = append(addr.Ipv4_dns, dns)
		}
		if strings.Contains(line, "IP4.DNS[2]:") {
			dns := strings.TrimSpace(strings.TrimPrefix(line, "IP4.DNS[2]:"))
			addr.Ipv4_dns = append(addr.Ipv4_dns, dns)
		}
	}
	return addr, nil
}

func (addr AddressDetail) construct_commands() []string {
	return generate_commands(addr)
}

func (conn Connection) construct_commands() []string {
	return generate_commands(conn)
}

//*********************
// HELPERS
// ********************

func parseConnection(conn_line string) Connection {
	regex := regexp.MustCompile(`^([\S\s]+)\s{2}(\S+)\s{2}(\S+)\s+(\S+)\s*`)
	match := regex.FindStringSubmatch(conn_line)
	// fmt.Println(match)
	if len(match) != 5 {
		fmt.Println("Error. Incorrect number of fields returned. Aborting.")
	}

	return Connection{
		Name:      strings.TrimSpace(match[1]),
		Uuid:      strings.TrimSpace(match[2]),
		Conn_type: strings.TrimSpace(match[3]),
		Device:    strings.TrimSpace(match[4]),
	}
}
