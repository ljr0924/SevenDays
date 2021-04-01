package cache

import (
    "strconv"
    "testing"
)

func TestHashMap(t *testing.T) {
    hash := NewMap(3, func(data []byte) uint32 {
        i, _ := strconv.Atoi(string(data))
        return uint32(i)
    })

    hash.Add("2", "4", "6")

    testCases := map[string]string{
        "2": "2",
        "11": "2",
        "23": "4",
        "27": "2",
    }

    for k, v := range testCases {
        if hash.Get(k) != v {
            t.Errorf("Asking for %s, should have yielded %s", k, v)
        }
    }

    hash.Add("8")

    testCases["27"] = "8"

    for k, v := range testCases {
        if hash.Get(k) != v {
            t.Errorf("Asking for %s, should have yielded %s", k, v)
        }
    }

}
