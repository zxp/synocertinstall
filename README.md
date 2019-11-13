# synocertinstall
A shell command to update Synology NAS SSL certifications and deploy for Synology AppPortal and ReverseProxy.

## How to use
```SHELL
Usage:
  -ca string
    	new CA certification file path
  -cert string
    	new certification file path
  -cert-key string
    	certification key
  -chain string
    	new full chain certification file path
  -format string
    	list format [a|s|p] for all, service, subscriber path (default "a")
  -info-file string
    	certification information file path
  -install
    	install system certifates to AppPortal or ReverseProxy
  -key string
    	new key file path
  -list
    	list applications
  -test
    	test mode, not really do it
  -update
    	update system certifates
```

### List all installed certifacates.
```SHELL
/volume1/docker/acme.sh# ./synocertinstall -list
...
certifate infomation file: /usr/syno/etc/certificate/_archive/INFO

Certifation Key: nvyfz6
Certifation Description: Test Let's Encrypt
Service Name: mail.test.com Subscriber: AppPortal Service Path: MailClient
Service Name: spreadsheet.test.com Subscriber: AppPortal Service Path: Spreadsheet
...
```
You will find the Certifation Key (**nvyfz6**), then you can update the new certificate.

### Update the specified certificate
```SHELL
/volume1/docker/acme.sh# ./synocertinstall_lin -update -cert-key nvyfz6 -cert test.com/test.com.cer \
                         -key test.com/test.com.key -ca test.com/ca.cer -chain test.com/fullchain.cer
```

### Install the new certificate to AppPortal and ReserveProxy
```SHELL
/volume1/docker/acme.sh# ./synocertinstall_lin -install -cert-key nvyfz6
```

### Other command options
```SHELL
  -info-file <certificate information file path>
                          specify the certificate information file path, normaly will be at 
                          `/usr/syno/etc/certificate/_archive/INFO`, but you can copy this file to 
                          anywhere and use it.

  -format [a|s|p]
  a, all                  list certificates service name, subscriber and service path.
  s, service              only list certificates service name.
  p, subscriber, path     only list certificates subscriber and service path.
  
  -test                   test mode, only display what will be done, and where the files will be 
                          copied to.
```