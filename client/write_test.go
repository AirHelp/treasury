package client_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/AirHelp/treasury/types"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	mocks "github.com/AirHelp/treasury/backend/mocks"
	"github.com/AirHelp/treasury/client"
	test "github.com/AirHelp/treasury/test/backend"
)

func TestWrite(t *testing.T) {
	dummyClientOptions := &client.Options{
		Backend:      &test.MockBackendClient{},
		S3BucketName: "fake_s3_bucket",
	}
	treasury, err := client.New(dummyClientOptions)
	if err != nil {
		t.Error(err)
	}
	err = treasury.Write(test.Key1, test.KeyValueMap[test.Key1], false)
	if err != nil {
		t.Error(err)
	}
}

var _ = Describe("Write", func() {
	var (
		treasury       *client.Client
		mockCtrl       *gomock.Controller
		mockBackendAPI *mocks.MockBackendAPI
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockBackendAPI = mocks.NewMockBackendAPI(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("secret that will be written does not exist", func() {
		When("force flag is false", func() {
			It("writes the secret successfully", func() {
				treasury = &client.Client{
					Backend: mockBackendAPI,
				}
				mockBackendAPI.EXPECT().GetObject(gomock.Any()).Return(&types.GetObjectOutput{}, nil)
				mockBackendAPI.EXPECT().PutObject(gomock.Any()).Return(nil)
				err := treasury.Write(test.Key1, test.KeyValueMap[test.Key1], false)
				Expect(err).To(BeNil())
			})
		})

		When("force flag is false and GetObject returns unexpected err", func() {
			It("returns err", func() {
				treasury = &client.Client{
					Backend: mockBackendAPI,
				}
				mockBackendAPI.EXPECT().GetObject(gomock.Any()).Return(nil, errors.New("unexpected error"))
				err := treasury.Write(test.Key1, test.KeyValueMap[test.Key1], false)
				Expect(err).To(HaveOccurred())
			})
		})

		When("force flag is false and the error is ParameterNotFound", func() {
			It("writes the secret successfully", func() {
				treasury = &client.Client{
					Backend: mockBackendAPI,
				}
				mockBackendAPI.EXPECT().GetObject(gomock.Any()).Return(nil, fmt.Errorf("%w", &ssmTypes.ParameterNotFound{}))
				mockBackendAPI.EXPECT().PutObject(gomock.Any()).Return(nil)
				err := treasury.Write(test.Key1, test.KeyValueMap[test.Key1], false)
				Expect(err).To(BeNil())
			})
		})
	})

	Context("some secret exists under given key", func() {
		When("force flag is false and secret is the same as that in the parameter store", func() {
			It("returns nil", func() {
				treasury = &client.Client{
					Backend: mockBackendAPI,
				}
				mockBackendAPI.EXPECT().GetObject(gomock.Any()).Return(&types.GetObjectOutput{Value: test.KeyValueMap[test.Key1]}, nil)
				err := treasury.Write(test.Key1, test.KeyValueMap[test.Key1], false)
				Expect(err).To(BeNil())
			})
		})

		When("force flag is false and secret is different than that in the parameter store", func() {
			It("writes the secret successfully", func() {
				treasury = &client.Client{
					Backend: mockBackendAPI,
				}
				mockBackendAPI.EXPECT().GetObject(gomock.Any()).Return(&types.GetObjectOutput{Value: "different_secret"}, nil)
				mockBackendAPI.EXPECT().PutObject(gomock.Any()).Return(nil)
				err := treasury.Write(test.Key1, test.KeyValueMap[test.Key1], false)
				Expect(err).To(BeNil())

			})
		})
	})

	Context("force flag is true", func() {
		It("writes the secret successfully", func() {
			treasury = &client.Client{
				Backend: mockBackendAPI,
			}
			mockBackendAPI.EXPECT().PutObject(gomock.Any()).Return(nil)
			err := treasury.Write(test.Key1, test.KeyValueMap[test.Key1], true)
			Expect(err).To(BeNil())
		})
	})
})
