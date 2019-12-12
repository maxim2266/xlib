FILES := file.go strings.go
TEST_FILES := file_test.go strings_test.go

TEST := test-result

test: $(TEST)

$(TEST): $(FILES) $(TEST_FILES)
	goimports -w $?
	golint $?
	go test | tee $@

clean:
	rm -f $(TEST)
