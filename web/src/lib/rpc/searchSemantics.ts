import { ArticleWithScore } from ".";

export default async function searchSemantics(query: string) {
    "use server";
    const url = process.env.BACKEND_URL as string;
    const response = await fetch(`${url}/semantics?text=${query}`, {
        method: "GET",
    });
    const result=await response.json() as ArticleWithScore[];
    const sortedResult=result.sort((a,b)=>b.Score-a.Score)
    return await sortedResult;
}