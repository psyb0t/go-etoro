#!/bin/bash
# Downloads one eToro API spec into api/. The specs are not committed; each
# generated package fetches its spec as the first step of `go generate`.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

readonly SPEC_BASE_URL="https://api-portal.etoro.com/api-reference"
# The portal answers non-browser user agents with 403.
readonly USER_AGENT="Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0 Safari/537.36"
readonly DOWNLOAD_TIMEOUT_SECONDS=60
readonly SPEC_OPENAPI="openapi.json"
readonly SPEC_ASYNCAPI="asyncapi.json"
SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

usage() {
	printf 'Usage: %s <%s|%s>\n' "$(basename "$0")" "$SPEC_OPENAPI" "$SPEC_ASYNCAPI" >&2
}

main() {
	local spec="${1:-}"
	if [[ "$spec" != "$SPEC_OPENAPI" && "$spec" != "$SPEC_ASYNCAPI" ]]; then
		log ERROR "unknown spec name=${spec}"
		usage
		exit 2
	fi

	mkdir -p "$SPEC_DIR"
	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.${spec}.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp now, the local is gone when EXIT fires
	trap "rm -f -- '${tmp}'" EXIT

	log INFO "downloading spec=${spec}"
	wget -q \
		-T "$DOWNLOAD_TIMEOUT_SECONDS" \
		-U "$USER_AGENT" \
		--header 'Accept: application/json' \
		-O "$tmp" \
		"${SPEC_BASE_URL}/${spec}"

	if ! jq -e 'type == "object" and ((.openapi // .asyncapi) | type == "string")' "$tmp" >/dev/null; then
		log ERROR "downloaded file is not an OpenAPI or AsyncAPI document spec=${spec}"
		exit 1
	fi

	chmod 0644 "$tmp"
	mv -- "$tmp" "${SPEC_DIR}/${spec}"
	log INFO "saved spec=${spec} version=$(jq -r '.info.version' "${SPEC_DIR}/${spec}")"
}

main "$@"
