import { defineStore } from 'pinia'

export const useCounterStore = defineStore('counter', {
  state: () => ({ n: 0 }),
  actions: {
    inc() {
      this.n++
    },
  },
})
