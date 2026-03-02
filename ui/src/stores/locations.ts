import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { mande } from "mande";

const multiLocationAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/locations') : mande('/api/locations')
const singleLocationAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/location') : mande('/api/location')

export interface Location {
    id: number | null
    label: string
    description: string
    printable: boolean
    tally: boolean
    capacity: number
}

export const useLocationsStore = defineStore('locations', () => {

    const sorted = ref<Location[]>([])
    // const editError = ref<string[]>([])

    async function find(sortBy: string = "label", sortDir: string = "asc") {
        if (sortDir === "desc") {
            sortBy = "-" + sortBy
        }

        await multiLocationAPI.get<Array<Location>>({ query: { "sort_by": sortBy } }).then((results: Location[]) => {
            sorted.value = results
        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }

    async function update(id: number, record: Location) {

        record.capacity = record.capacity / 1

        alert("update: " + JSON.stringify(record))
        await singleLocationAPI.patch<Location>(id, record).then((result: Location) => {
            const idx = indexOfID(result.id as number)
            sorted.value[idx] = result
        }).catch(err => {
            alert("find: " + err)
        })

    }

    async function create(record: Location) {
        record.capacity = record.capacity / 1
        alert("create: " + JSON.stringify(record))
        multiLocationAPI.post<Location>(record).then((result: Location) => {
            sorted.value.push(result)
        }).catch(err => {
            alert("find: " + err)
        })
    }

    async function kill(locationID: number) {
        singleLocationAPI.delete(locationID).then(() => {
            const idx = indexOfID(locationID)
            sorted.value.splice(idx, 1)
        }).catch(err => {
            alert("find: " + err)
        })
    }

    const count = computed(() => {
        return sorted.value.length
    })

    function findByID(id: number): Location {
        const found = sorted.value.find((b: { id: number | null }) => b.id === id)
        if (found === null || found === undefined) {
            return {
                id: null,
                label: "",
                description: "",
                printable: false,
                tally: false,
                capacity: 0,
            }
        }
        return {
            id: found.id,
            label: found.label,
            description: found.description,
            printable: found.printable,
            tally: found.tally,
            capacity: found.capacity,
        }
    }

    function indexOfID(id: number): number {
        return sorted.value.findIndex((b: { id: number | null }) => b.id === id)
    }

    return {
        sorted,
        count,
        find,
        update,
        create,
        kill,

        findByID,
        indexOfID,
    }

})
