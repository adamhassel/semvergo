package semver

import "testing"

func Test_pre_Max(t *testing.T) {
	type args struct {
		q pre
	}
	tests := []struct {
		name string
		p    pre
		args args
		want string
	}{
		{
			name: "simple, numeric",
			p:    "1",
			args: args{
				q: "2",
			},
			want: "2",
		},
		{
			name: "simple, letters",
			p:    "a",
			args: args{
				q: "b",
			},
			want: "b",
		},
		{
			name: "string vs numbers, strings win!",
			p:    "a",
			args: args{
				q: "1",
			},
			want: "a",
		},
		{
			name: "multiple numbers",
			p:    "1.1",
			args: args{
				q: "1.2",
			},
			want: "1.2",
		},
		{
			name: "multiple numbers, longer",
			p:    "1.2.2.3.4.5",
			args: args{
				q: "1.2.2.3.4.6",
			},
			want: "1.2.2.3.4.6",
		},
		{
			name: "strings and numbers combined, in same places",
			p:    "1.2.a.b.1.2",
			args: args{
				q: "1.2.a.b.1.3",
			},
			want: "1.2.a.b.1.3",
		},
		{
			name: "strings and numbers combined, not in same places, strings win",
			p:    "1.2.3.a.b.1.2",
			args: args{
				q: "1.2.a.3.b.1.3",
			},
			want: "1.2.a.3.b.1.3",
		},
		{
			name: "different length, but equal until final, in which case fewer components win",
			p:    "1.2.3.a.b.1.2",
			args: args{
				q: "1.2.3.a.b.1.2.moar",
			},
			want: "1.2.3.a.b.1.2.moar",
		},
		{
			name: "different length, not equal. Do regular comparison, above case is a tie breaker",
			p:    "1.2.4.a.b.1.2",
			args: args{
				q: "1.2.3.a.b.1.2.moar",
			},
			want: "1.2.4.a.b.1.2",
		},
		{
			name: "nothing is always more prevalent than something",
			p:    "",
			args: args{
				q: "blah",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Max(tt.args.q); got != tt.want {
				t.Errorf("MaxSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
