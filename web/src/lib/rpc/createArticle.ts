export default async function createArticle(title: string, content: string) {
    "use server"
    const url = process.env.BACKEND_URL as string
    const response = await fetch(`${url}/article/1`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ title, content }),
    })
    
    console.log(response)
    return await response.json();
}
