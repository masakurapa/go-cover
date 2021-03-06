package filter_test

import (
	"testing"

	"github.com/masakurapa/gover-html/internal/cover/filter"
	"github.com/masakurapa/gover-html/internal/option"
	"github.com/masakurapa/gover-html/test/helper"
)

func TestFilter_IsOutputTarget(t *testing.T) {
	type newArgs struct {
		option option.Option
	}
	type args struct {
		path     string
		fileName string
	}

	tests := []struct {
		name    string
		newArgs newArgs
		args    args
		want    bool
	}{
		{
			name: "include, excludeの指定なし、pathが/で始まらない場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForDefault(t),
			},

			args: args{
				path: "path/to/dir1",
			},

			want: true,
		},
		{
			name: "include, excludeの指定なし、pathが/で始まる場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForDefault(t),
			},
			args: args{
				path: "/path/to/dir1",
			},
			want: false,
		},

		// validate "include"
		{
			name: "include = path の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: true,
		},
		{
			name: "include != path の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir2",
			},
			want: false,
		},
		{
			name: "includeの接頭語(./)を除いた値 = path の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"./path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: true,
		},
		{
			name: "includeの接頭語(./)を除いた値 != path の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"./path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir2",
			},
			want: false,
		},
		{
			name: "includeの接尾語(/)を除いた値 = path の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1/"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: true,
		},
		{
			name: "includeの接尾語(/)を除いた値 != path の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1/"}),
			},
			args: args{
				path: "path/to/dir2",
			},
			want: false,
		},
		{
			name: "include = pathの接頭語(./)を除いた値 の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "./path/to/dir1",
			},
			want: true,
		},
		{
			name: "include != pathの接頭語(./)を除いた値 の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "./path/to/dir2",
			},
			want: false,
		},
		{
			name: "include = pathの接尾語(/)を除いた値 の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir1/",
			},
			want: true,
		},
		{
			name: "include != pathの接尾語(/)を除いた値 の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir2/",
			},
			want: false,
		},
		{
			name: "include = pathの接頭語 の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: true,
		},
		{
			name: "include != pathの接頭語 の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to"}),
			},
			args: args{
				path: "path/tooo/dir1",
			},
			want: false,
		},
		{
			name: "includeが複数あり、いずれかと同じパスの場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir2", "path/to/dir3", "path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: true,
		},
		{
			name: "includeが複数あり、いずれかと同じパスにもマッチしない場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir2", "path/to/dir3", "path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir4",
			},
			want: false,
		},
		{
			name: "include = path + fileName の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1/file.go"}),
			},
			args: args{
				path:     "path/to/dir1",
				fileName: "file.go",
			},
			want: true,
		},
		{
			name: "include = path + fileName の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForInclude(t, []string{"path/to/dir1/file.go"}),
			},
			args: args{
				path:     "path/to/dir1",
				fileName: "file2.go",
			},
			want: false,
		},

		// validate "exclude"
		{
			name: "exclude = path の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: false,
		},
		{
			name: "exclude != path の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir2",
			},
			want: true,
		},
		{
			name: "excludeの接頭語(./)を除いた値 = path の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"./path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: false,
		},
		{
			name: "excludeの接頭語(./)を除いた値 != path の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"./path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir2",
			},
			want: true,
		},
		{
			name: "excludeの接尾語(/)を除いた値 = path の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1/"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: false,
		},
		{
			name: "excludeの接尾語(/)を除いた値 != path の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1/"}),
			},
			args: args{
				path: "path/to/dir2",
			},
			want: true,
		},
		{
			name: "exclude = pathの接頭語(./)を除いた値 の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "./path/to/dir1",
			},
			want: false,
		},
		{
			name: "exclude != pathの接頭語(./)を除いた値 の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "./path/to/dir2",
			},
			want: true,
		},
		{
			name: "exclude = pathの接尾語(/)を除いた値 の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir1/",
			},
			want: false,
		},
		{
			name: "exclude != pathの接尾語(/)を除いた値 の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir2/",
			},
			want: true,
		},
		{
			name: "exclude = pathの接頭語 の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: false,
		},
		{
			name: "exclude != pathの接頭語 の場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to"}),
			},
			args: args{
				path: "path/tooo/dir1",
			},
			want: true,
		},
		{
			name: "excludeがカンマ区切りで複数あり、いずれかと同じパスの場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir2", "path/to/dir3", "path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir1",
			},
			want: false,
		},
		{
			name: "excludeがカンマ区切りで複数あり、いずれかと同じパスにもマッチしない場合はtrueが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir2", "path/to/dir3", "path/to/dir1"}),
			},
			args: args{
				path: "path/to/dir4",
			},
			want: true,
		},
		{
			name: "exclude = path + fileName の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1/file.go"}),
			},
			args: args{
				path:     "path/to/dir1",
				fileName: "file.go",
			},
			want: false,
		},
		{
			name: "exclude = path + fileName の場合はfalseが返る",
			newArgs: newArgs{
				option: helper.GetOptionForExclude(t, []string{"path/to/dir1/file.go"}),
			},
			args: args{
				path:     "path/to/dir1",
				fileName: "file2.go",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filter.New(tt.newArgs.option)
			if got := f.IsOutputTarget(tt.args.path, tt.args.fileName); got != tt.want {
				t.Errorf("filter.IsOutputTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}
