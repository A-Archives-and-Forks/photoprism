#!/usr/bin/env bash

# Installs a Chromium-based browser on Linux:
#   - AMD64: Google Chrome stable (apt repo at dl.google.com).
#   - ARM64: Chromium (native apt package; no snap).
#       - Debian hosts: native apt chromium.
#       - Ubuntu hosts: chromium from PhotoPrism's internal mirror at
#         https://dl.photoprism.app/dist/chromium/. The mirror is a small
#         reprepro tree of the XtraDeb PPA's chromium / chromium-driver /
#         chromium-sandbox .debs (BSD-1-Clause); we mirror them so CI builds
#         survive Launchpad outages. Ubuntu's own chromium-browser is a
#         snap-shim transitional package and unusable in Docker.
#         Upstream PPA: https://launchpad.net/~xtradeb/+archive/ubuntu/apps
#
# Override the ARM64 Ubuntu chromium source via env var:
#   CHROMIUM_SOURCE=internal  (default) — https://dl.photoprism.app/dist/chromium/
#   CHROMIUM_SOURCE=xtradeb             — https://ppa.launchpadcontent.net/xtradeb/apps/ubuntu
# Use 'xtradeb' only when our mirror has fallen behind a chromium update
# (e.g. CVE patch landed at XtraDeb but cron hasn't refreshed yet).
#
# This script must run as root. Use one of these invocations:
#
#   # Pipe via stdin (recommended one-liner — works everywhere, incl. SSH):
#   curl -fsSL https://raw.githubusercontent.com/photoprism/photoprism/develop/scripts/dist/install-chrome.sh | sudo bash
#
#   # Or download first and run:
#   curl -fsSL https://raw.githubusercontent.com/photoprism/photoprism/develop/scripts/dist/install-chrome.sh -o /tmp/install-chrome.sh
#   sudo bash /tmp/install-chrome.sh
#
# Do NOT use process substitution under sudo (`sudo bash <(curl …)`): the
# substitution opens /dev/fd/63 in the *parent* (unprivileged) shell, which
# the elevated bash cannot read once sudo re-execs — the script aborts
# immediately with `bash: /dev/fd/63: No such file or directory`.

PATH="/usr/local/sbin:/usr/sbin:/sbin:/usr/local/bin:/usr/bin:/bin:/scripts:$PATH"

# Abort if not executed as root.
if [[ $(id -u) != "0" ]]; then
  echo "Usage: ${0##*/} must run as root. Try:" 1>&2
  echo "  curl -fsSL <url-to-this-script> | sudo bash" 1>&2
  exit 1
fi

# Determine target architecture.
if [[ $PHOTOPRISM_ARCH ]]; then
  SYSTEM_ARCH=$PHOTOPRISM_ARCH
else
  SYSTEM_ARCH=$(uname -m)
fi

DESTARCH=${BUILD_ARCH:-$SYSTEM_ARCH}

# shellcheck source=/dev/null
. /etc/os-release

# Adds PhotoPrism's internal chromium mirror as an apt source and installs
# chromium from it. Mirror is rebuilt from XtraDeb upstream weekly via
# `make -C services/downloads chromium` on web2 (services/downloads/scripts/).
install_chromium_from_internal_mirror() {
  local keyring=/etc/apt/keyrings/photoprism-apt.gpg
  local src=/etc/apt/sources.list.d/photoprism-chromium.sources

  install -m 0755 -d /etc/apt/keyrings
  # PhotoPrism APT signing key (master): 99F2 9643 6E34 A5DD 1782  A5C9 7FE2 4EBF 235B EBF9
  curl -fsSL "https://dl.photoprism.app/dist/chromium/photoprism-apt.gpg.asc" \
    | gpg --no-tty --batch --yes --dearmor -o "$keyring"

  cat > "$src" <<EOF
Types: deb
URIs: https://dl.photoprism.app/dist/chromium
Suites: ${VERSION_CODENAME}
Components: main
Signed-By: ${keyring}
EOF

  apt-get update
  apt-get -qq install -y --no-install-recommends \
    chromium chromium-driver chromium-sandbox
}

# Adds the XtraDeb PPA apt source and installs chromium from it. Fallback
# path used only when CHROMIUM_SOURCE=xtradeb.
install_chromium_from_xtradeb_ppa() {
  local keyring=/etc/apt/keyrings/xtradeb-apps.gpg
  local src=/etc/apt/sources.list.d/xtradeb-apps.sources

  install -m 0755 -d /etc/apt/keyrings
  # PPA signing key fingerprint: 5301FA4FD93244FBC6F6149982BB6851C64F6880
  curl -fsSL "https://keyserver.ubuntu.com/pks/lookup?op=get&search=0x5301FA4FD93244FBC6F6149982BB6851C64F6880" \
    | gpg --no-tty --batch --yes --dearmor -o "$keyring"

  cat > "$src" <<EOF
Types: deb
URIs: https://ppa.launchpadcontent.net/xtradeb/apps/ubuntu
Suites: ${VERSION_CODENAME}
Components: main
Signed-By: ${keyring}
EOF

  apt-get update
  apt-get -qq install -y --no-install-recommends \
    chromium chromium-driver chromium-sandbox
}

case $DESTARCH in
  amd64 | AMD64 | x86_64 | x86-64)
    echo "Installing Google Chrome (stable) on ${ID} for ${DESTARCH^^}..."
    set -e
    curl -fsSL https://dl-ssl.google.com/linux/linux_signing_key.pub | gpg --no-tty --batch --yes --dearmor -o /etc/apt/trusted.gpg.d/dl-ssl.google.com.gpg
    sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list'
    apt-get update
    apt-get -qq install google-chrome-stable
    ;;

  arm64 | ARM64 | aarch64)
    set -e
    case $ID in
      debian)
        echo "Installing Chromium on ${ID} for ${DESTARCH^^}..."
        apt-get update
        apt-get -qq install chromium chromium-driver chromium-sandbox
        ;;

      ubuntu)
        apt-get -qq install -y --no-install-recommends ca-certificates curl gnupg
        case "${CHROMIUM_SOURCE:-internal}" in
          internal)
            echo "Installing Chromium (via PhotoPrism internal mirror) on ${ID} ${VERSION_CODENAME:-} for ${DESTARCH^^}..."
            install_chromium_from_internal_mirror
            ;;
          xtradeb)
            echo "Installing Chromium (via XtraDeb PPA — fallback) on ${ID} ${VERSION_CODENAME:-} for ${DESTARCH^^}..."
            install_chromium_from_xtradeb_ppa
            ;;
          *)
            echo "Unknown CHROMIUM_SOURCE='${CHROMIUM_SOURCE}' (expected 'internal' or 'xtradeb')" 1>&2
            exit 1
            ;;
        esac
        ;;

      *)
        echo "Unsupported distribution \"${ID}\" for ARM64 Chromium install" 1>&2
        exit 1
        ;;
    esac
    ;;

  *)
    echo "Unsupported Machine Architecture: \"$DESTARCH\"" 1>&2
    exit 0
    ;;
esac

echo "Done."
