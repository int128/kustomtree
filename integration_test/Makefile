OUTPUT_DIR=output

.PHONY: test
test:
	rm -fr $(OUTPUT_DIR)
	mkdir -p $(OUTPUT_DIR)
	cp -a testdata/ $(OUTPUT_DIR)/testdata/
	go run .. $(OUTPUT_DIR)/testdata/
	kustomize build testdata > $(OUTPUT_DIR)/original.yaml
	kustomize build $(OUTPUT_DIR)/testdata > $(OUTPUT_DIR)/kustomtree.yaml
	diff -u $(OUTPUT_DIR)/original.yaml $(OUTPUT_DIR)/kustomtree.yaml
