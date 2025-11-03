<script>
export default {
  props: ["deviceName", "ip", "visibility", "deviceStatus"],
};
</script>

<template>
  <div class="device-container">
    <p>{{ deviceName }}</p>
    <p>{{ ip }}</p>
    <p>{{ visibility }}</p>
    <p>{{ deviceStatus }}</p>
    <button>View</button>
  </div>
</template>

<style scoped>
.device-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
}
</style>
