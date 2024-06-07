import { Title } from "@solidjs/meta";
import { createAsync,cache } from "@solidjs/router";
import { onMount } from "solid-js";
import Counter from "~/components/Counter";


const getSample = async () => {
  "use server"
  const url = import.meta.env.VITE_BACKEND_URL as string
  const res = await fetch(`${url}/articles`)
  const articles = await res.json() as string
  return import.meta.env.VITE_BACKEND_URL as string
}
export default function Home() {
  const sample = createAsync(() => getSample())
  onMount(() => {
    console.log(sample())
  })
  return (
    <main>
      <div>
        {sample()}
      </div>
    </main>
  );
}
