package util

import (
	"reflect"
	"testing"
)

func TestCSVToMap(t *testing.T) {
	type args struct {
		csvData [][]string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "Positive Test",
			args: args{
				csvData: [][]string{
					{"Title1", "Title2"},
					{"A1", "A2"},
					{"B1", "B2"},
				},
			},
			want: map[string][]string{
				"Title1": {"A1", "A2"},
				"Title2": {"B1", "B2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CSVSlicesToMap(tt.args.csvData)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSVToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMaxLengthOfSlices(t *testing.T) {
	type args struct {
		item map[string][]string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Default",
			args: args{
				item: map[string][]string{
					"Title1": {"A1", "A2"},
					"Title2": {"B1", "B2"},
				},
			},
			want: 2,
		}, {
			name: "Irregular Length",
			args: args{
				item: map[string][]string{
					"Title1": {"A1"},
					"Title2": {"B1", "B2"},
				},
			},
			want: 2,
		}, {
			name: "Empty",
			args: args{
				item: map[string][]string{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMaxLengthOfSlices(tt.args.item); got != tt.want {
				t.Errorf("GetMaxLengthOfSlices() = %v, want %v", got, tt.want)
			}
		})
	}
}
