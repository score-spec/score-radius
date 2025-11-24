
## `score-radius`

- `--version`|`-v`: version for `score-radius`

## `score-radius init`

Initialize the local state directory and sample Score file.

- `--file`|`-f` - The score file to initialize (default `score.yaml`).
- `--no-sample` - Disable generation of the sample score file.
- `--provisioners` - 

## `score-radius generate`

Run the conversion from Score file to output manifests.

- `--image`|`-i` - An optional container image to use for any container with image == '.'.
- `--output`|`-o` - The output manifests file to write the manifests to (default `app.bicep`).
- `--override-property` - An optional set of path=key overrides to set or remove.
- `--overrides-file` - `An optional file of Score overrides to merge in.

## `score-radius provisioners`

### `list`

The list command will list out the resources provisioners. This requires an active `score-redius` state after `init` has been run. The list of resources provisioners will be empty if no provisioners are defined.

- `--format`|`-f` - Format of the output: `table` (default), `json`.
