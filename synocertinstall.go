package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	auroraPackage "github.com/logrusorgru/aurora"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

var (
	version, tagName, branch, commitID, buildTime string
	aurora                                        auroraPackage.Aurora
)

type certObject struct {
	Desc     string `json:"desc"`
	Services []struct {
		DisplayName string `json:"display_name"`
		Service     string `json:"service"`
		Subscriber  string `json:"subscriber"`
	} `json:"services"`
}

type certInfoObject map[string]certObject

func init() {
	aurora = auroraPackage.NewAurora(isatty.IsTerminal(os.Stdout.Fd()))
	log.SetOutput(colorable.NewColorableStdout())
	log.SetFlags(0)
}

func main() {
	version = fmt.Sprintf("Version: %s, Branch: %s, Build: %s, Build time: %s",
		aurora.BrightCyan(tagName),
		aurora.BrightCyan(branch),
		aurora.BrightCyan(commitID),
		aurora.BrightCyan(buildTime))

	log.Println(aurora.Cyan("Synology NAS certification install tool"))
	log.Println(version)
	log.Println()

	flag.Usage = func() {
		log.Println("Usage:")
		flag.PrintDefaults()
		log.Println()
	}

	var listFlag, testFlag, updateFlag, installFlag bool
	var infoFile, certKey, newCertPath, newKeyPath, newChainPath, newFullChainPath string
	var listFormat string
	var originFile string = "/usr/syno/etc/certificate/_archive/INFO"

	flag.BoolVar(&listFlag, "list", false, "list applications")
	flag.BoolVar(&updateFlag, "update", false, "update system certifates")
	flag.BoolVar(&installFlag, "install", false, "install system certifates to AppPortal or ReverseProxy")
	flag.BoolVar(&testFlag, "test", false, "test mode, not really do it")
	flag.StringVar(&infoFile, "info-file", "", "certification information file path")
	flag.StringVar(&certKey, "cert-key", "", "certification key")
	flag.StringVar(&listFormat, "format", "a", "list format [a|s|p] for all, service, subscriber path")

	flag.StringVar(&newCertPath, "cert", "", "new certification file path")
	flag.StringVar(&newKeyPath, "key", "", "new key file path")
	flag.StringVar(&newChainPath, "ca", "", "new CA certification file path")
	flag.StringVar(&newFullChainPath, "chain", "", "new full chain certification file path")

	flag.CommandLine.SetOutput(os.Stdout)
	flag.Parse()

	if infoFile == "" {
		infoFile = originFile
	}
	log.Println("certifate infomation file:", aurora.BrightYellow(infoFile))
	log.Println()

	jsonFile, err := os.Open(infoFile)
	if err != nil {
		log.Fatalln(aurora.Red(err.Error()))
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var certInfo certInfoObject
	err = json.Unmarshal(byteValue, &certInfo)
	if err != nil {
		log.Fatalln(aurora.Red(err.Error()))
	}

	if listFlag {
		listCertObjs(&certInfo, listFormat)
		os.Exit(0)
	}

	if updateFlag {
		if certKey == "" || newCertPath == "" || newKeyPath == "" || newChainPath == "" || newFullChainPath == "" {
			log.Fatalln(aurora.Red("need cert-key, cert, key, ca and chain options"))
		}

		var findCert bool
		for k := range certInfo {
			if k == certKey {
				findCert = true
				certPath := fmt.Sprintf("/usr/syno/etc/certificate/_archive/%s", k)
				destCertPath := fmt.Sprintf("%s/cert.pem", certPath)
				destKeyPath := fmt.Sprintf("%s/privkey.pem", certPath)
				destChainPath := fmt.Sprintf("%s/chain.pem", certPath)
				destFullChainPath := fmt.Sprintf("%s/fullchain.pem", certPath)

				if testFlag {
					log.Println(aurora.BrightRed("test mode actived, not really update system."))
					log.Printf("copy %s %s\n", aurora.BrightGreen(newKeyPath), aurora.BrightYellow(destKeyPath))
					log.Printf("copy %s %s\n", aurora.BrightGreen(newCertPath), aurora.BrightYellow(destCertPath))
					log.Printf("copy %s %s\n", aurora.BrightGreen(newChainPath), aurora.BrightYellow(destChainPath))
					log.Printf("copy %s %s\n", aurora.BrightGreen(newFullChainPath), aurora.BrightYellow(destFullChainPath))
					os.Exit(0)
				}

				log.Println(aurora.BrightBlue("begin update system..."))

				err = copyFile(newKeyPath, destKeyPath)
				if err != nil {
					log.Fatalln(aurora.Red(err.Error()))
				}

				err = copyFile(newCertPath, destCertPath)
				if err != nil {
					log.Fatalln(aurora.Red(err.Error()))
				}

				err = copyFile(newChainPath, destChainPath)
				if err != nil {
					log.Fatalln(aurora.Red(err.Error()))
				}

				err = copyFile(newFullChainPath, destFullChainPath)
				if err != nil {
					log.Fatalln(aurora.Red(err.Error()))
				}

				log.Println(aurora.BrightBlue("update system successfully"))
				os.Exit(0)
			}
		}

		if !findCert {
			log.Fatalln(aurora.Red("not find certification key:"), aurora.BrightYellow(certKey))
		}

		log.Println(aurora.BrightBlue("update certification"), aurora.BrightYellow(certKey), aurora.BrightBlue("OK!"))
	}

	if installFlag {
		if certKey == "" {
			log.Fatalln(aurora.Red("need cert-key option"))
		}

		var findCert bool
		for k, co := range certInfo {
			if k == certKey {
				findCert = true

				log.Println(aurora.BrightBlue("Certificate key:"), aurora.BrightYellow(k),
					aurora.BrightBlue("Certificate description:"), aurora.BrightGreen(co.Desc))
				log.Println()

				certPath := fmt.Sprintf("/usr/syno/etc/certificate/_archive/%s", k)
				srcCertPath := fmt.Sprintf("%s/cert.pem", certPath)
				srcKeyPath := fmt.Sprintf("%s/privkey.pem", certPath)
				srcChainPath := fmt.Sprintf("%s/chain.pem", certPath)
				srcFullChainPath := fmt.Sprintf("%s/fullchain.pem", certPath)

				if testFlag {
					log.Println(aurora.BrightRed("test mode actived, not really install certifications.\n"))
				} else {
					log.Println(aurora.BrightBlue("begin install certifications...\n"))
				}

				for _, s := range co.Services {
					var logBuff []string
					logStr := fmt.Sprintf("%s %s", aurora.BrightBlue("install to service:"), aurora.BrightYellow(s.DisplayName))

					servicePath := fmt.Sprintf("/usr/syno/etc/certificate/%s/%s", s.Subscriber, s.Service)
					destCertPath := fmt.Sprintf("%s/cert.pem", servicePath)
					destKeyPath := fmt.Sprintf("%s/privkey.pem", servicePath)
					destChainPath := fmt.Sprintf("%s/chain.pem", servicePath)
					destFullChainPath := fmt.Sprintf("%s/fullchain.pem", servicePath)

					var installFail bool

					if testFlag {
						logBuff = append(logBuff,
							fmt.Sprintf("copy %s %s", aurora.BrightGreen(srcKeyPath), aurora.BrightYellow(destKeyPath)),
							fmt.Sprintf("copy %s %s", aurora.BrightGreen(srcCertPath), aurora.BrightYellow(destCertPath)),
							fmt.Sprintf("copy %s %s", aurora.BrightGreen(srcChainPath), aurora.BrightYellow(destChainPath)),
							fmt.Sprintf("copy %s %s", aurora.BrightGreen(srcFullChainPath), aurora.BrightYellow(destFullChainPath)))
					} else {
						err = copyFile(srcKeyPath, destKeyPath)
						if err != nil {
							logBuff = append(logBuff, fmt.Sprint(aurora.Red(err.Error())))
							installFail = true
						}

						err = copyFile(srcCertPath, destCertPath)
						if err != nil {
							logBuff = append(logBuff, fmt.Sprint(aurora.Red(err.Error())))
							installFail = true
						}

						err = copyFile(srcChainPath, destChainPath)
						if err != nil {
							logBuff = append(logBuff, fmt.Sprint(aurora.Red(err.Error())))
							installFail = true
						}

						err = copyFile(srcFullChainPath, destFullChainPath)
						if err != nil {
							logBuff = append(logBuff, fmt.Sprint(aurora.Red(err.Error())))
							installFail = true
						}
					}

					if !testFlag && !installFail {
						logStr = fmt.Sprintf("%s %s", logStr, aurora.BrightGreen("Ok"))
					} else if !testFlag && installFail {
						logStr = fmt.Sprintf("%s %s", logStr, aurora.BrightRed("Fail"))
					} else {
						logStr = fmt.Sprintf("%s %s", logStr, aurora.BrightMagenta("Test mode"))
					}

					log.Println(logStr)
					for k := range logBuff {
						log.Println(logBuff[k])
					}
				}

				if !testFlag {
					log.Println(aurora.BrightBlue("install certifications successfully"))
				}
			}
		}

		if !findCert {
			log.Fatalln(aurora.Red("not find certification key:"), aurora.BrightYellow(certKey))
		}

		log.Println(aurora.BrightBlue("install services certification from system certification"), aurora.BrightYellow(certKey), aurora.BrightBlue("OK!"))
	}
}

func listCertObjs(ci *certInfoObject, listFormat string) {
	for k, co := range *ci {
		if len(co.Services) > 0 {
			log.Println("Certifation Key:", aurora.BrightYellow(k))
			log.Println("Certifation Description:", aurora.BrightYellow(co.Desc))
			for _, s := range co.Services {
				switch listFormat {
				case "a", "all":
					log.Println("Service Name:", aurora.BrightGreen(s.DisplayName),
						"Subscriber:", aurora.BrightGreen(s.Subscriber),
						"Service Path:", aurora.BrightGreen(s.Service))
				case "s", "service":
					log.Println("Service Name:", aurora.BrightGreen(s.DisplayName))
				case "p", "path", "subscriber":
					log.Println("Subscriber:", aurora.BrightGreen(s.Subscriber), "Service Path:", aurora.BrightGreen(s.Service))
				default:
					log.Println(aurora.Red("Wrong format string, use:"))
					log.Println(aurora.Red("a, all for all informations"))
					log.Println(aurora.Red("s, service for service name"))
					log.Println(aurora.Red("p, path, subscriber for subscriber and service path"))
					os.Exit(1)
				}
			}
			log.Println()
		}
	}
}

func copyFile(srcFile, destFile string) error {
	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(destFile, input, 0400)
	return err
}
