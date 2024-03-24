# Netmon

Network Monitor (Netmon)

## Getting Started

### Backend

Build the application,

```
go build
```

Start the collector (to collect network packets) and visualise the packet information in CLI.

```
sudo ./netmon watch packets wlp0s20f3
```

Replace wlp0s20f3 with your system's interface.

Otherwise, if you want to visualise a web page,

```
sudo ./netmon watch packets wlp0s20f3 web
```

Note: Remember to start the frontend if you use the above command.

### Frontend

```
cd web/netmon-ui
npm run dev
```

After executing the above commands, go to port 5173.



