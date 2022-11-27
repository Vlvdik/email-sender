## Simple email sender with support for delayed sending, html templates and open message tracking
#### Also, service saves data about the tasks, which allows you to perform tasks later, even if the service fails
---
#### Mailing is based on data in JSON format, which are stored in the file ***receivers.json***
#### **Important**: This implementation of the service **does not** know how to check the existence of email addresses and **does not** delete non-existent addresses from the ***receivers.json*** file
---
#### ***receivers.json*** example
```json
[
  {
    "email": ["someEmail@gmail.com"],
    "personalInfo":{
      "Name":"Alex",
      "lastName":"Alexovich",
      "birthday":"2002.02.02"
    }
  },
  {
  "email": ["anotherOne@gmail.com"],
    "personalInfo":{
      "Name":"Ivan",
      "lastName":"Ivanov",
      "birthday":"2002.03.02"
    }
  }
]
```
---
### To start the service, you need to configure ***config.toml***
```toml
sender_email = "..."
password = "..."
host = "..."
port = "..."
```
