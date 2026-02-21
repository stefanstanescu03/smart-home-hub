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
      <button class="delete-btn" @click="$emit('delete')" title="Delete Widget">
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="3"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
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
  gap: 1rem;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #333;
}

.header-title {
  font-size: 1rem;
  font-weight: 600;
  color: #eeeeee;
}

.delete-btn {
  background: #303030;
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 50%;
  border: none;
  color: #666;
  cursor: pointer;
  padding: 0;
  margin: 0;

  display: flex;
  align-items: center;
  justify-content: center;

  backface-visibility: hidden;
  transform: translateZ(0);

  transition: all 0.2s ease;
}

.delete-btn svg {
  width: 0.75rem;
  height: 0.75rem;
}

.delete-btn:hover {
  background: #ff4d4d;
  color: white;
  transform: scale(1.1);
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
