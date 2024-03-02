<template>
  <main>
    <div>Packets</div>
    <ul>
      <li v-for="(packet, index) in packets" :key="index">{{ packet }}</li>
    </ul>
  </main>
</template>

<script>
export default {
  data() {
    return {
        packets: []
    };
  },
  created() {
    this.setupEventSource();
  },
  methods: {
    setupEventSource() {
      const eventSource = new EventSource('http://localhost:4444/packets');

      eventSource.onmessage = (event) => {
        this.packets.push(JSON.parse(event.data));
      };

      eventSource.onerror = (error) => {
        console.error('EventSource failed:', error);
        eventSource.close();
      };
    }
  }
};
</script>


