# Jobicolet

## What is a Jobicolet?

A **Jobicolet** is a specialized WebAssembly (WASM) program designed to process an event and generate a result within the Jobico platform. It represents the executable logic that is dynamically loaded and executed by the Job Executors when handling specific 

## Key Characteristics:

1. **WASM Execution:**
   - A Jobicolet is implemented as a WebAssembly module, allowing it to be written in any programming language that compiles to WebAssembly. This flexibility empowers users to express their event processing logic in a language of their choice.

2. **Event Processing:**
   - The primary function of a Jobicolet is to process events. It takes as input the event data, performs the specified logic defined within the WASM module, and produces a result based on the defined processing rules.

3. **Result Generation:**
   - Upon processing an event, a Jobicolet generates a result. The nature of the result depends on the specific logic implemented in the WASM module. It could be a computation outcome, a transformed dataset, or any other relevant output.

4. **Language Agnostic:**
   - Jobicolets are language-agnostic in the sense that they can be written in any programming language that supports compilation to WebAssembly. This feature provides developers with the freedom to choose a language that aligns with their expertise and the requirements of their event processing tasks.

## Benefits:

- **Flexibility:**
  - The language-agnostic nature of Jobicolets provides developers with flexibility, allowing them to choose the most suitable programming language for expressing their event processing logic.

- **Scalability:**
  - As Jobicolets are executed within the scalable and isolated environment of Job Executors, the platform can efficiently scale to handle a large number of concurrent event processing tasks.

- **Interoperability:**
  - Jobicolets can interact with other components within the Jobico platform, facilitating seamless integration with queues, event definitions, and runtime environments.

A Jobicolet, at its core, represents the embodiment of programmable and scalable event processing within the Jobico platform, offering developers the freedom to innovate using the power of WebAssembly.

## Getting Started

### Docker

#### Prerequisites:

1. Ensure Docker are installed on your local machine.

#### Steps:

1. Clone the jobicolet-examples repository from GitHub using the following command:

     ```bash
     git clone https://github.com/andrescosta/jobico
     ```

2. **Navigate to the `/compose` directory:**

    ```bash
    cd compose
    ```

2. **Run Docker Compose:**

    ```bash
    docker compose up
    ```

    This command will start the Docker Compose stack based on the configuration defined in the `compose.yml` file.

3. **Verify Listener on Port 8080:**

    Once the Docker Compose stack is up and running, you can verify that the Listener is running on port 8080 by making a request. You can use a tool like `curl` or a web browser.

    - Using `curl`:

        ```bash
        curl http://localhost:8080
        ```

    - If using a web browser, navigate to `http://localhost:8080` in your browser.

4. **Stop and Cleanup:**

    When you are done testing, you can stop the Docker Compose stack using:

    ```bash
    docker compose down
    ```

    This will stop and remove the containers defined in the `compose.yml` file.

### Open Telemetry


A Docker Compose file with the OpenTelemetry stack enabled is provided. You can initiate it by executing the following command:

```bash
docker compose -f compose-otel.yml --profile obs up
```

To shut it down, use the following command:

```bash
docker compose -f compose-otel.yml --profile obs down
```

The Prometheus console is reachable at: http://localhost:9090/, while the Jaeger console can be accessed at: http://localhost:16686/search

### Rust

#### Prerequisites:

1. **Install Rust:**
   - Download and install Rust from [https://www.rust-lang.org/](https://www.rust-lang.org/).

2. **Install wasm-pack:**
   - After installing Rust, install `wasm-pack` by running the following command:

     ```bash
     cargo install wasm-pack
     ```

3. **Clone the Examples Repository:**
   - Clone the jobicolet-examples repository from GitHub using the following command:

     ```bash
     git clone https://github.com/andrescosta/jobicolet-examples.git
     ```

#### Build the Rust Example:

1. **Navigate to the Rust Example Directory:**
   - Change your working directory to the location of the Rust example in the jobicolet-examples repository:

     ```bash
     cd jobicolet-examples/rust/greet
     ```

2. **Compile the Example using Cargo:**
   - Use the following command to compile the Rust program (`greet.rs`) using Cargo:

     ```bash
     cargo build --release --target wasm32-unknown-unknown
     ```

   This command instructs Cargo to build the Rust program in release mode (`--release`) for the WebAssembly target (`--target wasm32-unknown-unknown`).

3. **Verify the Output:**
   - After a successful build, you should find the compiled WebAssembly module in the `target/wasm32-unknown-unknown/release/` directory. The file will be named `greet.wasm`.

4. **Upload the WASM file to Jobico:**
   - Following the compilation of the file, it is imperative to upload it to Jobico by executing the following command:

    ```bash
     cli upload wasm demorust greet-wasm-rust.wasm target\wasm32-unknown-unknown\release\greet.wasm
     ```
5. **Upload the schema file to Jobico:**
   - Executing this command will upload the schema file, facilitating the validation of the associated event:

    ```bash
     cli upload json demorust  greet-schema-rust.json schema.json
     ```
6. **Deploy the job:**
   - Executing this command will initiate the deployment of the Job:

    ```bash
     cli deploy job-rust-greet.yml
     ```
7. **Start streaming results:**
   - Executing this command will initiate the streaming of results from the Recorder component:

    ```bash
     cli recorder
     ```
8. **Send an event:**
   - Executing this command will dispatch an event to Jobico:

    ```bash
          curl --request POST \ 
          --url http://localhost:8080/events/demogo/evgo \  
          --header 'content-type: application/json' \
          --data '{"data": [{"firstName": "Andres","lastName": "C"}]}'
     ```
   - Return to the terminal where the results are currently being streamed and review the log.

### Tinygo

#### Prerequisites:

1. **Install TinyGo:**
   - Download and install TinyGo from [https://tinygo.org/](https://tinygo.org/).

2. **Clone the Examples Repository:**
   - Clone the jobicolet-examples repository from GitHub using the following command:

     ```bash
     git clone https://github.com/andrescosta/jobicolet-examples.git
     ```

3. ** The Jobico platform is up and running **
   - The Docker section taches you how to start Jobico

#### Build the Go Example:

1. **Navigate to the Go Example Directory:**
   - Change your working directory to the location of the Go example in the jobicolet-examples repository:

     ```bash
     cd jobicolet-examples/go/greet
     ```

2. **Compile the Example using TinyGo:**
   - Use the following command to compile the greet.go example using TinyGo:

     ```bash
     tinygo build -scheduler=none --no-debug -target=wasi greet.go
     ```

   This command instructs TinyGo to build the Go program (`greet.go`) for the WebAssembly System Interface (WASI) target.

3. **Verify the Output:**
   - After a successful build, you should see an executable file named `greet.wasm` in the same directory.

4. **Upload the WASM file to Jobico:**
   - Following the compilation of the file, it is imperative to upload it to Jobico by executing the following command:

    ```bash
     cli upload wasm demogo greet-wasm-go.wasm greet.wasm
     ```
5. **Upload the schema file to Jobico:**
   - Executing this command will upload the schema file, facilitating the validation of the associated event:

    ```bash
     cli upload json demogo greet-schema-go.json schema.json
     ```
6. **Deploy the job:**
   - Executing this command will initiate the deployment of the Job:

    ```bash
     cli deploy job-go-greet.yml
     ```
7. **Start streaming results:**
   - Executing this command will initiate the streaming of results from the Recorder component:

    ```bash
     cli recorder
     ```
8. **Send an event:**
   - Executing this command will dispatch an event to Jobico:

    ```bash
          curl --request POST \ 
          --url http://localhost:8080/events/demogo/evgo \  
          --header 'content-type: application/json' \
          --data '{"data": [{"firstName": "Andres","lastName": "C"}]}'
     ```
   - Return to the terminal where the results are currently being streamed and review the log.