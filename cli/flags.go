package cli

import "flag"

type options struct {
	address string
	schema  string
}

func ParseArgs() *options {
	address := flag.String("addr", "", "Address of ethereum transaction.")
	schema := flag.String("schema", "", "Filename of db schema structure.")

	return &options{
		address: *address,
		schema:  *schema,
	}
}
