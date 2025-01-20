<script setup lang="ts">
</script>

<template>
  <div class="box-border p-10">
    <p class="font-bold text-center text-3xl my-5">
      Exam marks gathering
    </p>
    <p class="font-normal text-center text-2xl">Made with love</p>
  </div>
  <div class="flex flex-1 flex-col items-center">
    <p class="my-2">Powered by</p>
    <div
      class="self-center justify-center items-center flex md:flex-auto md:flex-row flex-col"
    >
      <div
        class="box-content h-32 w-64 border-4 rounded-2xl my-2 mx-2 flex flex-row drop-shadow-2xl hover:drop-shadow-md"
      >
        <img class="box-content size-1/2 m-auto" src="../assets/vue.svg" />
        <div class="flex flex-col m-auto justify-center">
          <p class="font-bold text-xl mx-2">Vue3</p>
          <p class="font-normal text-xs mx-2">
            The Progressive JavaScript Framework
          </p>
        </div>
      </div>
      <div
        class="box-content h-32 w-64 border-4 rounded-2xl my-2 mx-2 flex flex-row drop-shadow-2xl hover:drop-shadow-md"
      >
        <img class="box-content size-1/2 m-auto" src="../assets/vite.svg" />
        <div class="flex flex-col m-auto justify-center">
          <p class="font-bold text-xl mx-2">Vite</p>
          <p class="font-normal text-xs mx-2">
            Next Generation Frontend Tooling
          </p>
        </div>
      </div>
    </div>
    <div
      class="self-center justify-center items-center flex md:flex-auto md:flex-row flex-col"
    >
      <div
        class="box-content h-32 w-64 border-4 rounded-2xl my-2 mx-2 flex flex-row drop-shadow-2xl hover:drop-shadow-md"
      >
        <img
          class="box-content p-4 m-auto size-2/6"
          src="../assets/tailwindcss.svg"
        />
        <div class="flex flex-col mx-auto justify-center">
          <p class="font-bold text-xl mx-2">TailwindCSS</p>
          <p class="font-normal text-xs mx-2">A utility-first CSS framework</p>
        </div>
      </div>
      <div
        class="box-content h-32 w-64 border-4 rounded-2xl my-2 mx-2 flex flex-row drop-shadow-2xl hover:drop-shadow-md"
      >
        <img class="box-content size-1/2  m-auto" src="../assets/vue.svg" />
        <div class="flex flex-col mx-auto justify-center">
          <p class="font-bold text-xl mx-2">VueRouter</p>
          <p class="font-normal text-xs mx-2">
            The official Router for Vue.js.
          </p>
        </div>
      </div>
    </div>
    <p class="my-4">Licensed under</p>
    <div
      class="box-content md:h-32 md:w-72 rounded-2xl mx-auto flex flex-col items-center justify-center"
    >
      <img class="size-1/2" src="../assets/MIT_logo.svg" />
      <p class="text-xl font-bold">MIT Licence</p>
    </div>
  </div>
</template>

<style scoped></style>