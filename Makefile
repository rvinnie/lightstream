build_streaming_service:
	$(MAKE) -C ./services/streaming

run: build_streaming_service


.DEFAULT_GOAL := run
.PHONY: build_streaming_service, run