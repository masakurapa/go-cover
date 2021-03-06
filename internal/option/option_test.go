package option_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/masakurapa/gover-html/internal/option"
	"github.com/masakurapa/gover-html/internal/reader"
	"github.com/masakurapa/gover-html/test/helper"
)

type mockReader struct {
	reader.Reader
	mockRead   func(string) (io.Reader, error)
	mockExists func(string) bool
}

func (m *mockReader) Read(s string) (io.Reader, error) {
	return m.mockRead(s)
}

func (m *mockReader) Exists(s string) bool {
	return m.mockExists(s)
}

func TestNew(t *testing.T) {
	type args struct {
		input   *string
		output  *string
		theme   *string
		include *string
		exclude *string
	}
	type testCase []struct {
		name     string
		settings string
		args     args
		want     *option.Option
		wantErr  bool
	}

	t.Run("設定ファイルが存在しない", func(t *testing.T) {
		tests := testCase{
			{
				name: "全項目に設定値が存在(theme=dark)",
				args: args{
					input:   helper.StringP("example.out"),
					output:  helper.StringP("example.html"),
					theme:   helper.StringP("dark"),
					include: helper.StringP("path/to/dir1,path/to/dir2"),
					exclude: helper.StringP("path/to/dir3,path/to/dir4"),
				},
				want: &option.Option{
					Input:   "example.out",
					Output:  "example.html",
					Theme:   "dark",
					Include: []string{"path/to/dir1", "path/to/dir2"},
					Exclude: []string{"path/to/dir3", "path/to/dir4"},
				},
				wantErr: false,
			},
			{
				name: "全項目に設定値が存在(theme=light)",
				args: args{
					input:   helper.StringP("example.out"),
					output:  helper.StringP("example.html"),
					theme:   helper.StringP("light"),
					include: helper.StringP("path/to/dir1,path/to/dir2"),
					exclude: helper.StringP("path/to/dir3,path/to/dir4"),
				},
				want: &option.Option{
					Input:   "example.out",
					Output:  "example.html",
					Theme:   "light",
					Include: []string{"path/to/dir1", "path/to/dir2"},
					Exclude: []string{"path/to/dir3", "path/to/dir4"},
				},
				wantErr: false,
			},
			{
				name: "全項目に空文字を指定",
				args: args{
					input:   helper.StringP(""),
					output:  helper.StringP(""),
					theme:   helper.StringP(""),
					include: helper.StringP(""),
					exclude: helper.StringP(""),
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "全項目にnilを指定",
				args: args{
					input:   nil,
					output:  nil,
					theme:   nil,
					include: nil,
					exclude: nil,
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "includeに空の値を持つ",
				args: args{
					include: helper.StringP("path/to/dir1,,path/to/dir2,,"),
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{"path/to/dir1", "path/to/dir2"},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "include./で始まるパスを指定",
				args: args{
					include: helper.StringP("./path/to/dir1"),
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{"path/to/dir1"},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "includeに/で終わるパスを指定",
				args: args{
					include: helper.StringP("path/to/dir1/"),
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{"path/to/dir1"},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "includeに/で始まるパスを指定",
				args: args{
					include: helper.StringP("/path/to/dir1"),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "excludeに空の値を持つ",
				args: args{
					exclude: helper.StringP("path/to/dir3,,path/to/dir4,,"),
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{"path/to/dir3", "path/to/dir4"},
				},
				wantErr: false,
			},
			{
				name: "excludeに./で始まるパスを指定",
				args: args{
					exclude: helper.StringP("./path/to/dir3"),
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{"path/to/dir3"},
				},
				wantErr: false,
			},
			{
				name: "excludeに/で終わるパスを指定",
				args: args{
					exclude: helper.StringP("path/to/dir3/"),
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{"path/to/dir3"},
				},
				wantErr: false,
			},
			{
				name: "excludeに/で始まるパスを指定",
				args: args{
					exclude: helper.StringP("/path/to/dir3"),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "themeに期待値以外を設定",
				args: args{
					theme: helper.StringP("unknown"),
				},
				want:    nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				readerMock := &mockReader{
					mockExists: func(string) bool { return false },
				}

				got, err := option.New(readerMock).
					Generate(tt.args.input, tt.args.output, tt.args.theme, tt.args.include, tt.args.exclude)
				if (err != nil) != tt.wantErr {
					t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("New() = %#v, want %#v", got, tt.want)
				}
			})
		}
	})

	t.Run("設定ファイルが存在する", func(t *testing.T) {
		tests := testCase{
			{
				name: "全項目に設定値が存在(theme=dark)",
				settings: `
input: example.out
output: example.html
theme: dark
include:
  - path/to/dir1
  - path/to/dir2
exclude:
  - path/to/dir3
  - path/to/dir4
`,
				args: args{},
				want: &option.Option{
					Input:   "example.out",
					Output:  "example.html",
					Theme:   "dark",
					Include: []string{"path/to/dir1", "path/to/dir2"},
					Exclude: []string{"path/to/dir3", "path/to/dir4"},
				},
				wantErr: false,
			},
			{
				name: "全項目に設定値が存在(theme=light)",
				settings: `
input: example.out
output: example.html
theme: light
include:
  - path/to/dir1
  - path/to/dir2
exclude:
  - path/to/dir3
  - path/to/dir4
`,
				args: args{},
				want: &option.Option{
					Input:   "example.out",
					Output:  "example.html",
					Theme:   "light",
					Include: []string{"path/to/dir1", "path/to/dir2"},
					Exclude: []string{"path/to/dir3", "path/to/dir4"},
				},
				wantErr: false,
			},
			{
				name: "全項目に設定値が存在し、引数に全項目に設定値が存在",
				settings: `
input: example.out
output: example.html
theme: dark
include:
  - path/to/dir1
  - path/to/dir2
exclude:
  - path/to/dir3
  - path/to/dir4
`,
				args: args{
					input:   helper.StringP("example2.out"),
					output:  helper.StringP("example2.html"),
					theme:   helper.StringP("light"),
					include: helper.StringP("path/to/dir5"),
					exclude: helper.StringP("path/to/dir6"),
				},
				want: &option.Option{
					Input:   "example2.out",
					Output:  "example2.html",
					Theme:   "light",
					Include: []string{"path/to/dir5"},
					Exclude: []string{"path/to/dir6"},
				},
				wantErr: false,
			},
			{
				name: "全項目に設定値が存在し、引数に全項目に空文字を設定",
				settings: `
input: example.out
output: example.html
theme: light
include:
  - path/to/dir1
  - path/to/dir2
exclude:
  - path/to/dir3
  - path/to/dir4
`,
				args: args{
					input:   helper.StringP(""),
					output:  helper.StringP(""),
					theme:   helper.StringP(""),
					include: helper.StringP(""),
					exclude: helper.StringP(""),
				},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{},
				},
				wantErr: false,
			},

			{
				name: "全項目のキーのみが存在する",
				settings: `
input:
output:
theme:
include:
exclude:
`,
				args: args{},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "全項目のキーが存在しない",
				settings: `
# empty settings
`,
				args: args{},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "includeに空の値を持つ",
				settings: `
include:
  - path/to/dir1
  -
  - path/to/dir2
  -
  -
`,
				args: args{},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{"path/to/dir1", "path/to/dir2"},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "include./で始まるパスを指定",
				settings: `
include:
  - ./path/to/dir1
`,
				args: args{},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{"path/to/dir1"},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "includeに/で終わるパスを指定",
				settings: `
include:
  - path/to/dir1/
`,
				args: args{},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{"path/to/dir1"},
					Exclude: []string{},
				},
				wantErr: false,
			},
			{
				name: "includeに/で始まるパスを指定",
				settings: `
include:
  - /path/to/dir1
`,
				args:    args{},
				want:    nil,
				wantErr: true,
			},
			{
				name: "excludeに空の値を持つ",
				settings: `
exclude:
  - path/to/dir3
  -
  - path/to/dir4
  -
  -
`,
				args: args{},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{"path/to/dir3", "path/to/dir4"},
				},
				wantErr: false,
			},
			{
				name: "excludeに./で始まるパスを指定",
				settings: `
exclude:
  - ./path/to/dir3
`,
				args: args{},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{"path/to/dir3"},
				},
				wantErr: false,
			},
			{
				name: "excludeに/で終わるパスを指定",
				settings: `
exclude:
  - path/to/dir3/
`,
				args: args{},
				want: &option.Option{
					Input:   "coverage.out",
					Output:  "coverage.html",
					Theme:   "dark",
					Include: []string{},
					Exclude: []string{"path/to/dir3"},
				},
				wantErr: false,
			},
			{
				name: "excludeに/で始まるパスを指定",
				settings: `
exclude:
  - /path/to/dir3
`,
				args:    args{},
				want:    nil,
				wantErr: true,
			},
			{
				name: "themeに期待値以外を設定",
				args: args{
					theme: helper.StringP("unknown"),
				},
				want:    nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				readerMock := &mockReader{
					mockExists: func(string) bool { return true },
					mockRead: func(string) (io.Reader, error) {
						return strings.NewReader(tt.settings), nil
					},
				}

				got, err := option.New(readerMock).
					Generate(tt.args.input, tt.args.output, tt.args.theme, tt.args.include, tt.args.exclude)
				if (err != nil) != tt.wantErr {
					t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("New() = %#v, want %#v", got, tt.want)
				}
			})
		}

	})
}
