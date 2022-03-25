test:
	go test --cover --race -count=100 -failfast ./...
	go test . -fuzz=FuzzLogWriterJSONInPlaceTagsWithInt -fuzztime 20s
	go test . -fuzz=FuzzLogWriterJSONInPlaceTagsWithFloat -fuzztime 20s
	go test . -fuzz=FuzzLogWriterJSONInPlaceTagsWithBool -fuzztime 20s
	go test . -fuzz=FuzzLogWriterJSONInPlaceTagsWithString -fuzztime 20s
	go test . -fuzz=FuzzLogWriterJSONInPlaceTagsWithMessage -fuzztime 20s
