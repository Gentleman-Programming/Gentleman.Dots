export def main [
  path?: string
  --export: list
  --shellvars (-s)
  --fn (-f): list
] {
  let fn_args = if ($fn | is-not-empty) {
    ['--shellfns' ($fn | str join ',')]
  } else {
    []
  }

  let path_args = if $path != null {
    [($path | path expand)]
  } else {
    []
  }

  let raw = ($in | str join "\n") | bash-env-json ...($fn_args ++ $path_args) | complete
  let raw_json = $raw.stdout | from json

  let error = $raw_json | get -i error
  if $error != null {
    error make { msg: $error }
  } else if $raw.exit_code != 0 {
    error make { msg: $"unexpected failure from bash-env-json ($raw.stderr)" }
  }

  if ($export | is-not-empty) {
    print "warning: --export is deprecated, use --shellvars(-s) instead"
    let exported_shellvars = ($raw_json.shellvars | select -i ...$export)
    $raw_json.env | merge ($exported_shellvars)
  } else if $shellvars or ($fn | is-not-empty) {
    $raw_json
  } else {
    $raw_json.env
  }
}
