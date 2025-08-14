# makefile listing from: http://stackoverflow.com/questions/4219255/how-do-you-get-the-list-of-targets-in-a-makefile
default-goal:
	@echo "viable targets:"
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

# security, performance check
# https://staticcheck.io/docs/getting-started/#distribution-packages
# go install directly instead of using a potentially outdated package management version
# go install honnef.co/go/tools/cmd/staticcheck@latest
exists_go_static: ; @which staticcheck > /dev/null

# You can use staticcheck -explain <check> to get a helpful description of a check.
do_go_static: exists_go_static
	staticcheck ./...

checks: do_go_static
	@# correctness check
	go vet ./...

