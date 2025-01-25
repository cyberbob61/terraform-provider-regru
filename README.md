# Terraform provider for managing DNS records on reg.ru

This project contains a Terraform provider for managing DNS records using the reg.ru API. The provider allows you to create, read, and delete various types of DNS records, including A, AAAA, CNAME, MX, and TXT.

## Installation

To use this provider, you need to have Terraform version 0.12 or higher installed. You can download Terraform from the official website(https://www.terraform.io/downloads.html).

## Development and Build

1. **Build**:
Create a file with a name `Dockerfile` 
```bash
mkdir ~/temp && cd ~/temp
touch Dockerfile
```
Open a file in any text editor and add:
```dockerfile
FROM golang:1.20

RUN apt-get update && apt-get install -y make git && rm -rf /var/lib/apt/lists/*

RUN git clone -b dev https://github.com/cyberbob61/terraform-provider-regru.git /app && \
    cd /app && \
    make install-deps && \
    make build 
```

2. **Build**
```bash
docker build --no-cache -t terraform-provider-regru .
docker run -d --name terraform-provider-regru terraform-provider-regru && docker cp terraform-provider-regru:/app/out/terraform-provider-regru $(pwd) && docker rm terraform-provider-regru
```

3. **Copy the provider**
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/murtll/regru/0.3.0/linux_amd64/
mv terraform-provider-regru ~/.terraform.d/plugins/registry.terraform.io/murtll/regru/0.3.0/linux_amd64/
```
   
## Configuration

1. **terraformrc configuration**:

Edit or create a file if it doesn't exist `~/.terraformrc`:

```hcl
provider_installation {
  filesystem_mirror {
    path    = "/home/<<USER>>/.terraform.d/plugins"
    include = ["murtll/regru"]
  }
  direct {
    exclude = ["murtll/regru"]
  }
}
```


```sh
sed -i "s/<<USER>>/$USER/g" ~/.terraformrc
```

Create a configuration file for example `main.tf`:
```hcl
terraform {
  required_providers { 
    regru = { 
      version = "~>0.3.0"
      source  = "murtll/regru"
    } 
  }
}
```

2. **Testing provider**:

```sh
terraform init
```

3. **Terraform configuration examples**:
Add to `main.tf`

```hcl
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

4. **Terraform actions**:

```sh
terraform plan
terraform apply
```
## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for more details.