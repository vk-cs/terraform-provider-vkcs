#!/usr/bin/env bash
set -euo pipefail

detect_build_target() {
    case "$(uname -s)" in
        Darwin) printf 'build_darwin\nterraform-provider-vkcs_darwin\n' ;;
        Linux)  printf 'build_linux\nterraform-provider-vkcs_linux\n' ;;
        *) printf 'Unsupported OS: %s\n' "$(uname -s)" >&2; return 2 ;;
    esac
}

build_provider_if_stale() {
    local repo_root="$1"
    local binary="$2"
    local target="$3"
    if [[ -x "${repo_root}/${binary}" ]] \
        && [[ -z "$(find "${repo_root}" -name '*.go' -newer "${repo_root}/${binary}" -not -path "${repo_root}/vendor/*" -not -path "${repo_root}/bin/*" -print -quit 2>/dev/null)" ]]; then
        echo "==> Reusing existing ${binary}"
        return 0
    fi
    echo "==> Building ${binary}"
    (cd "${repo_root}" && make "${target}" >/dev/null)
}

write_cli_config() {
    local plugin_dir="$1"
    cat > "${plugin_dir}/terraformrc" <<EOF
provider_installation {
  dev_overrides {
    "vk-cs/vkcs" = "${plugin_dir}"
  }
  direct {}
}
EOF
}

list_dirs() {
    local roots=("$@")
    if (( ${#roots[@]} == 0 )); then
        roots=(examples)
    fi
    find "${roots[@]}" -name '*.tf' -not -path '*/.terraform/*' -print0 \
        | xargs -0 -n1 dirname | sort -u
}

list_tf_files() {
    local dir="$1"
    local exclude="$2"
    find "${dir}" -maxdepth 1 -name '*.tf' \
        -not -name "${exclude}" -not -path '*/.terraform/*' | sort
}

dir_has_required_providers() {
    local dir="$1"
    local exclude="$2"
    find "${dir}" -maxdepth 1 -name '*.tf' -not -name "${exclude}" -print0 \
        | xargs -0 grep -l "required_providers" >/dev/null 2>&1
}

print_dir_files() {
    local repo_root="$1"
    shift
    local f
    for f in "$@"; do
        if [[ -L "${f}" ]]; then
            printf '      %s -> %s\n' "${f#"${repo_root}/"}" "$(readlink "${f}")"
        else
            printf '      %s\n' "${f#"${repo_root}/"}"
        fi
    done
}

validate_dir() {
    local dir="$1"
    local inject_name="$2"
    local inject_body="$3"

    local inject_file=""
    if ! dir_has_required_providers "${dir}" "${inject_name}"; then
        inject_file="${dir}/${inject_name}"
        printf '%s' "${inject_body}" > "${inject_file}"
    fi

    local out rc=0
    out=$(
        cd "${dir}" || exit 1
        rm -rf .terraform .terraform.lock.hcl
        terraform init -backend=false -input=false -no-color 2>&1
        terraform validate -no-color 2>&1
    ) || rc=$?

    [[ -n "${inject_file}" ]] && rm -f "${inject_file}"

    if (( rc != 0 )); then
        printf '%s' "${out}"
        return 1
    fi
    return 0
}

providers_tf_body() {
    cat <<'TF'
terraform {
  required_providers {
    vkcs = {
      source = "vk-cs/vkcs"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}
TF
}

run_loop() {
    local repo_root="$1"
    local inject_name="$2"
    local inject_body="$3"
    shift 3
    local -a dirs=("$@")

    echo "==> Validating ${#dirs[@]} example dirs"

    local -a failed=()
    local passed=0
    local total_files=0
    local dir out
    for dir in "${dirs[@]}"; do
        local -a tf_files=()
        mapfile -t tf_files < <(list_tf_files "${dir}" "${inject_name}")
        total_files=$((total_files + ${#tf_files[@]}))

        if out=$(validate_dir "${dir}" "${inject_name}" "${inject_body}"); then
            printf '  %s  [ok, %d file(s)]\n' "${dir}" "${#tf_files[@]}"
            print_dir_files "${repo_root}" "${tf_files[@]}"
            passed=$((passed + 1))
        else
            failed+=("${dir}")
            printf '  %s  [FAIL]\n' "${dir}"
            print_dir_files "${repo_root}" "${tf_files[@]}"
            printf '%s\n' "${out}" | sed 's/^/      /'
        fi
    done

    echo
    echo "==> Summary: ${passed} ok / ${#failed[@]} failed / ${total_files} tf files validated"
    if (( ${#failed[@]} > 0 )); then
        printf '  - %s\n' "${failed[@]}"
        return 1
    fi
}

main() {
    local repo_root
    repo_root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
    cd "${repo_root}"

    local target_pair build_target binary
    target_pair=$(detect_build_target)
    build_target=$(printf '%s' "${target_pair}" | sed -n '1p')
    binary=$(printf '%s' "${target_pair}" | sed -n '2p')

    build_provider_if_stale "${repo_root}" "${binary}" "${build_target}"

    local plugin_dir
    plugin_dir=$(mktemp -d -t tf-validate-examples-XXXXXX)

    cp "${repo_root}/${binary}" "${plugin_dir}/terraform-provider-vkcs"
    write_cli_config "${plugin_dir}"

    export TF_CLI_CONFIG_FILE="${plugin_dir}/terraformrc"
    export TF_IN_AUTOMATION=1
    export TF_PLUGIN_CACHE_DIR="${plugin_dir}/cache"
    mkdir -p "${TF_PLUGIN_CACHE_DIR}"

    local inject_name="_validate-providers.tf"
    local inject_body
    inject_body=$(providers_tf_body)

    local -a dirs
    mapfile -t dirs < <(list_dirs "$@")

    local rc=0
    run_loop "${repo_root}" "${inject_name}" "${inject_body}" "${dirs[@]}" || rc=$?
    rm -rf "${plugin_dir}"
    return "${rc}"
}

main "$@"
