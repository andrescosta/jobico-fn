

New Test()
.Service(Ctl.Service)
.Error(Connect,1)
.Run()


plt, err := NewPlatform(ctx)
.Service(Ctl)
.Error(Connect,1)
.Service(Listener)
.Error(Stop)
.Build()

plt, err := NewPlatformBuilder(ctx)
.WithService(Ctl)
.WithError(Connect,1)
.WithService(Listener)
.WithError(Stop)
.Build()


err := plt.StopAndDispose()
