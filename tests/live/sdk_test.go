package live

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// SDK compatibility: the OpenAI Python SDK's live client, run from python/ with pytest against
// the gateway's OpenAI drop-in route. The fake upstream, the pricing and the keys are this
// suite's, so the Python tests only need the gateway's address and a credential.

func TestSDK_OpenAIPython(t *testing.T) {
	requireGateway(t)
	uv, err := exec.LookPath("uv")
	if err != nil {
		t.Skip("uv is not installed; the OpenAI Python SDK suite needs it")
	}
	env := append(os.Environ(),
		"BIFROST_BASE_URL="+gatewayURL,
		"LIVE_UPSTREAM="+upstreamMode,
		"LIVE_VOICE_MODEL="+voiceModel,
		"LIVE_BACKEND_MODEL="+backendModel,
		"LIVE_BACKEND_MODEL_2="+backendModel2,
	)
	if upstreamMode == fakeUpstream {
		vk := createVirtualKey(t, virtualKeySpec{})
		restricted := createVirtualKey(t, virtualKeySpec{allowedModels: []string{"gpt-4o-mini"}})
		env = append(env, "BIFROST_VK="+vk.Value, "BIFROST_VK_RESTRICTED="+restricted.Value)
	} else if envVirtualKey != "" {
		env = append(env, "BIFROST_VK="+envVirtualKey)
	}

	cmd := exec.Command(uv, "run", "--quiet", "pytest", "-q", "-x", "--timeout=120")
	cmd.Dir = "python"
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	t.Logf("pytest:\n%s", out)
	require.NoError(t, err, "the OpenAI Python SDK suite failed")
}
