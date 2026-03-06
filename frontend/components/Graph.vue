<script>
import { Line } from "vue-chartjs";
import axios from "axios";
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  CategoryScale,
} from "chart.js";

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  LinearScale,
  PointElement,
  CategoryScale,
);

export default {
  props: ["filename", "param"],
  components: { Line },
  data() {
    return {
      values: [],
      num: 10,
      steps: 3,
      lag: 7,
      is_loading: false,
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: { legend: { display: false } },
        scales: {
          x: {
            grid: { color: "#222" },
            ticks: { color: "#999", font: { size: 10 } },
          },
          y: {
            grid: { color: "#222" },
            ticks: { color: "#999", font: { size: 10 } },
          },
        },
      },
    };
  },
  computed: {
    chartData() {
      return {
        labels: this.values.map((item) => item.timestamp),
        datasets: [
          {
            label: this.param,
            backgroundColor: "#eeeeee",
            borderColor: "#eeeeee",
            data: this.values.map((item) => item[this.param]),
            tension: 0.3,
            pointRadius: 2,
            borderWidth: 2,
            segment: {
              borderColor: (ctx) =>
                this.values[ctx.p1DataIndex]?.isForecast
                  ? "#42b883"
                  : "#eeeeee",

              borderDash: (ctx) =>
                this.values[ctx.p1DataIndex]?.isForecast ? [5, 5] : undefined,
            },
            pointBackgroundColor: (ctx) =>
              this.values[ctx.dataIndex]?.isForecast ? "#42b883" : "#eeeeee",
          },
        ],
      };
    },
  },
  methods: {
    // Change this
    async handleFetchValues(path) {
      try {
        const res = await axios.get(
          `http://localhost:8000/${path}/?filename=${this.filename}&num=${this.num}&param=${this.param}`,
        );
        this.values = res.data.values;
      } catch (err) {
        console.log(err);
      }
    },
    async handleForecast(path) {
      try {
        this.is_loading = true;
        const res = await axios.get(
          `http://localhost:8000/predict/${path}/?filename=${this.filename}&lag=${this.lag}&steps=${this.steps}&param=${this.param}`,
        );
        this.is_loading = false;
        const forecasted = res.data.values;

        if (this.values.length > 1 && forecasted.length > 0) {
          const last_entry = this.values[this.values.length - 1];
          const second_last_entry = this.values[this.values.length - 2];

          const last_time = new Date(last_entry.timestamp);
          const prev_time = new Date(second_last_entry.timestamp);

          const interval = last_time.getTime() - prev_time.getTime();

          const new_values = forecasted.map((val, index) => {
            const next_time = new Date(
              last_time.getTime() + interval * (index + 1),
            );

            const year = next_time.getFullYear();
            const month = String(next_time.getMonth() + 1).padStart(2, "0");
            const day = String(next_time.getDate()).padStart(2, "0");
            const hours = String(next_time.getHours()).padStart(2, "0");
            const minutes = String(next_time.getMinutes()).padStart(2, "0");
            const seconds = String(next_time.getSeconds()).padStart(2, "0");

            return {
              timestamp: `${year}-${month}-${day}T${hours}:${minutes}:${seconds}`,
              [this.param]: val,
              isForecast: true,
            };
          });

          this.values = [...this.values, ...new_values];
        }
      } catch (err) {
        console.log(err);
      }
    },
  },
  async mounted() {
    await this.handleFetchValues("tail");
  },
};
</script>

<template>
  <div class="chart-card">
    <div class="title-container">
      <h2 class="chart-title">{{ param }}</h2>
      <button @click="this.handleFetchValues('tail')" class="title-button">
        Last values
      </button>
      <button @click="this.handleFetchValues('hour')" class="title-button">
        Last hours
      </button>
      <button @click="this.handleFetchValues('day')" class="title-button">
        Last days
      </button>
      <div class="field">
        <label for="value">Number of values</label>
        <input type="number" id="value" v-model="this.num" />
      </div>
    </div>

    <div class="chart-container">
      <Line :data="chartData" :options="chartOptions" />
    </div>

    <div class="bottom-container">
      <button @click="this.handleForecast('values')" class="title-button">
        Forecast
      </button>
      <div class="field">
        <label for="value">Historical Window</label>
        <input type="number" id="value" v-model="this.lag" />
      </div>
      <div class="field">
        <label for="value">Prediction Horizon</label>
        <input type="number" id="value" v-model="this.steps" />
      </div>
      <span v-if="is_loading" class="spinner"></span>
    </div>
  </div>
</template>

<style scoped>
.chart-card {
  background: #1a1a1a;
  border: 1px solid #333;
  padding: 0.75rem;
  width: 50%;
}

.title-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
  align-items: baseline;
}

.chart-title {
  color: #888;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.05rem;
  margin: 0 0 0.5rem 0;
}

.bottom-container {
  display: flex;
  flex-direction: row;
  gap: 1rem;
  align-items: center;
  padding: 1rem;
}

.chart-container {
  position: relative;

  height: 300px;
}

.title-button {
  margin-bottom: 1rem;
  border: 1px solid #eeeeee;
  background: transparent;
  color: #eeeeee;
  padding: 0.5rem 0.75rem;
  border-radius: 0.4rem;
  cursor: pointer;
  transition: all 200ms ease;
}

.title-button:hover {
  color: #8c8c8c;
  border-color: #8c8c8c;
}

input,
select {
  color: #eeeeee;
  background-color: #252525;
  border: none;
  outline: none;
  padding: 0.5rem;
  font-size: 0.95rem;
  border-radius: 0.35rem;
  width: 5rem;
}

.field {
  display: flex;
  flex-direction: row;
  gap: 0.4rem;
  align-items: center;
}

label {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #b0b0b0;
}

@media (max-width: 600px) {
  .chart-container {
    height: 150px;
  }
  .chart-card {
    width: 90%;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.spinner {
  width: 18px;
  height: 18px;
  border: 2px solid #333;
  border-top-color: #eeeeee;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  display: inline-block;
}
</style>
