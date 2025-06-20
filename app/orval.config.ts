 import { defineConfig } from 'orval';

  export default defineConfig({
    'flashcard-api': {
        input: {
            target: '../api-schema/api.yaml'
        },
        output: {
            // mode: 'split',
            mode: 'single',
            client: 'fetch',
            workspace: 'src/orval-client/',
            target: './api-client.ts'
        },
    },

 });