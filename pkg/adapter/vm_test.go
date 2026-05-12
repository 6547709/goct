package adapter

import (
	"reflect"
	"testing"

	"github.com/openlyinc/pointy"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

// TestToVM_IPSplitFiltersEmpty 锁定 Bug 14 修复：
// CloudTower ips 字段在多网卡时是逗号分隔，可能出现 "1.2.3.4,," 这类带空段的值。
// toVM 必须 trim + 过滤空字符串，否则下游 vm.ip 会输出空行。
func TestToVM_IPSplitFiltersEmpty(t *testing.T) {
	tests := []struct {
		name string
		ips  string
		want []string
	}{
		{"single", "10.0.0.1", []string{"10.0.0.1"}},
		{"two", "10.0.0.1,10.0.0.2", []string{"10.0.0.1", "10.0.0.2"}},
		{"with-spaces", "10.0.0.1 , 10.0.0.2", []string{"10.0.0.1", "10.0.0.2"}},
		{"trailing-empty", "10.0.0.1,,", []string{"10.0.0.1"}},
		{"all-empty", ",,", nil},
		{"only-spaces", " , , ", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &models.VM{
				ID:   pointy.String("vm-1"),
				Name: pointy.String("test"),
				Ips:  pointy.String(tt.ips),
			}
			got := toVM(v).IPs
			if len(tt.want) == 0 && len(got) != 0 {
				t.Fatalf("want empty, got %v", got)
			}
			if len(tt.want) != 0 && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("want %v, got %v", tt.want, got)
			}
		})
	}
}

// TestMapBus 验证 bus 字符串到 SDK enum 的映射，未识别值回退 SCSI。
func TestMapBus(t *testing.T) {
	cases := map[string]models.Bus{
		"":         models.BusSCSI,
		"SCSI":     models.BusSCSI,
		"scsi":     models.BusSCSI,
		"IDE":      models.BusIDE,
		"ide":      models.BusIDE,
		"VIRTIO":   models.BusVIRTIO,
		"virtio":   models.BusVIRTIO,
		"NVMe":     models.BusVIRTIO, // CloudTower 暂用 VIRTIO 表示 NVMe
		"NVME":     models.BusVIRTIO,
		"unknown":  models.BusSCSI, // fallback
		" SCSI  ":  models.BusSCSI, // trim
	}
	for in, want := range cases {
		if got := mapBus(in); got != want {
			t.Errorf("mapBus(%q) = %v, want %v", in, got, want)
		}
	}
}
