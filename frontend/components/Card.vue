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
    <div class="title-container">
      <h1 class="name-h">{{ deviceName }}</h1>
      <button @click="$emit('delete')" class="delete-button">Delete</button>
    </div>
    <div class="values-container">
      <div v-for="(val, key) in values" :key="key">{{ key }}: {{ val }}</div>
    </div>
  </div>
</template>

<style scoped>
.name-h {
  font-size: large;
}
.card-container {
  border: 1px solid #eeeeee;
  border-radius: 1rem;
}
.values-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 0.5rem;
}

.title-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
  align-items: center;
  border-bottom: 1px solid #eeeeee;
  padding: 0.5rem;
}

.delete-button {
  border: none;
  background-color: transparent;
  border: 1px solid #eeeeee;
  text-decoration: none;
  cursor: pointer;
  color: #eeeeee;
  padding: 0.5rem;
  border-radius: 0.3rem;
  transition-duration: 300ms;
}

.delete-button:hover {
  border: 1px solid #9e2323;
  color: #9e2323;
}
</style>
