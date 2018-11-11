package kubeless

import (
	"fmt"
	"math"
	"strconv"

	"github.com/kubeless/kubeless/pkg/functions"
	"github.com/sirupsen/logrus"
)

func IsPrime(event functions.Event, context functions.Context) (string, error) {
	num, err := strconv.Atoi(event.Data)
	if err != nil {
		return "", fmt.Errorf("Failed to parse %s as int! %v", event.Data, err)
	}
	logrus.Infof("Checking if %s is prime", event.Data)
	if num <= 1 {
		return fmt.Sprintf("%d is not prime", num), nil
	}
	for i := 2; i <= int(math.Floor(float64(num)/2)); i++ {
		if num%i == 0 {
			return fmt.Sprintf("%d is not prime", num), nil
		}
	}
	return fmt.Sprintf("%d is prime", num), nil
}
