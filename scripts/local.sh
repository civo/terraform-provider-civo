#! /bin/bash

RESET=`tput sgr0`
RED=`tput setaf 1`
GREEN=`tput setaf 2`
YELLOW=`tput setaf 3`

get_os_and_arch() {
  # Check if running on Linux
  if [[ "$(uname)" == "Linux" ]]; then
    # Check if running on ARM
    if [[ "$(uname -m)" == "aarch64" ]]; then
      echo "linux_arm64"
    else
      echo "linux_amd64"
    fi
  # Check if running on macOS
  elif [[ "$(uname)" == "Darwin" ]]; then
    # Check if running on ARM
    if [[ "$(uname -m)" == "arm64" ]]; then
      echo "darwin_arm64"
    else
      echo "darwin_amd64"
    fi
  else
    echo "Unknown operating system"
  fi
}

# verify civo cli is installed
if [ -z `which civo` ]; then
   printf "${RED}Civo CLI is not installed! Please install\n\n${RESET}"
   exit
fi

civo_api_key=`civo apikey show | tail -n +4 | head -1 | awk '{print $4}'`

# verify terraform is installed
if [ -z `which terraform` ]; then
   printf "${RED}terraform is not installed! Please install\n\n${RESET}"
   exit
fi

# verify golang is installed
if [ -z `which go` ]; then
   printf "${RED}go is not installed! Please install\n\n${RESET}"
   exit
fi

read -p "Provide the path of the folder containing the ${GREEN}.tf${RESET} files you want to apply: " manifests_folder

echo "Creating the plugin folder to allow installing the ${GREEN}civo provider locally...${RESET}"
os_arch=$(get_os_and_arch)
plugin_folder=$manifests_folder/.terraform.d/plugins/registry.terraform.io/civo/civo/99.0.0/$os_arch
mkdir -p $plugin_folder

read -p "Provide the ${GREEN}region${RESET} you want to create resources in ${YELLOW}[LON1/FRA1/NYC1/PHX1]${RESET}: " region

echo "Overriding ${GREEN}~/.terraformrc${RESET} file with the $manifests_folder/.terraform.d/plugins"
cat > ~/.terraformrc << EOF
provider_installation {
  filesystem_mirror {
    path = "$manifests_folder/.terraform.d/plugins"
    include = ["registry.terraform.io/civo/civo"]
  }
  direct {
    exclude = ["registry.terraform.io/civo/civo"]
  }
}
EOF

# verify that ~/provider.tf is in place else create it
if [ -z `ls $manifests_folder/provider.tf 2>/dev/null` ]; then
cat > $manifests_folder/provider.tf << EOF
terraform {
  required_providers {
    civo = {
      source  = "civo/civo"
      version = "99.0.0"
    }
  }
}
provider "civo" {
  token = "$civo_api_key"
  region = "$region"
}
EOF
echo "${GREEN}${manifests_folder}/provider.tf${RESET} file created for civo provider using your civo apikey"
else 
   echo "${GREEN}${manifests_folder}/provider.tf${RESET} file found"
fi

echo "Building the terraform civo provider binary and moving it into the manifests plugin folder"
go build -o $plugin_folder/terraform-provider-civo_v99.0.0 main.go

cd $manifests_folder
echo "${GREEN}Init${RESET} civo provider..."
terraform init
terraform plan
printf "To apply desired resources, you can now use \n ${YELLOW}cd manifests_folder \n terraform apply${RESET}\n"


