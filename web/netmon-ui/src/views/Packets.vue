<template>
  <main>
    <h1>Packets</h1>
    <div>
      <ul>
        <li v-for="(packet, index) in packets" :key="index">{{ prettyPrint(packet) }}</li>
      </ul>
    </div>
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
    },
    prettyPrint(packet) {
      // { "SrcAddress": { "MAC": "74:d8:3e:b8:f0:4c", "IP": "192.168.29.239", "PORT": 60656 }, 
      // "DstAddress": { "MAC": "f0:ed:b8:f1:04:9b", "IP": "51.116.246.106", "PORT": 443 }, 
      // "Protocol": "TCP", 
      // "Timestamp": "2024-03-10T12:31:59.121113+05:30", 
      // "CaptureLength": 66 }
      return `${packet.Timestamp}, ${packet.SrcAddress.IP}, ${packet.DstAddress.IP}`
    }
  }
};
</script>

<style scoped>
ul {
  list-style: none;
}
</style>