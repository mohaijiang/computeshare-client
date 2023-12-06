package vm

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	queueTaskV1 "github.com/mohaijiang/computeshare-server/api/queue/v1"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func getVirtManager() *VirtManager {
	logger := log.With(log.NewStdLogger(os.Stdout),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", "1",
		"service.name", "Name",
		"service.version", "Version",
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	manage, err := NewVirtManager(logger)
	if err != nil {
		panic(err)
	}
	return manage
}

func TestCreateVm(t *testing.T) {
	manage := getVirtManager()
	param := queueTaskV1.ComputeInstanceTaskParamVO{
		Name:      "ubuntu1",
		Image:     "ubuntu-20.04",
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC2mLWYddGeahdk6i3muy72XDbppnG4LIDhyj/rSuzLstdVLI7mF7efkwCZgyYcYRJoIjNI5mnb17o7/qVWdgGSiMnSgiPcw4r0Dp1pghWXBEog3o7pI3gicY6//Y4+liqypBEDmBSJnDsMJqVARzFV0rjJLhYSCbYk99LPB1ZLj0mDvIY/1SjRR9bfPuW9Ht6QjkS9DEWIdTrJ0dAaGwJkc+a5pCVzcopq4ycvBVLEnEq4xCrhbNx/LrpYxytA7WXg6kUcN+4Me63QVPxUExcn14qXr5uYxo+ePkoBCNdbqFsm0Z1rxrEX8oGDHvAfsoCpQr/OV8J5WwO7i/QIOyK7 mohaijiang110@163.com",
		Cpu:       1,
		Memory:    1024,
	}
	id, err := manage.Create(param)
	if err != nil {
		panic(err)
	}

	fmt.Println(id)
}

func TestVirtManager_Shutdown(t *testing.T) {
	manage := getVirtManager()

	err := manage.Shutdown("ubuntu1")

	assert.NoError(t, err)
}

func TestVirtManager_Start(t *testing.T) {
	manage := getVirtManager()

	err := manage.Start("ubuntu1")

	assert.NoError(t, err)
}

func TestVirtManager_Status(t *testing.T) {
	manage := getVirtManager()

	status, err := manage.Status("ubuntu1")
	if err != nil {
		return
	}

	assert.NoError(t, err)
	fmt.Println(status)
}

func TestVirtManager_GetIp(t *testing.T) {
	manage := getVirtManager()
	ip, err := manage.GetIp("my-vm")
	if err != nil {
		return
	}

	assert.NoError(t, err)
	fmt.Println(ip)
}

func TestVirtManager_Reboot(t *testing.T) {
	manage := getVirtManager()

	err := manage.Reboot("ubuntu1")

	assert.NoError(t, err)
}

func TestVirtManager_Destroy(t *testing.T) {
	manage := getVirtManager()

	err := manage.Destroy("ubuntu1")
	assert.NoError(t, err)
}

func TestVirtManager_Init(t *testing.T) {
	manage := getVirtManager()

	manage.initBaseData()
}

func TestConsole(t *testing.T) {
	manage := getVirtManager()

	vncPort := manage.GetVncPort("my-vm")
	fmt.Println(vncPort)

}

func TestVirtManager_GetMaxVncPort(t *testing.T) {
	manage := getVirtManager()
	port := manage.GetMaxVncPort()
	fmt.Println(port)
}

func TestVirtManager_VncOpen(t *testing.T) {
	manage := getVirtManager()
	err := manage.VncOpen("my-vm")
	fmt.Println("vnc open", err)

	time.Sleep(time.Second * 10)

	err = manage.VncClose("my-vm")
	fmt.Println("vnc close", err)

	time.Sleep(time.Second * 20)

	fmt.Println("end")
}