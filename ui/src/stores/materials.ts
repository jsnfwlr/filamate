import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {mande} from "mande";

const multiMaterialAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/materials') : mande('/api/materials')
const singleMaterialAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/material') : mande('/api/material')

export interface Material {
  id: number | null
  label: string
  description: string
  class: string
  special: boolean
}



export const useMaterialsStore = defineStore('materials', () => {

    const sorted = ref<Material[]>([])
    // const editError = ref<string[]>([])

    async function find(sortBy: string = "label", sortDir: string = "asc") {
        if (sortDir === "desc") {
            sortBy = "-" + sortBy
        }

        await multiMaterialAPI.get<Array<Material>>({ query: { "sort_by": sortBy } }).then((results: Material[]) => {
            sorted.value = results
        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }

    async function update(id: number, record: Material) {
        await singleMaterialAPI.patch<Material>(id, record).then((result: Material) => {
            const idx = indexOfID(result.id as number)
            sorted.value[idx] = result
        }).catch(err => {
            alert("find: " + err)
        })

    }

    async function create(record: Material) {
        multiMaterialAPI.post<Material>(record).then( (result: Material) => {
            sorted.value.push( result)
        }).catch(err => {
            alert("find: " + err)
        })
    }

    async function kill(materialID: number) {
        singleMaterialAPI.delete(materialID).then(resp => {
            const idx = indexOfID(materialID)
            sorted.value.splice(idx, 1)
        }).catch(err => {
            alert("find: " + err)
        })
    }

    const count = computed( () => {
        return sorted.value.length
    })

    function findByID(id: number): Material | undefined {
        const found = sorted.value.find((b: { id: number | null }) => b.id === id)
        if (found === null || found === undefined) {
            return {
                id: null,
                label: "",
                description: "",
                class: "",
                special: false,
            }
        }
        return {
            id: found.id,
            label: found.label,
            description: found.description,
            class: found.class,
            special: found.special,
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
