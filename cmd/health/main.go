package main

import (
	"context"
	"crypto/tls"
	"fmt"
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

type healthCheckJob struct {
	vasp *pb.VASP
	hc   *models.HealthCheck
}

func main() {
	// open database
	db, err := store.Open(config.DatabaseConfig{URL: "leveldb:///fixtures/db"}) // TODO replace
	if err != nil {
		panic(err)
	}
	ctx := context.TODO()
	now := time.Now()

	all := make(chan *pb.VASP)
	check := make(chan *healthCheckJob)
	var cntCheckedVASPS int
	var cntCheckedErrs int

	// retrieve the health info for each vasp and determine if it needs to be checked
	go func() {
		for vasp := range all {
			if vasp == nil {
				close(all)
				check <- nil
				return // done
			}
			healthCheck, err := models.GetHealthCheckInfo(vasp)
			if err != nil {
				panic(err)
			} else if healthCheck.DelayCheck() {
				continue
			}
			// TODO if we want to check the last check time, we would do that here but
			// it is probably not necessary yet
			check <- &healthCheckJob{
				vasp: vasp,
				hc:   healthCheck,
			}
			cntCheckedVASPS++
		}
	}()

	// call the vasp endpoint and save the results
	go func() {
		for v := range check {
			if v == nil {
				close(check)
				return // done
			}
			client, err := initClient(ctx, v.vasp.TrisaEndpoint)
			if err != nil {
				cntCheckedErrs++
				log.Warn().Err(err).Str("health_check", "could not init client")
				continue
			}
			var state api.ServiceState
			if err := client.Invoke(ctx, "/Status", v.hc, &state); err != nil {
				cntCheckedErrs++
				log.Warn().Err(err).Str("health_check", "could not retrieve status")
				continue
			}
			attempts := int32(v.hc.Attempts + 1)
			if state.Status == api.ServiceState_Status(pb.ServiceState_HEALTHY) {
				attempts = 0
			}
			if err := db.UpdateStatus(v.vasp.Id, int32(state.Status)); err != nil {
				cntCheckedErrs++
				log.Warn().Err(err).Str("health_check", "could not update status")
				continue
			}
			if err := models.SetHealthCheckInfo(v.vasp, models.HealthCheck{
				CheckAfter:  state.NotBefore,
				CheckBefore: state.NotAfter,
				Attempts:    attempts,
				LastChecked: now.Format(time.RFC3339),
			}); err != nil {
				cntCheckedErrs++
				log.Warn().Err(err).Str("health_check", "could not save extra data")
				continue
			}
		}
	}()

	// retrieve all the vasps that are verified
	verificationStatus := pb.VerificationState_VERIFIED
	if err := db.RetrieveAll(&models.RetrieveAllOpts{
		VerificationStatus:  &verificationStatus,
		TrisaEndpointExists: true,
	}, all); err != nil {
		log.Warn().Err(err).Str("health_check", "could not retrieve vasps")
		// return nil, status.Error(codes.NotFound, "could not find VASP by ID")
		panic(err)
	}

	fmt.Println("checked ", cntCheckedVASPS, " vasps with ", cntCheckedErrs, " errors")
}

func initClient(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	config := &tls.Config{}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(config))}
	return grpc.Dial(endpoint, opts...)
}
