package activity_test

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCanFillBuyOrder(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard
}
