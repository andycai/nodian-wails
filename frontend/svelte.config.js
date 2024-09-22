import { vitePreprocess } from '@sveltejs/vite-plugin-svelte'
import preprocess from 'svelte-preprocess'

export default {
  // Consult https://github.com/sveltejs/svelte-preprocess
  // for more information about preprocessors
  preprocess: [
    preprocess({
      postcss: true,
    }),
    vitePreprocess(),
  ],
}
