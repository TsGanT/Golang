# Golang Ping
This program can be used as CLI ping application
You must use this application based on Golang
Install Golang: https://golang.org/dl/

How to use it:
# Because this program need to use ICMP package, you need to run this with root
sudo go run ping.go
With no argument, you can get help of this program

Here are some usage of this prigram:

1. Ping IPV4 address:
sudo go run ping.go --Host 8.8.8.8

2. Ping IPV4 address with setting TTL by yourself:
sudo go run ping.go --Host 8.8.8.8 --TTL 50

3. Ping IPV6 address:
In Mac OS machine, you need to run it in this way:
sudo go run ping.go --Host fe80::dd8b:3601:bb86:293a%en0 --IPV6 On

You need to pay attention to %en0. This is used to set the hardware used to ping IPV6 address.

Author: Shi Tang
Email: tangshi6666@gmail.com
