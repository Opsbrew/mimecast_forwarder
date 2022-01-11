# mimecast_forwarder


###### Mimecast is a SaaS-based email management platform enabling companies to administer business communications and data.


Mimecast forwarder can forward your mimecast logs to Any syslog forwarders easily.

add a .env to your directory

- MM_APP_ID: ""
- MM_APP_KEY: ""
- MM_URI: "/api/audit/get-siem-logs"
- MM_EMAIL_ADDRESS: ""
- MM_ACCESS_KEY: ""
- MM_SECRET_KEY: ""
- REMOTE_SYSLOG_SERVER: ""
- PORT: "5003"
