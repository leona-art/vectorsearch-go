export interface Article {
    id: string;
    title: string;
    content: string;
}

export interface ArticleWithScore{
    Article: Article;
    Score: number;
}