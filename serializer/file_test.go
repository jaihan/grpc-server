package serializer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/techschool/pcbook/pb"
	"gitlab.com/techschool/pcbook/sample"
	"gitlab.com/techschool/pcbook/serializer"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &pb.Laptop{}
	err2 := serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err2)
	require.True(t, proto.Equal(laptop1, laptop2))

	err3 := serializer.WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err3)
}
