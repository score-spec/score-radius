# Known limitations from Score to Radius

Here is a list of known limitations about the mapping between Score and Radius. If you want to see one of these implemented, please file an issue with your example here: https://github.com/mathieu-benoit/score-radius/issues. It will influence the roadmap of this project, thanks!

## On Score

- Just the first container from Score is mapped to an `Applications.Core/containers`.
  - Note: Most of the case it's one container per Workload per Score file.
- In `containers`'s, `resources.cpu` and `resources.memory` are not in `Applications.Core/containers`.
  - Note: Maybe to map as [PodSpecTemplate](https://docs.radapp.io/guides/author-apps/kubernetes/patch-podspec/)?

## On Radius

- In `Applications.Core/containers`'s `env`, secret reference is not yet taken into account like illustrated [here](https://docs.radapp.io/reference/resource-schema/core-schema/container-schema/#container).
  - Note: Something that we can easily do like we do with `score-k8s`'s `encodeSecretRef` function.
- In `Applications.Core/containers`, `workingDir`.
- In `Applications.Core/containers`'s `extensions`, [`manualscaling`](https://docs.radapp.io/reference/resource-schema/core-schema/container-schema/#manualscaling).
  - Note: Something that we can easily do with a Workload's `annotation`.