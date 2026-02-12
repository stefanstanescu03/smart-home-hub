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
      handler(message) {
        if (!message || typeof message !== "string") return;

        const newValues = { ...this.values };
        const parts = message.split(",");

        parts.forEach((p) => {
          const [k, v] = p.split(":");
          if (k && v && k !== "timestamp" && k !== "id") {
            newValues[k] = v;
          }
        });
        this.values = newValues;
      },
    },
  },
};
</script>

<template>
  <div class="device-card">
    <div class="card-header">
      <span class="header-title">{{ deviceName }}</span>
      <button @click="$emit('delete')" class="minimal-delete" title="Remove">
        &times;
      </button>
    </div>

    <div class="card-body">
      <div v-if="Object.keys(values).length === 0" class="waiting-state">
        Waiting for data...
      </div>

      <div v-for="(val, key) in values" :key="key" class="data-row">
        <span class="data-key">{{ key }}</span>
        <span class="data-value">{{ val }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.device-card {
  background-color: #1a1a1a;
  border: 1px solid #333;
  width: 100%;
  max-width: 400px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #333;
}

.header-title {
  font-size: 1rem;
  font-weight: 600;
  color: #eeeeee;
}

.minimal-delete {
  background: none;
  border: none;
  color: #666;
  font-size: 1.4rem;
  cursor: pointer;
  line-height: 1;
}

.minimal-delete:hover {
  color: #ff4444;
}

.card-body {
  padding: 0.5rem 0;
}

.data-row {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.8rem 1rem;
  border-bottom: 1px solid #252525;
}

.data-row:last-child {
  border-bottom: none;
}

.data-key {
  font-size: 0.9rem;
  color: #888;
  text-transform: capitalize;
}

.data-value {
  font-size: 1rem;
  color: #eeeeee;
  font-weight: 500;
  font-family: monospace;
}

.waiting-state {
  padding: 2rem;
  text-align: center;
  color: #555;
  font-size: 0.9rem;
}
</style>
