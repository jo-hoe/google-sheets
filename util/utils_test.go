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
					{"Title1", "Title2", "Title3"},
					{"A1", "B1", "C1"},
					{"A2", "B2", "C2"},
				},
			},
			want: map[string][]string{
				"Title1": {"A1", "A2"},
				"Title2": {"B1", "B2"},
				"Title3": {"C1", "C2"},
			},
		}, {
			name: "Headers only",
			args: args{
				csvData: [][]string{
					{"Title1", "Title2", "Title3"},
				},
			},
			want: map[string][]string{
				"Title1": {},
				"Title2": {},
				"Title3": {},
			},
		}, {
			name: "Empty",
			args: args{
				csvData: [][]string{
					{},
				},
			},
			want: map[string][]string{},
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

func TestValidatedKeys(t *testing.T) {
	type args struct {
		items          map[string][]string
		keysToValidate []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "All item in map",
			args: args{
				items: map[string][]string{
					"Title1": {"A1"},
					"Title2": {"B1", "B2"},
				},
				keysToValidate: []string{"Title1", "Title2"},
			},
			wantErr: false,
		}, {
			name: "One missing key",
			args: args{
				items: map[string][]string{
					"Title1": {"A1"},
				},
				keysToValidate: []string{"Title1", "Title2"},
			},
			wantErr: true,
		}, {
			name: "Validate subset of keys",
			args: args{
				items: map[string][]string{
					"Title1": {"A1"},
					"Title2": {"B1", "B2"},
				},
				keysToValidate: []string{"Title2"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatedKeys(tt.args.items, tt.args.keysToValidate...); (err != nil) != tt.wantErr {
				t.Errorf("ValidatedKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
