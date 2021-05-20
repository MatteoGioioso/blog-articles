package main

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/jackc/pgx/v4/pgxpool"
	"io/ioutil"
	"os"
)

// You need to pass the pool configuration to this method
// In thi case I am using pgx
func getTSLConfig() error {
	// Get the certificate file
	certPath := "/certs/rds-combined-ca-bundle.pem"
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return err
	}

	// Alternatively you can use x509.SystemCertPool()
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Initialize the tsl config
	tlsConfig := &tls.Config{}
	// The server name is the RDS endpoint,
	// should be something like: <name>.<id>.<aws_region>.rds.amazonaws.com
	// You can find the endpoint in the RDS console
	tlsConfig.ServerName = os.Getenv("POSTGRES_HOST")
	tlsConfig.RootCAs = caCertPool

	// Attach the tsl config to the pool
	config, _ := pgxpool.ParseConfig("")
	config.ConnConfig.TLSConfig = tlsConfig
	return nil
}
