Python2でAWSLambdaからSlack通知をするための最低限のコード---2017-06-19 10:22:59

検証のために、Lambdaから最低限のコードでSlack通知を行うようにした。

### コード

```
#-*- coding:utf-8 -*-
from urllib import urlencode
from urllib2 import Request, urlopen
import json

SLACK_POST_URL = "https://hooks.slack.com/services/XXXXXXXXXXXXXXXXXXXXXXXXXXXX"

def lambda_handler(event, context):
    post()

def post():
    slack_message = {
        ''channel'': ''test_channel'',
        ''text'': ''test_message'',
        ''username'': ''test_username'',
        ''icon_emoji'': '':smile:''
    }
    req = Request(SLACK_POST_URL, json.dumps(slack_message))
    response = urlopen(req)
```

これで、SLACK_POST_URL と slack_message だけいい感じに変えれば動く。
