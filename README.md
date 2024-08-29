# nmcli 

The nmcli package is a simple golang wrapper for the linux network-manager cli client (`nmcli`).

On Linux, `nmcli` is a command-line tool for controlling NetworkManager and reporting network status. It is used to create, display, edit, delete, activate, and deactivate network connections, as well as control and display network device status.

Inspired by the wonderful golang package [netlink](https://github.com/vishvananda/netlink) - the nmcli wrapper has been written to help users complete simple networking tasks in Linux that do not require the full power of the netlink package (or that may not yet be implemented).

## Dependencies 

## install 

* NetworkManager
  * `sudo apt -y install network-manager` (Ubuntu)

## enable

* 启动NetworkManager

```bash
systemctl restart NetworkManager
systemctl enable NetworkManager
```

* netplan管理网络的系统
- 如果是netplan管理网络的系统（如ubuntu22.04）需要在netplan中指定NetworkManager接管网络：

```
# 每个系统的文件名都不一样，我这里叫00-installer-config.yaml 
vim /etc/netplan/00-installer-config.yaml

# 在version下添加,注意开头对齐: 
renderer: NetworkManager
保存退出


# 示例
# This is the network config written by 'subiquity'
network:
  ethernets:
    ens33:
      dhcp4: true
  version: 2
  renderer: NetworkManager

```

* netplan应用
```bash
netplan apply
```




## Usage 

`go get github.com/BleetNie/nmcli`


### Check network manager (`nmcli`) version on system 
```golang
ver, err := nmcli.ValidateNmcliInstalled()
if err != nil {
  // nmcli is not installed
}
fmt.Println(ver)  // => nmcli version string
```

### Check for network-manager 
### Radios 
Get current state of radios:
```golang
radios, err := nmcli.Radios()
```

Change radio state:
```golang
// First get radios object
radios, err := nmcli.Radios()

// Set WIFI state
msg, err := radios.ChangeState(nmcli.WIFI, nmcli.OFF)

// Set state of all radios
msg, err := radios.ChangeState(nmcli.ALL, nmcli.ON)
```

### Connections 

Get all connections:
```golang
list_of_connections := nmcli.Connections()
```

Get connection(s) by connection name:
```golang
list_of_connections, err := nmcli.GetConnectionByName("{con-name}")
```

Add new connection:
```golang
```

#### Connection Object 
Delete connection:
```golang
c_list, _ := nmcli.GetConnectionByName("{con-name}")
// take first instance of connection with name {con-name}
// typically only one connection object is returned, but it is possible to have multiple connections with the same con-name
c := c_list[0]
// This will delete all connections that have name {con-name}
msg, err := c.Delete()
```

Modify connection:
```golang
c_list, _ := nmcli.GetConnectionByName("{con-name}")
// take first instance of connection with name {con-name}
// typically only one connection object is returned, but it is possible to have multiple connections with the same con-name
c := c_list[0]

// updates desired to connection details
// See docs for supported fields
c_updates := Connection{
    Name: "new-name",
    Device: "wlp58s0",
    Addr: &nmcli.AddressDetails{
      Ipv4_method:  "manual",
      Ipv4_address: "192.168.2.1",
      Ipv4_dns:     []string{"8.8.8.8", "1.1.1.1"},
    }
}
msg, err := c.Modify(c_updates)
```


### Devices 




## Compatibility Table 

| Object     | Command      | Status        |
| ---------- | ------------ | ------------- |
| general    |              | not supported |
| general    | status       | not supported |
| general    | hostname     | not supported |
| general    | permissions  | not supported |
| general    | logging      | not supported |
| networking |              | not supported |
| networking | on           | not supported |
| networking | off          | not supported |
| networking | connectivity | not supported |
| radio      |              | supported     |
| radio      | all          | supported     |
| radio      | wifi         | supported     |
| radio      | wwan         | supported     |
| connection |              | supported     |
| connection | show         | supported     |
| connection | up           | supported |
| connection | down         | supported |
| connection | add          | supported     |
| connection | modify       | supported     |
| connection | clone        | supported     |
| connection | edit         | not supported |
| connection | delete       | supported     |
| connection | reload       | not supported |
| connection | load         | not supported |
| connection | import       | not supported |
| connection | export       | not supported |
| device     |              | not supported |
| device     | status       | not supported |
| device     | show         | not supported |
| device     | set          | not supported |
| device     | connect      | not supported |
| device     | reapply      | not supported |
| device     | modify       | not supported |
| device     | disconnect   | not supported |
| device     | delete       | not supported |
| device     | monitor      | not supported |
| device     | wifi         | not supported |
| device     | wifi connect | not supported |
| device     | wifi rescan  | not supported |
| device     | wifi hotspot | not supported |
| device     | lldp         | not supported |
| agent      |              | not supported |
| agent      | secret       | not supported |
| agent      | polkit       | not supported |
| agent      | all          | not supported |
| monitor    |              | not supported |
