package activity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rotationalio/go-ensign"
	api "github.com/rotationalio/go-ensign/api/v1beta1"
	mimetype "github.com/rotationalio/go-ensign/mimetype/v1beta1"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	NetworkActivityMimeType = mimetype.ApplicationMsgPack
	NetworkActivityWindow   = 5 * time.Minute
)

var NetworkActivityEventType = api.Type{
	Name:         "NetworkActivity",
	MajorVersion: 1,
}

// NetworkActivity represents a time-aggregated collection of events on the GDS or rVASP
// that are a proxy for TRISA transfer usage; e.g. GDS lookups or rVSAP transfers. This
// event is published as msgpack data to an Ensign topic so that the BFF can render a
// timeseries of network activity.
type NetworkActivity struct {
	Network      Network                     `msgpack:"network"`       // The network refers to TestNet or MainNet and possibly also rVASP
	Activity     ActivityCount               `msgpack:"activity"`      // A count of activity events by name
	VASPActivity map[uuid.UUID]ActivityCount `msgpack:"vasp_activity"` // Per-vasp activity count should be less than or equal to activity counts
	Timestamp    time.Time                   `msgpack:"timestamp"`     // The start time of the aggregation window
	Window       time.Duration               `msgpack:"window"`        // The window size of the aggregation window
}

func New(network Network, ts time.Time) *NetworkActivity {
	if ts.IsZero() {
		ts = time.Now()
	}

	return &NetworkActivity{
		Network:      network,
		Activity:     make(ActivityCount),
		VASPActivity: make(map[uuid.UUID]ActivityCount),
		Timestamp:    ts,
		Window:       NetworkActivityWindow,
	}
}

func Parse(event *ensign.Event) (_ *NetworkActivity, err error) {
	if event.Mimetype != NetworkActivityMimeType {
		return nil, errors.New("unhandled mimetype")
	}

	if !event.Type.Equals(&NetworkActivityEventType) {
		return nil, errors.New("unhandled event type")
	}

	activity := &NetworkActivity{}
	if err = msgpack.Unmarshal(event.Data, activity); err != nil {
		return nil, err
	}
	return activity, nil
}

type ActivityCount map[Activity]uint64

type Activity uint16

const (
	UnknownActivity Activity = iota
	LookupActivity
	SearchActivity
)

func (a Activity) String() string {
	switch a {
	case LookupActivity:
		return "lookup"
	case SearchActivity:
		return "search"
	default:
		return "unknown"
	}
}

type Network uint8

const (
	UnknownNetwork Network = iota
	TestNet
	MainNet
	RVASP
)

func (n Network) String() string {
	switch n {
	case TestNet:
		return "testnet"
	case MainNet:
		return "mainnet"
	case RVASP:
		return "rvasp"
	default:
		return "unknown"
	}
}

func (a *NetworkActivity) Incr(activity Activity) {
	if a.Activity == nil {
		a.Activity = make(ActivityCount)
	}
	a.Activity[activity]++
}

func (a *NetworkActivity) IncrVASP(vaspID uuid.UUID, activity Activity) {
	if a.VASPActivity == nil {
		a.VASPActivity = make(map[uuid.UUID]ActivityCount)
	}

	if _, ok := a.VASPActivity[vaspID]; !ok {
		a.VASPActivity[vaspID] = make(ActivityCount)
	}

	a.VASPActivity[vaspID][activity]++
	a.Incr(activity)
}

func (a *NetworkActivity) Event() (_ *ensign.Event, err error) {
	// TODO: what do we want to add for the event metadata?
	event := &ensign.Event{
		Metadata: make(ensign.Metadata),
		Type:     &NetworkActivityEventType,
		Mimetype: NetworkActivityMimeType,
		Created:  time.Now(),
	}

	if event.Data, err = msgpack.Marshal(a); err != nil {
		return nil, err
	}
	return event, nil
}

func (a *NetworkActivity) WindowStart() time.Time {
	return a.Timestamp
}

func (a *NetworkActivity) WindowEnd() time.Time {
	return a.Timestamp.Add(a.Window)
}
