FILES := file.go strings.go parallel.go
TEST_FILES := file_test.go strings_test.go

TEST := test-result

test: $(TEST)

$(TEST): $(FILES) $(TEST_FILES)
	goimports -w $?
	go test | tee $@

clean:
	rm -f $(TEST)
