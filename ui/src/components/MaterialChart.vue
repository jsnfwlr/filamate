<script lang="ts" setup>
import { computed } from 'vue'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
import type { ChartData, ChartOptions, TooltipItem } from 'chart.js'
import { Pie } from 'vue-chartjs'

ChartJS.register(ArcElement, Tooltip, Legend)

const props = defineProps<{
    chartData: ChartData<'pie'>,
}>()

ChartJS.defaults.color = '#ffffff'
ChartJS.defaults.backgroundColor = '#222222'
ChartJS.defaults.borderColor = '#ffffff47'

const options = computed<ChartOptions<'pie'>>(() => {
    const options: ChartOptions<'pie'> = {
        responsive: true,
        plugins: {
            legend: {
                display: false,
            },
            tooltip: {
                callbacks: {
                    title: function (context: TooltipItem<'pie'>[]): string {
                        if (context[0] === undefined) {
                            return ''
                        }

                        let datasetColors = context[0].chart.data.datasets.map(function (e) {
                            return e.backgroundColor;
                        });
                        datasetColors = datasetColors.flat();

                        let labelIndex = 0;

                        if (context[0].dataset === undefined) {
                            return ''
                        }

                        if (context[0].dataset.backgroundColor instanceof Array) {
                            const li = datasetColors.indexOf(context[0].dataset.backgroundColor[context[0].dataIndex]);
                            if (li !== -1) {
                                labelIndex = li;
                            }
                        }

                        return (context[0].chart.data.labels === undefined ? '' : context[0].chart.data.labels[labelIndex]) as string;

                    }
                }
            }
        }
    }

    return options
})

</script>

<template>
    <div class="q-dark">
        <Pie :data="props.chartData" :options="options" />
    </div>
</template>

<style scoped>
div {
    /* Ensure the chart takes up the full width of its container */
    width: 100%;
    border: 1px solid #ffffff47;
    border-radius: 4px;
    padding: 16px;
    background: #222222;

    canvas {
        border-radius: 4px;
    }
}
</style>
