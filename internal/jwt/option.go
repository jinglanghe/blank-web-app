package jwt

type Options struct {
	// signing algorithm - possible values are HS256, HS384, HS512, RS256, RS384 or RS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Public key file for asymmetric algorithms
	PubKeyFile string
}

type Option func(*Options)

func SigningAlgorithm(algorithm string) Option {
	return func(opts *Options) {
		opts.SigningAlgorithm = algorithm
	}
}

func SecretKey(secret []byte) Option {
	return func(opts *Options) {
		opts.Key = secret
	}
}

func PublicKey(pub string) Option {
	return func(opts *Options) {
		opts.PubKeyFile = pub
	}
}