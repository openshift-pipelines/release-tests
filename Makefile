##########------------------------------------------------------------##########
##########- Download Latest Tkn binary ------------------------------------------##########
##########------------------------------------------------------------##########

STABLE_DOWNLOAD_URL := https://github.com/tektoncd/cli/releases/download/v${TKN_VERSION}/tkn_${TKN_VERSION}_Linux_x86_64.tar.gz
DOWNLOAD_PATH=build/tkn/v${TKN_VERSION}

.PHONY: download-tkn
download-tkn:
ifndef TKN_VERSION
	@echo TKN_VERSION not set
	@exit 1
endif
	[[ -d "${DOWNLOAD_PATH}" ]] | mkdir -p ${DOWNLOAD_PATH}
	curl -L -o ${DOWNLOAD_PATH}/tkn.tar.gz ${STABLE_DOWNLOAD_URL}
	tar xvzf ${DOWNLOAD_PATH}/tkn.tar.gz -C ${DOWNLOAD_PATH} tkn
	rm -rf ${DOWNLOAD_PATH}/tkn.tar.gz


.PHONY: lint
lint: ## run linter(s)
	@echo "Linting..."
	@golangci-lint run ./... --timeout 5m
