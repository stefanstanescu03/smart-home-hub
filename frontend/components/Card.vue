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
          if (k !== "timestamp" && k !== "id") {
            this.values[k] = v;
          }
        });
      },
    },
  },
};
</script>

<template>
  <div class="device-card">
    <div class="card-header">
      <div class="title">
        <h3>{{ deviceName }}</h3>
      </div>

      <button
        @click="$emit('delete')"
        class="icon-button"
        title="Delete device"
      >
        ✕
      </button>
    </div>

    <div class="card-body">
      <div v-for="(val, key) in values" :key="key" class="value-row">
        <span class="key">{{ key }}</span>
        <span class="value">{{ val }}</span>
      </div>

      <div v-if="Object.keys(values).length === 0" class="empty">
        Waiting for data…
      </div>
    </div>
  </div>
</template>

<style scoped>
.device-card {
  background: #1b1b1b;
  border: 1px solid #2a2a2a;
  border-radius: 1rem;
  overflow: hidden;
  width: 100%;
  max-width: 420px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.9rem 1rem;
  border-bottom: 1px solid #2a2a2a;
}

.title h3 {
  margin: 0;
  font-size: large;
  font-weight: bold;
  color: #eeeeee;
}

.icon-button {
  background: transparent;
  border: 1px solid #333;
  color: #aaa;
  border-radius: 0.4rem;
  cursor: pointer;
  padding: 0.25rem 0.45rem;
  font-size: 0.8rem;

  transition:
    color 150ms ease,
    border-color 150ms ease,
    background 150ms ease;
}

.icon-button:hover {
  color: #e55353;
  border-color: #e55353;
  background: rgba(229, 83, 83, 0.1);
}

.card-body {
  padding: 0.9rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.45rem;
}
.value-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  padding: 0.35rem 0.45rem;
  border-radius: 0.4rem;
  background: #202020;
}

.key {
  font-size: normal;
  text-transform: uppercase;
  color: #9aa0a6;
}

.value {
  font-size: normal;
  font-weight: 500;
  color: #eeeeee;
}

.empty {
  font-size: normal;
  color: #777;
  text-align: center;
  padding: 0.75rem 0;
}
</style>
