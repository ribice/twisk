twirp:
	./twirp
	rm -rf openapi/statik
	retool do statik -m -f -src openapi/ -dest openapi 