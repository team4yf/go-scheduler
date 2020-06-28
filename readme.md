## Go-Scheduler

> A  Job scheduler.

### Content
- Introduce
- Features
- TODO
- API List
- Deploy Manual

#### Introduce

it's build with `gin`, base the postgresql/mysql, send http request accroding to the defined CRON.

#### Features

- [x] add job
- [x] list jobs
- [x] run once
- [x] start/puase the single job
- [x] notify the result by topics.
- [x] save all the logs for each job
- [x] see the job execute report

#### TODO

- [ ] use redis/etcd to make it support cluster distribute.
- [ ] use MQ to support more notify methods.
- [x] support notify by email.

#### API List

- [x] GET `/api/v1/job/list`
  
  fetch all job list.
- [x] GET `/api/v1/job/get/:code`

  get the specific job detail by the code.
- [x] GET `/api/v1/job/execute/:code`

  run once by the specific job code.
- [x] POST `/api/v1/job/create`

  create a job.
- [x] POST `/api/v1/job/update`

  update the job, basiclly use it to define the job's status. auto start or not.
- [x] GET `/api/v1/task/list/:code?p=1&l=10`

  get the specific job's tasks.

- [x] GET `/api/v1/task/export/:code?p=1&l=10`

  export the tasks log by the specific job.

- [x] GET `/api/v1/task/detail/:id`

  get the detail infomation of the taskid

- [x] GET `/api/v1/task/report/:code`

  get the tasks report by the specific job, include the succes/fail/total, last execute time & last execute result.
- [x] POST `/api/v1/subscribe/sub/:code`

  subscribe the job's notification after the task finished; the `envent` should be `success`, `error`, `timeout`, `all`.

- [x] POST `/api/v1/subscribe/unSub/:code`

  unSubscribe the topic.

- [x] GET `/api/v1/subscribe/list/:code`

  list all the subscriber of the specific job by code.

  

#### Deploy Manual

- defined the config.*.yaml

    ```yaml
    mode: debug
    addr: ':8080'
    name: go-scheduler
    logger: 
        file: api-access.log
        level: info

    db:
        engine: postgres
        user: devuser
        password: DevPass123
        host: localhost
        port: 5432
        database: devdb
        charset: utf8
        showsql: true
    ```

- define the enviroment `GS_DEPLOY_MODE=PROD`

  it will load the config file `conf/config.prod.yaml`, 
  
  `config.local.yaml` for default.

- define the email config

  - GS_MODE
  - GS_EMAIL_HOST 
  - GS_EMAIL_PORT
  - GS_EMAIL_USERNAME
  - GS_EMAIL_PASSWORD
  - GS_EMAIL_NAME
  - GS_EMAIL_ADDRESS
  - GS_EMAIL_REPLY

  - GS_LOG_LEVEL
  
  - GS_DB_USER
  - GS_DB_PASSWORD
  - GS_DB_HOST
  - GS_DB_PORT
  - GS_DB_DATABASE
  - GS_DB_SHOWSQL

