package tests

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/n25a/eavesdropper/cmd"
	"github.com/n25a/eavesdropper/internal/app"
	"github.com/n25a/eavesdropper/internal/config"
	"github.com/stretchr/testify/suite"
)

type subjectOne struct {
	Data   int    `json:"data" db:"data"`
	Status string `json:"status" db:"status"`
}

type subjectTwo struct {
	Var1 string  `json:"var1" db:"var1"`
	Var2 int     `json:"var2" db:"var2"`
	Var3 float32 `json:"var3" db:"var3"`
	Var4 bool    `json:"var4" db:"var4"`
}

type EavesdropperTimeScaleDBNatsTestSuite struct {
	suite.Suite
	shutdown chan os.Signal
}

func (suite *EavesdropperTimeScaleDBNatsTestSuite) SetupSuite() {
	require := suite.Require()
	cmd.ConfigPath = "./config_timescaledb_nats_test.yaml"
	suite.shutdown = make(chan os.Signal, 2)
	signal.Notify(suite.shutdown, os.Interrupt, syscall.SIGTERM)

	// migrate database

	// stat service
	go cmd.Eavesdropping(&suite.shutdown)

	// initialize config and app
	err := config.LoadConfig(cmd.ConfigPath)
	require.NoError(err)

	err = app.InitApp()
	require.NoError(err)
}

func (suite *EavesdropperTimeScaleDBNatsTestSuite) TearDownSuite() {
	suite.shutdown <- os.Interrupt
}

func (suite *EavesdropperTimeScaleDBNatsTestSuite) TestEavesdropping() {
	require := suite.Require()

	// publish message to nats
	subjectOneData1 := subjectOne{1, "ok"}
	err := app.A.MQ.Publish("subject-one", subjectOneData1)
	require.NoError(err)

	subjectOneData2 := subjectOne{13, "fail"}
	err = app.A.MQ.Publish("subject-one", subjectOneData2)
	require.NoError(err)

	subjectTwoData1 := subjectTwo{"test1", 2, 3.14, true}
	err = app.A.MQ.Publish("subject-two", subjectTwoData1)
	require.NoError(err)

	subjectTwoData2 := subjectTwo{"test2", 558, 6789.14, false}
	err = app.A.MQ.Publish("subject-two", subjectTwoData2)
	require.NoError(err)

	// check database
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	data, err := app.A.DB.Get(
		ctx,
		"SELECT * FROM subject_one WHERE data = $1",
		subjectOne{},
		subjectOneData1.Data,
	)
	require.NoError(err)
	require.Equal(subjectOneData1, data)

	data, err = app.A.DB.Get(
		ctx,
		"SELECT * FROM subject_one WHERE data = $1",
		subjectOne{},
		subjectOneData2.Data,
	)
	require.NoError(err)
	require.Equal(subjectOneData2, data)

	data, err = app.A.DB.Get(
		ctx,
		"SELECT * FROM subject_two WHERE var1 = $1",
		subjectTwo{},
		subjectTwoData1.Var1,
	)
	require.NoError(err)
	require.Equal(subjectTwoData1, data)

	data, err = app.A.DB.Get(
		ctx,
		"SELECT * FROM subject_two WHERE var1 = $1",
		subjectTwo{},
		subjectTwoData2.Var1,
	)
	require.NoError(err)
	require.Equal(subjectTwoData2, data)
}

func TestEavesdropperTimeScaleDBNats(t *testing.T) {
	suite.Run(t, new(EavesdropperTimeScaleDBNatsTestSuite))
}
