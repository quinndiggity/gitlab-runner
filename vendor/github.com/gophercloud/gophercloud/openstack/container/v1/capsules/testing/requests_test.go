package testing

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/container/v1/capsules"
	th "github.com/gophercloud/gophercloud/testhelper"
	fakeclient "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGetCapsule(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleCapsuleGetSuccessfully(t)

	actualCapsule, err := capsules.Get(fakeclient.ServiceClient(), "cc654059-1a77-47a3-bfcf-715bde5aad9e").Extract()

	th.AssertNoErr(t, err)

	uuid := "cc654059-1a77-47a3-bfcf-715bde5aad9e"
	status := "Running"
	id := 1
	userID := "d33b18c384574fd2a3299447aac285f0"
	projectID := "6b8ffef2a0ac42ee87887b9cc98bdf68"
	cpu := float64(1)
	memory := "1024M"
	metaName := "test"

	createdAt, _ := time.Parse(gophercloud.RFC3339ZNoT, "2018-01-12 09:37:25+00:00")
	updatedAt, _ := time.Parse(gophercloud.RFC3339ZNoT, "2018-01-12 09:37:25+01:00")
	links := []interface{}{
		map[string]interface{}{
			"href": "http://10.10.10.10/v1/capsules/cc654059-1a77-47a3-bfcf-715bde5aad9e",
			"rel":  "self",
		},
		map[string]interface{}{
			"href": "http://10.10.10.10/capsules/cc654059-1a77-47a3-bfcf-715bde5aad9e",
			"rel":  "bookmark",
		},
	}
	capsuleVersion := "beta"
	restartPolicy := "always"
	metaLabels := map[string]string{
		"web": "app",
	}
	containersUUIDs := []string{
		"1739e28a-d391-4fd9-93a5-3ba3f29a4c9b",
		"d1469e8d-bcbc-43fc-b163-8b9b6a740930",
	}
	addresses := map[string][]capsules.Address{
		"b1295212-64e1-471d-aa01-25ff46f9818d": []capsules.Address{
			{
				PreserveOnDelete: false,
				Addr:             "172.24.4.11",
				Port:             "8439060f-381a-4386-a518-33d5a4058636",
				Version:          float64(4),
				SubnetID:         "4a2bcd64-93ad-4436-9f48-3a7f9b267e0a",
			},
		},
	}
	volumesInfo := map[string][]string{
		"67618d54-dd55-4f7e-91b3-39ffb3ba7f5f": []string{
			"4b725a92-2197-497b-b6b1-fb8caa4cb99b",
		},
	}

	expectedCapsule := capsules.Capsule{
		UUID:            uuid,
		ID:              id,
		UserID:          userID,
		ProjectID:       projectID,
		CPU:             cpu,
		Status:          status,
		Memory:          memory,
		MetaName:        metaName,
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		Links:           links,
		CapsuleVersion:  capsuleVersion,
		RestartPolicy:   restartPolicy,
		MetaLabels:      metaLabels,
		ContainersUUIDs: containersUUIDs,
		Addresses:       addresses,
		VolumesInfo:     volumesInfo,
	}

	th.AssertDeepEquals(t, &expectedCapsule, actualCapsule)
}
