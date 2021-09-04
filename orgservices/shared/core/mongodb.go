package core

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB is an instance of the MongoDB client for managing off-chain persistent data.
var MongoDB *mongo.Client

func InitMongoDB() {
	var (
		opts = options.Client()
		err  error
	)

	if !viper.GetBool("mongo_enabled") {
		return
	}

	opts.ApplyURI(viper.GetString("mongo_address"))
	opts.SetConnectTimeout(viper.GetDuration("mongo_connection_timeout"))

	if viper.GetBool("mongo_auth") {
		opts.SetAuth(options.Credential{
			Username: viper.GetString("mongo_username"),
			Password: viper.GetString("mongo_password"),
		})
	}

	if viper.GetBool("mongo_tls") {
		tlsConfig, err := getMongoTlsConfig(viper.GetString("mongo_ca_cert_path"))
		if err != nil {
			Logger.Fatal(err)
		}

		opts.SetTLSConfig(tlsConfig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), viper.GetDuration("mongo_connection_timeout"))
	defer cancel()

	if MongoDB, err = mongo.Connect(ctx, opts); err != nil {
		Logger.Fatal(errors.Wrap(err, "failed to create MongoDB client"))
	}

	ctx, cancel = context.WithTimeout(context.Background(), viper.GetDuration("mongo_query_timeout"))
	defer cancel()

	if err = MongoDB.Ping(ctx, nil); err != nil {
		Logger.Fatal(errors.Wrap(err, "failed to ping MongoDB"))
	}

	Logger.Info("Successfully connected to MongoDB")
}

func getMongoTlsConfig(path string) (*tls.Config, error) {
	certs := x509.NewCertPool()
	pem, err := ioutil.ReadFile(path); if err != nil {
		return nil, errors.Wrapf(err, "failed to read certificate from path, %s", path)
	}
	certs.AppendCertsFromPEM(pem)
	return &tls.Config{
		RootCAs: certs,
	}, nil
}
