# code-exec
An Online Code Execution API written in Golang, using echo as web framework, isolate as sandbox and Google Cloud Storage as permanent data storage.


## Configuration
* `ADDRESS`: The address that api server listens on, defaults to `:8000`
* `BODY_LIMIT`: The maximum size of the request body, defaults to `4M`
* `MAX_SANDBOX`: Maximum amount of sandbox isolate can create, defaults to `1000`

* `SUBMISSION_BUCKET`: The storage bucket to store submissions
* `TASK_BUCKET`: The storage bucket to store tasks

* `MAX_TASK`: The maximum amount of tasks a submission can have
* `MAX_TIME`: Maximum time limit, share by all the process / thread
* `MAX_MEMORY`: Maximum memory limit, share by all the process / thread
* `MAX_PROCESS`: Max process / thread limit
* `MAX_FILESIZE`: Max filesize limit