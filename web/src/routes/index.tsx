import { Title } from "@solidjs/meta";
import { createAsync, redirect,revalidate, action } from "@solidjs/router";
import { createEffect, onMount } from "solid-js";
import Counter from "~/components/Counter";
import { Button } from "~/components/ui/button";


const getArticles = async () => {
  "use server"
  const url = import.meta.env.VITE_BACKEND_URL as string
  const res = await fetch(`${url}/articles`)
  const articles = await res.json() as string
  return articles
}
const postArticleAction = action(async (title: string, content: string) => {
  "use server"
  console.log(title, content)
  const url = import.meta.env.VITE_BACKEND_URL as string
  // const title = formData.get("title") as string
  // const content = formData.get("content") as string
  const res = await fetch(
    `${url}/article/1`,
    {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title, content }),
    }
  )
  const article = await res.json()
  
  throw redirect("/")
})
export default function Home() {
  const sample = createAsync(() => getArticles())
  createEffect(() => {
    console.log(sample())
  })
  return (
    <main>
      <div>
        <form action={postArticleAction.with("title","content")} method="post">
          <Button type="submit">Submit</Button>
        </form>
      </div>
    </main>
  );
}
