// Package parser provides functionality for reading and parsing .env files
// into structured key-value maps.
//
// Supported .env syntax:
//   - KEY=VALUE              plain key-value pairs
//   - KEY="VALUE"            double-quoted values (quotes stripped)
//   - KEY='VALUE'            single-quoted values (quotes stripped)
//   - # comment              lines starting with '#' are ignored
//   - (empty lines)          empty lines are ignored
//   - KEY=a=b=c              values may contain '=' characters
//
// Example usage:
//
//	env, err := parser.ParseFile(".env.production")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(env["DATABASE_URL"])
package parser
