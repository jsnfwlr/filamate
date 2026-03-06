<script lang="ts" setup>
import { computed } from 'vue'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, Filler, LinearScale } from 'chart.js'
import type { ChartData, ChartOptions, TooltipItem } from 'chart.js'
import { Bar } from 'vue-chartjs'

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Filler, Legend)


const props = defineProps<{
    labels: string[],
    added: Array<number | [number, number] | null>,
    emptied: Array<number | [number, number] | null>,
    stored: Array<number | [number, number] | null>
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


const options = computed<ChartOptions<'bar'>>(() => {
    return {
        interaction: {
            mode: 'index'
        },
        responsive: true,
        plugins: {
            legend: {
                position: 'bottom',
                labels: {
                    color: '#ffffff',
                }

            },


            tooltip: {
                mode: 'index',
                intersect: true,
                displayColors: true,
                callbacks: {
                    beforeBody: function (tooltipItems: TooltipItem<'bar'>[]) {
                        tooltipItems.forEach(function (tooltipItem: TooltipItem<'bar'>) {
                            if (tooltipItem.parsed.y !== null && tooltipItem.parsed.y < 0) {
                                tooltipItem.formattedValue = String(tooltipItem.parsed.y * -1);
                            }
                        });
                    }
                }
            },
        },
        scales: {
            x: {
                stacked: false,

                title: {
                    display: true,
                    text: 'Date'
                }
            },
            y: {
                stacked: false,
                title: {
                    display: true,
                    text: 'Stored spools'
                },
                suggestedMin: 0,
                suggestedMax: props.stored != undefined && props.stored.length > 0 ? Math.max(...props.stored.map(value => Array.isArray(value) ? value[1] : value) as number[]) + 10 : 0
            }
        },
    }
})

const stock = computed(() => {
    const stock: number[] = []

     props.stored.map((value, index) => {
        const emptied = (props.emptied[index] !== null ? (Array.isArray(props.emptied[index]) ? props.emptied[index][1] : props.emptied[index]) : 0) as number
        const added = (props.added[index] !== null ? (Array.isArray(props.added[index]) ? props.added[index][1] : props.added[index]) : 0) as number
        const stored = (props.stored[index] !== null ? (Array.isArray(props.stored[index]) ? props.stored[index][1] : props.stored[index]) : 0) as number
        stock.push(stored + emptied + added)
    })
    return stock
})

const chartData = computed<ChartData<'bar'>>(() => {
    const data: ChartData<'bar'> = {

        labels: props.labels,
        datasets: [
            {
                label: 'Emptied',
                backgroundColor: colors.green,
                data: props.emptied,
                order: 0,
                // stack: 'Stack 1'

            },
            {
                label: 'Added',
                backgroundColor: colors.red,
                data: props.added,
                order: 1,
                // stack: 'Stack 2'
            },
            {
                label: 'Stored',
                backgroundColor: colors.blue,
                data: props.stored,
                order: 2,
                // stack: 'Stack 3'
            },

            {
                label: 'Stock',
                backgroundColor: colors.purple,
                data: stock.value,
                order: 3,
                // stack: 'Stack 4'
            }

        ]
    }

    return data
})
</script>

<template>
    <div class="q-dark">
        <Bar :data="chartData" :options="options" />
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
