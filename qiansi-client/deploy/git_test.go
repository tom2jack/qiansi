package deploy

import (
	"qiansi/common/models"
	"testing"
)

func TestGit(t *testing.T) {
	config := &models.Deploy{
		Id:         4,
		Uid:        2,
		Title:      "纸喵应用",
		DeployType: 1,
		RemoteUrl:  "git@gitee.com:zhimiao/qiansi_test.git",
		LocalPath:  `D:\go\dist`,
		Branch:     "develop",
		DeployKeys: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEA4SkiNHu2rTiZsTUVFyYMt14Uoh6r2TZem6DrY+1Lh68xMk01DKBx
5e+jI+EKUJOKOittV03Nmsl4xpj3S7/V6pbo6x+iVQu57/FQ5ZNUKztP37o0Tzc+Lt9jZZ
tvqzxLT8iLmXX5/bQT4MyaLhiltY/TM5uyCg+P5WFV4TLdLzJxWrHzun5hksKgCYS7hmKo
7vIsmJgiKaXv+IhbbKCJTRzzjqUJgFZ3rGBu+DzYpxo8p2ecE7PopR7IZmjyrkfIV/NSTU
Vnpzokx94pfGGl0jucb6CWpysNnI0KsSIbG/r4VQnWLixyM5hbNd67ZQ+x2Rh8+x1YIG7V
c3ofWN7dxwAAA8DSCez30gns9wAAAAdzc2gtcnNhAAABAQDhKSI0e7atOJmxNRUXJgy3Xh
SiHqvZNl6boOtj7UuHrzEyTTUMoHHl76Mj4QpQk4o6K21XTc2ayXjGmPdLv9XqlujrH6JV
C7nv8VDlk1QrO0/fujRPNz4u32Nlm2+rPEtPyIuZdfn9tBPgzJouGKW1j9Mzm7IKD4/lYV
XhMt0vMnFasfO6fmGSwqAJhLuGYqju8iyYmCIppe/4iFtsoIlNHPOOpQmAVnesYG74PNin
GjynZ5wTs+ilHshmaPKuR8hX81JNRWenOiTH3il8YaXSO5xvoJanKw2cjQqxIhsb+vhVCd
YuLHIzmFs13rtlD7HZGHz7HVggbtVzeh9Y3t3HAAAAAwEAAQAAAQEAlZq4IHkm6reV3xmv
Fr9waZH4UbPhaSTn/a4RWUb9DX2JSavlGKuuoiH0ms1XBizSBk5+iyil+TfuqL5QaiNfpk
x5HGjbenidJeGIZ9HZdhQlwTi+sve4uHozV/rMWtFoFO3iW6f43+p73rzzoLc9u4KByWOl
C2xFpxpiboxWTJmPWiExUZqE6HkgDNU321vGdOLZUcGPlITNVS4YqDO1juwYtZxjJFqLJ1
RPS3QEYFNABGBpyHLuRK/kqevgse4ugcVj5MC+Xr/9+XMY35f5NLFo5vGX4/yWGCEnwg5e
Sju5tFLQZh1P4PYK+bXHBqNObys7KY6+VGPLdL05fM0cAQAAAIBDThtH4ZsY4pABU955/J
8RYSujBPt+dJmAAjgqm2KevOWyE9P0YPItXCZTJ32cZsdSUiUBPDLvR9d9qKYsXlNLnJMq
koVrEYcIdiqFK60BW62csz20BLng2Uq3mTgfqIoZP4911YUvBXrvMu9Mje+uRI50bpJOHV
NGE80oFnfGBwAAAIEA+MWA/w8V452W8lrI4i3X664bexNkoQy+7NEt05J/8WolU/YEKTQz
VuRe7pd2XKsqRqCKVWGHsYDjMIr4lqE5vXuJbFze6oJQ+5Y5hTNnWr2KiWkMY4C/sGYjOl
QwRbi1Fbobw4Dl+vgj0/8CKdwKChon4oGxo7+8I7xcykxWN8cAAACBAOez/+8oAhVIVx14
tdYt+XWIIZWh/cG78C2zMb5Io13c/wt7pn93YQHVyJDYr9zuq7LAAKYl2WBcfN2KsaFV+n
f/op1Zuu8I4LDgygY1BP0CjqxFWId4BEseZXDMj+lKW5i6zX/Kg1q2M+mowaZ0wVvYcd/P
6yxX4/2Fg/pRjioBAAAACDI0LlBDLTAxAQI=
-----END OPENSSH PRIVATE KEY-----`,
	}
	if err := Git(config); err != nil {
		t.Error(err)
	}
}
