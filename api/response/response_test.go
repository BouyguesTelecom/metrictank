package response

import (
	"math"

	"github.com/grafana/metrictank/api/models"
	"github.com/grafana/metrictank/schema"
)

type series struct {
	in  []models.Series
	out string
}

func testSeries() []series {
	cases := []series{
		{
			in:  []models.Series{},
			out: `[]`,
		},
		{
			in: []models.Series{
				{
					Target:     "a",
					Datapoints: []schema.Point{},
					Interval:   60,
				},
			},
			out: `[{"target":"a","datapoints":[]}]`,
		},
		{
			in: []models.Series{
				{
					Target: "a",
					Datapoints: []schema.Point{
						{123, 60},
						{10000, 120},
						{0, 180},
						{1, 240},
					},
					Interval: 60,
				},
			},
			out: `[{"target":"a","datapoints":[[123,60],[10000,120],[0,180],[1,240]]}]`,
		},
		{
			in: []models.Series{
				{
					Target: "a",
					Datapoints: []schema.Point{
						{123, 60},
						{10000, 120},
						{0, 180},
						{1, 240},
					},
					Interval: 60,
				},
				{
					Target: "foo(bar)",
					Datapoints: []schema.Point{
						{123.456, 10},
						{math.NaN(), 20},
						{124.10, 30},
						{125.0, 40},
						{126.0, 50},
					},
					Interval: 10,
				},
			},
			out: `[{"target":"a","datapoints":[[123,60],[10000,120],[0,180],[1,240]]},{"target":"foo(bar)","datapoints":[[123.456,10],[null,20],[124.1,30],[125,40],[126,50]]}]`,
		},
	}
	return cases
}
