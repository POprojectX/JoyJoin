<script setup lang="ts">
import { computed, ref } from 'vue'

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')

const showPw = ref(false)
const showConfirm = ref(false)

const pwType = computed(() => (showPw.value ? 'text' : 'password'))
const confirmType = computed(() => (showConfirm.value ? 'text' : 'password'))

const uiInput = {
  base:
    'w-full rounded-2xl bg-white/90 border border-slate-200 text-slate-900 placeholder:text-slate-400 ring-0 focus:ring-0 focus-visible:ring-2 focus-visible:ring-violet-300 focus-visible:ring-offset-0'
}

const socialBtnClass =
  'w-full justify-center rounded-2xl border border-slate-200 bg-white/70 hover:bg-violet-50 active:bg-violet-100 text-violet-500'

function onSubmit() {
  console.log({
    name: name.value,
    email: email.value,
    password: password.value,
    confirmPassword: confirmPassword.value
  })
}
</script>

<template>
  <div class="min-h-screen relative overflow-hidden">
    <div class="absolute inset-0 bg-cover bg-center" style="background-image:url('/images/login-bg.jpg')"></div>

    <div class="relative min-h-screen grid place-items-center p-5">
      <div class="w-full max-w-sm">
        <div class="flex items-center justify-center gap-2 mb-5">
          <span class="h-6 w-6 rounded-full bg-linear-to-br from-violet-300 to-sky-300"></span>
          <span class="text-lg font-semibold text-slate-800">JoyJoin</span>
        </div>

        <div class="rounded-3xl bg-white/70 backdrop-blur-xl border border-white/70 shadow-xl p-6">
          <h1 class="text-2xl font-semibold text-center text-slate-900">Sign up</h1>
          <p class="text-sm text-center text-slate-600 mt-1">
            Create an account to plan your perfect event
          </p>

          <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
            <div>
              <label class="block text-sm text-slate-700 mb-1">Name</label>
              <UInput
                v-model="name"
                type="text"
                placeholder="Your name"
                size="lg"
                class="w-full"
                :ui="uiInput"
              />
            </div>

            <div>
              <label class="block text-sm text-slate-700 mb-1">Email</label>
              <UInput
                v-model="email"
                type="email"
                placeholder="Your email"
                size="lg"
                class="w-full"
                :ui="uiInput"
              >
                <template #trailing>
                  <UIcon name="i-heroicons-envelope" class="text-slate-400" />
                </template>
              </UInput>
            </div>

            <div>
              <label class="block text-sm text-slate-700 mb-1">Password</label>
              <UInput
                v-model="password"
                :type="pwType"
                placeholder="Password"
                size="lg"
                class="w-full"
                :ui="uiInput"
              >
                <template #trailing>
                  <button type="button" class="p-1" @click="showPw = !showPw" aria-label="Toggle password">
                    <UIcon :name="showPw ? 'i-heroicons-eye-slash' : 'i-heroicons-eye'" class="text-slate-400" />
                  </button>
                </template>
              </UInput>
            </div>

            <div>
              <label class="block text-sm text-slate-700 mb-1">Confirm password</label>
              <UInput
                v-model="confirmPassword"
                :type="confirmType"
                placeholder="Confirm password"
                size="lg"
                class="w-full"
                :ui="uiInput"
              >
                <template #trailing>
                  <button type="button" class="p-1" @click="showConfirm = !showConfirm" aria-label="Toggle confirm password">
                    <UIcon :name="showConfirm ? 'i-heroicons-eye-slash' : 'i-heroicons-eye'" class="text-slate-400" />
                  </button>
                </template>
              </UInput>
            </div>

            <UButton
              type="submit"
              size="xl"
              class="w-full rounded-2xl font-semibold text-slate-900 justify-center"
              :ui="{ base: 'justify-center' }"
              style="background:linear-gradient(90deg,#c4b5fd 0%, #bae6fd 100%); border:none;"
            >
              Sign up
            </UButton>

            <div class="text-center">
              <div class="text-xs text-slate-600">Already have an account?</div>
              <NuxtLink to="/login" class="text-sm font-medium text-violet-500 hover:underline">Sign in</NuxtLink>
            </div>

            <div class="space-y-3 pt-2">
              <UButton variant="soft" size="lg" :class="socialBtnClass">
                <template #leading>
                  <UIcon name="i-simple-icons-google" class="text-violet-500" />
                </template>
                Sign up with Google
              </UButton>

              <UButton variant="soft" size="lg" :class="socialBtnClass">
                <template #leading>
                  <UIcon name="i-simple-icons-apple" class="text-violet-500" />
                </template>
                Sign up with Apple
              </UButton>
            </div>

            <p class="text-[11px] text-center text-slate-500 leading-4 pt-1">
              By signing up, you agree to our
              <NuxtLink to="#" class="text-violet-500 hover:underline">Terms of Service</NuxtLink>
              and
              <NuxtLink to="#" class="text-violet-500 hover:underline">Privacy Policy</NuxtLink>.
            </p>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>
