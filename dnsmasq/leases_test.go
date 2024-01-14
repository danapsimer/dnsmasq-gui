package dnsmasq

import (
	"github.com/stretchr/testify/assert"
	"net"
	"strings"
	"testing"
	"time"
)

var leasesText = `1704926889 00:e0:4c:68:01:b2 192.168.215.5 thd 01:00:e0:4c:68:01:b2
1704890089 30:d0:42:eb:42:f6 192.168.215.2 bart 01:30:d0:42:eb:42:f6
1704871294 5a:b3:8f:64:a9:44 192.168.215.72 * 01:5a:b3:8f:64:a9:44
1704882639 64:12:69:57:00:15 192.168.215.132 DIRECTV-HR54-69570013 01:64:12:69:57:00:15
1704875340 c4:dd:57:28:06:33 192.168.215.118 ESP_280633 *
1704876600 76:99:a5:2d:0f:7e 192.168.215.117 * 01:76:99:a5:2d:0f:7e
1704879738 3a:dd:15:c4:ab:f3 192.168.215.69 * 01:3a:dd:15:c4:ab:f3
1704878221 78:ca:39:be:1a:ca 192.168.215.130 NancyAPimersMBP 01:78:ca:39:be:1a:ca
1704883206 7a:eb:bb:ec:bf:db 192.168.215.147 * 01:7a:eb:bb:ec:bf:db
1704890049 14:87:6a:f0:93:23 192.168.215.3 moe 01:14:87:6a:f0:93:23
1704875873 22:0f:75:1d:5d:b6 192.168.215.85 * 01:22:0f:75:1d:5d:b6
1704865545 3c:9b:d6:81:5b:9c 192.168.215.230 viziocastdisplay *
1704916541 2c:f0:5d:d2:f5:e1 192.168.215.6 DESKTOP-GF9E2EF 01:2c:f0:5d:d2:f5:e1
1704884161 a8:bb:50:86:45:5e 192.168.215.225 wiz_86455e *
1704875246 d0:c2:4e:27:27:e8 192.168.215.95 Samsung 01:d0:c2:4e:27:27:e8
1704874910 e8:31:cd:b7:0b:e8 192.168.215.227 Litter-Robot4 01:e8:31:cd:b7:0b:e8
1704877973 f8:0f:f9:8f:bd:9e 192.168.215.182 Google-Nest-Mini *
1704880741 c4:dd:57:29:d1:6e 192.168.215.75 ESP_29D16E *
1704880750 70:03:9f:90:bd:fa 192.168.215.97 ESP_90BDFA *
1704865388 70:48:f7:d5:a9:91 192.168.215.186 * *
1704882422 00:16:cb:a9:50:11 192.168.215.64 * 01:00:16:cb:a9:50:11
1704881756 3c:e4:41:07:89:3c 192.168.215.168 * *
1704864744 78:6c:84:74:d0:24 192.168.215.170 * *
1704866846 ee:fa:cb:25:4e:0f 192.168.215.223 * 01:ee:fa:cb:25:4e:0f
1704867571 44:42:01:30:f5:30 192.168.215.169 * *
1704869935 94:91:7f:78:87:f2 192.168.215.83 * *
1704853206 ec:b5:fa:89:35:72 192.168.215.172 ecb5fa893572 *
1704860954 a8:bb:50:79:ea:a8 192.168.215.87 wiz_79eaa8 *
1704861351 a8:bb:50:ee:64:ae 192.168.215.152 wiz_ee64ae *
1704861699 6c:29:90:59:97:56 192.168.215.222 wiz_599756 *
1704861738 a8:bb:50:ed:45:38 192.168.215.151 wiz_ed4538 *
1704862240 6c:29:90:59:8d:38 192.168.215.219 wiz_598d38 *
1704864074 44:4f:8e:a1:b5:36 192.168.215.218 wiz_a1b536 01:44:4f:8e:a1:b5:36
1704865191 64:16:66:1f:3e:de 192.168.215.78 09AA01AC051803UM 01:64:16:66:1f:3e:de
1704865219 64:16:66:1e:c8:bc 192.168.215.79 09AA01AC041807KL 01:64:16:66:1e:c8:bc
1704865233 1c:f2:9a:1e:a5:05 192.168.215.76 * *
1704865266 7c:78:b2:1b:d1:71 192.168.215.91 * *
1704865268 64:16:66:6c:2a:12 192.168.215.77 Nest-Hello-2a12 01:64:16:66:6c:2a:12
1704865269 e0:e2:e6:af:44:b0 192.168.215.90 * *
1704865269 1c:9d:c2:51:ec:08 192.168.215.89 espressif *
1704873964 64:db:a0:0d:83:e0 192.168.215.88 500-64dba00d83e0 *
1704865559 ea:ce:9c:8b:4c:bd 192.168.215.206 * 01:ea:ce:9c:8b:4c:bd
1704866378 6c:29:90:59:a8:96 192.168.215.220 wiz_59a896 01:6c:29:90:59:a8:96
1704870327 8a:b8:45:b4:e1:24 192.168.215.139 Chelsea-s-S21 01:8a:b8:45:b4:e1:24
1704870914 00:1d:4b:03:4e:51 192.168.215.191 * *
1704874597 14:d4:fe:db:69:cb 192.168.215.143 * 01:14:d4:fe:db:69:cb
1704876496 46:3a:5f:fe:ac:2b 192.168.215.200 * 01:46:3a:5f:fe:ac:2b`

func TestReadLeases(t *testing.T) {
	leases, err := ReadLeases(strings.NewReader(leasesText))
	if assert.NoError(t, err) {
		assert.Equal(t, 47, len(leases))
		checkCount := 0
		for _, lease := range leases {
			t.Logf("%s %s %s %s %s", lease.ExpireTime, lease.MacAddress, lease.IPAddress, lease.Name, lease.ClientId)
			switch lease.MacAddress.String() {
			case "64:db:a0:0d:83:e0":
				assert.Equal(t, "500-64dba00d83e0", lease.Name)
				assert.Equal(t, time.Date(2024, 1, 10, 3, 6, 4, 0, time.Local), lease.ExpireTime)
				assert.Equal(t, &net.IPAddr{net.IP([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 192, 168, 215, 88}), "ip"}, lease.IPAddress)
				assert.Equal(t, "", lease.ClientId)
				checkCount += 1
			case "46:3a:5f:fe:ac:2b":
				assert.Equal(t, "", lease.Name)
				assert.Equal(t, time.Date(2024, 1, 10, 3, 48, 16, 0, time.Local), lease.ExpireTime)
				assert.Equal(t, &net.IPAddr{net.IP([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 192, 168, 215, 200}), "ip"}, lease.IPAddress)
				assert.Equal(t, "01:46:3a:5f:fe:ac:2b", lease.ClientId)
				checkCount += 1
			case "7c:78:b2:1b:d1:71":
				assert.Equal(t, "", lease.Name)
				assert.Equal(t, time.Date(2024, 1, 10, 0, 41, 6, 0, time.Local), lease.ExpireTime)
				assert.Equal(t, &net.IPAddr{net.IP([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 192, 168, 215, 91}), "ip"}, lease.IPAddress)
				assert.Equal(t, "", lease.ClientId)
				checkCount += 1
			}
		}
		assert.Equal(t, 3, checkCount)
	}
}
