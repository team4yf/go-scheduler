## Go-Scheduler

> A  Job scheduler.

### Content
- Introduce
- Features
- TODO
- API List
- Deploy Manual

#### Introduce

it's build with `yf-fpm-server`, base the postgresql/mysql, send http request accroding to the defined CRON.

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

- [x] GET `/biz/job/list`
  
  fetch all job list.

- [x] GET `/biz/job/execute?code={code}`

  run once by the specific job code.
- [x] POST `/biz/job/add`

  create a job.
- [x] POST `/biz/job/update`

  update the job, basiclly use it to define the job's status. auto start or not.
- [x] GET `/biz/job/remove?code={code}`

  remove the job. it can not be restart.
- [x] GET `/biz/job/get?code={code}`

  get the job detail.
- [x] GET `/biz/job/pause?code={code}`

  pause the job, it can be restart.
- [x] GET `/biz/job/tasks?code={code}&skip={skip}&limit={limit}`

  get the tasks of the job.

#### Deploy Manual

- defined the config.*.yaml

    ```yaml
    mode: debug
    addr: ':8080'
    name: go-scheduler

    cron:
        store: db
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