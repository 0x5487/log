package main

// import (
// 	"errors"
// 	"fmt"

// 	"github.com/jasonsoft/log"
// 	"github.com/jasonsoft/log/handlers/console"
// 	pkgErr "github.com/pkg/errors"
// )

// var ErrNotFound = pkgErr.New("record not found")

// func main() {
// 	clog := console.New()
// 	log.RegisterHandler(clog, log.AllLevels...) // use console handler to log all level log

// 	err := http()
// 	if err != nil {
// 		if errors.Is(err, ErrNotFound) {
// 			fmt.Printf("err is ErrNotFound: %v\n", err)
// 		}

// 		var errDup error
// 		if errors.As(err, &errDup) {
// 			fmt.Printf("err as : %v\n", errDup)
// 		}

// 		if pkgErr.Cause(err) == ErrNotFound {
// 			fmt.Printf("root ErrNotFound error: %v\n", err)
// 		}

// 		// format the error to a string and print it
// 		//formattedStr := eris.ToString(err, true)

// 		log.WithError(err).Error("oops")
// 	}
// }

// func repo() error {
// 	return pkgErr.Wrap(ErrNotFound, "id 6 was not found")
// }

// func service() error {
// 	err := repo()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func http() error {
// 	err := service()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
