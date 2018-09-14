package model

type Options struct {
	FallBackUser string `hcl:"fallback_user"`
	Meta         Meta   `hcl:"meta"`
}

type Meta struct {
	User           string `hcl:"user"`
	Groups         string `hcl:"groups"`
	HostVarsPrefix string `hcl:"hostvars_prefix"`
	Env            string `hcl:"environment"`
}
