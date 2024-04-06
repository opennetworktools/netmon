<template>
  <main>
    <RealTimeChart :packets="packets"/>
    <div class="wrapper">
      <div class="container-packets">
        <div>Timestamp, SrcIP, SrcPort, DstIP, DstPort</div>
          <div class="container">
            <ul>
              <li v-for="(packet, index) in packets" :key="index">{{ prettyPrint(packet) }}</li>
            </ul>
          </div>
      </div>
      <Hosts />
    </div>
  </main>
</template>

<script>
import { watch, ref, onMounted, computed } from 'vue'
import RealTimeChart from '@/components/RealTimeChart.vue'
import Hosts from '@/components/Hosts.vue'
import dayjs from 'dayjs'
import { v4 as uuidv4 } from 'uuid'

// import { usePacketsStore } from '@/stores/packets'
// const store = usePacketsStore()
// console.log(store.packets)

export default {
  components: {
    RealTimeChart,
    Hosts
  },
  setup() {
    let packets = ref([])

    function setupEventSource() {
      const eventSource = new EventSource('http://localhost:4444/packets');

      eventSource.onmessage = (event) => {
        let packet = JSON.parse(event.data)
        packet.ID = uuidv4()
        packets.value.push(packet)
      };

      eventSource.onerror = (error) => {
        console.error('EventSource failed:', error);
        eventSource.close();
      };
    }

    function prettyPrint(packet) {
      let formattedTimestamp = dayjs(packet.Timestamp).format('DD-MM-YYYY hh:mm:ss');
      return `${formattedTimestamp}, ${packet.SrcAddress.IP}, ${packet.SrcAddress.PORT}, ${packet.DstAddress.IP}, ${packet.DstAddress.PORT}, ${packet.Protocol}`
    }

    onMounted(() => {
      setupEventSource()
    })

    const packetsForChart = computed(() => packets.value.slice()) // Make a reactive computed property

    watch(packets, (currentValue, oldValue) => {
      console.log(currentValue);
    });

    return {
      packets: packetsForChart,
      prettyPrint
    }
  }
}
</script>

<style scoped>

main {
  padding: 0 2rem;
}

.wrapper {
  display: flex;
  gap: 2rem;
}

.container-packets {
  margin-top: 1rem;
}

.container {
  margin-top: 0.5rem;
  max-height: 400px;
  overflow-x: hidden;
  overflow-y: auto;
}
</style>
