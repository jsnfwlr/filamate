<script setup lang="ts">
import { RouterLink, RouterView, useRoute } from 'vue-router'

import { ref } from 'vue'

const route = useRoute()
</script>

<template>
  <q-layout view="hHr lpr fFr">
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <RouterLink text-white to="/">
          <q-btn text-white dense flat icon="mdi-view-dashboard" label="Dashboard" />
        </RouterLink>
        <q-separator dark vertical class="q-mx-sm" />
        <RouterLink text-white to="/spools">
          <q-btn dense flat icon="mdi-printer-3d-nozzle" label="Spools" />
        </RouterLink>
        <RouterLink text-white to="/brands">
          <q-btn dense flat icon="mdi-copyright" label="Brands" />
        </RouterLink>
        <RouterLink text-white to="/colors">
          <q-btn dense flat icon="mdi-palette" label="Colors" />
        </RouterLink>
        <RouterLink text-white to="/locations">
          <q-btn dense flat icon="mdi-crosshairs-gps" label="Locations" />
        </RouterLink>
        <RouterLink text-white to="/materials">
          <q-btn dense flat icon="mdi-movie-roll" label="Materials" />
        </RouterLink>
        <RouterLink text-white to="/stores">
          <q-btn dense flat icon="mdi-cart" label="Stores" />
        </RouterLink>
        <q-space />
        <q-toolbar-title shrink>
          <RouterLink to="/">
            <span class="q-pr-sm">Filamate</span>
            <q-avatar><img src="/logo.png" /></q-avatar>
          </RouterLink>
        </q-toolbar-title>
      </q-toolbar>
    </q-header>

    <q-page-container>
      <RouterView />
    </q-page-container>
  </q-layout>
</template>

<style scoped>
.q-page-container {
  padding-top: calc(50px + 1.5rem) !important;
  min-height: 100vh;
}

.q-header a {
  color: white;
  text-decoration: none;
}
</style>
