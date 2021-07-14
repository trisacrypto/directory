package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/trisacrypto/directory/pkg/gds/config"
	"github.com/trisacrypto/directory/pkg/gds/models/v1"
	"github.com/trisacrypto/directory/pkg/gds/store"

	api "github.com/trisacrypto/trisa/pkg/trisa/gds/api/v1beta1"
	pb "github.com/trisacrypto/trisa/pkg/trisa/gds/models/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// open database
	fmt.Println("hi", os.Getenv("GDS_DATABASE_URL"))
	db, err := store.Open(config.DatabaseConfig{URL: "leveldb:///fixtures/db"}) // TODO replace
	if err != nil {
		panic(err)
	}

	run(db, 6*time.Hour)
}

type healthCheckJob struct {
	vasp *pb.VASP
	hc   *models.HealthCheckExtra
}

func run(db store.Store, duration time.Duration) {
	ctx := context.TODO()
	now := time.Now()

	all := make(chan *pb.VASP)
	check := make(chan *healthCheckJob)

	// retrieve the health info for each vasp and determine if it needs to be checked
	go func() {
		for vasp := range all {
			fmt.Println("vasp", vasp)
			healthCheck, err := models.GetHealthCheckInfo(vasp)
			if err != nil {
				log.Warn().Err(err).Str("health_check", fmt.Sprintf("could not retrieve info for vasp id %s", vasp.Id))
			} else if healthCheck.DelayCheck() {
				log.Info().Err(err).Str("health_check", fmt.Sprintf("delay for vasp id %s", vasp.Id))
				continue
			}
			check <- &healthCheckJob{
				vasp: vasp,
				hc:   healthCheck,
			}
		}
	}()

	// call the vasp endpoint and save the results
	go func() {
		for v := range check {
			client, err := initClient(ctx, v.vasp.TrisaEndpoint)
			if err != nil {
				log.Warn().Err(err).Str("health_check", fmt.Sprintf("could not init client for vasp id %s", v.vasp.Id))
				continue
			}
			var state api.ServiceState
			if err := client.Invoke(ctx, "/Status", v.hc, &state); err != nil {
				log.Warn().Err(err).Str("health_check", fmt.Sprintf("could not retrieve status for vasp id %s", v.vasp.Id))
				continue
			}
			attempts := int32(v.hc.Attempts + 1)
			if state.Status == api.ServiceState_Status(pb.ServiceState_HEALTHY) {
				attempts = 0
			}
			if err := db.UpdateStatus(v.vasp.Id, int32(state.Status)); err != nil {
				log.Warn().Err(err).Str("health_check", fmt.Sprintf("could not update status for vasp id %s", v.vasp.Id))
				continue
			}
			if err := models.SetHealthCheckInfo(v.vasp, models.HealthCheckExtra{
				CheckAfter:  state.NotBefore,
				CheckBefore: state.NotAfter,
				Attempts:    attempts,
				LastChecked: now.Format(time.RFC3339),
			}); err != nil {
				log.Warn().Err(err).Str("health_check", fmt.Sprintf("could not save extra data for vasp id %s", v.vasp.Id))
				continue
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

loop:
	for {
		fmt.Println("tik")
		// retrieve all the vasps that are verified
		verificationStatus := pb.VerificationState_VERIFIED
		if err := db.RetrieveAll(&models.RetrieveAllOpts{
			VerificationStatus:  &verificationStatus,
			TrisaEndpointExists: true,
		}, all); err != nil {
			log.Warn().Err(err).Str("health_check", "could not retrieve vasps")
		}

		select {
		case <-ticker.C:
			continue
		case <-interrupt:
			break loop
		}
	}
}

func initClient(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	config := &tls.Config{}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(config))}
	return grpc.Dial(endpoint, opts...)
}
