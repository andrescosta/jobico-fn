$file=$args[0]
$name = (Get-Item $file ).Basename 
$name = $name + ".wasm"
tinygo build -o ../cmd/executor/target/$name -target=wasi $file
