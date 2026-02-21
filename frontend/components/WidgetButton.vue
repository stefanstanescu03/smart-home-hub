<script>
import axios from "axios";
export default {
  props: ["deviceId", "deviceName", "label", "payload", "widgetID"],
  data() {
    return {
      active: false,
    };
  },
  methods: {
    getToken() {
      const cookies = document.cookie;
      let token = cookies.split("=")[1];
      if (token === undefined) {
        return undefined;
      }
      return token;
    },
    async isActive() {
      try {
        await axios
          .get(`/api/device/ping/${this.deviceId}`)
          .then((response) => {
            if (response.data.status == true) {
              this.active = true;
            } else {
              this.status = false;
            }
          });
      } catch (err) {
        this.status = false;
      }
    },
    async handleClick() {
      try {
        await axios.post(
          "/api/widget/cmd",
          {
            Widget: this.widgetID,
            DeviceId: this.deviceId,
            Payload: this.payload,
          },
          { headers: { Authorization: `Bearer ${this.getToken()}` } },
        );
      } catch (err) {
        console.log(err);
      }
    },
  },
  async mounted() {
    await this.isActive();
  },
};
</script>

<template>
  <div
    class="smart-tile"
    :class="{ 'is-active': active, 'is-disabled': !active }"
    @click="handleClick"
  >
    <div class="title-header">
      <span class="device-name">{{ deviceName }}</span>
      <button
        class="delete-btn"
        @click.stop="$emit('delete')"
        title="Delete Widget"
      >
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

    <div class="tile-body">
      <h2 class="label">{{ label }}</h2>
    </div>

    <div class="tile-footer">
      <div class="status-row">
        <div class="indicator"></div>
        <span class="status-text">
          {{ active ? "Online" : "Offline" }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.smart-tile {
  width: 8rem;
  height: 8rem;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 24px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  user-select: none;
}

.is-disabled {
  opacity: 0.5;
  filter: grayscale(0.8);
  border-style: dashed;
}

.is-disabled:hover {
  transform: none;
  background: #1a1a1a;
}

.is-active .indicator {
  background: #2ba618;
}

.is-active .status-text {
  color: #2ba618;
}

.smart-tile:active:not(.is-disabled) {
  transform: scale(0.95);
  background: #252525 !important; /* The color you wanted */
  border-color: #555; /* Slight border highlight */
  box-shadow: inset 0 2px 8px rgba(0, 0, 0, 0.4); /* Adds 3D "pushed in" depth */
}

.title-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.delete-btn {
  background: #303030;
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 50%;
  border: none;
  color: #666;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 5;
}

.delete-btn svg {
  width: 0.75rem;
  height: 0.75rem;
}
.delete-btn:hover {
  background: #ff4d4d;
  color: white;
}

.status-row {
  display: flex;
  align-items: center;
  gap: 6px;
}

.indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #444;
}

.device-name {
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
  color: rgba(255, 255, 255, 0.3);
}

.label {
  font-size: 1.1rem;
  font-weight: 500;
  color: #ffffff;
  margin: 0;
}

.status-text {
  font-size: 0.65rem;
  color: rgba(255, 255, 255, 0.3);
}

.smart-tile:hover:not(.is-disabled) {
  background: #1d1d1d;
  border-color: #444;
  transform: translateY(-4px);
}
</style>
