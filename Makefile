version = $(shell cat CHANGES.rst | awk '/^[0-9]+\.[0-9]+(\.[0-9]+)?/' | head -n1)
tag_name = v$(version)

tag:
	@git rev-parse --abbrev-ref HEAD | grep '^main$$'
	@git tag -a $(tag_name) -m "New version of memory-dheck."
	@git push origin $(tag_name)
	@echo "successfully push new version of memory-dheck"