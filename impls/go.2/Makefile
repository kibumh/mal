all: step0_repl

common: core/core.go env/env.go printer/printer.go reader/reader.go types/types.go

all: step0_repl step1_read_print step2_eval step3_env step4_if_fn_do step5_tco step6_file step7_quote step8_macros step9_try stepA_mal

step0_repl: common step0/main.go
	go build -o build/step0_repl step0/main.go
step1_read_print: common step1/main.go
	go build -o build/step1_read_print step1/main.go
step2_eval: common step2/main.go
	go build -o build/step2_eval step2/main.go
step3_env: common step3/main.go
	go build -o build/step3_env step3/main.go
step4_if_fn_do: common step4/main.go
	go build -o build/step4_if_fn_do step4/main.go
step5_tco: common step5/main.go
	go build -o build/step5_tco step5/main.go
step6_file: common step6/main.go
	go build -o build/step6_file step6/main.go
step7_quote: common step7/main.go
	go build -o build/step7_quote step7/main.go
step8_macros: common step8/main.go
	go build -o build/step8_macros step8/main.go
step9_try: common step9/main.go
	go build -o build/step9_try step9/main.go
stepA_mal: common stepA/main.go
	go build -o build/stepA_mal stepA/main.go

