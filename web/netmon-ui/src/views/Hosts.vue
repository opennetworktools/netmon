<template>
    <main>
      <h1>Hosts</h1>
      <div class="container">
        <ul>
          <li v-for="(hostsObj, index) in hosts" :key="index">{{ aggregate(hostsObj) }}</li>
        </ul>
      </div>
    </main>
</template>
  
<script>
  export default {
    data() {
      return {
          hosts: [],
          aggregation: {}
      };
    },
    created() {
      this.setupEventSource();
    },
    methods: {
      setupEventSource() {
        const eventSource = new EventSource('http://localhost:4444/hosts');
  
        eventSource.onmessage = (event) => {
          this.hosts.push(JSON.parse(event.data));
        };
  
        eventSource.onerror = (error) => {
          console.error('EventSource failed:', error);
          eventSource.close();
        };
      },
      aggregate(host) {
        // { "IP": "65.0.200.43", 
        // "HostName": "ec2-65-0-200-43.ap-south-1.compute.amazonaws.com.", 
        // "HostNames": [ "ec2-65-0-200-43.ap-south-1.compute.amazonaws.com." ], 
        // "ASNumber": 16509, 
        // "ASName": "AMAZON-02", 
        // "Bytes": 120 }
        return `${host.IP}, ${host.HostName}, ${host.ASNumber}, ${host.ASName}, ${host.Bytes}`
      }
    }
  };
  </script>
  
  <style scoped>
  .container {
    margin-top: 1rem;
  }
  
  ul {
    list-style: none;
  }
  </style>