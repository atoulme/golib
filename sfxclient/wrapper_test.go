package sfxclient

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"golang.org/x/net/context"
	"github.com/signalfx/golib/datapoint"
	"sync/atomic"
)

func TestNewMultiCollector(t *testing.T) {
	Convey("a NewMultiCollector", t, func() {
		c1 := GoMetricsSource
		c2 := GoMetricsSource
		Convey("should return itself for one item", func() {
			So(NewMultiCollector(c1), ShouldEqual, c1)
		})
		Convey("should wrap multiple items", func() {
			c3 := NewMultiCollector(c1, c2)
			So(len(c3.Datapoints()), ShouldEqual, 2*len(c1.Datapoints()))
		})
	})
}

func TestWithDimensions(t *testing.T) {
	Convey("a WithDimensions should work", t, func() {
		c1 := GoMetricsSource
		c2 := WithDimensions{
			Collector:  c1,
			Dimensions: map[string]string{"name": "jack"},
		}
		dp0 := c2.Datapoints()[0]
		So(dp0.Dimensions["name"], ShouldEqual, "jack")

		c2.Dimensions = nil
		dp0 = c2.Datapoints()[0]
		So(dp0.Dimensions["name"], ShouldNotEqual, "jack")
	})
}

func ExampleNewMultiCollector() {
	var a Collector
	var b Collector
	c := NewMultiCollector(a, b)
	c.Datapoints()
}

func ExampleCumulativeP() {
	client := NewHTTPDatapointSink()
	ctx := context.Background()
	var countThing int64
	go func() {
		atomic.AddInt64(&countThing, 1)
	}()
	client.AddDatapoints(ctx, []*datapoint.Datapoint {
		CumulativeP("server.request_count", nil, &countThing),
	})
}

func ExampleWithDimensions() {
	sched := NewScheduler()
	sched.AddCallback(&WithDimensions{
		Collector: GoMetricsSource,
		Dimensions: map[string]string {
			"extra": "dimension",
		},
	})
}
