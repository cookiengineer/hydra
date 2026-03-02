
# TODO

- State.VirtualScreen is necessary due to XTest API, so there needs to be a GlobalState and a LocalState?
  Maybe it makes sense to rename types.State into types.LocalState?

- types.State currently is unfinished. Backport listeners.Init() into types.NewState()
  with properties that make sense without adding redundant listeners

