ls /foo
envsubst '$PORT_API' < /foo/swagger.yaml > /foo/swagger-ready.yaml
