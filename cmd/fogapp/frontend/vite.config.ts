import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { resolve } from "path";

export default defineConfig({
    plugins: [tailwindcss(), svelte()],
    resolve: {
        alias: {
            $lib: resolve(__dirname, "./src/lib"),
        },
    },
    server: {
        port: 5173,
        strictPort: false,
    },
    build: {
        outDir: "dist",
        emptyOutDir: true,
        sourcemap: false,
    },
});
