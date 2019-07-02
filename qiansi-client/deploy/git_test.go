package deploy

import (
	"qiansi/models"
	"testing"
)

func TestGit(t *testing.T) {
	config := &models.Deploy{
		Id:         4,
		Uid:        2,
		Title:      "纸喵应用",
		DeployType: 1,
		RemoteUrl:  "ssh://git@gitee.com/273000727/go-test.git",
		LocalPath:  `D:\go\tools-client\dist`,
		Branch:     "master",
		DeployKeys: `-----BEGIN OPENSSH PRIVATE KEY-----
		b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
		NhAAAAAwEAAQAAAQEAyslyjGtcG9FlAUwfo56w28tSdgyZDiWfBQiXL4kTj92dyS1puQRl
		KR9cJtpMbD+TF19CUqGKBzo8GWhFF80CW2v5JzyMC/YdJifOnDYKbIhOULO/q5ZCoMtTO9
		R/FP/oikddWkl4s3W9EjuBBNpua9L3GHLVDZi2Nqfra1aU6tr/6eYovwEUptSEhMmGiJfv
		ZGYbH3cO3YHEFg0lSE8LBLPe32KQommAqvzbi5xStrLaoReMVI7kHikxwWz5pscb+GRNRW
		Xn0zV3Ib+H8L0pt4OT0/yPO0rUSq8V1cEGvqUV3pnJ7rDGlzyYcJkakHuwnpfEjFIdgMT0
		nGKbPRAVcQAAA9h9Ok1afTpNWgAAAAdzc2gtcnNhAAABAQDKyXKMa1wb0WUBTB+jnrDby1
		J2DJkOJZ8FCJcviROP3Z3JLWm5BGUpH1wm2kxsP5MXX0JSoYoHOjwZaEUXzQJba/knPIwL
		9h0mJ86cNgpsiE5Qs7+rlkKgy1M71H8U/+iKR11aSXizdb0SO4EE2m5r0vcYctUNmLY2p+
		trVpTq2v/p5ii/ARSm1ISEyYaIl+9kZhsfdw7dgcQWDSVITwsEs97fYpCiaYCq/NuLnFK2
		stqhF4xUjuQeKTHBbPmmxxv4ZE1FZefTNXchv4fwvSm3g5PT/I87StRKrxXVwQa+pRXemc
		nusMaXPJhwmRqQe7Cel8SMUh2AxPScYps9EBVxAAAAAwEAAQAAAQA2/zdv2dYbPUj1dx3F
		lE5G7fepSHViHtXn2ZKXM8f4ZpRacVSQ9x4wbu7hIqdDXGKaHh2wp1r15tdR1LOYZuNSxA
		/IkmUxAUiahoVEXGurT7RdssIy2Qes8DfcrB7jJRx+FCi/SdnQYggrH7Q4Cr2TxJ17Jfme
		PGJ+pD/21n6AyvSygH0x1D/4jzmJvLX3pFRmHrScxq08N12h21/6pksB1VV09jMEEz+tRI
		0cHtlGwYuNj7RsrD55xI2l2RyG5X0KcjPCO1ALsW08HIPahfaYUaEYs+n+xB+Z9clhZC/j
		mgkF9v5sBiq76u1GJUGNQg69gMdQqD5f/tJkWIDtXMbZAAAAgQCPYKT+gi9cNT+6wBhsus
		dz7g/Bu/W2EHiDPaXwItUUEotY56b+4n6LI/v1w5ygTYY9PqUu8VR866Bz6yh3P94SXnv0
		tGbXwkR742H18fY0T0qImjpoQ+9CuDO0Hnx2ecrMzitEnURPnin1spRi1TmDcki9R5gXZu
		BKsT7EupoRvQAAAIEA5Gz+qmwMeClIlLWktcjEELQLLQC0ozR342r8SUBOgpr53Dm6Oxpy
		q8uCvrXmdAX3zNM3vHUMik/dGAmjO87616GkubptwffPBRD4RqxYRjlzPwbUUd0W2qJpgW
		RrPupDXj1PAKfDXhlXJGfyvj+/ZlqAWtQ3UIfVdOGPdoAEmGMAAACBAONEI52/N4fh49Kk
		vNRfsoMFdhzShReXV8zjveZ1bsAByglpi0xAZ+2NZ7M4zdhjM5pkFH+kgS6Ys/vXFuED3v
		3oqYlll6S0Z39+78HBne/9f7RClc7i/161Pk0ozrgTYbyzLjR6kXoXUpUd1/C3lDiHGfcF
		dJpxahA7uCJDV+EbAAAAHUFkbWluaXN0cmF0b3JAU0MtMjAxOTAzMDIxNjE3AQIDBAU=
		-----END OPENSSH PRIVATE KEY-----`,
	}
	if err := Git(config); err != nil {
		t.Error(err)
	}
}
