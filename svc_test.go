package sign

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dromara/carbon/v2"
	"github.com/prashantv/gostub"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

var rewards = []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

func mockDB() *gorm.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err)
	}
	mock.MatchExpectationsInOrder(false)

	signAtRows := sqlmock.NewRows([]string{"sign_at"}).AddRow(20210102)
	mock.ExpectQuery("SELECT `sign_at` FROM `sign_daily`.+").WillReturnRows(signAtRows)

	cumulativeRows := sqlmock.NewRows([]string{"maxnum"})
	mock.ExpectQuery("SELECT max\\(num\\) as maxnum FROM `sign_cumulative`.+").WillReturnRows(cumulativeRows)

	mock.ExpectExec("INSERT INTO `sign_daily`.+").
		WillReturnResult(sqlmock.NewResult(1, 1))

	return gdb
}
func TestCumulativeLists(t *testing.T) {
	type args struct {
		db  *gorm.DB
		cid int64
	}

	stubs := gostub.StubFunc(&timeNow, carbon.Parse("2021-01-02"))
	defer stubs.Reset()
	db := mockDB()

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "good",
			args:    args{db, 1},
			want:    5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CumulativeLists(tt.args.db, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("CumulativeLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CumulativeLists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCumulativeSignIn(t *testing.T) {
	type args struct {
		db  *gorm.DB
		cid int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CumulativeSignIn(tt.args.db, tt.args.cid); (err != nil) != tt.wantErr {
				t.Errorf("CumulativeSignIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDailyLists(t *testing.T) {
	type args struct {
		db      *gorm.DB
		rewards []interface{}
		cid     int64
	}

	stubs := gostub.StubFunc(&timeNow, carbon.Parse("2021-01-02"))
	defer stubs.Reset()
	db := mockDB()

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "good",
			args: args{db, rewards, 0},
			want: map[string]interface{}{
				"list":     rewards,
				"sign_num": 1,
				"max_num":  31,
				"status":   1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DailyLists(tt.args.db, tt.args.rewards, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("DailyLists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DailyLists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDailySignIn(t *testing.T) {
	type args struct {
		db       *gorm.DB
		cid      int64
		mocktime string
	}

	db := mockDB()

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "good",
			args:    args{db, 1, "2021-01-01"},
			wantErr: false,
		},
		{
			name:    "bad",
			args:    args{db, 1, "2021-01-02"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubs := gostub.StubFunc(&timeNow, carbon.Parse(tt.args.mocktime))
			defer stubs.Reset()
			if err := DailySignIn(tt.args.db, tt.args.cid); (err != nil) != tt.wantErr {
				t.Errorf("DailySignIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
