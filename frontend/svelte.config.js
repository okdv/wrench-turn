import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: preprocess({
    postcss: true,
    typescript: true
  }),

	kit: {
		prerender: {
      entries: ['/join', '/login', '/', '/jobs/create', '/jobs', '/jobs/[id]']
    },
		adapter: adapter({
      pages: 'build', 
      assets: 'build', 
      fallback: 'index.html', 
      precompress: false, 
      struct: true 
    })
	}
};

export default config;
