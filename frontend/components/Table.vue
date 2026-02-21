<script>
export default {
  props: ["deviceId", "deviceName", "streamData"],
  data() {
    return {
      valuesArray: [],
    };
  },
  watch: {
    streamData: {
      immediate: true,
      deep: true,
      handler(message) {
        if (!message) return;

        const parts = message.split(",");
        let value = {};

        parts.forEach((p) => {
          const [k, v] = p.split(":");
          if (k !== "id") {
            value[k] = v;
          }
        });

        this.valuesArray.push(value);

        if (this.valuesArray.length > 100) {
          this.valuesArray.shift();
        }
      },
    },
  },
};
</script>

<template>
  <div class="stream-card">
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

    <div class="table-scroll-area">
      <table v-if="valuesArray.length">
        <thead>
          <tr>
            <th v-for="(value, key) in valuesArray[0]" :key="key">
              {{ key }}
            </th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="(row, index) in valuesArray.slice().reverse()"
            :key="index"
          >
            <td v-for="(cell, key) in row" :key="key">
              {{ cell }}
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else class="waiting-state">Waiting for data...</div>
    </div>
  </div>
</template>

<style scoped>
.stream-card {
  background-color: #1a1a1a;
  border: 1px solid #333;
  width: 100%;
  display: flex;
  flex-direction: column;
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

.table-scroll-area {
  max-height: 300px;
  overflow-y: auto;
  background-color: #1a1a1a;
}

table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

thead {
  position: sticky;
  top: 0;
  z-index: 10;
  background-color: #252525;
}

th {
  padding: 0.75rem 1rem;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  color: #888;
  border-bottom: 1px solid #333;
}

td {
  padding: 0.7rem 1rem;
  font-size: 0.95rem;
  color: #eeeeee;
  border-bottom: 1px solid #252525;
  font-family: monospace;
}

tbody tr:hover {
  background-color: #222222;
}

.waiting-state {
  padding: 3rem 1rem;
  text-align: center;
  color: #555;
  font-size: 0.9rem;
}

.table-scroll-area::-webkit-scrollbar {
  width: 6px;
}
.table-scroll-area::-webkit-scrollbar-track {
  background: #1a1a1a;
}
.table-scroll-area::-webkit-scrollbar-thumb {
  background: #333;
}
</style>
