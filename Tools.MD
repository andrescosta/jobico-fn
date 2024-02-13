# Tools

## Jobico Command Line Tool

### Overview:

The **Jobico Command Line Tool** is a management interface designed to facilitate the deployment, rollback, or redeployment of jobs within the Jobico platform. It offers a range of commands for uploading WebAssembly (WASM) and schema files, as well as streaming information from the Executions Recorder. This tool is an indispensable component for administrators and tenants alike, providing a streamlined and efficient way to interact with the platform.

### Key Features:

1. **Deployment Commands:**
   - **Deploy:**
     - The `deploy` command is employed to add a job definition to the Jobico platform.If the `-update` flag is provided and the job has already been deployed, the command will redeploy it.
 
     ```bash
     jobico deploy [-update] my-job-definition.yaml
     ```

   - **Rollback:**
     - The `rollback` command allows for the rollback of a deployed job to a previous state. It's a useful feature for reverting to a stable configuration in case of issues. The executor associated with the JOB are stopped.

     ```bash
     jobico rollback my-job-definition.yaml
     ```

2. **File Upload Commands:**
   - **Upload WASM:**
     - The `upload wasm` command enables the upload of a WebAssembly file to the Job Repository. This file can then be referenced by Jobicolets during deployment.

     ```bash
     jobico upload wasm <tenant id> <file id> <my-job-logic.wasm>
     ```

   - **Upload Schema:**
     - The `upload json` command allows the upload of JSON schema files to the Job Repository. These files define the structure of events processed by the platform.

     ```bash
     jobico upload json <tenant id> <file id> <my-job-logic.json>
     ```

3. **Streaming Information:**
   - **Stream from Recorder:**
     - The `recorder` command allows users to stream information from the Executions Recorder. This feature is valuable for real-time monitoring and analysis of job executions. Using the '-lines <NM>' flag outputs the last NM lines produced by the Jobs for the latest executions.

     ```bash
     jobico recorder [-lines NM]
     ```
4. **Information Commands:**
   - **Deployments**
     -  The `show deploy` command prints information about a Job Definition deployed previously. It offers details on the configuration, queues, runtims, and associated schema of a deployed job. 

     ```bash
     jobico show deploy <tenant id> <definition id>
     ```
   - **Environment(experimental)** 
     -  The `show env` command prints information about the nodes that composed a Jobico's cluster. This information provided is currently not used by the platform at this moment. 

     ```bash
     jobico show env
     ```

## Dashboard - Terminal GUI

### Overview:

The **Jobico Dashboard** is a terminal-based graphical user interface (GUI) designed to offer an interactive and visual representation of the Jobico platform. This GUI allows users to seamlessly visualize deployed jobs, explore files in the repository, and stream real-time results produced by executed jobs, all within the convenience of the terminal.

### Key Features:

1. **Visualizing Deployed Jobs:**
   - The Dashboard GUI presents a visual overview of deployed jobs, displaying relevant details such as job names, configurations, and status indicators. Users can easily navigate and interact with job-related information.

2. **Browsing Repository Files:**
   - Users can explore files stored in the Job Repository through an intuitive graphical interface. This includes the ability to inspect WebAssembly (WASM) files, JSON schema definitions, and other artifacts crucial for job execution.

3. **Streaming Job Results:**
   - The Dashboard supports real-time streaming of results produced by executed jobs. The GUI provides a dynamic display of outcomes, offering users immediate visibility into the status and performance of their jobs.

#### Launching the Dashboard:

```bash
dashboard [-debug] [-sync]
```

Executing this command launches the Dashboard GUI, initiating an interactive environment for users to visually explore deployed jobs and related information.

### Screenshots

1. **Job definitions**
   
![alt](docs/screenshots/gen-jobd.png?)

2. **Schema**
   
![alt](docs/screenshots/gen-schema.png?)

3. **Recorder output**
   
![alt](docs/screenshots/gen-recorder.png?)