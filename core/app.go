package core

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"res-downloader/core/shared"
	"strconv"
	"syscall"
	"time"

	"github.com/vrischmann/userdir"
)

type App struct {
	ctx         context.Context
	assets      embed.FS
	AppName     string `json:"AppName"`
	Version     string `json:"Version"`
	Description string `json:"Description"`
	Copyright   string `json:"Copyright"`
	UserDir     string `json:"-"`
	LockFile    string `json:"-"`
	PublicCrt   []byte `json:"-"`
	PrivateKey  []byte `json:"-"`
	IsProxy     bool   `json:"IsProxy"`
	IsReset     bool   `json:"-"`
}

var (
	appOnce        *App
	globalConfig   *Config
	globalLogger   *Logger
	resourceOnce   *Resource
	systemOnce     *SystemSetup
	proxyOnce      *Proxy
	httpServerOnce *HttpServer
	ruleOnce       *RuleSet
)

func GetApp(assets embed.FS, wjs string) *App {
	if appOnce == nil {
		matches := regexp.MustCompile(`"productVersion":\s*"([\d.]+)"`).FindStringSubmatch(wjs)
		version := "1.0.1"
		if len(matches) > 0 {
			version = matches[1]
		}

		appOnce = &App{
			assets:      assets,
			AppName:     "res-downloader",
			Version:     version,
			Description: "res-downloader是一款集网络资源嗅探 + 高速下载功能于一体的软件，高颜值、高性能和多样化，提供个人用户下载自己上传到各大平台的网络资源功能！",
			Copyright:   "Copyright © 2023~" + strconv.Itoa(time.Now().Year()),
			IsReset:     false,
			PublicCrt: []byte(`-----BEGIN CERTIFICATE-----
MIIFVzCCAz+gAwIBAgIUFL4AWoydofrOGk9g7Ua4bt5vtLMwDQYJKoZIhvcNAQEL
BQAwQjELMAkGA1UEBhMCQ04xFzAVBgNVBAoMDnJlcy1kb3dubG9hZGVyMRowGAYD
VQQDDBFyZXMtZG93bmxvYWRlciBDQTAeFw0yNjA2MTIxMDExNDBaFw00NjA2MDcx
MDExNDBaMEIxCzAJBgNVBAYTAkNOMRcwFQYDVQQKDA5yZXMtZG93bmxvYWRlcjEa
MBgGA1UEAwwRcmVzLWRvd25sb2FkZXIgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4IC
DwAwggIKAoICAQCpRcJqbxZPvGNYD5BUFMsnyjvSV1Xs3zaNgoUoZCn99GE+452F
43YdSw9hjFTn18WaolVSZSjYL8ZcXNsQ8N2xLi5Y+D3lzy5RjFSucOEqmSoPblHI
xDHx49gWLn16ci/rndkUIiH8pUitYMWthH8JHKwOSxbudl0kiG+iFOH30veq5UKN
4dJ7z3gDGUSUIyMI99a620n7kLkYRyNRGc8vcDMrlUBmfjnA8t9eMYzDbVIbUttl
KAKuhxGBtPMXLk+OTmfgK/cxYxTMoRBhkD4VD8Fu6hKFudUUtwOVdVgqpvs4+9h6
a43WQsj1Tf+0b+e5lov2keEziK2AoHWjxvLXW//hbqtMOjSpJgwsyzLRJhXHhAW/
8y3QX1mG3O4+ypWnafSqH/yY1PKswIYGhsCBhrNf5nSZLHkcbtj8ieh/JOvyXZEw
lYks44iPgtckZYNJCCDU8A0vyIsEwBMfmP+yKqmGEE3Mw/rZcNTtfD5hkRgBkWgB
b87walCSB3Iz+3U8IBlkG98JVYKWNN/CYMiJHSRzuOpAOjRikJFzLOtz6OaLxX5b
dfOWVxWgCqZmnLoLdwKm54GNwRdmVSRclVAmHv7SJR36f5zYtfYGtWSGW2Ih8f6e
niz3M1e29nlx/LD+6qFK99B8M7F1J/f7N2vTeot71zyyPrrFRWy6actlTQIDAQAB
o0UwQzASBgNVHRMBAf8ECDAGAQH/AgEAMA4GA1UdDwEB/wQEAwIBBjAdBgNVHQ4E
FgQUwx0j3zH9Cz560cfoK+ACEMb3cnEwDQYJKoZIhvcNAQELBQADggIBAJ+lg3+p
2P2SKO38loP33sJU6UJCb0P42zgrI0APBn9SOB4l6CKK1vMSgVICJoRNZ4brl0SK
WIr7azIUwSCX5C1VFH8KgDJi9t5C2O+o+L6e8I2mnHNTFWV2MgDOVNM5d/D1yXLf
O3qfOC7VeiH4ongKpvYO8z3z9vgu9rXWAvQhmFBSbOC3mN2oBC8llnFfqHMJxFyI
Obi53jUAMACqVUlNN4Qrfea1KjFBwhwQPaVH8YePjVt8UGwUWQrqMbzojvkH2hw/
ka8ucp9TgPxHQDbpGoIYzHjmeoFfTVqTP5xC7jpsRG1Va6G5R2Qaa/rmKLvNPOUH
RikF2Er9lxKTm3OfelU/H55TV8KqOje3KtshsOt9HuxLs2Cs3qAlp9p9quoPeMyZ
DqK7sN/efxzxgL76ZG3/CJ6G6HQHdY2u6IDE1MmSScE6YxNzM3DkEQoqXVU8WaY4
UH3XM1YoCcbNJDmG9ZeEYcCe+EzYgDpYqj1Zu34ZIUBQJIA/vzI9/Rsi7L6n/ksv
5N1LfEBJf6D9jntJYv4r9J3DC0D2SuCHvQ1IM6CtwEdJ1JUP9iaejX5ZFtzxJ5xz
0uTN+rG5EUzh4kQWyOKdFI7NFmuW+wNJll32hkvKDF3s0sg7BJk0fCdzrMTpTUZ+
V3p454kzcjbly543qLJRHhoBKoPWB+2lpy2V
-----END CERTIFICATE-----
`),
			PrivateKey: []byte(`-----BEGIN PRIVATE KEY-----
MIIJQwIBADANBgkqhkiG9w0BAQEFAASCCS0wggkpAgEAAoICAQCpRcJqbxZPvGNY
D5BUFMsnyjvSV1Xs3zaNgoUoZCn99GE+452F43YdSw9hjFTn18WaolVSZSjYL8Zc
XNsQ8N2xLi5Y+D3lzy5RjFSucOEqmSoPblHIxDHx49gWLn16ci/rndkUIiH8pUit
YMWthH8JHKwOSxbudl0kiG+iFOH30veq5UKN4dJ7z3gDGUSUIyMI99a620n7kLkY
RyNRGc8vcDMrlUBmfjnA8t9eMYzDbVIbUttlKAKuhxGBtPMXLk+OTmfgK/cxYxTM
oRBhkD4VD8Fu6hKFudUUtwOVdVgqpvs4+9h6a43WQsj1Tf+0b+e5lov2keEziK2A
oHWjxvLXW//hbqtMOjSpJgwsyzLRJhXHhAW/8y3QX1mG3O4+ypWnafSqH/yY1PKs
wIYGhsCBhrNf5nSZLHkcbtj8ieh/JOvyXZEwlYks44iPgtckZYNJCCDU8A0vyIsE
wBMfmP+yKqmGEE3Mw/rZcNTtfD5hkRgBkWgBb87walCSB3Iz+3U8IBlkG98JVYKW
NN/CYMiJHSRzuOpAOjRikJFzLOtz6OaLxX5bdfOWVxWgCqZmnLoLdwKm54GNwRdm
VSRclVAmHv7SJR36f5zYtfYGtWSGW2Ih8f6eniz3M1e29nlx/LD+6qFK99B8M7F1
J/f7N2vTeot71zyyPrrFRWy6actlTQIDAQABAoICADt7WYGYLrFvREOeGHwLYIZH
cPUNnpKhr2RTsKIMbJsiZIe6aVCyqP9LMIj5TJ65umUH1U6iYJNzWlN7h7lwwp5v
2XaHKQ0X3DFozBsObHlRIDAS9qdDlj9nbrgAtzQvavpzWeRSyDmlBSBzyJMcY52a
lzVgmprKOhnL3dqJVwyEdGZ3sIb2C0ZZldUU7H3XyQhuOuUniNxfM0O+P8FQffcw
CLMwe9RoV7gfQHGznMDRqhAS8iieQi79JKA9K00Ch488qxDhsjgHFrh/gqeeDcrN
4g2tMDwQnLluRFKhTQ2T5uTTzvLZ15oqlr5oncnUpwuWxPEsuwOmbD2uZRhboVqM
puRqsjLffU6M6ZoOhaaLqXZKsMm4D0I8/4Jk9hWwz3np9Vb3tAxYRJbz859EeRN/
jwYBJ2JDWolxTP0E2x3qqITMfFCW2FHNSbBW0xWniddZVL3sYZhuvMGC8UCwne+r
R+fXfARbc42N0T3NS+BRZSE6XP8O/6MxzwSRo3qQEj9Sy5H0U0uIMT1yjTRTbCg8
HDzYSyw1I5J7RO154vNxvadI4KvjzQo6UQCAtrv+MEVmsRrO6bm0Zr4LUc8HTynO
OTgXtYZstQkeTZ8kDkVrzqN3zoA/RnQgf4M5cWPZkrbGE09zXGfTc/ZaD657BjPI
0MSmxpFzrlstsnnaQ+b7AoIBAQDbcj2iEMln4oHBk9F1vvAT3uKsm3rbBzY4BKog
ocm+VsWSxbIpPG93SQBYF10TjQK5+fQedxZnpSK0F9i852h8cEkM13zPp+4bpq7N
rZHiASM/uSpjf/MTAZelJeUCVPvLySoA7Ntvc/5s5JhNbjmHQadmu5gjdb8hzkEa
OB2ipjK1V47Bbocfm4tFKZqDky0IUrtbRSQ/3NNrkKVNH20XxmhwdCbPSJVfJdZP
ZLnvaSi29J/jUHLN0Wj7hkjvSPXwA3AJDkIofldqkrYXXsD/dc8e8xeImMOYlJ2y
4LCUqtKWcZCKuzltHvIRasX8REUCtEQ0sio7Ee5Dk2HUA8ZvAoIBAQDFd/qtPliV
ZCJD9Y2VRhz2LemdOgKQ0pxOrKi049R/gHnr9O/7YspstZIlQDXlcwv8Geii2jTu
RXfnzbQAZ2judF01WFFxa4/sEMrMKw0OnFgXvXat+Oeqqy19m4u0aY6qsYaaV/T/
thrS3PeP3z1iT8GUn/TeEXUqLod/8uQ5N/rwDtfwXyhqCjfjD0UZJaUtpy9DmukB
uy5V5KIRQlw2/M9+Xsv7PwbHPsmBuknOFs1/GlgY+ou3iSBcGoS7a0V9KtWuA/gC
3VJdIG9frxqaRH52KJn3tR+Lwl+oYVf5sgI79p5bGUK15oCka+MAnskC4wXn/hkd
1a8gmaZYeg4DAoIBAQCnZ6gC2UF74XxQ+v3gEA+/aNmNCXMYYZEH/O73w9ROQo2o
IO5/rJ3v5p2/ldsoTfsVesuy7fAGkyA9OK/bs8CupU3k4QJSu23WZDqXpuBSA8Ir
G4ttqi75gc54ascgF0qatFQ5rnbbuCYQVfalov954ijdIyC1dF6hYGGjqclZyeWH
F0tM3o9wGk2NLma0FvNUlSBeSQmVOlWCii6//chQSchkeQceO+XPVuL9X/7D13n6
z/SlCTr7LdQjhNZgzEzpkwXFsr3ffDodj9wfSeZ7OxkNKC78wmT8IeuHiQbL8uCB
ahL08kylpOTPSp/MiRwIKxKZDI1Q0KXtoSIRBew7AoIBAEBGcV36sTLPSSf9wXwZ
OkwXXbdDrpodM6uYH8HhdsWZuBXJwGN/IIyJ/WwKnoB1Fi2U8Vgw1pHIIuNc3X5U
Kp+TWNOIT9ovPMWGIbybNsDOuw3fKcYvAplW5pPAEZVD1qBQ8JNElga1670/F4XJ
EF3zv/r0peuymwSD6K8JDKhjRFbnPfqLvsflU5Og4Mjyq/VUOdozjix2FPr4VJhx
lTqAx7lGefbp6PbpxQVo8aUXXwREOvDGfLvc8p+IMjQfEptPWgYuXIfyKmENsbLf
eDVGnjMvvA4Y8o+3UROpu80W2AtUlu7oJgK3aKAcTeNhy+QYqU9+Mga4Cyz9Vlvu
PqkCggEBAL8rEaP9j2sEvzRt4x9T1W8LkgRVHrXo1g12ZUE7mqzDXeYoTxDJ2ot+
buno6sZDXTYUUAGj2yNfMKQAmMEZz4CE84vsY0qyZW3P/GRUQOA49rMi/odWM4xc
4itK91UzreVYWeWke+2mvCy3C0E4x55oNu+M5GzFpiV0Yg5OSMLcgprXNXeTZ89v
/rQEZgIzzF4cZtOr7VkuDx3of9ePabXMZQpionOkdiLtabHR0vAPEbOlRT6CD/tj
dVou10XCAWmkUDuW35E/17mYB/AGzZ5QZhgTFllW6mTjKwlc7TSidDkaFyPRAjvz
p6WW5V1LqYJHLtRNKpHwTNkrTdVS9vE=
-----END PRIVATE KEY-----
`),
		}
		appOnce.UserDir = filepath.Join(userdir.GetConfigHome(), appOnce.AppName)
		err := os.MkdirAll(appOnce.UserDir, 0750)
		if err != nil {
			fmt.Println("Mkdir UserDir err: ", err.Error())
		}
		appOnce.LockFile = filepath.Join(appOnce.UserDir, "install.lock")
		initLogger()
		initConfig()
		initProxy()
		initResource()
		initHttpServer()
		initSystem()
		initRule()
	}
	return appOnce
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	go httpServerOnce.run()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		a.UnsetSystemProxy()
		os.Exit(0)
	}()
}

func (a *App) OnExit() {
	a.UnsetSystemProxy()
	globalLogger.Close()
	if appOnce.IsReset {
		err := a.ResetApp()
		fmt.Println("err:", err)
	}
}

func (a *App) installCert() (string, error) {
	out, err := systemOnce.installCert()
	if err != nil {
		globalLogger.Esg(err, out)
		return out, err
	} else {
		if err := a.lock(); err != nil {
			globalLogger.Err(err)
		}
	}
	return out, nil
}

func (a *App) OpenSystemProxy() error {
	if a.IsProxy {
		return nil
	}
	err := systemOnce.setProxy()
	if err == nil {
		a.IsProxy = true
		return nil
	}
	return err
}

func (a *App) UnsetSystemProxy() error {
	if !a.IsProxy {
		return nil
	}
	err := systemOnce.unsetProxy()
	if err == nil {
		a.IsProxy = false
		return nil
	}
	return err
}

func (a *App) isInstall() bool {
	return shared.FileExist(a.LockFile)
}

func (a *App) lock() error {
	err := os.WriteFile(a.LockFile, []byte("success"), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) ResetApp() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return err
	}

	_ = os.Remove(filepath.Join(appOnce.UserDir, "install.lock"))
	_ = os.Remove(filepath.Join(appOnce.UserDir, "pass.cache"))
	_ = os.Remove(filepath.Join(appOnce.UserDir, "config.json"))
	_ = os.Remove(filepath.Join(appOnce.UserDir, "cert.crt"))

	cmd := exec.Command(exePath)
	cmd.Start()
	return nil
}
