package semver

import (
	"reflect"
	"testing"
)

func Test_version(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		want_major uint
		want_minor uint
		want_patch uint
		wantErr    bool
	}{
		{
			name: "zeroes",
			args: args{
				s: "0.0.0",
			},
			want_major: 0,
			want_minor: 0,
			want_patch: 0,
			wantErr:    false,
		},
		{
			name: "not zeroes",
			args: args{
				s: "1.2.3",
			},
			want_major: 1,
			want_minor: 2,
			want_patch: 3,
			wantErr:    false,
		},
		{
			name: "not numbers",
			args: args{
				s: "foo.bar.baz",
			},
			want_major: 0,
			want_minor: 0,
			want_patch: 0,
			wantErr:    true,
		},
		{
			name: "incomplete",
			args: args{
				s: "2.3",
			},
			want_major: 0,
			want_minor: 0,
			want_patch: 0,
			wantErr:    true,
		},
		{
			name: "not a semver",
			args: args{
				s: "blabla",
			},
			want_major: 0,
			want_minor: 0,
			want_patch: 0,
			wantErr:    true,
		},
		{
			name: "with noise",
			args: args{
				s: "bla3.4.5blah",
			},
			want_major: 3,
			want_minor: 4,
			want_patch: 5,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := version(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("version() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want_major {
				t.Errorf("version() got = %v, want_major %v", got, tt.want_major)
			}
			if got1 != tt.want_minor {
				t.Errorf("version() got1 = %v, want_major %v", got1, tt.want_minor)
			}
			if got2 != tt.want_patch {
				t.Errorf("version() got2 = %v, want_major %v", got2, tt.want_patch)
			}
		})
	}
}

func Test_prefix(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "stnd",
			args: args{
				s: "foo-1.2.3",
			},
			want:    "foo-",
			wantErr: false,
		},
		{
			name: "noprefix",
			args: args{
				s: "1.2.3",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "not semver",
			args: args{
				s: "blabla",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "incomplete semver",
			args: args{
				s: "0.1",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "all wrong and not semver",
			args: args{
				s: "foo.1.blah.0.1",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "prefixed incomplete semver",
			args: args{
				s: "foo.1.1",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "pre- and suffixed semver",
			args: args{
				s: "foo-1.1.0-bar",
			},
			want:    "foo-",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := prefix(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("prefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("prefix() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_suffix(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				s: "1.2.3-foo",
			},
			want:    "-foo",
			wantErr: false,
		},
		{
			name: "noprefix",
			args: args{
				s: "1.2.3",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "not semver",
			args: args{
				s: "blabla",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "incomplete semver",
			args: args{
				s: "0.1",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "all wrong and not semver",
			args: args{
				s: "1.blah.0.1.foo",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "suffixed incomplete semver",
			args: args{
				s: "1.1.foo",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "pre- and suffixed semver",
			args: args{
				s: "foo-1.1.0-bar",
			},
			want:    "-bar",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := suffix(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("suffix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("suffix() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    SemVer
		wantErr bool
	}{
		{
			name: "regular",
			args: args{
				s: "v1.2.3-dev",
			},
			want: SemVer{
				major:  1,
				minor:  2,
				patch:  3,
				prefix: "v",
				presep: "",
				suffix: "-dev",
				sufsep: "",
			},
			wantErr: false,
		},
		{
			name: "regular",
			args: args{
				s: "v1.2.3-dev",
			},
			want: SemVer{
				major:  1,
				minor:  2,
				patch:  3,
				prefix: "v",
				presep: "",
				suffix: "-dev",
				sufsep: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemVer_separators(t *testing.T) {
	type fields struct {
		prefix string
		presep string
		suffix string
		sufsep string
	}
	type args struct {
		pre string
		suf string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "default",
			fields: fields{
				prefix: "foo-",
				presep: "",
				suffix: "-bar",
				sufsep: "",
			},
			args: args{
				pre: "-",
				suf: "-",
			},
			want: fields{
				prefix: "foo",
				presep: "-",
				suffix: "bar",
				sufsep: "-",
			},
		},
		{
			name: "multi-char separators",
			fields: fields{
				prefix: "foo--",
				presep: "",
				suffix: "--bar",
				sufsep: "",
			},
			args: args{
				pre: "--",
				suf: "--",
			},
			want: fields{
				prefix: "foo",
				presep: "--",
				suffix: "bar",
				sufsep: "--",
			},
		},
		{
			name: "some fuckery",
			fields: fields{
				prefix: "foo-",
				presep: "",
				suffix: "-bar",
				sufsep: "",
			},
			args: args{
				pre: "",
				suf: "-",
			},
			want: fields{
				prefix: "foo-",
				presep: "",
				suffix: "bar",
				sufsep: "-",
			},
		},
		{
			name: "only separator",
			fields: fields{
				prefix: "-",
				presep: "",
				suffix: "-",
				sufsep: "",
			},
			args: args{
				pre: "-",
				suf: "-",
			},
			want: fields{
				prefix: "",
				presep: "-",
				suffix: "",
				sufsep: "-",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SemVer{
				prefix: tt.fields.prefix,
				presep: tt.fields.presep,
				suffix: tt.fields.suffix,
				sufsep: tt.fields.sufsep,
			}
			s.separators(tt.args.pre, tt.args.suf)
			got := fields{
				prefix: s.prefix,
				presep: s.presep,
				suffix: s.suffix,
				sufsep: s.sufsep,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	type args struct {
		v []SemVer
	}
	tests := []struct {
		name string
		args args
		want SemVer
	}{
		{
			name: "Major only",
			args: args{
				v: []SemVer{
					{
						major: 1,
					},
					{
						major: 2,
					},
				},
			},
			want: SemVer{
				major: 2,
			},
		},
		{
			name: "Minor only",
			args: args{
				v: []SemVer{
					{
						minor: 1,
					},
					{
						minor: 2,
					},
				},
			},
			want: SemVer{
				minor: 2,
			},
		},
		{
			name: "Patch only",
			args: args{
				v: []SemVer{
					{
						patch: 1,
					},
					{
						patch: 2,
					},
				},
			},
			want: SemVer{
				patch: 2,
			},
		},
		{
			name: "With suffix is smaller than no suffix",
			args: args{
				v: []SemVer{
					{
						suffix: "foo",
					},
					{},
				},
			},
			want: SemVer{
				suffix: "",
			},
		},
		{
			name: "Major and minor",
			args: args{
				v: []SemVer{
					{
						major: 1,
						minor: 2,
					},
					{
						major: 1,
						minor: 3,
					},
				},
			},
			want: SemVer{
				major: 1,
				minor: 3,
			},
		},
		{
			name: "Major and minor and patch and suffix",
			args: args{
				v: []SemVer{
					{
						major:  0,
						minor:  0,
						patch:  1,
						suffix: "staging",
					},
					{
						major:  0,
						minor:  0,
						patch:  3,
						suffix: "dev",
					},
					{
						major:  0,
						minor:  0,
						patch:  2,
						suffix: "staging",
					},
				},
			},
			want: SemVer{
				major:  0,
				minor:  0,
				patch:  3,
				suffix: "dev",
			},
		},
		{
			name: "Major and minor and patch",
			args: args{
				v: []SemVer{
					{
						major: 1,
						minor: 2,
						patch: 3,
					},
					{
						major: 1,
						minor: 2,
						patch: 4,
					},
				},
			},
			want: SemVer{
				major: 1,
				minor: 2,
				patch: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Max(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}
