# Terraform provider for managing DNS records on reg.ru

This project contains a Terraform provider for managing DNS records using the reg.ru API. The provider allows you to create, read, and delete various types of DNS records, including A, AAAA, CNAME, MX, and TXT.

## Installation

To use this provider, you need to have Terraform version 0.12 or higher installed. You can download Terraform from the official website(https://www.terraform.io/downloads.html).

## Development and Build

1. **Create a Dockerfile**:
```dockerfile
FROM golang:1.20

RUN apt-get update && apt-get install -y make git && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/cyberbob61/terraform-provider-regru.git /app && \
    cd /app && \
    make install-deps && \
    make build

```

2. **Check**:

    Verify that the built provider is located in the correct directory:

    ```sh
    ls ~/.terraform.d/plugins/registry.terraform.io/cyberbob61/regru/0.2.1/linux_amd64/
    ```
   
## Configuration

1. **Configuration**:

    Create a configuration file for example `main.tf`:

    ```hcl
    terraform {
      required_providers {
        regru = {
          version = "~>0.2.0"
          source  = "letenkov/regru"
        }
      }
    }

    provider "regru" {
      api_username = "<username>"
      api_password = "<password>"
      cert_file    = "cert.crt"
      key_file     = "key.crt"
    }

    resource "regru_dns_record" "example_com" {
      zone   = "example.com"
      name   = "@"
      type   = "A"
      record = "1.1.1.1"
    }

    resource "regru_dns_record" "example_com_ipv6" {
      zone   = "example.com"
      name   = "@"
      type   = "AAAA"
      record = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
    }

    resource "regru_dns_record" "example_com_mx" {
      zone   = "example.com"
      name   = "@"
      type   = "MX"
      record = "10 mail.example.com"
    }

    resource "regru_dns_record" "example_com_txt" {
      zone   = "example.com"
      name   = "@"
      type   = "TXT"
      record = "v=spf1 include:example.com ~all"
    }
    ```

3. **Terraform actions**:

    ```sh
    terraform init
    terraform plan
    terraform apply
    ```

## Лицензия

Этот проект лицензируется на условиях лицензии Apache 2.0. Подробнее см. в файле [LICENSE](LICENSE).
