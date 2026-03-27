#!/usr/bin/env bash
set -euo pipefail

REPO="interdev7/fcm-cli"
BINARY_NAME="fcm"

log() {
  printf '%s\n' "$1"
}

detect_platform() {
  OS="$(uname -s)"
  ARCH="$(uname -m)"

  case "$OS" in
    Linux) OS="linux" ;;
    Darwin) OS="darwin" ;;
    *)
      log "Unsupported OS: $OS"
      exit 1
      ;;
  esac

  case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
      log "Unsupported architecture: $ARCH"
      exit 1
      ;;
  esac

  PLATFORM="${OS}-${ARCH}"
}

get_latest_tag() {
  LATEST_TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')"

  if [ -z "${LATEST_TAG:-}" ]; then
    log "Failed to fetch latest release"
    exit 1
  fi
}

download_binary() {
  FILE_NAME="${BINARY_NAME}-${PLATFORM}"
  URL="https://github.com/${REPO}/releases/download/${LATEST_TAG}/${FILE_NAME}"

  TMP_DIR="$(mktemp -d)"
  TARGET_FILE="${TMP_DIR}/${BINARY_NAME}"

  log "Downloading ${URL}"
  curl -fL "${URL}" -o "${TARGET_FILE}"
  chmod +x "${TARGET_FILE}"
}

install_binary() {
  INSTALL_DIR="/usr/local/bin"

  if [ ! -w "$INSTALL_DIR" ]; then
    log "Installing to ${INSTALL_DIR} with sudo"
    sudo mv "${TARGET_FILE}" "${INSTALL_DIR}/${BINARY_NAME}"
  else
    mv "${TARGET_FILE}" "${INSTALL_DIR}/${BINARY_NAME}"
  fi

  log "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
}

verify_install() {
  "${BINARY_NAME}" --version || "${BINARY_NAME}" -v || true
}

main() {
  detect_platform
  get_latest_tag
  download_binary
  install_binary
  verify_install
}

main "$@"