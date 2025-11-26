<script>
export default {
  props: ["deviceId", "deviceName", "streamData"],
  data() {
    return {
      values: {},
    };
  },
  watch: {
    streamData: {
      immediate: true,
      deep: true,
      handler(message) {
        if (!message) return;
        const parts = message.split(",");
        parts.forEach((p) => {
          const [k, v] = p.split(":");
          if (k != "timestamp" && k != "id") {
            this.values[k] = v;
          }
        });
      },
    },
  },
};
</script>

<template>
  <div class="card-container">
    <h1 class="name-h">{{ deviceName }}</h1>
    <div v-for="(val, key) in values" :key="key">{{ key }}: {{ val }}</div>
  </div>
</template>

<style scoped>
.name-h {
  font-size: large;
}
.card-container {
  border: 1px solid #a6a6a6;
  border-radius: 1rem;
  padding: 0.5rem;
}
</style>
