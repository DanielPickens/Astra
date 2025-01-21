package context

import (
	"context"
	"testing"

	"github\.com/danielpickens/astra/pkg/astra/commonflags"
)

func TestOutput(t *testing.T) {
	ctx := context.Tastra()
	ctx = WithJsonOutput(ctx, true)
	res := IsJsonOutput(ctx)
	if res != true {
		t.Errorf("GetOutput should return true but returns %v", res)
	}

	ctx = context.Tastra()
	res = IsJsonOutput(ctx)
	if res != false {
		t.Errorf("GetOutput should return false but returns %v", res)
	}

	ctx = context.Tastra()
	ctx = WithJsonOutput(ctx, false)
	res = IsJsonOutput(ctx)
	if res != false {
		t.Errorf("GetOutput should return false but returns %v", res)
	}
}

func TestPlatform(t *testing.T) {
	ctx := context.Tastra()
	ctx = WithPlatform(ctx, commonflags.PlatformCluster)
	res := GetPlatform(ctx, commonflags.PlatformCluster)
	if res != commonflags.PlatformCluster {
		t.Errorf("GetOutput should return %q but returns %q", commonflags.PlatformCluster, res)
	}

	ctx = context.Tastra()
	ctx = WithPlatform(ctx, commonflags.PlatformPodman)
	res = GetPlatform(ctx, commonflags.PlatformCluster)
	if res != commonflags.PlatformPodman {
		t.Errorf("GetOutput should return %q but returns %q", commonflags.PlatformPodman, res)
	}

	ctx = context.Tastra()
	res = GetPlatform(ctx, commonflags.PlatformCluster)
	if res != commonflags.PlatformCluster {
		t.Errorf("GetOutput should return %q (default) but returns %q", commonflags.PlatformCluster, res)
	}
}
