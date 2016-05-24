#!/bin/bash

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

detect_go_dep() {
  [[ -f "${__proj_dir}/Godeps/Godeps.json" ]] && _dep='godep'
  [[ -f "${__proj_dir}/glide.yaml" ]] && _dep='glide'
  [[ -f "${__proj_dir}/vendor/vendor.json" ]] && _dep='govendor'
  _info "golang dependency tool: ${_dep}"
  echo "${_dep}"
}

install_go_dep() {
  local _dep=${_dep:=$(_detect_dep)}
  _info "ensuring ${_dep} is available"
  case $_dep in
    godep)
      _go_get github.com/tools/godep
      ;;
    glide)
      _go_get github.com/Masterminds/glide
      ;;
    govendor)
      _go_get github.com/kardianos/govendor
      ;;
  esac
}

restore_go_dep() {
  local _dep=${_dep:=$(_detect_dep)}
  _info "restoring dependency with ${_dep}"
  case $_dep in
    godep)
      godep restore
      ;;
    glide)
      glide install
      ;;
    govendor)
      govendor sync
      ;;
  esac
}

detect_go_dep
install_go_dep
restore_go_dep
