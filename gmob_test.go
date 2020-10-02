package gmob

import (
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestBuild(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		request interface{}
		wantErr bool
		want    bson.M
	}{
		{
			name: "struct",
			request: struct {
				ID      string                 `bson:"id"`
				Name    string                 `bson:"name"`
				Active  bool                   `bson:"active"`
				Publish bool                   `bson:"publish"`
				Count   int                    `bson:"count"`
				Amount  float64                `bson:"amount"`
				List    [2]string              `bson:"list"`
				Meta    map[string]interface{} `bson:"meta"`
			}{
				ID:      "abc",
				Name:    "ABC",
				Active:  false,
				Publish: true,
				Count:   10,
				Amount:  2.5,
				List:    [2]string{"l1", "l2"},
				Meta:    map[string]interface{}{"type": "test"},
			},
			want: bson.M{
				"id":      "abc",
				"name":    "ABC",
				"publish": true,
				"count":   10,
				"amount":  2.5,
				"list":    bson.A{"l1", "l2"},
				"meta": bson.M{
					"type": "test",
				},
			},
		},
		{
			name: "struct inside struct",
			request: struct {
				ID        string     `bson:"id"`
				CreatedAt time.Time  `bson:"createdAt"`
				DeletedAt *time.Time `bson:"deletedAt"`
				Notes     []struct {
					Content   string    `bson:"content"`
					CreatedAt time.Time `bson:"createdAt"`
				} `bson:"notes"`
			}{
				ID:        "abc",
				CreatedAt: now,
				DeletedAt: &now,
				Notes: []struct {
					Content   string    `bson:"content"`
					CreatedAt time.Time `bson:"createdAt"`
				}{
					{
						Content:   "test content",
						CreatedAt: now,
					},
				},
			},
			want: bson.M{
				"id":        "abc",
				"createdAt": now,
				"deletedAt": now,
				"notes": bson.A{
					bson.M{
						"content":   "test content",
						"createdAt": now,
					},
				},
			},
		},
		{
			name: "map",
			request: map[string]interface{}{
				"active": false,
				"type":   "test",
				"list":   []string{"l1", "l2"},
				"meta": map[string]interface{}{
					"name": "",
					"time": now,
				},
			},
			want: bson.M{
				"active": false,
				"type":   "test",
				"list": bson.A{
					"l1", "l2",
				},
				"meta": bson.M{
					"name": "",
					"time": now,
				},
			},
		},
		{
			name: "map contains struct",
			request: map[string]interface{}{
				"active": false,
				"meta": struct {
					Status string `bson:"status"`
				}{
					Status: "test",
				},
			},
			want: bson.M{
				"active": false,
				"meta": bson.M{
					"status": "test",
				},
			},
		},
		{
			name:    "not map/struct input with error",
			request: 100,
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Build(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkBuild(b *testing.B) {
	input := struct {
		ID   string                 `bson:"id"`
		List [2]string              `bson:"list"`
		Meta map[string]interface{} `bson:"meta"`
	}{
		ID:   "abc",
		List: [2]string{"l1", "l2"},
		Meta: map[string]interface{}{"type": "test"},
	}

	for i := 0; i < b.N; i++ {
		Build(input)
	}
}
