#!/bin/bash
# Makes the downloaded AsyncAPI document describe what eToro actually sends.
#
# Some content schemas, such as InstrumentRate, declare number fields but say
# in their own description that every value arrives as a JSON string. Strict
# decoding of such a field fails on real frames, so in every schema whose
# description says so, number, integer and boolean properties become strings.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR
readonly SOURCE_SPEC="${SPEC_DIR}/asyncapi.json"
readonly OUTPUT_SPEC="${SPEC_DIR}/asyncapi.wire.json"
# Matched against the lowercased schema description.
readonly WIRE_NOTE="on the wire every value is a json string"

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

main() {
	if [[ ! -f "$SOURCE_SPEC" ]]; then
		log ERROR "source spec missing path=${SOURCE_SPEC}"
		exit 1
	fi

	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.asyncapi.wire.json.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp now, the local is gone when EXIT fires
	trap "rm -f -- '${tmp}'" EXIT

	jq --arg note "$WIRE_NOTE" '
		.components.schemas |= map_values(
			if ((.description // "") | ascii_downcase | contains($note)) then
				.properties |= map_values(
					if (.type == "number" or .type == "integer" or .type == "boolean") then
						.type = "string" | del(.format)
					else
						.
					end
				)
			else
				.
			end
		)
	' "$SOURCE_SPEC" >"$tmp"

	chmod 0644 "$tmp"
	mv -- "$tmp" "$OUTPUT_SPEC"
	log INFO "matched wire types schemas=$(jq -r --arg note "$WIRE_NOTE" '[.components.schemas | to_entries[] | select((.value.description // "") | ascii_downcase | contains($note)) | .key] | join(",")' "$OUTPUT_SPEC")"
}

main "$@"
