<script lang="ts" setup>
import { ref, computed } from 'vue'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
import type { ChartData, ChartOptions, TooltipItem } from 'chart.js'
import { Pie } from 'vue-chartjs'

ChartJS.register(ArcElement, Tooltip, Legend)

const props = defineProps<{
    chartData: ChartData<'pie'>,
}>()

const colors = {
    red: 'rgb(255, 99, 132)',
    orange: 'rgb(255, 159, 64)',
    yellow: 'rgb(255, 205, 86)',
    green: 'rgb(75, 192, 192)',
    blue: 'rgb(54, 162, 235)',
    purple: 'rgb(153, 102, 255)',
    grey: 'rgb(201, 203, 207)'
};


ChartJS.defaults.color = '#ffffff'
ChartJS.defaults.backgroundColor = '#222222'
ChartJS.defaults.borderColor = '#ffffff47'

const options = computed<ChartOptions<'pie'>>(() => {
    const options: ChartOptions<'pie'> = {
        responsive: true,
        plugins: {
            legend: {
                display: false,
                /*
                labels: {
                    generateLabels: function (chart) {
                        // Get the default label list
                        const original = ChartJS.overrides.pie.plugins.legend.labels.generateLabels;
                        const labelsOriginal = original.call(this, chart);

                        // Build an array of colors used in the datasets of the chart
                        let datasetColors = chart.data.datasets.map(function (e) {
                            return e.backgroundColor;
                        });
                        datasetColors = datasetColors.flat();

                        // loop over the datasets, using the color index to figure out which dataset the label comes from
                        // set the datasetIndex of the nth label when the nth color is found within the dataset making sure
                        // the label is never reused across datasets. Also make sure the hidden state of the label is updated
                        // to match the dataset hidden state, and the fillStyle is updated to match the dataset backgroundColor
                        for (let i = 0; i < chart.data.datasets.length; i++) {
                            const dataset = chart.data.datasets[i];
                            if (dataset === undefined) continue;

                            if (dataset.backgroundColor instanceof Array) {
                                for (let j = 0; j < dataset.backgroundColor.length; j++) {
                                    const bgColor = dataset.backgroundColor[j];
                                    const colorIndex = datasetColors.indexOf(bgColor);
                                    if (colorIndex !== -1) {
                                        if (labelsOriginal[colorIndex] !== undefined) {
                                            labelsOriginal[colorIndex].datasetIndex = i;
                                            labelsOriginal[colorIndex].hidden = !chart.isDatasetVisible(i);
                                            labelsOriginal[colorIndex].fillStyle = bgColor as Color;
                                            labelsOriginal[colorIndex].text = chart.data.labels ? chart.data.labels[colorIndex] + " (" + dataset.data[j] + ")" : '';

                                        }
                                    }
                                }
                            }
                        }

                        return labelsOriginal;
                    }
                },
                onClick: function (mouseEvent, legendItem, legend) {
                    if (legendItem.datasetIndex !== undefined) {
                        // toggle the visibility of the dataset from what it currently is
                        legend.chart.getDatasetMeta(
                            legendItem.datasetIndex
                        ).hidden = legend.chart.isDatasetVisible(legendItem.datasetIndex);
                        legend.chart.update();
                    }
                }
                */
            },
            tooltip: {
                callbacks: {
                    title: function (context: TooltipItem<'pie'>[]): any {
                        if (context[0] === undefined) {
                            return
                        }





                        let datasetColors = context[0].chart.data.datasets.map(function (e) {
                            return e.backgroundColor;
                        });
                        datasetColors = datasetColors.flat();

                        let labelIndex = 0;

                        if (context[0].dataset === undefined) {
                            return
                        }

                        if (context[0].dataset.backgroundColor instanceof Array) {
                            const li = datasetColors.indexOf(context[0].dataset.backgroundColor[context[0].dataIndex]);
                            if (li !== -1) {
                                labelIndex = li;
                            }
                        }


                        return (context[0].chart.data.labels === undefined ? '' : context[0].chart.data.labels[labelIndex]);

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
        <Pie :data="chartData" :options="options" />
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
