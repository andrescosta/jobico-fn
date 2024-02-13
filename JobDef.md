# Job Definition

"Job Definition" provides a comprehensive specification within the Jobico framework, outlining how Events are validated, and processed by a Job, and how the results are managed. Utilizing YAML for this purpose ensures a structured and human-readable format, allowing clear articulation of rules and conditions governing event validation, processing, and result management by Jobs.

### Anatomy of a Job Definition YAML:
A Job Definition YAML file includes various attributes that define the job's characteristics:

#### `name`

- **Description:** A friendly name for the job, used for identification and management.

- **Example:**

```yaml
name: Customer Event Processing Definitions 
```

#### `id`

- **Description:** The "id" attribute represents the unique identifier for the job. It serves as a distinct reference to identify and manage the job within the Jobico platform.

- **Example:**

  ```yaml
  id: customer-proc-jobs
  ```

#### `tenant`

- **Description:** The "tenant" attribute represents the ID of the tenant associated with the job. It ensures that the job is attributed to a specific tenant within the multi-tenancy architecture of Jobico.

- **Example:**

  ```yaml
  tenant: pritty-tenant
  ```

#### `queues`

- **Description:** The "queues" section describes the queues associated with the job. This section allows for future expansion where queue environments and capabilities can be defined.

  - `queues.id`: ID of the queue.
  - `queues.name`: Friendly name of the queue.

- **Example:**

  ```yaml
     queues:
       - id: default-queue
         name: Default to all events
  ```

#### `jobs`

- **Description:** The "jobs" section is where the jobs and events are defined and how they will be validated and processed.

  - `jobs.event`: An event definition.
  - `jobs.event.name`: Friendly name for the event.
  - `jobs.event.id`: ID for the event, used by the REST API and executors to determine the schema and WASM file.
  - `jobs.event.datatype`: Specifies the data type of the event. "0" represents JSON.

  - **`jobs.event.schema`: Schema file definition:**

    - `jobs.event.schema.id`: ID of the schema file.
    - `jobs.event.schema.name`: Name of the schema file.
    - `jobs.event.schema.schemaref`: Reference used to retrieve the file from the repository.

  - `jobs.event.supplierqueue`: Specifies the ID of the queue where this event will be published.
  - `jobs.event.runtime`: ID of the runtime that will process this event.
  - `jobs.event.result`: Specifies how the result of the execution will be treated (Under Construction).

- **Example:**

  ```yaml
  jobs:
    - event:
        name: New customer
        id: customer-registration
        datatype: 0
        schema:
           id: customer-registration-schema
           name: Customer registration schema
           schemaref: customer-registration-schema.json
        supplierqueue: 1
        runtime: 1
  ```
 
#### `runtimes`

- **Description:** The "runtimes" section specifies the runtimes available to process the events.

  - `runtimes.id`: ID of the runtime, used to reference a specific runtime.
  - `runtimes.name`: Friendly name.
  - `runtimes.moduleref`: Reference used to retrieve the file from the repository.
  - `runtimes.mainfuncname`: Future usage.
  - `runtimes.type`: "0" represents WASM as the runtime type.

- **Example:**

  ```yaml
  runtimes:
    - id: wasm-runtime-customer-ev
      name: Wasm runtime for Customer events
      moduleref: wasm-runtime-customer-ev.wasm
      mainfuncname: event
      type: 0
  ```

  In this example, a runtime named "wasm-runtime-customer-ev" is defined with the associated WASM file and runtime type.

These attributes collectively form a comprehensive YAML file, capturing the essential details for defining and deploying jobs within the Jobico platform. 

### Example

```yaml
name: Customer Event Processing Definitions
id: customer-proc-jobs
tenant: pritty-tenant
queues:
  - id: queue-default
    name: Default to all events
jobs:
  - event:
      name: New customer
      id: customer-registration
      datatype: 0
      schema:
        id: customer-registration-schema
        name: Customer registration schema
        schemaref: customer-registration-schema.json
      supplierqueue: queue-default
      runtime: wasm-runtime-customer-ev
runtimes:
  - id: wasm-runtime-customer-ev
    name:  Wasm runtime for Customer events
    moduleref: wasm-runtime-customer-ev.wasm
    mainfuncname: event
    type: 0
```


