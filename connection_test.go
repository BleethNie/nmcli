package nmcli

import (
	"errors"
	"fmt"
	"testing"
)

func Test_CreateNewConnection(t *testing.T) {

	// initialise connection details
	newConn := Connection{
		Name:      "wcrd-go-nmcli-wrapper-test-connection",
		Conn_type: "dummy",
		Device:    "eth10",
		Addr: &AddressDetail{
			Ipv4_method:  "manual",
			Ipv4_address: "192.168.2.1",
			Ipv4_dns:     []string{"8.8.8.8", "1.1.1.1"},
		},
	}

	// create connection
	msg, err := AddConnection(&newConn)
	if err != nil {
		t.Errorf("Failed to add connection with message:\n%v\n", msg)
	}

	// Verify new connection exists
	_, err = GetConnectionByName("wcrd-go-nmcli-wrapper-test-connection")
	if err == errors.New("no connection found") {
		t.Errorf("New connection not found in nmcli connection list")
		t.Errorf("%v", err)
	}

}
func Test_ModifyConnection(t *testing.T) {
	c, _ := GetConnectionByName("有线连接 1")
	if len(c) == 0 {
		t.Skipf("Test connection has not been created. This may be due to a prior test failure. Skipping this test.")
	}
	conn:=Connection{Name: "eno1"}
	c[0].Modify(conn)
}


func Test_GetAddrDetail(t *testing.T) {
	addr,err:=GetAddrDetail("有线连接 1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(addr)
}



func Test_CloneConnection(t *testing.T) {
	// get connection
	c, _ := GetConnectionByName("eno1")
	if len(c) == 0 {
		t.Skipf("Test connection has not been created. This may be due to a prior test failure. Skipping this test.")
	}

	// clone
	msg, err := c[0].Clone("eno2")
	if err != nil {
		t.Errorf("failed to clone connection.\nmsg: %v", msg)
	}

	// verify creation
	c, _ = GetConnectionByName("wcrd-go-nmcli-wrapper-test-connection-clone")

	// clean-up
	_, err = c[0].Delete()
	if err != nil {
		fmt.Printf("Failed to delete cloned connection: %v\nPlease delete manually using nmcli.", c[0].Name)
	}
}

// TODO
func Test_ConnectionUp(t *testing.T) {
	// requires that the create new connection has run prior

	// get connection
	c, _ := GetConnectionByName("wcrd-go-nmcli-wrapper-test-connection")
	if len(c) == 0 {
		t.Skipf("Test connection has not been created. This may be due to a prior test failure. Skipping this test.")
	}

	// Verify state

}

// TODO
func Test_ConnectionDown(t *testing.T) {
	// requires that the create new connection has run prior

	// get connection
	c, _ := GetConnectionByName("wcrd-go-nmcli-wrapper-test-connection")
	if len(c) == 0 {
		t.Skipf("Test connection has not been created. This may be due to a prior test failure. Skipping this test.")
	}
}

func Test_DeleteConnection(t *testing.T) {
	// requires that the create new connection has run prior

	// get connection
	c, _ := GetConnectionByName("wcrd-go-nmcli-wrapper-test-connection")
	if len(c) == 0 {
		t.Skipf("Test connection has not been created. This may be due to a prior test failure. Skipping this test.")
	}

	// delete
	msg, err := c[0].Delete()
	if err != nil {
		t.Errorf("Failed to delete connection\n")
		t.Errorf("msg: %v", msg)
	}

}
