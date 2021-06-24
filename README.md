# Traefik SSL Certificate Exporter
[![Go](https://github.com/rafi0101/traefik-ssl-certificate-exporter/actions/workflows/go.yml/badge.svg)](https://github.com/rafi0101/traefik-ssl-certificate-exporter/actions/workflows/go.yml)
![](https://img.shields.io/docker/cloud/build/rafi0101/traefik-ssl-certificate-exporter.svg)
![Language](https://img.shields.io/badge/language-Golang-29BEB0.svg)
[![PRWelcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/rafi0101/traefik-ssl-certificate-exporter/pulls)
<div align="center">
  <sub>Built with ❤︎ by
  <a href="https://github.com/rafi0101">Raphael Ebner</a>
</div>
<br/>

Go tool to export traefik ssl certificates

The Docker container exports periodically the traefik ssl certificates and make them available in the ```certs/``` folder

Content
-----------
* [Important](#Important)
* [Usage](#Usage)
* [Variables](#Variables)
* [Input](#Input)
* [Output](#Output)
* [Developed by](#Developed-by)
* [License](#License)

Important:
----------
This tool only works with [Acmev2](https://letsencrypt.org/docs/client-options/)

Usage:
----------

### Command line:
```bash
go build && ./traefik-ssl-certificate-exporter --source traefik/acme.json --dest certs/
```

### Docker:

Docker image: [rafi0101/traefik-ssl-certificate-exporter](https://hub.docker.com/r/rafi0101/traefik-ssl-certificate-exporter)

See [docker-compose.yml](docker-compose.yml)

Variables:
----------
| CLI      | Docker        | Default Value           | Description  |
| -------- | ------------- | ------------------------| --------|
| --source | -             | ```traefik/acme.json``` | [Input](#Input) |
| --dest   | -             | ```certs/```            | [Output](#Output) |
| --owner  | CERT_OWNER_ID | ```0 (root)```          | owner for the extracted cert/keys |
| --group  | CERT_GROUP_ID | ```0 (root)```          | group for the extracted cert/keys |
| -        | CRON_TIME     | ```* * * * *```         | cron time for the container to extract certs |

Input:
----------
Your [acme.json](https://doc.traefik.io/traefik/https/acme/) from traefik tls configuration


Output:
----------
```
certs/
    example.com/
        cert.pem
        chain.pem
        fullchain.pem
        privkey.pem
    sub.example.com/
        cert.pem
        chain.pem
        fullchain.pem
        privkey.pem
```

Developed by
----------

* Raphael Ebner


License
----------

    MIT License

    Copyright (c) 2021 Raphael Ebner

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
