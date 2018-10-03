#!/usr/bin/env bash

OS=$(uname | tr '[:upper:]' '[:lower:]')
if [ -z ${VERSION} ]
then
    VERSION="v0.1.0"
fi
PROVIDER_DOWNLOAD_URL="https://github.com/Sedicii/terraform-provider-ssh/releases/download/${VERSION}/terraform-provider-ssh_${OS}-amd64_${VERSION}"
TF_PLUGINS_PATH="${HOME}/.terraform.d/plugins/${OS}_amd64/"

mkdir -p "${TF_PLUGINS_PATH}"
PLUGIN_DEST="${TF_PLUGINS_PATH}/terraform-provider-ssh_${VERSION}"
curl -L ${PROVIDER_DOWNLOAD_URL} --output ${PLUGIN_DEST} && chmod +x ${PLUGIN_DEST}
