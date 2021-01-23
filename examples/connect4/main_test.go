package main

import (
	"testing"
)

func TestGame_didTheyWin(t *testing.T) {
	type args struct {
		circle CircleType
		column int
		row    int
	}
	tests := []struct {
		name  string
		board [7][6]CircleType
		args  args
		want  bool
	}{
		{
			name:  "no win",
			board: [7][6]CircleType{},
			args: args{
				circle: CircleTypeRed,
				column: 0,
				row:    0,
			},
			want: false,
		},
		{
			name: "across columns",
			board: [7][6]CircleType{
				{CircleTypeRed},
				{CircleTypeRed},
				{CircleTypeRed},
				{CircleTypeRed},
			},
			args: args{
				circle: CircleTypeRed,
				column: 0,
				row:    0,
			},
			want: true,
		},
		{
			name: "down a column",
			board: [7][6]CircleType{
				{},
				{CircleTypeBlack, CircleTypeBlack, CircleTypeBlack, CircleTypeBlack},
			},
			args: args{
				circle: CircleTypeBlack,
				column: 1,
				row:    0,
			},
			want: true,
		},
		{
			name: "diag1",
			board: [7][6]CircleType{
				{},
				{CircleTypeBlack},
				{CircleTypeNone, CircleTypeBlack},
				{CircleTypeNone, CircleTypeNone, CircleTypeBlack},
				{CircleTypeNone, CircleTypeNone, CircleTypeNone, CircleTypeBlack},
			},
			args: args{
				circle: CircleTypeBlack,
				column: 1,
				row:    0,
			},
			want: true,
		},
		{
			name: "diag2",
			board: [7][6]CircleType{
				{},
				{CircleTypeNone, CircleTypeNone, CircleTypeNone, CircleTypeBlack},
				{CircleTypeNone, CircleTypeNone, CircleTypeBlack},
				{CircleTypeNone, CircleTypeBlack},
				{CircleTypeBlack},
			},
			args: args{
				circle: CircleTypeBlack,
				column: 1,
				row:    0,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{
				board: tt.board,
			}
			if got := g.didTheyWin(tt.args.circle, tt.args.column, tt.args.row); got != tt.want {
				t.Errorf("Game.didTheyWin() = %v, want %v", got, tt.want)
			}
		})
	}
}
