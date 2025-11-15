#!/bin/bash

# Install Go manually, since oss-fuzz ships with an outdated Go version.
# See https://github.com/google/oss-fuzz/pull/10643.
export CXX="${CXX} -lresolv" # required by Go 1.20
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz \
  && mkdir temp-go \
  && rm -rf /root/.go/* \
  && tar -C temp-go/ -xzf go1.25.0.linux-amd64.tar.gz \
  && mv temp-go/go/* /root/.go/ \
  && rm -rf temp-go go1.25.0.linux-amd64.tar.gz

(
# fuzz qpack
compile_go_fuzzer github.com/quic-go/qpack/fuzzing Fuzz qpack_fuzzer
)

(
# fuzz quic-go
compile_go_fuzzer github.com/quic-go/quic-go/fuzzing/frames Fuzz frame_fuzzer
compile_go_fuzzer github.com/quic-go/quic-go/fuzzing/header Fuzz header_fuzzer
compile_go_fuzzer github.com/quic-go/quic-go/fuzzing/transportparameters Fuzz transportparameter_fuzzer
compile_go_fuzzer github.com/quic-go/quic-go/fuzzing/tokens Fuzz token_fuzzer
compile_go_fuzzer github.com/quic-go/quic-go/fuzzing/handshake Fuzz handshake_fuzzer

if [ $SANITIZER == "coverage" ]; then
    # no need for corpora if coverage
    exit 0
fi

# generate seed corpora
cd $GOPATH/src/github.com/quic-go/quic-go/
go generate -x ./fuzzing/...

zip --quiet -r $OUT/header_fuzzer_seed_corpus.zip fuzzing/header/corpus
zip --quiet -r $OUT/frame_fuzzer_seed_corpus.zip fuzzing/frames/corpus
zip --quiet -r $OUT/transportparameter_fuzzer_seed_corpus.zip fuzzing/transportparameters/corpus
zip --quiet -r $OUT/handshake_fuzzer_seed_corpus.zip fuzzing/handshake/corpus
)

# for debugging
ls -al $OUT
