import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { mande } from "mande";
import type { Color } from './colors';

import { useBrandsStore } from './brands';
import { useStoresStore } from './stores';


const multiSpoolAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/spools') : mande('/api/spools')
const singleSpoolAPI = import.meta.env.DEV ? mande('http://bespin:9766/api/spool') : mande('/api/spool')

export interface Spool {
    id: number | null
    location: number
    material: number
    brand: number
    store: number
    colors: number[]
    weight: number
    combined_weight: number
    current_weight: number
    price: number
    empty: boolean
    created_at: Date
    updated_at: Date
    deleted_at: Date | null
    emptied_at: Date | null
}

export interface NewSpool {
    id: number | null
    location: number
    material: number
    brand: number
    store: number
    colors: number[]
    weight: number
    combined_weight: number
    current_weight: number
    price: number
    empty: boolean
}

export interface SpoolBrandLink {
    brand_link: string
    brand_label: string
    store_link: string
    store_label: string
}

const brandsStore = useBrandsStore()
const storesStore = useStoresStore()


export const useSpoolsStore = defineStore('spools', () => {




    const sorted = ref<Spool[]>([])

    const lastCreatedID = ref<number | null>(null)

    // const editError = ref<string[]>([])

    async function find(sortBy: string = "label", sortDir: string = "asc") {
        if (sortDir === "desc") {
            sortBy = "-" + sortBy
        }

        await multiSpoolAPI.get<Array<Spool>>({ query: { "sort_by": sortBy } }).then((results: Spool[]) => {
            if (!showEmpty.value) {
                results = results.filter((s: Spool) => !s.empty)
            }

            sorted.value = results
        }).catch(err => {
            alert("find: " + err)
        })

        return sorted.value
    }

    async function update(id: number, record: Spool) {
        await singleSpoolAPI.patch<Spool>(id, record).then((result: Spool) => {
            const idx = indexOfID(result.id as number)
            sorted.value[idx] = result
        })

    }

    async function create(record: NewSpool) {
        multiSpoolAPI.post<Spool>(record).then((result: Spool) => {
            sorted.value.push(result)
            lastCreatedID.value = result.id
        })
    }

    async function kill(spoolID: number) {
        singleSpoolAPI.delete(spoolID).then(resp => {
            const idx = indexOfID(spoolID)
            sorted.value.splice(idx, 1)
            // }).catch(err => {
            //     addNewStudentErrors.value = err.body
        })
    }

    const count = computed(() => {
        return sorted.value.length
    })

    function findByID(id: number): Spool | undefined {
        const found = sorted.value.find((b: { id: number | null }) => b.id === id)
        if (found === null || found === undefined) {
            return {
                id: null,
                location: 0,
                material: 0,
                brand: 0,
                store: 0,
                colors: [],
                weight: 0,
                combined_weight: 0,
                current_weight: 0,
                price: 0,
                empty: false,
                created_at: new Date(),
                updated_at: new Date(),
                deleted_at: null,
                emptied_at: null,
            }
        }
        return {
            id: found.id,
            location: found.location,
            material: found.material,
            brand: found.brand,
            store: found.store,
            colors: found.colors,
            weight: found.weight,
            combined_weight: found.combined_weight,
            current_weight: found.current_weight,
            price: found.price,
            empty: found.empty,
            created_at: found.created_at,
            updated_at: found.updated_at,
            deleted_at: found.deleted_at,
            emptied_at: found.emptied_at,
        }
    }

    function indexOfID(id: number): number {
        return sorted.value.findIndex((b: { id: number | null }) => b.id === id)
    }

    const showEmpty = computed(() => {
        return empty.value
    })


    const empty = ref(false)

    function toggleEmpty() {
        empty.value = !empty.value
    }

    function storeBrandLinks(id: number): SpoolBrandLink {
        const spool = findByID(id)
        const brand = brandsStore.findByID(spool?.brand as number)
        const store = storesStore.findByID(spool?.store as number)

        const result: SpoolBrandLink = {
            brand_link: "",
            brand_label: "",
            store_link: "skip",
            store_label: "skip",
        }

        if (brand?.store_id !== null && brand?.store_id !== undefined) {
            const storeForBrand = storesStore.findByID(brand?.store_id as number)
            if (storeForBrand?.url !== undefined && storeForBrand?.url !== null) {
                result.brand_link = storeForBrand?.url
                result.brand_label = brand.label
            }
        }

        if (store?.url !== null && store?.url !== undefined && store.url !== result.brand_link) {
            result.store_link = store.url
            result.store_label = store.label
        }

        return result
    }



    return {
        sorted,
        lastCreatedID,
        count,
        find,
        update,
        create,
        kill,

        showEmpty,
        toggleEmpty,
        findByID,
        indexOfID,
        storeBrandLinks
    }

})
