#!/bin/bash

#http://www.apache.org/licenses/LICENSE-2.0.txt
#
#
#Copyright 2015 Intel Corporation
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.

# Support travis.ci environment matrix:
TEST_TYPE="${TEST_TYPE:-$1}"
UNIT_TEST="${UNIT_TEST:-"gofmt goimports go_vet go_test go_cover"}"

set -e
set -u
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__proj_dir="$(dirname "$__dir")"

# shellcheck source=scripts/common.sh
. "${__dir}/common.sh"

_debug "script directory ${__dir}"
_debug "project directory ${__proj_dir}"

[[ "$TEST_TYPE" =~ ^(small|medium|large|legacy)$ ]] || _error "invalid TEST_TYPE (value must be 'small', 'medium', 'large', or 'legacy', recieved:${TEST_TYPE}"

_gofmt() {
  _info "running gofmt"
  test -z "$(gofmt -l -d $(find . -type f -name '*.go' -not -path "./vendor/*") | tee /dev/stderr)"
}

_goimports() {
  _info "running goimports"
  _go_get golang.org/x/tools/cmd/goimports
  test -z "$(goimports -l -d $(find . -type f -name '*.go' -not -path "./vendor/*") | tee /dev/stderr)"
}

_golint() {
  _info "running golint"
  _go_get github.com/golang/lint/golint
  golint ./...
}

_go_vet() {
  _info "running go_vet"
  go vet ${test_dirs}
}

_go_race() {
  _info "running go test -race"
  go test -race ./...
}

_go_test() {
  _info "running go test type: ${TEST_TYPE}"
  _go_get github.com/smartystreets/goconvey/convey
  _go_get github.com/stretchr/testify/mock
  # Standard go tooling behavior is to ignore dirs with leading underscors
  for dir in $test_dirs;
  do
    if [[ -z ${go_cover+x} ]]; then
      go test --tags="${TEST_TYPE}" -covermode=count -coverprofile="${dir}/profile.tmp" "${dir}"
      if [ -f "${dir}/profile.tmp" ]; then
        tail -n +2 "${dir}/profile.tmp" >> profile.cov
        rm "${dir}/profile.tmp"
      fi
    else
      go test -v --tags="${TEST_TYPE}" "${dir}"
    fi
  done
}

_go_cover() {
  _info "running go tool cover"
  go tool cover -func profile.cov
}

test_small() {
  # The script does automatic checking on a Go package and its sub-packages, including:
  # 1. gofmt         (http://golang.org/cmd/gofmt/)
  # 2. goimports     (https://github.com/bradfitz/goimports)
  # 3. golint        (https://github.com/golang/lint)
  # 4. go vet        (http://golang.org/cmd/vet)
  # 5. race detector (http://blog.golang.org/race-detector)
  # 6. go test
  # 7. test coverage (http://blog.golang.org/cover)
  local go_tests
  go_tests=(gofmt goimports golint go_vet go_race go_test go_cover)

  _debug "available unit tests: ${go_tests[*]}"
  _debug "user specified tests: ${UNIT_TEST}"

  ((n_elements=${#go_tests[@]}, max=n_elements - 1))

  for ((i = 0; i <= max; i++)); do
    if [[ "${UNIT_TEST}" =~ (^| )"${go_tests[i]}"( |$) ]]; then
      _"${go_tests[i]}"
    else
      _info "skipping ${go_tests[i]}"
    fi
  done
}

if [[ $TEST_TYPE == "small" ]]; then
  test_dirs=$(find . -type f -name '*.go' -not -path "./.*" -not -path "*/_*" -not -path "./Godeps/*" -not -path "./vendor/*" -print0 | xargs -0 -n1 dirname| sort -u)

  _debug "go code directories:
${test_dirs}"

  echo "mode: count" > profile.cov
  test_small
fi
