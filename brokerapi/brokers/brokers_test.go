package brokers_test

import (
	. "gcp-service-broker/brokerapi/brokers"

	"code.cloudfoundry.org/lager"
	"gcp-service-broker/brokerapi/brokers"
	"gcp-service-broker/brokerapi/brokers/models"
	"gcp-service-broker/brokerapi/brokers/models/modelsfakes"
	"gcp-service-broker/db_service"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
)

var _ = Describe("Brokers", func() {
	var (
		gcpBroker                *GCPAsyncServiceBroker
		err                      error
		logger                   lager.Logger
		serviceNameToId          map[string]string = make(map[string]string)
		bqProvisionDetails       models.ProvisionDetails
		cloudSqlProvisionDetails models.ProvisionDetails
	)

	BeforeEach(func() {
		logger = lager.NewLogger("brokers_test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))

		testDb, _ := gorm.Open("sqlite3", "test.db")
		testDb.CreateTable(models.ServiceInstanceDetails{})
		testDb.CreateTable(models.ServiceBindingCredentials{})
		testDb.CreateTable(models.PlanDetails{})
		testDb.CreateTable(models.ProvisionRequestDetails{})

		db_service.DbConnection = testDb

		os.Setenv("ROOT_SERVICE_ACCOUNT_JSON", `{
			"type": "service_account",
			"project_id": "foo",
			"private_key_id": "something",
			"private_key": "foobar",
			"client_email": "example@gmail.com",
			"client_id": "1",
			"auth_uri": "somelink",
			"token_uri": "somelink",
			"auth_provider_x509_cert_url": "somelink",
			"client_x509_cert_url": "somelink"
		      }`)
		os.Setenv("SECURITY_USER_NAME", "username")
		os.Setenv("SECURITY_USER_PASSWORD", "password")
		os.Setenv("SERVICES", `[
			{
			  "id": "b9e4332e-b42b-4680-bda5-ea1506797474",
			  "description": "A Powerful, Simple and Cost Effective Object Storage Service",
			  "name": "google-storage",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google Cloud Storage",
			    "longDescription": "A Powerful, Simple and Cost Effective Object Storage Service",
			    "documentationUrl": "https://cloud.google.com/storage/docs/overview",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "storage"]
			},
			{
			  "id": "628629e3-79f5-4255-b981-d14c6c7856be",
			  "description": "A global service for real-time and reliable messaging and streaming data",
			  "name": "google-pubsub",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google PubSub",
			    "longDescription": "A global service for real-time and reliable messaging and streaming data",
			    "documentationUrl": "https://cloud.google.com/pubsub/docs/",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "pubsub"]
			},
			{
			  "id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
			  "description": "A fast, economical and fully managed data warehouse for large-scale data analytics",
			  "name": "google-bigquery",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google BigQuery",
			    "longDescription": "A fast, economical and fully managed data warehouse for large-scale data analytics",
			    "documentationUrl": "https://cloud.google.com/bigquery/docs/",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "bigquery"]
			},
			{
			  "id": "4bc59b9a-8520-409f-85da-1c7552315863",
			  "description": "Google Cloud SQL is a fully-managed MySQL database service",
			  "name": "google-cloudsql",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google CloudSQL",
			    "longDescription": "Google Cloud SQL is a fully-managed MySQL database service",
			    "documentationUrl": "https://cloud.google.com/sql/docs/",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "cloudsql"]
			},
			{
			  "id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
			  "description": "Machine Learning Apis including Vision, Translate, Speech, and Natural Language",
			  "name": "google-ml-apis",
			  "bindable": true,
			  "plan_updateable": false,
			  "metadata": {
			    "displayName": "Google Machine Learning APIs",
			    "longDescription": "Machine Learning Apis including Vision, Translate, Speech, and Natural Language",
			    "documentationUrl": "https://cloud.google.com/ml/",
			    "supportUrl": "https://cloud.google.com/support/"
			  },
			  "tags": ["gcp", "ml"]
			}
		      ]`)
		os.Setenv("PRECONFIGURED_PLANS", `[
			{
			  "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
			  "name": "standard",
			  "display_name": "Standard",
			  "description": "Standard storage class",
			  "features": {"storage_class": "STANDARD"}
			},
			{
			  "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
			  "name": "nearline",
			  "display_name": "Nearline",
			  "description": "Nearline storage class",
			  "features": {"storage_class": "NEARLINE"}
			},
			{
			  "service_id": "b9e4332e-b42b-4680-bda5-ea1506797474",
			  "name": "reduced_availability",
			  "display_name": "Durable Reduced Availability",
			  "description": "Durable Reduced Availability storage class",
			  "features": {"storage_class": "DURABLE_REDUCED_AVAILABILITY"}
			},
			{
			  "service_id": "628629e3-79f5-4255-b981-d14c6c7856be",
			  "name": "default",
			  "display_name": "Default",
			  "description": "PubSub Default plan",
			  "features": ""
			},
			{ "service_id": "f80c0a3e-bd4d-4809-a900-b4e33a6450f1",
			  "name": "default",
			  "display_name": "Default",
			  "description": "BigQuery default plan",
			  "features": ""
			},
			{
			  "service_id": "5ad2dce0-51f7-4ede-8b46-293d6df1e8d4",
			  "name": "default",
			  "display_name": "Default",
			  "description": "Machine Learning api default plan",
			  "features": ""
			}
		      ]`)

		os.Setenv("PLANS", `{
			"test_plan": {
				"guid": "foo",
				"name": "bar",
				"description": "testplan",
				"tier": "4",
				"pricing_plan": "athing",
				"max_disk_size": "20",
				"display_name": "FOOBAR",
				"service": "4bc59b9a-8520-409f-85da-1c7552315863"
			}
		}`)

		gcpBroker, err = brokers.New(logger)
		if err != nil {
			logger.Error("error", err)
		}

		var someBigQueryPlanId string
		var someCloudSQLPlanId string
		for _, service := range *gcpBroker.Catalog {
			serviceNameToId[service.Name] = service.ID
			if service.Name == BigqueryName {
				someBigQueryPlanId = service.Plans[0].ID
			}
			if service.Name == CloudsqlName {
				someCloudSQLPlanId = service.Plans[0].ID
			}
		}

		for k, _ := range gcpBroker.ServiceBrokerMap {
			async := true
			if k == serviceNameToId[brokers.CloudsqlName] {
				async = true
			}
			gcpBroker.ServiceBrokerMap[k] = &modelsfakes.FakeServiceBrokerHelper{
				AsyncStub: func() bool { return async },
			}
		}

		bqProvisionDetails = models.ProvisionDetails{
			ServiceID: serviceNameToId[brokers.BigqueryName],
			PlanID:    someBigQueryPlanId,
		}

		cloudSqlProvisionDetails = models.ProvisionDetails{
			ServiceID: serviceNameToId[brokers.CloudsqlName],
			PlanID:    someCloudSQLPlanId,
		}

	})

	Describe("Broker init", func() {
		It("should have 5 services in sevices map", func() {
			Expect(len(gcpBroker.ServiceBrokerMap)).To(Equal(5))
		})

		It("should have a default client", func() {
			Expect(gcpBroker.GCPClient).NotTo(Equal(&http.Client{}))
		})

		It("should have loaded credentials correctly and have a project id", func() {
			Expect(gcpBroker.RootGCPCredentials.ProjectId).To(Equal("foo"))
		})
	})

	Describe("getting broker catalog", func() {
		It("should have 5 services available", func() {
			Expect(len(gcpBroker.Services())).To(Equal(5))
		})

		It("should have 3 storage plans available", func() {
			serviceList := gcpBroker.Services()
			for _, s := range serviceList {
				if s.ID == serviceNameToId[StorageName] {
					Expect(len(s.Plans)).To(Equal(3))
				}
			}

		})
	})

	Describe("provision", func() {
		Context("when the bigquery service id is provided", func() {
			It("should call bigquery provisioning", func() {
				bqId := serviceNameToId[brokers.BigqueryName]
				_, err := gcpBroker.Provision("something", bqProvisionDetails, true)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(gcpBroker.ServiceBrokerMap[bqId].(*modelsfakes.FakeServiceBrokerHelper).ProvisionCallCount()).To(Equal(1))
			})

		})

		Context("when too many services are provisioned", func() {
			It("should return an error", func() {
				gcpBroker.InstanceLimit = 0
				_, err := gcpBroker.Provision("something", bqProvisionDetails, true)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(models.ErrInstanceLimitMet))
			})
		})

		Context("when an unrecognized service is provisioned", func() {
			It("should return an error", func() {
				_, err = gcpBroker.Provision("something", models.ProvisionDetails{
					ServiceID: "nope",
					PlanID:    "nope",
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when an unrecognized plan is provisioned", func() {
			It("should return an error", func() {
				_, err = gcpBroker.Provision("something", models.ProvisionDetails{
					ServiceID: serviceNameToId[BigqueryName],
					PlanID:    "nope",
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when duplicate services are provisioned", func() {
			It("should return an error", func() {
				_, _ = gcpBroker.Provision("something", bqProvisionDetails, true)
				_, err := gcpBroker.Provision("something", bqProvisionDetails, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when async provisioning isn't allowed but the service requested requires it", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Provision("something", cloudSqlProvisionDetails, false)
				Expect(err).To(HaveOccurred())
			})
		})

	})

	// TODO(cbriant): wrap all the positive things in expect no errors
	Describe("deprovision", func() {
		Context("when the bigquery service id is provided", func() {
			It("should call bigquery deprovisioning", func() {
				bqId := serviceNameToId[brokers.BigqueryName]
				gcpBroker.Provision("something", bqProvisionDetails, true)
				gcpBroker.Deprovision("something", models.DeprovisionDetails{
					ServiceID: bqId,
				}, true)
				Expect(gcpBroker.ServiceBrokerMap[bqId].(*modelsfakes.FakeServiceBrokerHelper).DeprovisionCallCount()).To(Equal(1))
			})
		})

		Context("when the service doesn't exist", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Deprovision("something", models.DeprovisionDetails{
					ServiceID: serviceNameToId[brokers.BigqueryName],
				}, true)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when async provisioning isn't allowed but the service requested requires it", func() {
			It("should return an error", func() {
				_, err := gcpBroker.Deprovision("something", models.DeprovisionDetails{
					ServiceID: serviceNameToId[brokers.CloudsqlName],
				}, false)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	AfterEach(func() {
		os.Remove(brokers.AppCredsFileName)
		os.Remove("test.db")
	})
})