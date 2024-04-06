import { defineStore } from 'pinia'

export const usePacketsStore = defineStore('packets', {
    state: () => {
        packets: []
    },
    actions: {
        add(packet) {
          this.packets.append(packet)
        },
    },
})