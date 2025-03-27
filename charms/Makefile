.PHONY: interface, server
interface:
	cd ophelia-interface && \
	CHARMCRAFT_ENABLE_EXPERIMENTAL_EXTENSIONS=true charmcraft pack
	mv ophelia-interface/*.charm ../dist/

server:
	cd ophelia-server && \
	CHARMCRAFT_ENABLE_EXPERIMENTAL_EXTENSIONS=true charmcraft pack
	mv ophelia-server/*.charm ../dist/
