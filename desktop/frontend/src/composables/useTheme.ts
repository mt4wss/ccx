import { ref } from 'vue'

export type ThemeMode = 'light' | 'dark'

const STORAGE_KEY = 'ccx-desktop-theme'

const theme = ref<ThemeMode>('dark')

function applyTheme(mode: ThemeMode) {
  const html = document.documentElement
  if (mode === 'dark') {
    html.classList.add('dark')
  } else {
    html.classList.remove('dark')
  }
}

function init() {
  const stored = localStorage.getItem(STORAGE_KEY) as ThemeMode | null
  theme.value = stored === 'light' || stored === 'dark' ? stored : 'dark'
  applyTheme(theme.value)
}

function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  applyTheme(theme.value)
  localStorage.setItem(STORAGE_KEY, theme.value)
}

function setTheme(mode: ThemeMode) {
  theme.value = mode
  applyTheme(mode)
  localStorage.setItem(STORAGE_KEY, mode)
}

export function useTheme() {
  return { theme, init, toggleTheme, setTheme }
}
