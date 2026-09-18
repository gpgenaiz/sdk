package mgmt

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"

	"genaiz.com/genaiz-lib/lang/timez"
	"genaiz.com/genaiz/task/broker"
)

func TestUserLinkInstance_MarshalJSON(t *testing.T) {
	var testCreated = time.Now()
	var expectedCreated = timez.NewTodayFormatter().FormatMillis(testCreated.UnixMilli())
	var testUserLinkInstance = &UserLinkInstance{
		Id:          new(int64(37)),
		Oem:         "expectedOem",
		Handle:      "expectedHandle",
		Version:     "expectedVersion",
		Fqdn:        "expectedFqdn",
		Name:        "expectedName",
		Description: "expectedDescription",
		Created:     testCreated.UnixMilli(),
		Flags:       new(1337),
	}
	var bytes []byte
	var err error

	if bytes, err = testUserLinkInstance.MarshalJSON(); err == nil {
		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
		actual := string(bytes)
		assert.Contains(t, actual, fmt.Sprintf("\"id\":%d", *testUserLinkInstance.Id))
		assert.Contains(t, actual, fmt.Sprintf("\"oem\":\"%s\"", testUserLinkInstance.Oem))
		assert.Contains(t, actual, fmt.Sprintf("\"handle\":\"%s\"", testUserLinkInstance.Handle))
		assert.Contains(t, actual, fmt.Sprintf("\"version\":\"%s\"", testUserLinkInstance.Version))
		assert.Contains(t, actual, fmt.Sprintf("\"name\":\"%s\"", testUserLinkInstance.Name))
		assert.Contains(t, actual, fmt.Sprintf("\"description\":\"%s\"", testUserLinkInstance.Description))
		assert.Contains(t, actual, fmt.Sprintf("\"created\":\"%s\"", expectedCreated))
		assert.Contains(t, actual, fmt.Sprintf("\"flags\":%d", *testUserLinkInstance.Flags))
	} else {
		assert.Fail(t, err.Error())
	}
}

func TestUserLinkInstance_MarshalJSON_Modified(t *testing.T) {
	var testModified = time.Now()
	var expectedModified = timez.NewTodayFormatter().FormatMillis(testModified.UnixMilli())
	var testUserLinkInstance = &UserLinkInstance{
		Id:          new(int64(37)),
		Oem:         "expectedOem",
		Handle:      "expectedHandle",
		Version:     "expectedVersion",
		Fqdn:        "expectedFqdn",
		Name:        "expectedName",
		Description: "expectedDescription",
		Modified:    testModified.UnixMilli(),
		Flags:       new(1337),
	}
	var bytes []byte
	var err error

	if bytes, err = testUserLinkInstance.MarshalJSON(); err == nil {
		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
		actual := string(bytes)
		assert.Contains(t, actual, fmt.Sprintf("\"id\":%d", *testUserLinkInstance.Id))
		assert.Contains(t, actual, fmt.Sprintf("\"oem\":\"%s\"", testUserLinkInstance.Oem))
		assert.Contains(t, actual, fmt.Sprintf("\"handle\":\"%s\"", testUserLinkInstance.Handle))
		assert.Contains(t, actual, fmt.Sprintf("\"version\":\"%s\"", testUserLinkInstance.Version))
		assert.Contains(t, actual, fmt.Sprintf("\"name\":\"%s\"", testUserLinkInstance.Name))
		assert.Contains(t, actual, fmt.Sprintf("\"description\":\"%s\"", testUserLinkInstance.Description))
		assert.Contains(t, actual, fmt.Sprintf("\"modified\":\"%s\"", expectedModified))
		assert.Contains(t, actual, fmt.Sprintf("\"flags\":%d", *testUserLinkInstance.Flags))
	} else {
		assert.Fail(t, err.Error())
	}
}

func TestUserLinkInstance_MarshalSlice(t *testing.T) {
	var testCreated = time.Now()
	var expectedCreated = timez.NewTodayFormatter().FormatMillis(testCreated.UnixMilli())
	var testUserLinkInstance = &UserLinkInstance{
		Id:          new(int64(37)),
		Oem:         "expectedOem",
		Handle:      "expectedHandle",
		Version:     "expectedVersion",
		Fqdn:        "expectedFqdn",
		Name:        "expectedName",
		Description: "expectedDescription",
		Created:     testCreated.UnixMilli(),
		Active:      true,
	}
	var values []string
	var err error

	if values, err = testUserLinkInstance.MarshalSlice(); err == nil {
		assert.NoError(t, err)
		assert.NotEmpty(t, values)
		assert.Equal(t, values[0], cast.ToString(testUserLinkInstance.Id))
		assert.Equal(t, values[1], testUserLinkInstance.Name)
		assert.Equal(t, values[2], testUserLinkInstance.Fqdn)
		assert.Equal(t, values[3], testUserLinkInstance.Version)
		assert.Equal(t, values[4], expectedCreated)
		assert.Equal(t, values[5], "yes")
	} else {
		assert.Fail(t, err.Error())
	}
}

func TestUserLinkInstance_MarshalSlice_NoCreate(t *testing.T) {
	var testUserLinkInstance = &UserLinkInstance{
		Id:          new(int64(37)),
		Oem:         "expectedOem",
		Handle:      "expectedHandle",
		Version:     "expectedVersion",
		Sequence:    new(1),
		Fqdn:        "expectedFqdn",
		Name:        "expectedName",
		Description: "expectedDescription",
		Active:      false,
	}
	var values []string
	var err error

	if values, err = testUserLinkInstance.MarshalSlice(); err == nil {
		assert.NoError(t, err)
		assert.NotEmpty(t, values)
		assert.Equal(t, values[0], cast.ToString(testUserLinkInstance.Id))
		assert.Equal(t, values[1], testUserLinkInstance.Name)
		assert.Equal(t, values[2], testUserLinkInstance.Fqdn)
		assert.Equal(t, values[3], testUserLinkInstance.Version)
		assert.Equal(t, values[4], "-")
		assert.Equal(t, values[5], "no")
	} else {
		assert.Fail(t, err.Error())
	}
}

func TestToUserLinkInstance(t *testing.T) {
	var testLink = &broker.DataLinkInstance{
		Id:          new(int64(37)),
		Name:        "expectedName",
		Description: "expectedDescription",
		DataLink: broker.DataLink{
			Oem:     "expectedOem",
			Handle:  "expectedHandle",
			Version: "expectedVersion",
			Seq:     new(1),
			Flags:   new(broker.DataLinkFlags.Active),
		},
		Visibility: "visible",
		Properties: map[string]string{
			"prop": "value",
		},
		Created:  1,
		Modified: 2,
		Flags:    new(broker.DataLinkInstanceFlags.Active),
	}

	actual := ToUserLinkInstance(testLink)
	assert.NotEmpty(t, actual)
	assert.Equal(t, *testLink.Id, *actual.Id)
	assert.Equal(t, testLink.Name, actual.Name)
	assert.Equal(t, testLink.Description, actual.Description)
	assert.Equal(t, testLink.DataLink.Oem, actual.Oem)
	assert.Equal(t, testLink.DataLink.Handle, actual.Handle)
	assert.Equal(t, fmt.Sprintf("%s-rc-%d", testLink.DataLink.Version, *testLink.DataLink.Seq), actual.Version)
	assert.Equal(t, testLink.DataLink.Seq, actual.Sequence)
	actualPropValue, ok := actual.Properties["prop"]
	assert.True(t, ok)
	assert.Equal(t, testLink.Properties["prop"], actualPropValue)
	assert.Equal(t, testLink.Created, actual.Created)
	assert.Equal(t, testLink.Modified, actual.Modified)
	assert.Equal(t, *testLink.Flags, *actual.Flags)
	assert.True(t, actual.Active)
}

func TestToUserLinkInstance_Released(t *testing.T) {
	var testLink = &broker.DataLinkInstance{
		Id:          new(int64(37)),
		Name:        "expectedName",
		Description: "expectedDescription",
		DataLink: broker.DataLink{
			Oem:     "expectedOem",
			Handle:  "expectedHandle",
			Version: "expectedVersion",
			Seq:     new(1),
			Flags:   new(broker.DataLinkFlags.Active | broker.DataLinkFlags.Released),
		},
		Visibility: "visible",
		Properties: map[string]string{
			"prop": "value",
		},
		Created:  1,
		Modified: 2,
		Flags:    new(broker.DataLinkInstanceFlags.Active),
	}

	actual := ToUserLinkInstance(testLink)
	assert.NotEmpty(t, actual)
	assert.Equal(t, *testLink.Id, *actual.Id)
	assert.Equal(t, testLink.Name, actual.Name)
	assert.Equal(t, testLink.Description, actual.Description)
	assert.Equal(t, testLink.DataLink.Oem, actual.Oem)
	assert.Equal(t, testLink.DataLink.Handle, actual.Handle)
	assert.Equal(t, testLink.DataLink.Version, actual.Version)
	assert.Equal(t, testLink.DataLink.Seq, actual.Sequence)
	actualPropValue, ok := actual.Properties["prop"]
	assert.True(t, ok)
	assert.Equal(t, testLink.Properties["prop"], actualPropValue)
	assert.Equal(t, testLink.Created, actual.Created)
	assert.Equal(t, testLink.Modified, actual.Modified)
	assert.Equal(t, *testLink.Flags, *actual.Flags)
	assert.True(t, actual.Active)
}
