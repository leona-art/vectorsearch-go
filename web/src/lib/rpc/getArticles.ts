import type { Article } from "."
export default async function getArticles() {
    "use server"
    const url = process.env.BACKEND_URL as string
    const response = await fetch(`${url}/articles`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    })
    
    return await response.json() as Article[];
}