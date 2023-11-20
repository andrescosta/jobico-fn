$file=$args[0]
$name = (Get-Item $file ).Basename 
$name = $name + ".wasm"
tinygo build -o ./target/$name -target=wasi $file
