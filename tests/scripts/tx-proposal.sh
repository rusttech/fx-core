#!/usr/bin/env bash

set -eo pipefail

project_dir="$(git rev-parse --show-toplevel)"
readonly project_dir
readonly out_dir="${project_dir}/out"
readonly proposals_file="${project_dir}/tests/data/proposals.json"

## ARGS: <title> <summary>
## DESC: base64 encode metadata
function base64_metadata() {
  echo '{"title": "'"$1"'","summary": "'"$2"'","metadata":""}' | base64
}

## ARGS: <msg_type>
## DESC: get proposal template
function get_proposal_template() {
  local msg_type=$1
  jq -r --arg msg_type "$msg_type" '.[]|select(.msg_type == $msg_type)' "$proposals_file" > "$out_dir/${msg_type##*.}.json"
}

## ARGS: <msg_type> <amount>
## DESC: query min deposit
function query_min_deposit() {
  local msg_type=$1 amount=$2
  if [[ -z "$msg_type" ]]; then
    echo "$(cosmos_query gov params | jq -r '.params.min_deposit|select(.denom=="'"$STAKING_DENOM"'")|.amount')$STAKING_DENOM" && return
  fi

  base_deposit="$(cosmos_query gov params --msg-type="$msg_type" | jq -r '.params.min_deposit|select(.denom=="'"$STAKING_DENOM"'")|.amount')"
  if [[ "$msg_type" == "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" && -n "$amount" ]]; then
    deposit_threshold=$(cosmos_query gov egf-params | jq -r '.params.egf_deposit_threshold.amount')
    claim_ratio=$(cosmos_query gov egf-params | jq -r '.params.claim_ratio')

    amount_without=${amount%"$STAKING_DENOM"}
    if [[ $(echo "$amount_without - $deposit_threshold" | bc) -gt 0 ]]; then
      echo "$(echo "$amount_without * $claim_ratio" | bc)""$STAKING_DENOM"
    fi
  fi
  echo "${base_deposit}${STAKING_DENOM}"
}

## ARGS: <proposal_file>
## DESC: submit proposal
function submit_proposal() {
  local proposal_file=$1
  msg_type=$(jq -r '.msg_type' "$proposal_file")

  if [ "$msg_type" == "/cosmos.distribution.v1beta1.CommunityPoolSpendProposal" ]; then
    deposit=$(query_min_deposit "$msg_type" "$(jq -r '.amount' "$proposal_file")")
    json_processor "$proposal_file" '.deposit = "'"$deposit"'"'

    cosmos_tx gov submit-legacy-proposal community-pool-spend "$proposal_file" --from "$FROM" -y
  else
    title=$(jq -r '.title' "$proposal_file")
    summary=$(jq -r '.summary' "$proposal_file")
    metadata=$(base64_metadata "$title" "$summary")
    json_processor "$proposal_file" '.metadata = "'"$metadata"'"'

    deposit=$(query_min_deposit "$msg_type")
    json_processor "$proposal_file" '.deposit = "'"$deposit"'"'

    cosmos_tx gov submit-proposal "$proposal_file" --from "$FROM" -y
  fi
}

# shellcheck source=/dev/null
. "${project_dir}/tests/scripts/setup-env.sh"