#!/bin/bash

# Build and Run the different test permutations.
# This helps validate that nothing gets broken.

_build_proceed() {
    # return success (0) if we should, and 1 (fail) if not
    if [[ "${zforce}" ]]; then return 0; fi
    for a in "fastpath.generated.go" "json.mono.generated.go"; do
        if [[ ! -e "$a" ]]; then return 0; fi
        for b in `ls -1 *.go.tmpl gen.go gen_mono.go values_test.go`; do
            if [[ "$a" -ot "$b" ]]; then return 0; fi
        done
    done
    return 1
}

# _build generates fastpath.go
_build() {
    # if ! [[ "${zforce}" || $(_ng "fastpath.generated.go") || $(_ng "json.mono.generated.go") ]]; then return 0; fi
    _build_proceed
    if [ $? -eq 1 ]; then return 0; fi
    if [ "${zbak}" ]; then
        _zts=`date '+%m%d%Y_%H%M%S'`
        _gg=".generated.go"
        [ -e "fastpath${_gg}" ] && mv fastpath${_gg} fastpath${_gg}__${_zts}.bak
        [ -e "gen${_gg}" ] && mv gen${_gg} gen${_gg}__${_zts}.bak
    fi
    
    rm -f fast*path.generated.go *mono*generated.go *_generated_test.go gen-from-tmpl*.generated.go

    local btags="codec.build codec.notmono codec.safe codec.notfastpath"

    cat > gen-from-tmpl.codec.generated.go <<EOF
package codec
func GenTmplRun2Go(in, out string) { genTmplRun2Go(in, out) }
func GenMonoAll() { genMonoAll() }
EOF

    cat > gen-from-tmpl.generated.go <<EOF
//go:build ignore
package main
import "${zpkg}"
func main() {
codec.GenTmplRun2Go("fastpath.go.tmpl", "base.fastpath.generated.go")
codec.GenTmplRun2Go("fastpath.notmono.go.tmpl", "base.fastpath.notmono.generated.go")
codec.GenTmplRun2Go("mammoth_test.go.tmpl", "mammoth_generated_test.go")
codec.GenMonoAll()
}
EOF

    # explicitly return 0 if this passes, else return 1
    ${gocmd} run -tags "$btags" gen-from-tmpl.generated.go || return 1
    rm -f gen-from-tmpl*.generated.go
    return 0
}

_prebuild() {
    local d="$PWD"
    local zfin="test_values.generated.go"
    local zfin2="test_values_flex.generated.go"
    local zpkg="github.com/ugorji/go/codec"
    local returncode=1

    # zpkg=${d##*/src/}
    # zgobase=${d%%/src/*}
    # rm -f *_generated_test.go 
    # if [[ $zforce ]]; then ${gocmd} install ${zargs[*]} .; fi &&
    true &&
        _build &&
        cp $d/values_test.go $d/$zfin &&
        cp $d/values_flex_test.go $d/$zfin2 &&
        if [[ "$(type -t _codegenerators_external )" = "function" ]]; then _codegenerators_external ; fi &&
        returncode=0 &&
        echo "prebuild done successfully"
    rm -f $d/$zfin $d/$zfin2
    return $returncode
    # unset zfin zfin2 zpkg
}

_make() {
    _prebuild && ${gocmd} install ${zargs[*]} .
}

_clean() {
    rm -f \
       gen-from-tmpl.*generated.go \
       test_values.generated.go test_values_flex.generated.go
}

_tests_run_one() {
    local tt="alltests $i"
    local rr="TestCodecSuite"
    if [[ "x$i" == "xx" ]]; then tt="codec.notmono codec.notfastpath x"; rr='Test.*X$'; fi
    local g=( ${zargs[*]} ${ztestargs[*]} -count $nc -cpu $cpus -vet "$vet" -tags "$tt" -run "$rr" )
    [[ "$zcover" == "1" ]] && g+=( -cover )
    # g+=( -ti "$k" )
    g+=( -tdiff )
    [[ "$zcover" == "1" ]] && g+=( -test.gocoverdir $covdir )
    local -
    set -x
    ${gocmd} test "${g[@]}" &
}

_tests() {
    local vet="" # TODO: make it off
    local gover=$( ${gocmd} version | cut -f 3 -d ' ' )
    # go tool cover is not supported for gccgo, gollvm, other non-standard go compilers
    [[ $( ${gocmd} version ) == *"gccgo"* ]] && zcover=0
    [[ $( ${gocmd} version ) == *"gollvm"* ]] && zcover=0
    case $gover in
        go1.2[0-9]*|go2.*|devel*) true ;;
        *) return 1
    esac
    # we test the following permutations wnich all execute different code paths as below.
    echo "TestCodecSuite: (fastpath/unsafe), (!fastpath/unsafe), (fastpath/!unsafe), (!fastpath/!unsafe)"
    local nc=2 # count
    local cpus="1,$(nproc)"
    # if using the race detector, then set nc to
    if [[ " ${zargs[@]} " =~ "-race" ]]; then
        cpus="$(nproc)"
    fi
    local covdir=""
    local a=( "" "codec.safe" "codec.notfastpath" "codec.safe codec.notfastpath"
              "codec.notmono" "codec.notmono codec.safe"
              "codec.notmono codec.notfastpath" "codec.notmono codec.safe codec.notfastpath" )
    [[ "$zextra" == "1" ]] && a+=( "x" )
    [[ "$zcover" == "1" ]] && covdir=`mktemp -d`
    ${gocmd} vet -printfuncs "errorf" "$@" || return 1
    for i in "${a[@]}"; do
        local j=${i:-default}; j="${j// /-}"; j="${j//codec./}"
        [[ "$zwait" == "1" ]] && echo ">>>> TAGS: 'alltests $i'; RUN: 'TestCodecSuite'"
        _tests_run_one
        [[ "$zwait" == "1" ]] && wait
        # if [[ "$?" != 0 ]]; then return 1; fi
    done
    wait
    [[ "$zcover" == "1" ]] &&
        echo "go tool covdata output" &&
        ${gocmd} tool covdata percent -i $covdir &&
        ${gocmd} tool covdata textfmt -i $covdir -o __cov.out &&
        ${gocmd} tool cover -html=__cov.out
}

_usage() {
    # hidden args:
    # -pf [p=prebuild (f=force)]
    
    cat <<EOF
primary usage: $0
    -t[esow]   -> t=tests [e=extra, s=short, o=cover, w=wait]
    -[md]      -> [m=make, d=race detector]
    -v         -> v=verbose (more v's to increase verbose level)
EOF
    if [[ "$(type -t _usage_run)" = "function" ]]; then _usage_run ; fi
}

_main() {
    if [[ -z "$1" ]]; then _usage; return 1; fi
    local x # determines the main action to run in this build
    local zforce # force
    local zcover # generate cover profile and show in browser when done
    local zwait # run tests in sequence, not parallel ie wait for one to finish before starting another
    local zextra # means run extra (python based tests, etc) during testing
    
    local ztestargs=()
    local zargs=()
    local zverbose=()
    local zbenchflags=""

    local gocmd=${MYGOCMD:-go}
    
    OPTIND=1
    while getopts ":cetmnrgpfvldsowikxyz" flag
    do
        case "x$flag" in
            'xw') zwait=1 ;;
            'xv') zverbose+=(1) ;;
            'xo') zcover=1 ;;
            'xe') zextra=1 ;;
            'xf') zforce=1 ;;
            'xs') ztestargs+=("-short") ;;
            'xl') zargs+=("-gcflags"); zargs+=("-l=4") ;;
            'xn') zargs+=("-gcflags"); zargs+=("-m=2") ;;
            'xd') zargs+=("-race") ;;
            # 'xi') x='i'; zbenchflags=${OPTARG} ;;
            x\?) _usage; return 1 ;;
            *) x=$flag ;;
        esac
    done
    shift $((OPTIND-1))
    # echo ">>>> _main: extra args: $@"
    case "x$x" in
        'xt') _tests "$@" ;;
        'xm') _make "$@" ;;
        'xr') _release "$@" ;;
        'xg') _go ;;
        'xp') _prebuild "$@" ;;
        'xc') _clean "$@" ;;
    esac

    # handle from local run.sh
    case "x$x" in
        'xi') _check_inlining_one "$@" ;;
        'xk') _go_compiler_validation_suite ;;
        'xx') _analyze_checks "$@" ;;
        'xy') _analyze_debug_types "$@" ;;
        'xz') _analyze_do_inlining_and_more "$@" ;;
    esac
    # unset zforce zargs zbenchflags
}

[ "." = `dirname $0` ] && _main "$@"

# _xtrace() {
#     local -
#     set -x
#     "${@}"
# }
