<script>
export default {
  props: ["deviceId", "deviceName"],
};
</script>

<template>
  <div class="card-container">
    <h1 class="name-h">{{ deviceName }}</h1>
    <p class="value-h">Temperature: 23C</p>
    <p class="value-h">Humidity: 50%</p>
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
