# RedKing
RedKing is a simple tool for redirecting web requests.
It was created to help identify and exploit SSRF vulnerabilities similar to these:
 * [CVE-2021-21311](https://github.com/advisories/GHSA-x5r2-hj5c-8jx6)
 * [CVE-2021-21311 Writeup](https://github.com/vrana/adminer/files/5957311/Adminer.SSRF.pdf)
 * [Gitlab SSRF redirect vulnerability](https://gitlab.com/gitlab-org/gitlab-foss/-/issues/54649)

## How to Use it
Run RedKing with the `-h` flag to see available options and formats.

```
./RedKing -h
Usage of ./RedKing:
  -mode string
    	The mode RedKing should execute in. Select from:
    	single - redirect to a single URL
    	portscan - create a series of redirects at localhost/1,localhost/2,...
    		Each number will redirect to a different port on the target host.
    	The built in port scan ports are: 22,80,443,445,3389,8000,8080
    	 (default "single")
  -p int
    	The port used to host the redirect server (default 8080)
  -r int
    	Redirect status code - suggested 301, 302, or 307 (default 302)
  -url string
    	The URL used for redirects
  -v	Verbose
```

## Quickstart
```
./RedKing -url http://test.com


______         _   _   ___
| ___ \       | | | | / (_)
| |_/ /___  __| | | |/ / _ _ __   __ _
|    // _ \/ _' | |    \| | '_ \ / _' |
| |\ \  __/ (_| | | |\  \ | | | | (_| |
\_| \_\___|\__,_| \_| \_/_|_| |_|\__, |
                                  __/ |
                                 |___/


Mode: single
URL: http://test.com
Port: :8080

Starting server on localhost:8080
2021/05/15 21:17:04 Request from 127.0.0.1:12134
```

By default, RedKing opens a server on localhost:8080 and redirects all requests that hit it to the specified url.

### Single Mode
Single mode simply allows you to redirect all traffic to one specific host and path. See the example above for simple usage.

### Portscan Mode

**Note:** This mode is really designed to be used in conjunction with a tool like Burp's Intruder utility or something that will allow you to quickly
trigger requests and then grep through the output.

Portscan mode is designed to allow you to "scan" an internal IP for open ports.
By default it will create a redirect to allow you to test the following ports: 22, 80, 443, 445, 3389, 8000, 8080 \
It will create a series of redirects on your localhost, each of which corresponds to a specific port on the target server.
e.g. \
localhost:8080/0 -> Redirect To -> http://test.com:22 \
localhost:8080/1 -> Redirect To -> http://test.com:80 \
localhost:8080/2 -> Redirect To -> http://test.com:443

Depending on the specifics of the SSRF vulnerability you are exploiting, you may be able to glean information about running processes, or even access 
internal web pages and dump sensitive information. (e.g. the AWS metadata service)

Below is an example portscan mode and bash loop to demonstrate it.


**RedKing Output**
```
./RedKing -url http://test.com -mode portscan


______         _   _   ___
| ___ \       | | | | / (_)
| |_/ /___  __| | | |/ / _ _ __   __ _
|    // _ \/ _' | |    \| | '_ \ / _' |
| |\ \  __/ (_| | | |\  \ | | | | (_| |
\_| \_\___|\__,_| \_| \_/_|_| |_|\__, |
                                  __/ |
                                 |___/


Mode: portscan
URL: http://test.com
Port: :8080

Starting server on localhost:8080
Redirection from 127.0.0.1:13177 to http://test.com:22
Redirection from 127.0.0.1:13180 to http://test.com:80
Redirection from 127.0.0.1:13184 to http://test.com:443
Redirection from 127.0.0.1:13187 to http://test.com:445
Redirection from 127.0.0.1:13189 to http://test.com:3389
Redirection from 127.0.0.1:13192 to http://test.com:8000
Redirection from 127.0.0.1:13195 to http://test.com:8080
```

**TestScript Output**
```
for endpoint in {0..6}; do echo "Connecting to localhost:8080/$endpoint"; curl --retry 0 --connect-timeout 1  -L localhost:8080/$endpoint -Ss | head -5; echo ; done
Connecting to localhost:8080/0
curl: (28) Connection timed out after 1000 milliseconds

Connecting to localhost:8080/1
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=iso-8859-1">
<meta http-equiv="Content-Script-Type" content="text/javascript">

Connecting to localhost:8080/2
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>400 Bad Request - DOSarrest Internet Security</title>
curl: (23) Failed writing body (0 != 740)

Connecting to localhost:8080/3
curl: (28) Connection timed out after 1001 milliseconds

Connecting to localhost:8080/4
curl: (28) Connection timed out after 1000 milliseconds

Connecting to localhost:8080/5
curl: (28) Connection timed out after 1000 milliseconds

Connecting to localhost:8080/6
curl: (28) Connection timed out after 1000 milliseconds
```

Notice that for endpoints 0,3,4,5,and 6 the connection timed out. This indicates that there is likely _not_ a service running on those ports.

## Docker
Running from docker is easy, simply run the image and specify command line arguments.
```
git pull bpsizemore/redking
git run bpsizemore/redking -h
git run bpsizemore/redking -url test.com -mode simple 
```
**https://hub.docker.com/r/bpsizemore/redking**
