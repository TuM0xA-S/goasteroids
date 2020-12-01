package vobj

import "testing"

func TestSectionIntersects(t *testing.T) {
	type args struct {
		v0 Vec2
		v1 Vec2
		u0 Vec2
		u1 Vec2
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{Vec2{0, 0}, Vec2{3, 3}, Vec2{3, 0}, Vec2{0, 3}}, true},
		{"", args{Vec2{2, 3}, Vec2{2, 5}, Vec2{1, 4}, Vec2{5, 4}}, true},
		{"", args{Vec2{3, 6}, Vec2{4, 2}, Vec2{1, 4}, Vec2{5, 4}}, true},
		{"", args{Vec2{3, 6}, Vec2{4, 2}, Vec2{3, 3}, Vec2{5, 3}}, true},
		{"", args{Vec2{3, 6}, Vec2{4, 2}, Vec2{3, 2}, Vec2{4, 6}}, true},
		{"", args{Vec2{1, 1}, Vec2{6, 3}, Vec2{6, 2}, Vec2{1, 2}}, true},
		{"", args{Vec2{0, 5}, Vec2{10, 1.5}, Vec2{2, 0}, Vec2{5.2, 7}}, true},
		{"", args{Vec2{0, 5}, Vec2{10, 1.5}, Vec2{7.5, 0}, Vec2{5.2, 7}}, true},
		{"", args{Vec2{2, 5}, Vec2{9, 2}, Vec2{11, 3}, Vec2{6, 2}}, true},
		{"", args{Vec2{9, 2}, Vec2{9, 5}, Vec2{6, 2}, Vec2{11, 3},}, true},
		{"", args{Vec2{6, 2}, Vec2{11, 3}, Vec2{9, 2}, Vec2{9, 5},}, true},
		{"", args{Vec2{7, 5}, Vec2{11, 3}, Vec2{9, 2}, Vec2{9, 5},}, true},
		{"", args{Vec2{5, 6}, Vec2{6, 10}, Vec2{4, 7}, Vec2{9, 7},}, true},
		{"", args{Vec2{5, 6}, Vec2{6, 10}, Vec2{3, 9}, Vec2{10, 8},}, true},
		{"", args{Vec2{5, 6}, Vec2{6, 10}, Vec2{10, 8}, Vec2{3, 9},}, true},
		{"", args{Vec2{10, 8}, Vec2{3, 9}, Vec2{5, 6}, Vec2{6, 10},}, true},
		{"", args{Vec2{6, 10}, Vec2{5, 6}, Vec2{3, 9}, Vec2{10, 8},}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SectionIntersects(tt.args.v0, tt.args.v1, tt.args.u0, tt.args.u1); got != tt.want {
				t.Errorf("SectionIntersects() = %v, want %v", got, tt.want)
			}
		})
	}
}
