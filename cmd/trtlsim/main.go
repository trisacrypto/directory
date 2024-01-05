package main

import (
	"context"
	crand "crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	wr "github.com/mroth/weightedrand"
	"github.com/trisacrypto/directory/pkg/trtl/jitter"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"github.com/trisacrypto/trisa/pkg/trust"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	gprcInsecure "google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	interval  = 10 * time.Second                                       // the ticker interval, default to 5 seconds
	sigma     = 200 * time.Millisecond                                 // the amount of jitter, default to 100 ms
	accesses  = 1000                                                   // desired accesses per interval, default to 15
	regions   = 1                                                      // number of regions simultaneously using the accessor
	endpoint  = "localhost:4436"                                       // the endpoint of the running trtl server
	insecure  = true                                                   // connect without mTLS
	certPath  = "fixtures/certs/mtls.client.dev/client.pem"            // path to file on disk with certificate key-pair if using mTLS
	poolPath  = "fixtures/certs/mtls.client.dev/certificate.chain.pem" // path to file on disk with trust pool if using mTLS
	keyspace  = 6000                                                   // the number of keys the simulator operates on
	chunkSize = 512                                                    // size of each write
	strategy  = "weighted_random"                                      // simulation strategy
)

// dummy namespaces for simulation
var namespaces = []string{"catchfireBarons", "falselightCutters", "fullCrowns", ""}

// probabilities for reads/writes/deletes
var probabilities = map[string]uint{
	"read":   80,
	"write":  18,
	"delete": 2,
}

func main() {

	// assumes trtl is already being served (e.g. from the trtl cli)
	// sim needs the endpoint (e.g. localhost:port) + certs (just stubs for now)
	// run trtl in insecure mode
	sim := New(endpoint, insecure)
	simClient, err := sim.connect()
	if err != nil {
		log.Fatal(err)
	}

	// run the accessor, e.g. weighted random or trisa model
	sim.Accessor.Run(simClient)
}

type Simulator struct {
	Endpoint string   `yaml:"endpoint"`            // the replica endpoint to connect to
	Insecure bool     `yaml:"insecure,omitempty"`  // do not connect with TLS
	CertPath string   `yaml:"cert_path,omitempty"` // path to certificate key pair for client side mTLS
	PoolPath string   `yaml:"pool_path,omitempty"` // path to certificate trust chain for client side mTLS
	Accessor Accessor `yaml:"-"`                   // accessor to create accesses to the database
}

type Accessor interface {
	Run(client pb.TrtlClient)
}

func New(endpoint string, insecure bool) *Simulator {
	sim := &Simulator{
		Endpoint: endpoint,
		Insecure: insecure,
		CertPath: certPath,
		PoolPath: poolPath,
	}

	// Configure from the environment if envvars are supplied, otherwise use the defaults
	if endpoint := os.Getenv("TRTLSIM_ENDPOINT"); endpoint != "" {
		sim.Endpoint = endpoint
	}

	if insecure := os.Getenv("TRTLSIM_INSECURE"); insecure != "" {
		insecure = strings.ToLower(insecure)
		if insecure == "f" || insecure == "false" || insecure == "0" {
			sim.Insecure = false
		}
	}

	if certPath := os.Getenv("TRTLSIM_CERT_PATH"); certPath != "" {
		sim.CertPath = certPath
	}

	if poolPath := os.Getenv("TRTLSIM_POOL_PATH"); poolPath != "" {
		sim.PoolPath = poolPath
	}

	// Create simulator
	simulationStrategy := strings.TrimSpace(strings.ToLower(os.Getenv("TRTLSIM_STRATEGY")))
	if simulationStrategy == "" {
		simulationStrategy = strategy
	}

	switch simulationStrategy {
	case "weighted_random":
		sim.Accessor = NewWeightedRandom()
	case "trisa_model":
		sim.Accessor = NewTRISAModel()
	default:
		log.Fatal(fmt.Errorf("unknown simulation strategy %q use weighted_random or trisa_model", simulationStrategy))
	}

	return sim
}

// Connect to the trtl server and return a gRPC client
func (s *Simulator) connect() (_ pb.TrtlClient, err error) {
	var opts []grpc.DialOption
	if s.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(gprcInsecure.NewCredentials()))
	} else {
		var sz *trust.Serializer
		if sz, err = trust.NewSerializer(false); err != nil {
			return nil, err
		}

		var pool trust.ProviderPool
		if pool, err = sz.ReadPoolFile(s.PoolPath); err != nil {
			return nil, err
		}

		var provider *trust.Provider
		if provider, err = sz.ReadFile(s.CertPath); err != nil {
			return nil, err
		}

		var cert tls.Certificate
		if cert, err = provider.GetKeyPair(); err != nil {
			return nil, err
		}

		var certPool *x509.CertPool
		if certPool, err = pool.GetCertPool(false); err != nil {
			return nil, err
		}
		var u *url.URL
		if u, err = url.Parse(s.Endpoint); err != nil {
			return nil, err
		}
		conf := &tls.Config{
			ServerName:   u.Host,
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(conf)))
	}

	// Connect the replica client
	var cc *grpc.ClientConn
	if cc, err = grpc.Dial(s.Endpoint, opts...); err != nil {
		return nil, err
	}
	log.Printf("connected to trtl server at %s\n", s.Endpoint)
	return pb.NewTrtlClient(cc), nil
}

//===========================================================================
// Weighted Probability Accessor
//===========================================================================

type WeightedRandom struct {
	Selector *wr.Chooser `yaml:"chooser,omitempty"` // random selection helper
}

func NewWeightedRandom() Accessor {
	accessor := &WeightedRandom{}
	choices := make([]wr.Choice, len(namespaces))
	for i, w := range probabilities {
		choices = append(choices, wr.Choice{Item: i, Weight: w})
	}

	var err error
	if accessor.Selector, err = wr.NewChooser(choices...); err != nil {
		log.Fatal("error in chooser creation")
	}
	return accessor
}

func (wr *WeightedRandom) Run(client pb.TrtlClient) {
	ticker := jitter.New(interval, sigma)

	// start on the first tick
	for ; true; <-ticker.C {

		// multiple regions access the data store concurrently; 1 routine per region
		var wg sync.WaitGroup
		for r := 1; r <= regions; r++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				wr.accessor(client)
			}()
		}
	}
}

func (wr *WeightedRandom) accessor(client pb.TrtlClient) {
	var gets int
	var puts int
	var dels int

	// Make as many accesses as directed by the global variable/config
	for a := 1; a <= accesses; a++ {

		// randomly select a namespace from namespaces
		ns := namespaces[rand.Intn(len(namespaces))]

		// randomly select a key from keyspace and cast to bytes
		key := []byte(strconv.Itoa(rand.Intn(keyspace)))

		// select from read, write, delete using probabilities
		switch wr.Selector.Pick().(string) {
		case "read":
			// execute Get
			req := &pb.GetRequest{
				Key:       key,
				Namespace: ns,
				Options: &pb.Options{
					ReturnMeta: false,
				},
			}
			if _, err := client.Get(context.TODO(), req); err != nil {
				// Skip if the object doesn't exist
				if serr, ok := status.FromError(err); ok && serr.Code() == codes.NotFound {
					continue
				}
				log.Fatal(fmt.Errorf("could not read from database: %v", err))
			}
			gets++
		case "write":
			// create random data
			val := make([]byte, chunkSize)

			// execute Put
			req := &pb.PutRequest{
				Key:       key,
				Value:     val,
				Namespace: ns,
				Options: &pb.Options{
					ReturnMeta: false,
				},
			}
			if _, err := client.Put(context.TODO(), req); err != nil {
				fmt.Print(err)
				log.Fatal("could not write to database")
			}
			puts++
		case "delete":
			// First check that the object exists
			req := &pb.GetRequest{
				Key:       key,
				Namespace: ns,
				Options: &pb.Options{
					ReturnMeta: false,
				},
			}
			// Object doesn't exist so skip this one
			if _, err := client.Get(context.TODO(), req); err != nil {
				if serr, ok := status.FromError(err); ok && serr.Code() == codes.NotFound {
					continue
				}
			} else {
				// Object does exist so execute Delete
				req := &pb.DeleteRequest{
					Key:       key,
					Namespace: ns,
					Options: &pb.Options{
						ReturnMeta: false,
					},
				}
				if _, err := client.Delete(context.TODO(), req); err != nil {
					if serr, ok := status.FromError(err); ok && serr.Code() == codes.NotFound {
						continue
					}
					log.Fatal(fmt.Errorf("could not delete from database: %v", err))
				}
			}
			dels++
		default:
			log.Fatal(errors.New("unknown database operation"))
		}
	}

	log.Printf("%d gets %d puts %d dels", gets, puts, dels)
}

//===========================================================================
// TRISA Model
//===========================================================================

const (
	readers          = 30
	registerInterval = 15 * time.Minute
	registerSigma    = 30 * time.Second
	reissueAge       = 30 * time.Minute
	reissueInterval  = 2 * time.Minute
	reissueSigma     = 15 * time.Second
	nsVASPs          = "vasps"
	nsCertReqs       = "certreqs"
	nsUsers          = "users"
)

// The TRISA model attempts to model a production workload on the TRISA network. It runs
// multiple go routines that constantly read from the database, performing gets and
// iters at an aggregate rate of 10-30 reads per second. It also runs a write routine
// that performs a sequence of writes: first a create (registration), then multiple
// updates (e.g. for email verification, admin changes, reviewing, approval, rejection,
// etc.). It then starts the cycle over again. Finally we run two "updating" routines
// that will update already written records at a fairly constant rate based on their
// created timestamp to simulate certificate re-issuance and user profile changes.
type TRISAModel struct {
	sync.RWMutex
	client     pb.TrtlClient        // the trtl client to make accesses on
	namespaces []string             // a list of the namespaces to access
	records    map[string][]string  // maps namespace to a list of record UUIDs for random selection
	createdAt  map[string]time.Time // maps VASP UUIDs to their created timestamp and maintains the set
}

func NewTRISAModel() Accessor {
	return &TRISAModel{
		namespaces: []string{nsVASPs, nsCertReqs, nsUsers},
		records:    make(map[string][]string),
		createdAt:  make(map[string]time.Time),
	}
}

func (t *TRISAModel) Run(client pb.TrtlClient) {
	t.client = client
	var wg sync.WaitGroup
	log.Println("running TRISA model")

	// TODO: load keys from database periodically since other simulators are probably
	// running and adding keys from other regions - if we ignore them, we're just reading
	// data created by this simulator. To do this though, we need to track which keys
	// we're managing vs which keys the other servers are managing.

	// Launch read go routines, each go routine performs a read once per second with a
	// jitter, so the number of go routines determines the approximate read rate
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go t.reader(&wg)
	}

	// Launch the VASP registration routine
	wg.Add(1)
	go t.registrations(&wg)

	// Launch the certificate reissuance routine
	wg.Add(1)
	go t.reissuer(&wg)

	// Launch the user change password routine
	wg.Add(1)
	go t.userProfiles(&wg)

	wg.Wait()
}

// reader conducts a Get or an Iter every second or so with jitter
func (t *TRISAModel) reader(wg *sync.WaitGroup) {
	defer wg.Done()
	nErrors := 0
	ticker := jitter.New(1*time.Second, 300*time.Millisecond)

	for ; true; <-ticker.C {
		roulette := rand.Float64()
		namespace := t.namespaces[rand.Intn(3)]

		switch {
		case roulette < 0.1:
			// Perform a cursor iteration over all items in a random namespace with probability 0.1
			// NOTE: as the database gets bigger, these reads are going to get longer, to the point
			// where the simulator will eventually crash because the read takes longer than the timeout.
			if err := t.Cursor(namespace); err != nil {
				log.Println(err.Error())
				nErrors++
			}
		case roulette < 0.15:
			// Perform a batch iteration over all items in a random namespace with probability 0.05
			// NOTE: as the database gets bigger, these reads are going to get longer, to the point
			// where the simulator will eventually crash because the read takes longer than the timeout.
			if err := t.Iter(namespace); err != nil {
				log.Println(err.Error())
				nErrors++
			}
		default:
			// Get a random record from the namespace with probability 0.75
			// Flip a coin to determine if we fetch meta or not with the request
			key := t.pick(namespace)
			if key == "" {
				// No records have been created yet
				continue
			}

			withMeta := coinFlip()
			if err := t.Get(key, namespace, withMeta); err != nil {
				log.Println(err.Error())
				nErrors++
			}
		}

		if nErrors > 10 {
			// If 10 errors have occurred kill the program
			log.Fatal("too many errors occurred reading")
		}
	}
}

// select a random record from the namespace with uniform likelihood
// TODO: make read likelihood non-uniform, e.g. some VASPs should be more popular to lookup than others
func (t *TRISAModel) pick(namespace string) string {
	t.RLock()
	defer t.RUnlock()
	nRecords := len(t.records[namespace])
	if nRecords == 0 {
		return ""
	}
	return t.records[namespace][rand.Intn(nRecords)]
}

// simulates 1-3 registrations per jittered interval (e.g. 1-3 registrations a week)
func (t *TRISAModel) registrations(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := jitter.New(registerInterval, registerSigma)
	for ; true; <-ticker.C {
		roulette := rand.Float64()
		// every register interval kick off a new registration with probability 0.8
		if roulette < 0.8 {
			wg.Add(1)
			go t.register(wg)
		}

		// every register interval kick of a second registration with probability 0.1
		if roulette < 0.1 {
			wg.Add(1)
			go t.register(wg)
		}

		// every register interval kick off a third registration with probability 0.05
		if roulette < 0.05 {
			wg.Add(1)
			go t.register(wg)
		}
	}
}

// simulates a registration workflow with various writes
// NOTE: this is the process that would likely cause stomps as users in one region stomp
// admin edits in another region. However, only one process (e.g. this simulator) knows
// about the registration so stomps are currently impossible in this model. We either
// need to allow accesses from multiple regions in the simulator, or find a way to get
// multiple simulators to be able to write registrations simultaneously.
func (t *TRISAModel) register(wg *sync.WaitGroup) {
	defer wg.Done()
	// Write initial records: 1 VASP, 1 cert req, 1-4 users
	vasp := uuid.NewString()
	certreq := uuid.NewString()

	nUsers := rand.Intn(4) + 1
	users := make([]string, 0, nUsers)
	for i := 0; i < nUsers; i++ {
		users = append(users, uuid.NewString())
	}

	// Create VASP
	if err := t.Put(vasp, nsVASPs, rand.Intn(32768)+4096); err != nil {
		log.Fatal(fmt.Errorf("could not register vasp: %s", err))
	}
	t.upsert(vasp, nsVASPs, time.Now())

	// Create CertReq
	if err := t.Put(certreq, nsCertReqs, rand.Intn(8192)+228); err != nil {
		log.Fatal(fmt.Errorf("could not create certreq: %s", err))
	}
	t.upsert(certreq, nsCertReqs, time.Now())

	// Create Users
	for _, user := range users {
		if err := t.Put(user, nsUsers, rand.Intn(2048)+512); err != nil {
			log.Fatal(fmt.Errorf("could not create user: %s", err))
		}
		t.upsert(user, nsUsers, time.Now())
	}

	log.Printf("created vasp %s with 1 certreq and %d users", vasp, len(users))

	// Sleep and create random writes to users and VASP
	// note: this is the only writer so it is currently impossible to generate stomps
	// we'd have to have multiple regions in the simulator to simulate users from
	// different regions interacting with the record, or we'd have to have different
	// simulators write to records created by a different simulator
	nUpdates := rand.Intn(24) + 3
	for i := 0; i < nUpdates; i++ {
		time.Sleep(randSleep(1*time.Second, 48*time.Second))
		roulette := rand.Float64()
		switch {
		case roulette < 0.05:
			// Update CertReq
			if err := t.Put(certreq, nsCertReqs, rand.Intn(8192)+228); err != nil {
				log.Fatal(fmt.Errorf("could not update certreq: %s", err))
			}
			t.upsert(certreq, nsCertReqs, time.Now())
			log.Printf("updated certreq %s", certreq)
		case roulette < 0.45:
			// Update a random user
			user := users[rand.Intn(len(users))]
			if err := t.Put(user, nsUsers, rand.Intn(2048)+512); err != nil {
				log.Fatal(fmt.Errorf("could not update user: %s", err))
			}
			t.upsert(user, nsUsers, time.Now())
			log.Printf("updated user %s", user)
		default:
			if err := t.Put(vasp, nsVASPs, rand.Intn(32768)+4096); err != nil {
				log.Fatal(fmt.Errorf("could not update vasp: %s", err))
			}
			t.upsert(vasp, nsVASPs, time.Now())
			log.Printf("updated vasp %s", vasp)
		}
	}

	// Create last two updates to certreq/vasp to mimic issuance process
	if err := t.Put(vasp, nsVASPs, rand.Intn(32768)+4096); err != nil {
		log.Fatal(fmt.Errorf("could not update vasp: %s", err))
	}
	if err := t.Put(certreq, nsCertReqs, rand.Intn(8192)+228); err != nil {
		log.Fatal(fmt.Errorf("could not update certreq: %s", err))
	}
	t.upsert(vasp, nsVASPs, time.Now())
	t.upsert(certreq, nsCertReqs, time.Now())
	log.Printf("issuing certificates for vasp %s", vasp)

	// Trying to sleep long enough to let anti-entropy happen; assuming anti-entropy is 1 minute
	time.Sleep(90 * time.Second)

	// Create last two updates to certreq/vasp to mimic issuance process
	if err := t.Put(vasp, nsVASPs, rand.Intn(32768)+4096); err != nil {
		log.Fatal(fmt.Errorf("could not update vasp: %s", err))
	}
	if err := t.Put(certreq, nsCertReqs, rand.Intn(8192)+228); err != nil {
		log.Fatal(fmt.Errorf("could not update certreq: %s", err))
	}
	t.upsert(vasp, nsVASPs, time.Now())
	t.upsert(certreq, nsCertReqs, time.Now())
	log.Printf("certificates issued for vasp %s", vasp)
}

// Insert a new record into the simulation or update its timestamp
func (t *TRISAModel) upsert(key, namespace string, ts time.Time) {
	t.Lock()
	if _, ok := t.createdAt[key]; !ok {
		t.records[namespace] = append(t.records[namespace], key)
	}
	t.createdAt[key] = ts
	t.Unlock()
}

// Update VASPs that are older than the reissue age to make writes to older objects
func (t *TRISAModel) reissuer(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := jitter.New(registerInterval, registerSigma)
	for {
		// Wait for reissue interval
		<-ticker.C

		// Figure out which VASPs need certs reissued
		vasps := make([]string, 0)
		t.RLock()
		for _, vasp := range t.records[nsVASPs] {
			if time.Since(t.createdAt[vasp]) >= reissueAge {
				vasps = append(vasps, vasp)
			}
		}
		t.RUnlock()

		// Reissue certs for each VASP, simulating a request and issue process
		for _, vasp := range vasps {
			wg.Add(1)
			go func(vasp string) {
				defer wg.Done()
				// Create last two updates to certreq/vasp to mimic issuance process
				if err := t.Put(vasp, nsVASPs, rand.Intn(32768)+4096); err != nil {
					log.Fatal(fmt.Errorf("could not update vasp: %s", err))
				}
				certreq := uuid.NewString()
				if err := t.Put(certreq, nsCertReqs, rand.Intn(8192)+228); err != nil {
					log.Fatal(fmt.Errorf("could not update certreq: %s", err))
				}
				t.upsert(vasp, nsVASPs, time.Now())
				t.upsert(certreq, nsCertReqs, time.Now())
				log.Printf("reissuing certificates for vasp %s", vasp)

				// Trying to sleep long enough to let anti-entropy happen; assuming anti-entropy is 1 minute
				time.Sleep(90 * time.Second)

				// Create last two updates to certreq/vasp to mimic issuance process
				if err := t.Put(vasp, nsVASPs, rand.Intn(32768)+4096); err != nil {
					log.Fatal(fmt.Errorf("could not update vasp: %s", err))
				}
				if err := t.Put(certreq, nsCertReqs, rand.Intn(8192)+228); err != nil {
					log.Fatal(fmt.Errorf("could not update certreq: %s", err))
				}

				t.upsert(vasp, nsVASPs, time.Now())
				t.upsert(certreq, nsCertReqs, time.Now())
				log.Printf("certificates reissued for vasp %s", vasp)
			}(vasp)
		}
	}
}

// User updater randomly updates users' profiles
// This writer may generate stomps, but it's unlikely
func (t *TRISAModel) userProfiles(wg *sync.WaitGroup) {
	ticker := jitter.New(reissueInterval, reissueSigma)
	for {
		<-ticker.C
		// Update user with 0.2 probability every tick
		if roulette := rand.Float64(); roulette < 0.2 {
			// Randomly select a user to update
			user := t.pick(nsUsers)
			if user == "" {
				// No users in the database to update
				continue
			}

			if err := t.Put(user, nsUsers, rand.Intn(2048)+512); err != nil {
				log.Fatal(fmt.Errorf("could not update user: %s", err))
			}
			t.upsert(user, nsUsers, time.Now())
			log.Printf("updated user profile %s", user)
		}
	}
}

func (t *TRISAModel) Get(key, namespace string, withMeta bool) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Execute the Get request
	req := &pb.GetRequest{
		Key:       []byte(key),
		Namespace: namespace,
		Options: &pb.Options{
			ReturnMeta: withMeta,
		},
	}
	if _, err := t.client.Get(ctx, req); err != nil {
		// Ignore not found errors
		if serr, ok := status.FromError(err); ok && serr.Code() != codes.NotFound {
			return fmt.Errorf("could not read from database: %v", err)
		}
	}
	return nil
}

func (t *TRISAModel) Cursor(namespace string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.CursorRequest{
		Namespace: namespace,
		Options: &pb.Options{
			ReturnMeta: coinFlip(),
		},
	}

	var stream pb.Trtl_CursorClient
	if stream, err = t.client.Cursor(ctx, req); err != nil {
		return fmt.Errorf("could not open cursor stream: %v", err)
	}

	for {
		if _, err = stream.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error receiving from cursor: %v", err)
		}
	}
}

func (t *TRISAModel) Iter(namespace string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.IterRequest{
		Namespace: namespace,
		Options: &pb.Options{
			ReturnMeta: coinFlip(),
			PageSize:   100,
		},
	}

	var rep *pb.IterReply
	if rep, err = t.client.Iter(ctx, req); err != nil {
		return fmt.Errorf("could not get first page of requests: %v", err)
	}

	for rep.NextPageToken != "" {
		req.Options.PageToken = rep.NextPageToken
		if rep, err = t.client.Iter(ctx, req); err != nil {
			return fmt.Errorf("could not get next page of requests: %v", err)
		}
	}

	return nil
}

func (t *TRISAModel) Put(key, namespace string, nbytes int) (err error) {
	// Create a random value of the specified length
	value := make([]byte, nbytes)
	if _, err = crand.Read(value); err != nil {
		return fmt.Errorf("could not create random value: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.PutRequest{
		Key:       []byte(key),
		Value:     value,
		Namespace: namespace,
		Options: &pb.Options{
			ReturnMeta: coinFlip(),
		},
	}

	if _, err = t.client.Put(ctx, req); err != nil {
		return fmt.Errorf("could not put value: %v", err)
	}
	return nil
}

func coinFlip() bool {
	return rand.Float32() < 0.5
}

func randSleep(min, max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max))) + min
}
