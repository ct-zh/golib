package xmind

import "testing"

func TestCreate(t *testing.T) {
	var data = map[string]interface{}{
		"test1": map[string]interface{}{
			"test1 child": map[string]interface{}{
				"test1 child child": struct{}{},
			},
		},
		"test2": struct{}{},
		"test3": map[string]interface{}{
			"test3 222222": map[string]interface{}{
				"test3 33333333": map[string]interface{}{
					"test3 44444444": map[string]interface{}{
						"test3 555555555": struct{}{},
					},
				},
			},
		},
	}

	xmind := CreateXMindFromMap("tree for xmind", data)
	str := xmind.Output()
	t.Logf("xmind:\n%s", str)
}
