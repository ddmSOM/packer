package packer

// A Build represents a single job within Packer that is responsible for
// building some machine image artifact. Builds are meant to be parallelized.
type Build interface {
	Prepare()
	Run(ui Ui)
}

// A build struct represents a single build job, the result of which should
// be a single machine image artifact. This artifact may be comprised of
// multiple files, of course, but it should be for only a single provider
// (such as VirtualBox, EC2, etc.).
type coreBuild struct {
	name    string
	builder Builder
	rawConfig interface{}

	prepareCalled bool
}

// Implementers of Builder are responsible for actually building images
// on some platform given some configuration.
//
// Prepare is responsible for reading in some configuration, in the raw form
// of map[string]interface{}, and storing that state for use later. Any setup
// should be done in this method. Note that NO side effects should really take
// place in prepare. It is meant as a state setup step only.
//
// Run is where the actual build should take place. It takes a Build and a Ui.
type Builder interface {
	Prepare(config interface{})
	Run(build Build, ui Ui)
}

// Prepare prepares the build by doing some initialization for the builder
// and any hooks. This _must_ be called prior to Run.
func (b *coreBuild) Prepare() {
	b.prepareCalled = true
	b.builder.Prepare(b.rawConfig)
}

// Runs the actual build. Prepare must be called prior to running this.
func (b *coreBuild) Run(ui Ui) {
	if !b.prepareCalled {
		panic("Prepare must be called first")
	}

	b.builder.Run(b, ui)
}
