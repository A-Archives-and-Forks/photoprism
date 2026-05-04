#!/usr/bin/env bash

# Installs a Chromium-based browser on Linux:
#   - AMD64: Google Chrome stable (apt repo at dl.google.com).
#   - ARM64: Chromium (native apt package; no snap).
#       - Debian hosts: native apt chromium.
#       - Ubuntu hosts: Debian bookworm chromium, since Ubuntu's chromium-browser
#         is only a snap-shim transitional package and unusable in Docker.
#         Bookworm's userspace lib requirements (libjpeg62-turbo, libopenh264-7,
#         libminizip1, ...) are close enough to Ubuntu LTS that apt can resolve
#         them by pulling those few libs from the Debian repo; the chromium
#         binary itself runs fine on Ubuntu glibc from jammy (22.04, 2.35) up.
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

# Adds the Debian bookworm apt source and installs chromium from it.
install_chromium_from_debian_bookworm() {
  local keyring=/etc/apt/keyrings/debian-archive-bookworm.gpg
  local src=/etc/apt/sources.list.d/debian-bookworm-chromium.sources
  local pin=/etc/apt/preferences.d/debian-bookworm-chromium.pref

  install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://ftp-master.debian.org/keys/archive-key-12.asc \
    | gpg --no-tty --batch --yes --dearmor -o "$keyring"

  cat > "$src" <<EOF
Types: deb
URIs: http://deb.debian.org/debian
Suites: bookworm
Components: main
Signed-By: ${keyring}
EOF

  # Pin Debian bookworm to a low priority so apt only pulls Bookworm packages
  # that have no Ubuntu equivalent. Without this, apt's solver "upgrades"
  # unrelated Ubuntu packages such as libjpeg-dev to Bookworm's higher epoch
  # version (1:2.1.5-2 vs Ubuntu's 8c-*ubuntu*), which then pulls
  # libjpeg62-turbo-dev and collides with Ubuntu's already-installed
  # libjpeg-turbo8-dev (both ship /usr/include/<triplet>/jconfig.h).
  # The chromium-* packages are elevated to priority 990 so they install
  # from Bookworm; their runtime dep on libjpeg62-turbo (no Ubuntu equivalent)
  # is resolved as the only available candidate, which is harmless.
  #
  # Older Ubuntu LTS that can't satisfy current Bookworm chromium minreq
  # (Jammy 22.04 ships libharfbuzz0b 2.7.4, libopenjp2-7 2.4.0, libzstd1
  # 1.4.8 — all below the >=6.0.0/>=2.5.0/>=1.5.2 chromium asks for) need
  # an extended whitelist so apt picks the newer Bookworm versions of
  # those specific libs. Newer Ubuntu releases (Noble 24.04+) already
  # ship versions that satisfy chromium, so we keep the narrow pin there
  # to avoid downgrading native libs (e.g. Resolute's libharfbuzz0b
  # 12.3.2 → Bookworm 6.0.0). libc6 / libstdc++6 are intentionally never
  # whitelisted so the Ubuntu base stays intact; libjpeg-dev /
  # libjpeg62-turbo-dev are intentionally never whitelisted so the
  # Resolute libjpeg-dev/libjpeg-turbo8-dev collision stays prevented.
  local extra_pin=""
  if [[ $VERSION_CODENAME == "jammy" ]]; then
    extra_pin=$(cat <<EOF


Package: libharfbuzz-subset0 libharfbuzz0b libharfbuzz-icu0 libharfbuzz-gobject0
Pin: release o=Debian,n=bookworm
Pin-Priority: 990

Package: libopenjp2-7 libzstd1 libdav1d6 libjpeg62-turbo libminizip1 libopenh264-7
Pin: release o=Debian,n=bookworm
Pin-Priority: 990
EOF
)
  fi

  cat > "$pin" <<EOF
Package: *
Pin: release o=Debian,n=bookworm
Pin-Priority: 100

Package: chromium chromium-common chromium-driver chromium-sandbox
Pin: release o=Debian,n=bookworm
Pin-Priority: 990${extra_pin}
EOF

  apt-get update
  apt-get -qq install -y --no-install-recommends \
    chromium chromium-common chromium-driver chromium-sandbox
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
        echo "Installing Chromium (via Debian bookworm) on ${ID} ${VERSION_CODENAME:-} for ${DESTARCH^^}..."
        apt-get -qq install -y --no-install-recommends ca-certificates curl gnupg
        install_chromium_from_debian_bookworm
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
