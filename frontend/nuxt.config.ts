// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ["@nuxt/ui"],
  css: ["~/assets/css/main.css"],
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },

  typescript: {
    strict: true,
    typeCheck: false,
  },
  colorMode: {
    classSuffix: "", // uses 'dark' class, not 'dark-mode'
    preference: "light",
    fallback: "light",
    storageKey: "joyjoin-theme",
  },
});
