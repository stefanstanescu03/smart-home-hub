<script>
import axios from "axios";
export default {
  props: [
    "deviceId",
    "deviceName",
    "label",
    "widgetID",
    "payload",
    "payload2",
    "channel",
  ],
  data() {
    return {
      state: false,
      active: false,
    };
  },
  methods: {
    getToken() {
      const cookies = document.cookie;
      return cookies.split("=")[1] || undefined;
    },
    async isActive() {
      try {
        await axios
          .get(`/api/device/ping/${this.deviceId}`)
          .then((response) => {
            if (response.data.status == true) {
              this.active = true;
            } else {
              this.active = false;
            }
          });
      } catch (err) {
        this.active = false;
      }
    },
    async getState() {
      try {
        let channel_to_search =
          this.channel === "" ? "unregistered" : this.channel;
        const res = await axios.get(
          `/api/device/state/${this.deviceId}/${channel_to_search}`,
        );
        if (res.data.state === "not registered" || res.data.state === "false") {
          this.state = false;
        } else {
          this.state = true;
        }
      } catch (err) {
        console.log(err);
      }
    },
    async handleToggle() {
      if (this.active) {
        try {
          let payload_to_send = "";
          if (!this.state) {
            this.state = true;
            payload_to_send = this.payload;
          } else {
            this.state = false;
            payload_to_send = this.payload2;
          }

          await axios.post(
            "/api/widget/cmd",
            {
              Widget: this.widgetID,
              DeviceId: this.deviceId,
              Payload: payload_to_send,
            },
            { headers: { Authorization: `Bearer ${this.getToken()}` } },
          );
        } catch (err) {
          console.log(err);
        }
      }
    },
  },
  mounted() {
    this.isActive();
    this.getState();
  },
};
</script>

<template>
  <div
    class="smart-tile"
    :class="{ 'is-active': state, 'is-disabled': !active }"
    @click="handleToggle"
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
      <label class="accent-switch" @click.stop>
        <input type="checkbox" :checked="this.state" @change="handleToggle" />
        <span class="accent-slider"></span>
      </label>
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

.accent-switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 28px;
}

.accent-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.accent-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #2a2a2a;
  transition: 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  border-radius: 30px;
  border: 1px solid #333;
}

.accent-slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 3px;
  bottom: 3px;
  background-color: #555;
  transition: 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  border-radius: 50%;
}

input:checked + .accent-slider {
  background-color: #2ba618;
}

input:checked + .accent-slider:before {
  transform: translateX(22px);
  background-color: #fff;
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
  color: #eeeeee;
  margin: 0;
}

.smart-tile:active {
  transform: scale(0.94);
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
  transform: scale(1.1);
}

.title-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
