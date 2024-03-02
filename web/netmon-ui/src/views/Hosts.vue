<template>
    <main>
      <div>Hosts</div>
      <ul>
        <li v-for="(hostsObj, index) in hosts" :key="index">{{ hostsObj }}</li>
      </ul>
    </main>
</template>
  
<script>
  export default {
    data() {
      return {
          hosts: []
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
      }
    }
  };
  </script>
  
  
  